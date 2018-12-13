// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package tweedy

import (
	"github.com/ChimeraCoder/anaconda"
)

// API is a wrapper around anaconda.TwitterAPI it
// simplifies their interface, and gives us an interface
// that has only what we need.
type API struct {
	upstream *anaconda.TwitterApi
	opts     *Opts
}

// New gets the keys, and then returns an API
// created by Anaconda.  If you wish to do more advanced
// stuff than what I'm doing, hit up their docs
func New(opts *Opts) (*API, error) {
	k, sk := opts.Keys.Key, opts.Keys.SecretKey
	t, st := opts.Tokens.Token, opts.Tokens.TokenSecret
	a := anaconda.NewTwitterApiWithCredentials(t, st, k, sk)
	api := &API{upstream: a, opts: opts}
	return api, nil
}
