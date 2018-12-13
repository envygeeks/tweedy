// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package tweedy

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
	/**
	 * Private
	 */
	mapped   bool
	upstream *anaconda.Tweet
	opts     *Opts
	api      *API

	/**
	 * Public
	 */
	CreatedAt    time.Time `json:"-"`
	CreatedAtStr string    `json:"created_at"`
	FullText     string    `json:"full_text"`
	ID           int64     `json:",string"`
}

var (
	tweetTimefmt = time.RubyDate // For CreatedAtStr
	tweetMap     = map[string]string{
		"FullText":  "FullText",
		"CreatedAt": "CreatedAtStr",
		"Id":        "ID",
	}
)

// Tweets is obvious
type Tweets []*Tweet

// Init setups some stuff
// this is only necessary if you do &Tweet{} instead of
// using NewTweet, this mostly only happens in the
// case of serialization
func (t *Tweet) Init() (err error) {
	var emptyT time.Time
	if !t.mapped && t.upstream != nil {
		err = umap.Map(t, tweetMap)
		if err == nil {
			t.mapped = true
		}
	}

	if err == nil && t.CreatedAtStr != "" && t.CreatedAt == emptyT {
		t.CreatedAt, err = time.Parse(tweetTimefmt,
			t.CreatedAtStr)
	}

	return
}

// NewTweet creates a new Tweet
func NewTweet(t *anaconda.Tweet, a *API) (*Tweet, error) {
	tweet := &Tweet{}
	if t != nil {
		tweet.upstream = t
	}

	tweet.api = a
	tweet.opts = a.opts
	err := tweet.Init()
	return tweet, err
}

// IsRetweet tells if Retweet
func (t *Tweet) IsRetweet() bool {
	return strings.HasPrefix(t.FullText, "RT ")
}

// UnRetweet a Tweet
func (t *Tweet) UnRetweet() (err error) {
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
