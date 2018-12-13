// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/envygeeks/tweedy/twitter/keys"
)

// API is a wrapper around anaconda.TwitterAPI it
// simplifies their interface, and gives us an interface
// that has only what we need.
type API struct {
	upstream *anaconda.TwitterApi
}

// New gets the keys, and then returns an API
// created by Anaconda.  If you wish to do more advanced
// stuff than what I'm doing, hit up their docs
func New() (*API, error) {
	k, err := keys.Get()
	if err != nil {
		return nil, err
	}

	k1, k2, k3, k4 := k[0], k[1], k[2], k[3]
	upstream := anaconda.NewTwitterApiWithCredentials(k1, k2, k3, k4)
	return &API{
		upstream: upstream,
	}, nil
}
