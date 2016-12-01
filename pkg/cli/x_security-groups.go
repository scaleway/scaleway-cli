// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

var cmdSecurityGroups = &Command{
	Exec:        runSecurityGroups,
	UsageLine:   "_security-groups OPTIONS [ARGS]",
	Description: "Interacts with security-groups",
	Hidden:      true,
	Help:        "Interacts with security-groups",
	Examples: `

	SGID = SecurityGroupID
	ACTION = "{\"action\":"string", \"direction\": \"string\", \"ip_range\": \"string\", \"protocol\": \"string\", \"dest_port_from\": \"int\"}"

	$ scw _security-groups list-groups
	$ scw _security-groups show-group SGID
	$ scw _security-groups new-group --name=NAME --desc=DESC
	$ scw _security-groups update-group SGID --name=NAME --desc=DESC
	$ scw _security-groups delete-group SGID

	$ scw _security-groups list-rules SGID
	$ scw _security-groups show-rule SGID RULEID
	$ scw _security-gruops delete-rule SGID RULEID
	$ scw _security-groups new-rule SGID ACTION
	$ scw _security-groups update-rule SGID RULEID ACTION`,
}

// "show-rule"

func init() {
	cmdSecurityGroups.Flag.BoolVar(&securityGroupsHelp, []string{"h", "-help"}, false, "Print usage")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsName, []string{"n", "-name"}, "", "SecurityGroup's name")
	cmdSecurityGroups.Flag.StringVar(&securityGroupsDesc, []string{"d", "-description"}, "", "SecurityGroup's description")
	subCmdSecurityGroup = map[string]func(cmd *Command, args []string) error{
		"list-groups":  listSecurityGroup,
		"new-group":    newSecurityGroup,
		"update-group": updateSecurityGroup,
		"delete-group": deleteSecurityGroup,
		"show-group":   showSecurityGroup,
		"list-rules":   listSecurityGroupRule,
		"new-rule":     newSecurityGroupRule,
		"update-rule":  updateSecurityGroupRule,
		"delete-rule":  deleteSecurityGroupRule,
		"show-rule":    showSecurityGroupRule,
	}
}

// Flags
var securityGroupsHelp bool   // -h, --help flag
var securityGroupsName string // -n, --name flag
var securityGroupsDesc string // -d, --description flag

// SubCommand
var subCmdSecurityGroup map[string]func(cmd *Command, args []string) error

type rulesDefinition struct {
	Action       string `json:"action"`
	Direction    string `json:"direction"`
	IPRange      string `json:"ip_range"`
	Protocol     string `json:"protocol"`
	DestPortFrom *int   `json:"dest_port_from,omitempty"`
}

func printRawMode(out io.Writer, data interface{}) error {
	js, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Unable to parse the data: %v", err)
	}
	fmt.Fprintf(out, "%s\n", string(js))
	return nil
}

func newSecurityGroup(cmd *Command, args []string) error {
	fmt.Println(securityGroupsDesc)
	fmt.Println(securityGroupsName)
	if securityGroupsName == "" || securityGroupsDesc == "" {
		return cmd.PrintShortUsage()
	}
	return cmd.API.PostSecurityGroup(api.ScalewayNewSecurityGroup{
		Organization: cmd.API.Organization,
		Name:         securityGroupsName,
		Description:  securityGroupsDesc,
	})
}

func updateSecurityGroup(cmd *Command, args []string) error {
	fmt.Println(args)
	if securityGroupsName == "" || securityGroupsDesc == "" || len(args) != 1 {
		return cmd.PrintShortUsage()
	}
	return cmd.API.PutSecurityGroup(api.ScalewayUpdateSecurityGroup{
		Organization: cmd.API.Organization,
		Name:         securityGroupsName,
		Description:  securityGroupsDesc,
	}, args[0])
}

func deleteSecurityGroup(cmd *Command, args []string) error {
	if len(args) != 1 {
		return cmd.PrintShortUsage()
	}
	return cmd.API.DeleteSecurityGroup(args[0])
}

func showSecurityGroup(cmd *Command, args []string) error {
	if len(args) != 1 {
		return cmd.PrintShortUsage()
	}
	securityGroups, err := cmd.API.GetASecurityGroup(args[0])
	if err != nil {
		return err
	}
	return printRawMode(cmd.Streams().Stdout, *securityGroups)
}

func listSecurityGroup(cmd *Command, args []string) error {
	securityGroups, err := cmd.API.GetSecurityGroups()
	if err != nil {
		return err
	}
	return printRawMode(cmd.Streams().Stdout, *securityGroups)
}

func listSecurityGroupRule(cmd *Command, args []string) error {
	if len(args) != 1 {
		return cmd.PrintShortUsage()
	}
	GetSecurityGroupRules, err := cmd.API.GetSecurityGroupRules(args[0])
	if err != nil {
		return err
	}
	return printRawMode(cmd.Streams().Stdout, *GetSecurityGroupRules)
}

func newSecurityGroupRule(cmd *Command, args []string) error {
	var rule rulesDefinition
	var content api.ScalewayNewSecurityGroupRule

	if len(args) != 2 {
		return cmd.PrintShortUsage()
	}
	if err := json.Unmarshal([]byte(args[1]), &rule); err != nil {
		return err
	}
	content.Action = rule.Action
	content.Direction = rule.Direction
	content.IPRange = rule.IPRange
	content.Protocol = rule.Protocol
	if rule.DestPortFrom != nil {
		content.DestPortFrom = *rule.DestPortFrom
	}
	return cmd.API.PostSecurityGroupRule(args[0], content)
}

func updateSecurityGroupRule(cmd *Command, args []string) error {
	var rule rulesDefinition
	var content api.ScalewayNewSecurityGroupRule

	if len(args) != 3 {
		return cmd.PrintShortUsage()
	}
	if err := json.Unmarshal([]byte(args[2]), &rule); err != nil {
		return err
	}
	content.Action = rule.Action
	content.Direction = rule.Direction
	content.IPRange = rule.IPRange
	content.Protocol = rule.Protocol
	if rule.DestPortFrom != nil {
		content.DestPortFrom = *rule.DestPortFrom
	}
	return cmd.API.PutSecurityGroupRule(content, args[0], args[1])
}

func showSecurityGroupRule(cmd *Command, args []string) error {
	if len(args) != 2 {
		return cmd.PrintShortUsage()
	}
	GroupRuleID, err := cmd.API.GetASecurityGroupRule(args[0], args[1])
	if err != nil {
		return err
	}
	return printRawMode(cmd.Streams().Stdout, *GroupRuleID)
}

func deleteSecurityGroupRule(cmd *Command, args []string) error {
	if len(args) != 2 {
		return cmd.PrintShortUsage()
	}
	return cmd.API.DeleteSecurityGroupRule(args[0], args[1])
}

func runSecurityGroups(cmd *Command, args []string) error {
	if securityGroupsHelp || len(args) == 0 {
		return cmd.PrintUsage()
	}
	cmd.Flag.Parse(args[1:])
	if function, ok := subCmdSecurityGroup[args[0]]; ok {
		return function(cmd, cmd.Flag.Args())
	}
	return fmt.Errorf("subcommand not found: %s", args[0])
}
