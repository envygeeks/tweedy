// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package auth

// Auth holds the tokens,
// and the keys to the kingdom
// provided by the user.
type Auth struct {
	Tokens *Tokens
	Keys   *Keys
}

// New is a new Auth
func New() (*Auth, error) {
	a := &Auth{Tokens: &Tokens{}, Keys: &Keys{}}
	err := a.Tokens.FromEnv()
	if err == nil {
		err = a.Keys.FromEnv()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return a, nil
}
