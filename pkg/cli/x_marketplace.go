// Copyright (C) 2015 Scaleway. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE.md file.

package cli

import (
	"encoding/json"

	"github.com/scaleway/scaleway-cli/pkg/api"
)

var cmdMarketplace = &Command{
	Exec:        runMarketplace,
	UsageLine:   "_marketplace -r VERB [FIELD]+",
	Description: "",
	Hidden:      true,
	Help:        "List, read and write and delete marketplace",
	Examples: `
    $ scw _marketplace -r GET images
    $ scw _marketplace -r GET versions UUID_IMAGE
    $ scw _marketplace -r GET local_images UUID_IMAGE UUID_LOCAL_VERSION
    $ scw _marketplace -r POST images <data>
    $ scw _marketplace -r POST versions UUID_IMAGE <data>
    $ scw _marketplace -r POST local_images UUID_IMAGE UUID_LOCAL_VERSION <data>
    $ scw _marketplace -r POST images <data>
    $ scw _marketplace -r PUT versions UUID_IMAGE <data>
    $ scw _marketplace -r PUT local_images UUID_IMAGE UUID_LOCAL_VERSION <data>
    $ scw _marketplace -r PUT images <data>
    $ scw _marketplace -r PUT versions UUID_IMAGE <data>
    $ scw _marketplace -r DELETE local_images UUID_IMAGE UUID_LOCAL_VERSION
    $ scw _marketplace -r DELETE images
    $ scw _marketplace -r DELETE versions UUID_IMAGE
    $ scw _marketplace -r DELETE local_images UUID_IMAGE UUID_LOCAL_VERSION
`,
}

func init() {
	cmdMarketplace.Flag.BoolVar(&marketplaceHelp, []string{"h", "-help"}, false, "Print usage")
	cmdMarketplace.Flag.StringVar(&marketplaceRequestType, []string{"r", "-request"}, "", "Choice a request type")
}

// Flags
var marketplaceHelp bool          // -h, --help flag
var marketplaceRequestType string // -r, --request flag

func getMarketPlace(cmd *Command, args []string) error {
	ctx := cmd.GetContext(args)

	switch args[0] {
	case "images":
		if len(args) == 2 {
			marketplaceImages, err := ctx.API.GetMarketPlaceImages(args[1])
			if err != nil {
				return err
			}
			printRawMode(cmd.Streams().Stdout, marketplaceImages)
			return nil
		}
		marketplaceImages, err := ctx.API.GetMarketPlaceImages("")
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, marketplaceImages)
	case "versions":
		if len(args) != 3 && len(args) != 2 {
			return cmd.PrintUsage()
		}
		if len(args) == 3 {
			marketplaceVersions, err := ctx.API.GetMarketPlaceImageVersions(args[1], args[2])
			if err != nil {
				return err
			}
			printRawMode(cmd.Streams().Stdout, marketplaceVersions)
			return nil
		}
		marketplaceVersions, err := ctx.API.GetMarketPlaceImageVersions(args[1], "")
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, marketplaceVersions)
	case "current_versions":
		if len(args) < 2 {
			return cmd.PrintUsage()
		}
		marketplaceVersions, err := ctx.API.GetMarketPlaceImageCurrentVersion(args[1])
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, marketplaceVersions)
	case "local_images":
		if len(args) != 3 && len(args) != 4 {
			return cmd.PrintUsage()
		}
		if len(args) == 4 {
			marketplaceLocalVersions, err := ctx.API.GetMarketPlaceLocalImages(args[1], args[2], args[3])
			if err != nil {
				return err
			}
			printRawMode(cmd.Streams().Stdout, marketplaceLocalVersions)
			return nil
		}
		marketplaceLocalVersions, err := ctx.API.GetMarketPlaceLocalImages(args[1], args[2], "")
		if err != nil {
			return err
		}
		printRawMode(cmd.Streams().Stdout, marketplaceLocalVersions)
	default:
		return cmd.PrintUsage()
	}
	return nil
}

