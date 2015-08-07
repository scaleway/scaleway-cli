// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"github.com/scaleway/scaleway-cli/pkg/commands"
	log "github.com/scaleway/scaleway-cli/vendor/github.com/Sirupsen/logrus"
)

var cmdLogin = &Command{
	Exec:        cmdExecLogin,
	UsageLine:   "login [OPTIONS]",
	Description: "Log in to Scaleway API",
	Help: `Generates a configuration file in '/home/$USER/.scwrc'
containing credentials used to interact with the Scaleway API. This
configuration file is automatically used by the 'scw' commands.

You can get your credentials on https://cloud.scaleway.com/#/credentials
`,
}

func init() {
	cmdLogin.Flag.StringVar(&organization, []string{"o", "-organization"}, "", "Organization")
	cmdLogin.Flag.StringVar(&token, []string{"t", "-token"}, "", "Token")
	cmdLogin.Flag.BoolVar(&loginHelp, []string{"h", "-help"}, false, "Print usage")
}

// FLags
var organization string // -o flag
var token string        // -t flag
var loginHelp bool      // -h, --help flag

func cmdExecLogin(cmd *Command, rawArgs []string) {
	if loginHelp {
		cmd.PrintUsage()
	}
	if len(rawArgs) != 0 {
		cmd.PrintShortUsage()
	}

	args := commands.LoginArgs{
		Organization: organization,
		Token:        token,
	}
	ctx := cmd.GetContext(rawArgs)
	err := commands.RunLogin(ctx, args)
	if err != nil {
		log.Fatalf("Cannot execute 'login': %v", err)
	}
}
