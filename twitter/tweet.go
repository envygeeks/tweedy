// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package twitter

import (
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/envygeeks/tweedy/umap"
	"github.com/sirupsen/logrus"
)

// Tweet is an individual Tweet ID
// We don't care about anything but the specific
// Tweet ID that we plan to delete
type Tweet struct {
	upstream     *anaconda.Tweet
	CreatedAt    time.Time `json:"-"`
	FullText     string    `json:"full_text"`
	CreatedAtStr string    `json:"created_at"`
	ID           int64     `json:",string"`
	api          *API
}

var (
	tweetUpstreamMap = umap.Map{
		"FullText":     "FullText",
		"CreatedAtStr": "CreatedAt",
		"Id":           "ID",
	}
)

// Tweets are Tweet
type Tweets []*Tweet

// Setup does constructor stuff.
func (t *Tweet) Setup() *Tweet {
	var err error

	if t.upstream != nil {
		err = umap.MapValues(t.upstream, t, tweetUpstreamMap)
		if err == nil {
			t.CreatedAt, err = t.upstream.
				CreatedAtTime()
		}
	} else {
		t.CreatedAt, err = time.Parse(time.RubyDate,
			t.CreatedAtStr)
	}

	if err != nil {
		logrus.Fatalln(err)
	}

	return t
}

// IsRetweet tells if Retweet
func (t *Tweet) IsRetweet() bool {
	return strings.HasPrefix(t.FullText, "RT ")
}

// Unretweet a Tweet
func (t *Tweet) Unretweet() (err error) {
	_, err = t.api.upstream.UnRetweet(t.ID, true)
	if err == nil {
		c := t.CreatedAt
		m, d, y := c.Month(), c.Day(), c.Year()
		logrus.Infof("unretweeted %d from %d/%02d/%02d",
			t.ID, y, m, d)
	}

	return
}

// Delete the Tweet
func (t *Tweet) Delete() (err error) {
	_, err = t.api.upstream.DeleteTweet(t.ID, true)
	if err == nil {
		c := t.CreatedAt
		m, d, y := c.Month(), c.Day(), c.Year()
		logrus.Infof("deleted %d from %d/%02d/%02d",
			t.ID, y, m, d)
	}

	return
}
