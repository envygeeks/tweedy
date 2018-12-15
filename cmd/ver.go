// Copyright 2018 Jordon Bedwell. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	commit  string
	version string
	url     string

	Version = &cobra.Command{
		Use:   "version",
		Short: "Get the version",
		Run:   versionRun,
	}
)

func init() {
	Cmd.AddCommand(Version)
}

func versionRun(*cobra.Command, []string) {
	fmt.Printf("commit: %s\n", commit)
	fmt.Printf("commit-url: %s/show/%s\n", url, commit)
	fmt.Printf("full-source-url: %s\n", url)
	fmt.Printf("version: %s\n", version)
}
