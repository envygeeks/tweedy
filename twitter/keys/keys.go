// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package keys

// Get returns the keys to the kingdom
// in the future this will be via Env, via Keychain, or
// via the CLI itself if you wish
func Get() ([]string, error) {
	return FromEnv()
}
