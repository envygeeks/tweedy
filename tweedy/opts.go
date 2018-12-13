// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package tweedy

import (
	"github.com/envygeeks/tweedy/tweedy/auth"
)

// Opts hold the opts that we
// pass between instances of Tweet,
// API, and User, for stuff
type Opts struct {
	Keys   *auth.Keys
	Tokens *auth.Tokens
	DryRun bool
}
