// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package commands

import (
	"fmt"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

// CreateArgs are arguments passed to `RunCreate`
type CreateArgs struct {
	Volumes        []string
	Tags           []string
	Name           string
	Bootscript     string
	Image          string
	IP             string
	CommercialType string
	TmpSSHKey      bool
	IPV6           bool
}

// RunCreate is the handler for 'scw create'
func RunCreate(ctx CommandContext, args CreateArgs) error {
	if args.TmpSSHKey {
		err := AddSSHKeyToTags(ctx, &args.Tags, args.Image)
		if err != nil {
			return err
		}
	}

	env := strings.Join(args.Tags, " ")
	volume := strings.Join(args.Volumes, " ")
	config := api.ConfigCreateServer{
		ImageName:         args.Image,
		Name:              args.Name,
		Bootscript:        args.Bootscript,
		Env:               env,
		AdditionalVolumes: volume,
		DynamicIPRequired: false,
		IP:                args.IP,
		CommercialType:    args.CommercialType,
		EnableIPV6:        args.IPV6,
	}
	if args.IP == "dynamic" || args.IP == "" {
		config.DynamicIPRequired = true
		config.IP = ""
	} else if args.IP == "none" || args.IP == "no" {
		config.IP = ""
	}
	serverID, err := api.CreateServer(ctx.API, &config)
	if err != nil {
		return err
	}
	logrus.Debugf("Server created: %s", serverID)
	logrus.Debugf("PublicDNS %s", serverID+api.URLPublicDNS)
	logrus.Debugf("PrivateDNS %s", serverID+api.URLPrivateDNS)
	fmt.Fprintln(ctx.Stdout, serverID)
	return nil
}
