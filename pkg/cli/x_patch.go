// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	api "github.com/scaleway/scaleway-cli/pkg/api"
)

var cmdPatch = &Command{
	Exec:        runPatch,
	UsageLine:   "_patch [OPTIONS] IDENTIFIER FIELD=VALUE",
	Description: "",
	Hidden:      true,
	Help:        "PATCH an object on the API",
	Examples: `
    $ scw _patch myserver state_detail=booted
    $ scw _patch server:myserver state_detail=booted
`,
}

func init() {
	cmdPatch.Flag.BoolVar(&patchHelp, []string{"h", "-help"}, false, "Print usage")
}

// Flags
var patchHelp bool // -h, --help flag

func runPatch(cmd *Command, args []string) error {
	if patchHelp {
		return cmd.PrintUsage()
	}
	if len(args) != 2 {
		return cmd.PrintShortUsage()
	}

	// Parsing FIELD=VALUE
	updateParts := strings.SplitN(args[1], "=", 2)
	if len(updateParts) != 2 {
		return cmd.PrintShortUsage()
	}
	fieldName := updateParts[0]
	newValue := args[1][len(updateParts[0])+1:]

	changes := 0

	ident, err := api.GetIdentifier(cmd.API, args[0])
	if err != nil {
		return err
	}
	switch ident.Type {
	case api.IdentifierServer:
		currentServer, err := cmd.API.GetServer(ident.Identifier)
		if err != nil {
			return fmt.Errorf("Cannot get server %s: %v", ident.Identifier, err)
		}

		var payload api.ScalewayServerPatchDefinition

		switch fieldName {
		case "state_detail":
			log.Debugf("%s=%s  =>  %s=%s", fieldName, currentServer.StateDetail, fieldName, newValue)
			if currentServer.StateDetail != newValue {
				changes++
				payload.StateDetail = &newValue
			}
		case "name":
			log.Warnf("To rename a server, Use 'scw rename'")
			log.Debugf("%s=%s  =>  %s=%s", fieldName, currentServer.StateDetail, fieldName, newValue)
			if currentServer.Name != newValue {
				changes++
				payload.Name = &newValue
			}
		case "bootscript":
			log.Debugf("%s=%s  =>  %s=%s", fieldName, currentServer.Bootscript.Identifier, fieldName, newValue)
			if currentServer.Bootscript.Identifier != newValue {
				changes++
				payload.Bootscript = &newValue
			}
		case "security_group":
			log.Debugf("%s=%s  =>  %s=%s", fieldName, currentServer.SecurityGroup.Identifier, fieldName, newValue)
			if currentServer.SecurityGroup.Identifier != newValue {
				changes++
				payload.SecurityGroup = &api.ScalewaySecurityGroup{
					Identifier: newValue,
				}
			}
		case "tags":
			newTags := strings.Split(newValue, " ")
			log.Debugf("%s=%s  =>  %s=%s", fieldName, currentServer.Tags, fieldName, newTags)
			// fixme test equality with reflect.DeepEqual ?
			changes++
			payload.Tags = &newTags
		case "ipv6":
			log.Debugf("%s=%s  =>  %s=%s", fieldName, currentServer.Tags, fieldName, newValue)
			switch strings.ToLower(newValue) {
			case "true":
				t := true
				payload.EnableIPV6 = &t
				changes++
			case "false":
				f := false
				payload.EnableIPV6 = &f
				changes++
			}
		default:
			return fmt.Errorf("'_patch server %s=' not implemented", fieldName)
		}
		// FIXME: volumes, tags, dynamic_ip_required

		if changes > 0 {
			log.Debugf("updating server: %d change(s)", changes)
			err = cmd.API.PatchServer(ident.Identifier, payload)
		} else {
			log.Debugf("no changes, not updating server")
		}
		if err != nil {
			return fmt.Errorf("Cannot update server: %v", err)
		}
	default:
		return fmt.Errorf("_patch not implemented for this kind of object")
	}
	fmt.Println(ident.Identifier)
	return nil
}
