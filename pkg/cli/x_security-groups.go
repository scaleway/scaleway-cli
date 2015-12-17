// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

var cmdSecurityGroups = &Command{
	Exec:        runSecurityGroups,
	UsageLine:   "_security-groups [OPTIONS] [ARGS]",
	Description: "Interacts with security-groups",
	Hidden:      true,
	Help:        "Interacts with security-groups",
	Examples: `
    $ scw _security-groups
    $ scw _security-groups "SecurityGroupID"
    $ scw _security-groups --new "NAME:DESCRIPTION"
    $ scw _security-groups --update "NAME:DESCRIPTION" "SecurityGroupID"
    $ scw _security-groups --delete "SecurityGroupID"
    $ scw _security-groups --rules "SecurityGroupID"
    $ scw _security-groups --rule-id "SecurityGroupID:RuleID"
    $ scw _security-groups --rule-delete "SecurityGroupID:RuleID"
    $ scw _security-groups --rule-new "SecurityGroupID:ACTION:DIRECTION:IP_RANGE:PROTOCOL[:DEST_PORT_FROM]"
    $ scw _security-groups --rule-update "SecurityGroupID:RuleID:ACTION:DIRECTION:IP_RANGE:PROTOCOL[:DEST_PORT_FROM]"
`,
}

func init() {
	cmdSecurityGroups.Flag.BoolVar(&securityGroupsHelp, []string{"h", "-help"}, false, "Print usage")
	// cmdSecurityGroups.Flag.BoolVar(&securityGroupsRaw, []string{"r", "-raw"}, false, "Displays the output in raw mode")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsNew, []string{"n", "-new"}, "", "Adds a new security group")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsUpdate, []string{"u", "-update"}, "", "Updates a security group")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsDelete, []string{"d", "-delete"}, "", "Deteles a security group")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsRules, []string{"r", "-rules"}, "", "Displays the rules")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsRuleID, []string{"ri", "-rule-id"}, "", "Displays one rule")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsRuleDelete, []string{"rd", "-rule-delete"}, "", "Deletes one rule")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsRuleNew, []string{"rn", "-rule-new"}, "", "Adds a new rule")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsRuleUpdate, []string{"ru", "-rule-update"}, "", "Updates one rule")
}

// Flags
// var securityGroupsRaw bool       // -r, --raw flag
var securityGroupsHelp bool         // -h, --help flag
var securityGroupsNew string        // -n, --new flag
var securityGroupsUpdate string     // -u, --update flag
var securityGroupsDelete string     // -d, --delete flag
var securityGroupsRules string      // -r, --rules flag
var securityGroupsRuleID string     // -ri, --rule-id flag
var securityGroupsRuleDelete string // -rd, --rule-delete flag
var securityGroupsRuleNew string    // -rn, --rule-new flag
var securityGroupsRuleUpdate string // -ru, --rule-update flag

func printRawMode(out io.Writer, data interface{}) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Unable to parse the data: %v", err)
	}
	fmt.Fprintf(out, "%s\n", string(js))
	return nil
}

func runSecurityGroups(cmd *Command, args []string) error {
	if securityGroupsHelp {
		return cmd.PrintUsage()
	}

	if securityGroupsNew != "" {
		newParts := strings.SplitN(securityGroupsNew, ":", 2)
		if len(newParts) != 2 {
			return cmd.PrintShortUsage()
		}

		return cmd.API.PostSecurityGroup(api.ScalewayNewSecurityGroup{
			Organization: cmd.API.Organization,
			Name:         newParts[0],
			Description:  newParts[1],
		})
	} else if securityGroupsUpdate != "" {
		if len(args) != 1 {
			return cmd.PrintShortUsage()
		}
		newParts := strings.SplitN(securityGroupsUpdate, ":", 2)
		if len(newParts) != 2 {
			return cmd.PrintShortUsage()
		}
		return cmd.API.PutSecurityGroup(api.ScalewayNewSecurityGroup{
			Organization: cmd.API.Organization,
			Name:         newParts[0],
			Description:  newParts[1],
		}, args[0])
	} else if securityGroupsDelete != "" {
		return cmd.API.DeleteSecurityGroup(securityGroupsDelete)
	} else if securityGroupsRules != "" {
		GetSecurityGroupRules, err := cmd.API.GetSecurityGroupRules(securityGroupsRules)
		if err != nil {
			return err
		}
		return printRawMode(cmd.Streams().Stdout, *GetSecurityGroupRules)
	} else if securityGroupsRuleID != "" {
		newParts := strings.SplitN(securityGroupsRuleID, ":", 2)
		if len(newParts) != 2 {
			return cmd.PrintShortUsage()
		}
		GroupRuleID, err := cmd.API.GetASecurityGroupRule(newParts[0], newParts[1])
		if err != nil {
			return err
		}
		return printRawMode(cmd.Streams().Stdout, *GroupRuleID)
	} else if securityGroupsRuleDelete != "" {
		newParts := strings.SplitN(securityGroupsRuleDelete, ":", 2)
		if len(newParts) != 2 {
			return cmd.PrintShortUsage()
		}
		return cmd.API.DeleteSecurityGroupRule(newParts[0], newParts[1])
	} else if securityGroupsRuleNew != "" {
		newParts := strings.Split(securityGroupsRuleNew, ":")
		if len(newParts) != 5 && len(newParts) != 6 {
			return cmd.PrintShortUsage()
		}
		if len(newParts) == 6 {
			port, err := strconv.Atoi(newParts[5])
			if err != nil {
				return err
			}
			return cmd.API.PostSecurityGroupRule(newParts[0], api.ScalewayNewSecurityGroupRule{
				Action:       newParts[1],
				Direction:    newParts[2],
				IPRange:      newParts[3],
				Protocol:     newParts[4],
				DestPortFrom: port,
			})
		}
		return cmd.API.PostSecurityGroupRule(newParts[0], api.ScalewayNewSecurityGroupRule{
			Action:    newParts[1],
			Direction: newParts[2],
			IPRange:   newParts[3],
			Protocol:  newParts[4],
		})
	} else if securityGroupsRuleUpdate != "" {
		newParts := strings.Split(securityGroupsRuleUpdate, ":")
		if len(newParts) != 6 && len(newParts) != 7 {
			return cmd.PrintShortUsage()
		}
		if len(newParts) == 7 {
			port, err := strconv.Atoi(newParts[6])
			if err != nil {
				return err
			}
			return cmd.API.PutSecurityGroupRule(api.ScalewayNewSecurityGroupRule{
				Action:       newParts[2],
				Direction:    newParts[3],
				IPRange:      newParts[4],
				Protocol:     newParts[5],
				DestPortFrom: port,
			}, newParts[0], newParts[1])
		}
		return cmd.API.PutSecurityGroupRule(api.ScalewayNewSecurityGroupRule{
			Action:    newParts[2],
			Direction: newParts[3],
			IPRange:   newParts[4],
			Protocol:  newParts[5],
		}, newParts[0], newParts[1])
	}
	if len(args) == 1 {
		securityGroups, err := cmd.API.GetASecurityGroup(args[0])
		if err != nil {
			return err
		}
		return printRawMode(cmd.Streams().Stdout, *securityGroups)
	}
	securityGroups, err := cmd.API.GetSecurityGroups()
	if err != nil {
		return err
	}
	return printRawMode(cmd.Streams().Stdout, *securityGroups)
}
