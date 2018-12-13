// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package auth

import (
	"fmt"
	"os"
)

// Tokens are auth tokens
type Tokens struct {
	TokenSecret string
	Token       string
}

// FromEnv gets keys from the env
func (t *Tokens) FromEnv() error {
	keys := []string{"TWITTER_TOKEN", "TWITTER_TOKEN_SECRET"}
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

	t.Token = v[keys[0]]
	t.TokenSecret = v[keys[1]]
	return nil
}
