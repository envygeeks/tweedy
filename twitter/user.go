// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/envygeeks/tweedy/umap"
	"github.com/sirupsen/logrus"
)

// User is a user
type User struct {
	LikesCount int
	TweetCount int64
	upstream   *anaconda.User
	Handle     string
	EMail      string
	Name       string
	UID        int64
}

// Likes are likes
type Likes *[]Tweet

var (
	userUpstreamMap = umap.Map{
		"Email":           "EMail",
		"StatusesCount":   "TweetCount",
		"FavouritesCount": "LikesCount",
		"ScreenName":      "Handle",
		"Name":            "Name",
		"Id":              "UID",
	}
)

// UserQuery is used in cases where
// we will query an API that can take
// either a UID, or a Handle.
//
//   func MyFunc(u UserQuery) {
//     if u.Handle != ""
//       // Get UID, and do work.
//     }
//   }
type UserQuery struct {
	Handle string
	UID    int64
}

// Setup the user
func (u *User) Setup() *User {
	if u.upstream != nil {
		err := umap.MapValues(u.upstream, u, userUpstreamMap)
		if err != nil {
			logrus.Fatalln(err)
		}
	}

	return u
}

