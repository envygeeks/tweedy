// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/envygeeks/tweedy/umap"
)

// User is a user
type User struct {
	mapped   bool           // Whether we've mapped upstream
	upstream *anaconda.User // The upstream we are encapsulating
	api      *API           // The API we are working with

	/**
	 * Public
	 */
	EMail  string
	Handle string
	Name   string
	UID    int64

	/**
	 * Counts
	 * BLOCKED: The TODO requires umap.Map adjustments
	 * TODO: Move this onto a UserCount
	 */
	TweetCount int64
	LikesCount int
}

// Likes are likes
type Likes *[]Tweet

// UserQuery is used in cases where we will query an API
// that can take either a UID, or a Handle.
//
//   func MyFunc(query UserQuery) {
//     if query.Handle != ""
//       // Get UID, and do work.
//     }
//   }
type UserQuery struct {
	Handle string
	UID    int64
}

var (
	userMap = map[string]string{
		"Email":           "EMail",
		"StatusesCount":   "TweetCount",
		"FavouritesCount": "LikesCount",
		"ScreenName":      "Handle",
		"Name":            "Name",
		"Id":              "UID",
	}
)

// Init the user
func (u *User) Init() (err error) {
	if !u.mapped && u.upstream != nil {
		err := umap.Map(u, userMap)
		if err == nil {
			u.mapped = true
		}
	}

	return
}

// NewUser creates a new User
func NewUser(uu *anaconda.User, a *API) (u *User, err error) {
	u = &User{}
	if uu != nil {
		u.upstream = uu
	}

	u.api = a
	err = u.Init()
	return
}