func postMarketPlace(cmd *Command, args []string) error {
	ctx := cmd.GetContext(args)

	switch args[0] {
	case "images":
		var data api.MarketImage

		if len(args) < 2 {

		}
		err := json.Unmarshal([]byte(args[1]), &data)
		if err != nil {
			return err
		}
		err = ctx.API.PostMarketPlaceImage(data)
		if err != nil {
			return err
		}
	case "versions":
		var data api.MarketVersion

		if len(args) < 3 {
			return cmd.PrintUsage()
		}
		err := json.Unmarshal([]byte(args[2]), &data)
		if err != nil {
			return err
		}
		err = ctx.API.PostMarketPlaceImageVersion(args[1], data)
		if err != nil {
			return err
		}
	case "local_images":
		var data api.MarketLocalImage

		if len(args) < 4 {
			return cmd.PrintUsage()
		}
		err := json.Unmarshal([]byte(args[3]), &data)
		if err != nil {
			return err
		}
		err = ctx.API.PostMarketPlaceLocalImage(args[0], args[1], args[2], data)
		if err != nil {
			return err
		}
	default:
		return cmd.PrintUsage()
	}
	return nil
}

func putMarketPlace(cmd *Command, args []string) error {
	ctx := cmd.GetContext(args)

	switch args[0] {
	case "images":
		var data api.MarketImage

		if len(args) < 2 {
			return cmd.PrintUsage()
		}
		err := json.Unmarshal([]byte(args[1]), &data)
		if err != nil {
			return err
		}
		err = ctx.API.PostMarketPlaceImage(data)
		if err != nil {
			return err
		}
	case "versions":
		var data api.MarketVersion

		if len(args) < 3 {
			return cmd.PrintUsage()
		}
		err := json.Unmarshal([]byte(args[2]), &data)
		if err != nil {
			return err
		}
		err = ctx.API.PostMarketPlaceImageVersion(args[1], data)
		if err != nil {
			return err
		}
	case "local_images":
		var data api.MarketLocalImage

		if len(args) < 4 {
			return cmd.PrintUsage()
		}
		err := json.Unmarshal([]byte(args[3]), &data)
		if err != nil {
			return err
		}
		err = ctx.API.PostMarketPlaceLocalImage(args[0], args[1], args[2], data)
		if err != nil {
			return err
		}
	default:
		return cmd.PrintUsage()
	}
	return nil
}

func deleteMarketPlace(cmd *Command, args []string) error {
	ctx := cmd.GetContext(args)

	switch args[0] {
	case "images":
		if len(args) < 1 {
			return cmd.PrintUsage()
		}
		err := ctx.API.DeleteMarketPlaceImage(args[0])
		if err != nil {
			return err
		}
	case "versions":
		if len(args) < 3 {
			return cmd.PrintUsage()
		}
		err := ctx.API.DeleteMarketPlaceImageVersion(args[0], args[1])
		if err != nil {
			return err
		}
	case "local_images":
		if len(args) < 4 {
			return cmd.PrintUsage()
		}
		err := ctx.API.DeleteMarketPlaceLocalImage(args[0], args[1], args[2])
		if err != nil {
			return err
		}
	default:
		return cmd.PrintUsage()
	}
	return nil
}

func runMarketplace(cmd *Command, args []string) error {
	if marketplaceHelp {
		return cmd.PrintUsage()
	}
	if len(args) < 1 {
		return cmd.PrintShortUsage()
	}
	switch marketplaceRequestType {
	case "GET":
		return getMarketPlace(cmd, args)
	case "POST":
		return postMarketPlace(cmd, args)
	case "PUT":
		return putMarketPlace(cmd, args)
	case "DELETE":
		return deleteMarketPlace(cmd, args)
	default:
		return cmd.PrintUsage()
	}
}
