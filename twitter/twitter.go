// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package twitter

import (
	"encoding/json"
	"io/ioutil"
	"net/url"
	"strconv"
)

const (
	tweetlimit = 3200
)

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
		tweet.Init()
	}

	return
}

// GetFromTimeline pulls the Tweets from the Timeline
// UserQuery is just a struct that holds either UID, or a
// handle, and allows you to decide which you want
//
// 	GetFromTimeline(&UserQuery{UID: 123456789 })
// 	GetFromTimeline(&UserQuery{
//   		Handle: "myTwitterHandle"
// 	})
func (a *API) GetFromTimeline(uq UserQuery) (Tweets, error) {
	tweets := Tweets{}
	seen := map[int64]struct{}{}
	var last int64

	u, err := a.GetUser(uq)
	if err != nil {
		return Tweets{}, err
	}

	struid := strconv.FormatInt(u.UID, 10)
	q := url.Values{
		"user_id":         []string{struid},
		"include_rts":     []string{"1"},
		"exclude_replies": []string{"0"},
		"trim_user":       []string{"1"},
		"count": []string{
			"200",
		},
	}

	for {
		delete(q, "max_id")
		if last != 0 {
			strlid := strconv.FormatInt(last, 10)
			q["max_id"] = []string{
				strlid,
			}
		}

		timeline, err := timeline(a, q)
		if err != nil {
			return Tweets{}, err
		}

		if len(timeline) == 0 {
			break
		} else {
			slen := len(seen)
			for _, v := range timeline {
				if _, ok := seen[v.ID]; !ok {
					tweets = append(tweets, v)
					seen[v.ID] = struct{}{}
				}
			}

			last = tweets[len(tweets)-1].ID
			if slen == len(seen) || len(seen) >= tweetlimit {
				break
			}
		}
	}

	return tweets, nil
}

func timeline(a *API, q url.Values) (t Tweets, err error) {
	utweets, err := a.upstream.GetUserTimeline(q)
	if err != nil {
		return nil, err
	}

	for _, utweet := range utweets {
		tweet, err := NewTweet(&utweet, a)
		if err != nil {
			break
		}

		t = append(t, tweet)
	}

	return
}

// uByUID gets the user by their UID
func uByUID(a *API, id int64) (*User, error) {
	q := url.Values{}
	u, err := a.upstream.GetUsersShowById(id, q)
	if err != nil {
		return nil, err
	}

	return NewUser(&u, a)
}

// uByHandle gets a user based on the username, or
// handle that's really just whatever you decide to call it
// because it's not important enough to get upset about
func uByHandle(a *API, handle string) (*User, error) {
	q := url.Values{}
	u, err := a.upstream.GetUsersShow(handle, q)
	if err != nil {
		return nil, err
	}

	return NewUser(&u, a)
}

// GetUser gets the user
func (a *API) GetUser(q UserQuery) (u *User, err error) {
	if q.UID == 0 && q.Handle != "" {
		u, err = uByHandle(a, q.Handle)
	} else {
		if q.UID != 0 {
			u, err = uByUID(a, q.UID)
		}
	}

	return
}

// GetTweet retrieves a Tweet from the API
func (a *API) GetTweet(i int64) (*Tweet, error) {
	t, err := a.upstream.GetTweet(i, url.Values{})
	if err != nil {
		return nil, err
	}

	return NewTweet(&t, a)
}
