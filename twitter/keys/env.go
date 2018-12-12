// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package keys

import (
	"fmt"
	"os"
)

var (
	envKeys = [4]string{
		"TWITTER_ACCESS_TOKEN",
		"TWITTER_ACCESS_TOKEN_SECRET",
		"TWITTER_API_KEY",
		"TWITTER_API_SECRET",
	}
)

// FromEnv gets access, and api keys from the users env,
// you can visit Twitters developer.twitter.com to register
// to get your keys, I do not provide them for you
func FromEnv() (slice []string, err error) {
	for _, k := range envKeys {
		v, ok := os.LookupEnv(k)
		if !ok {
			err = fmt.Errorf("unable to find key %s in env", k)
			return
		}

		slice = append(slice, v)
	}

	return
}
