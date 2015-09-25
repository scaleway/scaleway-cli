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
		var newGroups api.ScalewayNewSecurityGroup

		newParts := strings.SplitN(securityGroupsNew, ":", 2)
		if len(newParts) != 2 {
			return cmd.PrintShortUsage()
		}
		newGroups.Organization = cmd.API.Organization
		newGroups.Name = newParts[0]
		newGroups.Description = newParts[1]
		resp, err := cmd.API.PostResponse("security_groups", newGroups)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Succeed POST code
		if resp.StatusCode == 201 {
			return nil
		}

		var error api.ScalewayAPIError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&error)
		if err != nil {
			return err
		}
		error.StatusCode = resp.StatusCode
		error.Debug()
		return error
	} else if securityGroupsUpdate != "" {
		var newGroups api.ScalewayNewSecurityGroup

		if len(args) != 1 {
			return cmd.PrintShortUsage()
		}
		newParts := strings.SplitN(securityGroupsUpdate, ":", 2)
		if len(newParts) != 2 {
			return cmd.PrintShortUsage()
		}
		newGroups.Organization = cmd.API.Organization
		newGroups.Name = newParts[0]
		newGroups.Description = newParts[1]
		resp, err := cmd.API.PutResponse(fmt.Sprintf("security_groups/%s", args[0]), newGroups)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Succeed PUT code
		if resp.StatusCode == 200 {
			return nil
		}

		var error api.ScalewayAPIError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&error)
		if err != nil {
			return err
		}
		error.StatusCode = resp.StatusCode
		error.Debug()
		return error
	} else if securityGroupsDelete != "" {
		resp, err := cmd.API.DeleteResponse(fmt.Sprintf("security_groups/%s", securityGroupsDelete))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Succeed PUT code
		if resp.StatusCode == 204 {
			return nil
		}

		var error api.ScalewayAPIError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&error)
		if err != nil {
			return err
		}
		error.StatusCode = resp.StatusCode
		error.Debug()
		return error
	} else if securityGroupsRules != "" {
		securityGroupRules, err := cmd.API.GetSecurityGroupRules(securityGroupsRules)
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, *securityGroupRules)
		return nil
	} else if securityGroupsRuleID != "" {
		newParts := strings.SplitN(securityGroupsRuleID, ":", 2)
		if len(newParts) != 2 {
			return cmd.PrintShortUsage()
		}
		GroupRuleID, err := cmd.API.GetASecurityGroupRule(newParts[0], newParts[1])
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, *GroupRuleID)
		return nil
	} else if securityGroupsRuleDelete != "" {
		newParts := strings.SplitN(securityGroupsRuleDelete, ":", 2)
		if len(newParts) != 2 {
			return cmd.PrintShortUsage()
		}
		resp, err := cmd.API.DeleteResponse(fmt.Sprintf("security_groups/%s/rules/%s", newParts[0], newParts[1]))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Succeed PUT code
		if resp.StatusCode == 204 {
			return nil
		}

		var error api.ScalewayAPIError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&error)
		if err != nil {
			return err
		}
		error.StatusCode = resp.StatusCode
		error.Debug()
		return error
	} else if securityGroupsRuleNew != "" {
		var newRule api.ScalewayNewSecurityGroupRule

		newParts := strings.Split(securityGroupsRuleNew, ":")
		if len(newParts) != 5 && len(newParts) != 6 {
			return cmd.PrintShortUsage()
		}
		newRule.Action = newParts[1]
		newRule.Direction = newParts[2]
		newRule.IPRange = newParts[3]
		newRule.Protocol = newParts[4]
		if len(newParts) == 6 {
			var err error

			newRule.DestPortFrom, err = strconv.Atoi(newParts[5])
			if err != nil {
				return err
			}
		}
		resp, err := cmd.API.PostResponse(fmt.Sprintf("security_groups/%s/rules", newParts[0]), newRule)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Succeed POST code
		if resp.StatusCode == 201 {
			return nil
		}

		var error api.ScalewayAPIError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&error)
		if err != nil {
			return err
		}
		error.StatusCode = resp.StatusCode
		error.Debug()
		return error
	} else if securityGroupsRuleUpdate != "" {
		var newRule api.ScalewayNewSecurityGroupRule

		newParts := strings.Split(securityGroupsRuleUpdate, ":")
		if len(newParts) != 6 && len(newParts) != 7 {
			return cmd.PrintShortUsage()
		}
		newRule.Action = newParts[2]
		newRule.Direction = newParts[3]
		newRule.IPRange = newParts[4]
		newRule.Protocol = newParts[5]
		if len(newParts) == 7 {
			var err error

			newRule.DestPortFrom, err = strconv.Atoi(newParts[6])
			if err != nil {
				return err
			}
		}
		resp, err := cmd.API.PutResponse(fmt.Sprintf("security_groups/%s/rules/%s", newParts[0], newParts[1]), newRule)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Succeed PUT code
		if resp.StatusCode == 200 {
			return nil
		}

		var error api.ScalewayAPIError
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&error)
		if err != nil {
			return err
		}
		error.StatusCode = resp.StatusCode
		error.Debug()
		return error
	}
	if len(args) == 1 {
		securityGroups, err := cmd.API.GetASecurityGroup(args[0])
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, *securityGroups)
		return nil
	}
	securityGroups, err := cmd.API.GetSecurityGroups()
	if err != nil {
		return err
	}
	printRawMode(cmd.Streams().Stdout, *securityGroups)
	return nil
}
