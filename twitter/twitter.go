// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
)

var (
	_envKeys = [4]string{
		"TWITTER_ACCESS_TOKEN",
		"TWITTER_ACCESS_TOKEN_SECRET",
		"TWITTER_API_KEY",
		"TWITTER_API_SECRET",
	}
)

// GetKeys gets access, and api keys from
// the current users environment, you can visit
// Twitters developer.twitter.com to register
func envKeys() (slice []string, err error) {
	for _, k := range _envKeys {
		v, ok := os.LookupEnv(k)
		if !ok {
			err = fmt.Errorf("unable to find key %s in env", k)
			return
		}

		slice = append(slice, v)
	}

	return
}

// API is a wrapper around anaconda.TwitterAPI
// it simplifies their interface, and gives us an
// interface that has only what we need.
type API struct {
	upstream *anaconda.TwitterApi
}

// GetFromFile pulls all the Tweets from a file.
func (a *API) GetFromFile(f string) (t Tweets, err error) {
	b, err := ioutil.ReadFile(f)
	if err == nil {
		err = json.Unmarshal(b, &t)
		if err != nil {
			return
		}
	} else {
		return
	}

	for _, tweet := range t {
		tweet.api = a
		tweet.Setup()
	}

	return
}

// GetFromTimeline pulls the Tweets from the Timeline
func (a *API) GetFromTimeline(handle string) (tweets Tweets, err error) {
	struid, err := a.GetUID(handle)
	if err != nil {
		return
	}

	uid := strconv.FormatInt(struid, 10)
	t, err := a.upstream.GetUserTimeline(url.Values{
		"trim_user":       []string{"1"},
		"exclude_replies": []string{"0"},
		"include_rts":     []string{"1"},
		"user_id":         []string{uid},
		"count": []string{
			"200",
		},
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	for _, tt := range t {
		tweets = append(tweets, (&Tweet{
			upstream: &tt,
			api:      a,
		}).Setup())
	}

	return
}

// GetUID gets the user
func (a *API) GetUID(handle string) (i int64, err error) {
	user, err := a.upstream.GetUsersLookup(handle, url.Values{})
	if err != nil {
		return
	}

	if len(user) > 1 {
		// I don't understand how this can happen?
		err = fmt.Errorf("too many users for %s", handle)
		return
	}

	i = user[0].Id
	return
}

// Get retrieves a Tweet from the API
func (a *API) Get(i int64) (*Tweet, error) {
	logrus.Infof("fetching %d", i)
	t, err := a.upstream.GetTweet(i, url.Values{})
	if err != nil {
		return nil, err
	}

	return (&Tweet{
		upstream: &t,
		api:      a,
	}).Setup(), nil
}

// Delete a Tweet
func (a *API) Delete(i int64) error {
	logrus.Infof("deleting %d", i)
	_, err := a.upstream.DeleteTweet(i, true)
	return err
}

// New gets the keys, and then returns an API
// created by Anaconda.  If you wish to do more advanced
// stuff than what I'm doing, hit up their docs
func New() (*API, error) {
	k, err := envKeys()
	if err != nil {
		return nil, err
	}

	a, b, c, d := k[0], k[1], k[2], k[3]
	upstream := anaconda.NewTwitterApiWithCredentials(a, b, c, d)
	return &API{
		upstream: upstream,
	}, nil
}