// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/envygeeks/tweedy/tweedy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
		Use:              "tweedy",
		Short:            "Delete your Tweets",
		PersistentPreRun: cmdPre,
		Run:              cmdRun,
	}
)

func init() {
	Cmd.Flags().Int("from", 0, "Delete Tweets after")
	Cmd.Flags().String("user", "", "User handle, or UID")
	Cmd.PersistentFlags().Bool("test", false, "Don't be destructive")
	Cmd.PersistentFlags().Bool("verbose", false, "Verbose output")
	Cmd.PersistentFlags().Bool("debug", false, "Debug output")
	Cmd.Flags().String("file", "", "Tweets JSON file")
	Cmd.Flags().Int("keep", 3, "Days to keep")

	Cmd.MarkFlagRequired("user")
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})
}

// pre does a setup on global stuff like loggers.
func cmdPre(cmd *cobra.Command, _ []string) {
	logrus.SetLevel(logrus.WarnLevel)
	if t, err := cmd.Flags().GetBool("verbose"); err == nil && t {
		logrus.SetLevel(logrus.DebugLevel)
		if t, err := cmd.Flags().GetBool("debug"); err == nil && t {
			logrus.SetReportCaller(true)
		}
	}
}

// from gets the from, w/out err
func from(cmd *cobra.Command) int64 {
	if i, err := cmd.Flags().GetInt("from"); err == nil {
		return int64(i)
	}

	return 0
}

// negative inverts the number so that we can
// use it to walk back the time object, to compare
// the Tweets date, to the days ago.
func keep(cmd *cobra.Command) int {
	// We don't care if there is an error tbqf, just ditch
	if i, err := cmd.Flags().GetInt("keep"); err == nil {
		if i > 0 {
			return -i
		}
	}

	return -3
}

// absPath expands, and checks the path
func file(cmd *cobra.Command) string {
	f, err := cmd.Flags().GetString("file")
	if err != nil {
		logrus.Fatalln(err)
	} else {
		if f == "" {
			return ""
		}
	}

	abs, err := filepath.Abs(f)
	if err == nil {
		if _, err = os.Stat(abs); os.IsNotExist(err) {
			msg := "unable to find the file %q"
			logrus.Fatalln(msg)
		}
	}

	return abs
}

func api(cmd *cobra.Command) *tweedy.API {
	t, _ := cmd.Flags().GetBool("test")
	o := &tweedy.Opts{Test: t}
	a, err := tweedy.New(o)
	if err != nil {
		logrus.Fatalln(err)
	}

	return a
}

// skipD logs if we are skipping because a Tweet
// is too young to be deleted, in that it's date is
// newer than the date the user wants to delete
func skipD(dDay time.Time, t *tweedy.Tweet) {
	c, i := t.CreatedAt, t.ID
	dy, dm, dd := dDay.Year(), dDay.Month(), dDay.Day()
	msg := "skipping %d from %d/%02d/%02d because after %d/%02d/%02d"
	ty, tm, td := c.Year(), c.Month(), c.Day()
	logrus.Infof(msg, i, ty, tm, td,
		dy, dm, dd)
}

// skipO logs skips because a Tweet ID is
// less than that of the ID of which the user
// wants to skip from, we don't delete older
func skipO(from int64, t *tweedy.Tweet) {
	c := t.CreatedAt
	msg := "skipping %d from %d/%02d/%02d because ID after %d"
	y, m, d := c.Year(), c.Month(), c.Day()
	logrus.Infof(msg, y, m, d, from)
}

// logE logs missing, or fatal messages
// We decide whether to throw based on the type of
// error, exp if it's a simple permission or
// missing error, we will keep going.
func check(e error, t *tweedy.Tweet) (ok bool) {
	if e == nil {
		return true
	}

	c := t.CreatedAt
	y, m, d := c.Year(), c.Month(), c.Day()
	err := e.(*anaconda.ApiError)

	if err.StatusCode == 404 {
		ok = true
		msg := "missing %d from %d/%02d/%02d"
		logrus.Infof(msg, t.ID, y, m, d)
	} else {
		if err.StatusCode == 403 {
			ok = true
			msg := "permission error on %d from %d/%02d/%02d"
			logrus.Warnf(msg, t.ID, y, m, d)
		}
	}

	if !ok {
		logrus.Fatalln(e)
	}

	return
}

// tweets loops on the Tweets, and deletes
func tweets(dDay time.Time, f int64, t []*tweedy.Tweet) {
	for _, tweet := range t {
		c := tweet.CreatedAt
		if c.After(dDay) {
			skipD(dDay, tweet)
			continue
		} else {
			if f > 0 && f > tweet.ID {
				skipO(f, tweet)
				continue
			}
		}

		if (tweet.IsRetweet() && check(tweet.UnRetweet(), tweet)) ||
			check(tweet.Delete(), tweet) {
			continue
		}
	}
}

func cmdRun(cmd *cobra.Command, _ []string) {
	var t tweedy.Tweets
	a, f := api(cmd), file(cmd)
	u, _ := cmd.Flags().GetString("user")
	var q tweedy.UserQuery
	var err error

	if f != "" {
		t, err = a.GetFromFile(f)
		if err != nil {
			logrus.Fatalln(err)
		}
	} else {
		q.Handle = u
		uid, err := strconv.Atoi(u)
		if err == nil {
			q.UID = int64(uid)
			q.Handle = ""
		}

		t, err = a.GetFromTimeline(q)
		if err != nil {
			logrus.Fatalln(err)
		}
	}

	k := keep(cmd)
	date := time.Now().AddDate(0, 0, k)
	tweets(date, from(cmd), t)
}
