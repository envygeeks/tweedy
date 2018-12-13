// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/envygeeks/tweedy/tweedy"
	"github.com/envygeeks/tweedy/tweedy/auth"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type cli struct {
	From    int
	File    string
	User    string
	DryRun  bool
	Verbose bool
	Silent  bool
	Keep    int
}

var (
	args    = &cli{}
	mainCmd = &cobra.Command{
		Use:   "tweedy",
		Short: "Tweedy deletes your Tweets, Retweets, and Likes",
		Run:   runCmd,
	}
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	mainCmd.Flags().IntVarP(&args.From, "from", "r", 0, "Delete Tweets from ID")
	mainCmd.Flags().BoolVarP(&args.Verbose, "verbose", "1", false, "Verbose output")
	mainCmd.Flags().BoolVarP(&args.DryRun, "dry-run", "d", false, "Disable destructive actions")
	mainCmd.Flags().StringVarP(&args.User, "user", "u", "", "User handle, or UID")
	mainCmd.Flags().IntVarP(&args.Keep, "keep", "k", 3, "Days to keep")
	mainCmd.Flags().StringVarP(&args.File, "file", "f",
		"", "Tweets JSON file")
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
		logrus.Infof("missing %d from %d/%02d/%02d",
			t.ID, y, m, d)
	} else {
		if err.StatusCode == 403 {
			ok = true
			logrus.Warnf("permission error on %d from %d/%02d/%02d",
				t.ID, y, m, d)
		}
	}

	if !ok {
		logrus.Fatalln(e)
	}

	return
}

// Wrap around okT and ship the first Tweet
func checkall(e error, t tweedy.Tweets) bool {
	return check(e, t[0])
}

func file() (f string) {
	f, err := filepath.Abs(args.File)
	if err != nil {
		logrus.Fatalln(err)
	} else {
		if _, err = os.Stat(f); os.IsNotExist(err) {
			msg := "unable to find the file %s"
			logrus.Fatalf(msg, args.File)
		}
	}

	return
}

func keep() (i int) {
	i = args.Keep
	if i > 0 {
		i = -i
	}

	return
}

// loopOn loops on the Tweets and deletes
func loopOn(deleteDate time.Time, from int64, t []*tweedy.Tweet) {
	for _, tweet := range t {
		c := tweet.CreatedAt
		if c.After(deleteDate) || (from > 0 && from > tweet.ID) {
			y, m, d := c.Year(), c.Month(), c.Day()
			msg := "skipping %d from %d/%02d/%02d"
			if tweet.CreatedAt.After(deleteDate) {
				msg = "not deleting %d from %d/%02d/%02d"
			}

			logrus.Infof(msg, tweet.ID, y, m, d)
			continue
		}

		if tweet.IsRetweet() {
			ok := false
			if check(tweet.UnRetweet(), tweet) {
				ok = true
			}

			// Double tap old RT's
			tweet.Delete()
			if ok {
				continue
			}
		} else {
			if check(tweet.Delete(), tweet) {
				continue
			}
		}
	}
}

func api() (api *tweedy.API) {
	t, k := &auth.Tokens{}, &auth.Keys{}
	err := t.FromEnv()
	if err != nil {
		log.Fatalln(err)
	}

	err = k.FromEnv()
	if err != nil {
		log.Fatalln(err)
	}

	o := &tweedy.Opts{
		Tokens: t,
		DryRun: args.DryRun,
		Keys:   k,
	}

	a, err := tweedy.New(o)
	if err != nil {
		logrus.Fatalln(err)
	}

	return a
}

// runWithFile is the runCmd for
// the command when ran with --file=tweets.json
// as it has it's own logic.
func runCmd(*cobra.Command, []string) {
	setupLogging()

	var (
		tweets tweedy.Tweets
		err    error
	)

	api := api()
	if args.File != "" {
		f := file()
		tweets, err = api.GetFromFile(f)
	} else {
		if u := args.User; u != "" {
			if uid, fail := strconv.Atoi(u); fail == nil {
				tweets, err = api.GetFromTimeline(
					tweedy.UserQuery{
						UID: int64(uid),
					})
			} else {
				tweets, err = api.GetFromTimeline(
					tweedy.UserQuery{
						Handle: u,
					})
			}
		} else {
			// We need one or the other.
			log.Fatalf("no user provided, must provide %s, or %s",
				"--user", "--file")
		}
	}
	if err != nil {
		log.Fatalln(err)
	}

	date := time.Now().AddDate(0, 0, keep())
	loopOn(date, int64(args.From),
		tweets)
}

func setupLogging() {
	logrus.SetLevel(logrus.WarnLevel)
	if args.Verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

// main
func main() {
	if err := mainCmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
