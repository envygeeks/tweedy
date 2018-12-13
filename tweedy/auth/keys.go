// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package auth

import (
	"fmt"
	"os"
)

// Keys holds the keys to the API
type Keys struct {
	SecretKey string
	Key       string
}

// FromEnv gets keys from the env
func (k *Keys) FromEnv() error {
	keys := []string{"TWITTER_KEY", "TWITTER_SECRET_KEY"}
	v := map[string]string{}
	for _, k := range keys {
		var ok bool

		v[k], ok = os.LookupEnv(k)
		if !ok {
			m := "unable to find %s in env"
			err := fmt.Errorf(m, k)
			return err
		}
	}

	k.Key = v[keys[0]]
	k.SecretKey = v[keys[1]]
	return nil
}
