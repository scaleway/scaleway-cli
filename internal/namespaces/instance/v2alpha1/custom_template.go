package instance

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
)

func templateMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp instance.Template
	template := tmp(i.(instance.Template))

	// Convert empty slices to nil so HideIfEmpty works correctly
	if len(template.Tags) == 0 {
		template.Tags = nil
	}
	if len(template.ServerTags) == 0 {
		template.ServerTags = nil
	}
	if len(template.Volumes) == 0 {
		template.Volumes = nil
	}
	if len(template.PrivateNetworks) == 0 {
		template.PrivateNetworks = nil
	}
	if len(template.FilesystemIDs) == 0 {
		template.FilesystemIDs = nil
	}

	opt.Sections = []*human.MarshalSection{
		{FieldName: "Tags", HideIfEmpty: true},
		{FieldName: "ServerTags", HideIfEmpty: true, Title: "Server Tags"},
		{FieldName: "Volumes", HideIfEmpty: true},
		{FieldName: "PrivateNetworks", HideIfEmpty: true, Title: "Private Networks"},
		{FieldName: "FilesystemIDs", HideIfEmpty: true, Title: "Filesystems"},
	}

	str, err := human.Marshal(template, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

type customCreateTemplateRequest struct {
	*instance.CreateTemplateRequest
	PublicIPv4Count uint32
	PublicIPv6Count uint32
}

func TemplateCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("name").Default = core.RandomValueGenerator("tmpl")
	c.ArgSpecs.GetByName("server-type").Required = true
	c.ArgSpecs.GetByName("public-ip-v4-count").Name = "public-ipv4-count"
	c.ArgSpecs.GetByName("public-ip-v6-count").Name = "public-ipv6-count"

	c.ArgsType = reflect.TypeFor[customCreateTemplateRequest]()

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		args := argsI.(*customCreateTemplateRequest)

		request := args.CreateTemplateRequest
		request.PublicIPV4Count = args.PublicIPv4Count
		request.PublicIPV6Count = args.PublicIPv6Count

		return runner(ctx, request)
	}

	return c
}

type customUpdateTemplateRequest struct {
	*instance.UpdateTemplateRequest
	PublicIPv4Count uint32
	PublicIPv6Count uint32
	Volumes         []*instance.CreateTemplateRequestVolumeTemplate
	PrivateNetworks []*instance.CreateTemplateRequestPrivateNetworkTemplate
}

func TemplateUpdateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true
	c.ArgSpecs.GetByName("public-ip-v4-count").Name = "public-ipv4-count"
	c.ArgSpecs.GetByName("public-ip-v6-count").Name = "public-ipv6-count"
	c.ArgSpecs.GetByName("update-volumes.volumes.{index}.volume-type").Name = "volumes.{index}.volume-type"
	c.ArgSpecs.GetByName("update-volumes.volumes.{index}.name").Name = "volumes.{index}.name"
	c.ArgSpecs.GetByName("update-volumes.volumes.{index}.tags.{index}").Name = "volumes.{index}.tags.{index}"
	c.ArgSpecs.GetByName("update-volumes.volumes.{index}.size").Name = "volumes.{index}.size"
	c.ArgSpecs.GetByName("update-volumes.volumes.{index}.base-snapshot-id").Name = "volumes.{index}.base-snapshot-id"
	c.ArgSpecs.GetByName("update-volumes.volumes.{index}.image-label").Name = "volumes.{index}.image-label"
	c.ArgSpecs.GetByName("update-volumes.volumes.{index}.perf-iops").Name = "volumes.{index}.perf-iops"
	c.ArgSpecs.GetByName("update-private-networks.private-networks.{index}.private-network-id").Name = "private-networks.{index}.private-network-id"

	c.ArgsType = reflect.TypeFor[customUpdateTemplateRequest]()

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		args := argsI.(*customUpdateTemplateRequest)

		request := args.UpdateTemplateRequest
		request.PublicIPV4Count = &args.PublicIPv4Count
		request.PublicIPV6Count = &args.PublicIPv6Count
		request.UpdateVolumes = &instance.UpdateTemplateRequestUpdateVolumes{Volumes: args.Volumes}
		request.UpdatePrivateNetworks = &instance.UpdateTemplateRequestUpdatePrivateNetworks{
			PrivateNetworks: args.PrivateNetworks,
		}

		return runner(ctx, request)
	}

	return c
}

func TemplateDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true

	return c
}

func TemplateGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true

	return c
}

func TemplateCheckBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true

	return c
}

func TemplateCreateServerBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true
	c.ArgSpecs.GetByName("name").Default = core.RandomValueGenerator("tmpl")

	return c
}

func TemplateListBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		templateList, ok := rawResp.(*instance.ListTemplatesResponse)
		if !ok {
			return "", fmt.Errorf(
				"expected response of type instance.ListTemplatesResponse, got %T",
				rawResp,
			)
		}

		return templateList.Templates, nil
	}

	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "NAME",
				FieldName: "Name",
			},
			{
				Label:     "SERVER TYPE",
				FieldName: "ServerType",
			},
			{
				Label:     "ZONE",
				FieldName: "Zone",
			},
			{
				Label:     "TAGS",
				FieldName: "Tags",
			},
			{
				Label:     "SERVER TAGS",
				FieldName: "ServerTags",
			},
			{
				Label:     "PUBLIC IPv4 COUNT",
				FieldName: "PublicIPV4Count",
			},
			{
				Label:     "PUBLIC IPv6 COUNT",
				FieldName: "PublicIPV6Count",
			},
			{
				Label:     "SECURITY GROUP ID",
				FieldName: "SecurityGroupID",
			},
			{
				Label:     "PLACEMENT GROUP ID",
				FieldName: "PlacementGroupID",
			},
			{
				Label:     "FILESYSTEM IDS",
				FieldName: "FilesystemIDs",
			},
			{
				Label:     "PROJECT ID",
				FieldName: "ProjectID",
			},
			{
				Label:     "CREATED AT",
				FieldName: "CreatedAt",
			},
			{
				Label:     "UPDATED AT",
				FieldName: "UpdatedAt",
			},
		},
	}

	return c
}

//
// User Data
//

func TemplateSetUserDataBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true
	c.ArgSpecs.GetByName("content").CanLoadFile = true

	return c
}

func TemplateGetUserDataBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true

	return c
}

func TemplateListUserDataKeysBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}

		keyList, ok := rawResp.(*instance.ListTemplateUserDataKeysResponse)
		if !ok {
			return "", fmt.Errorf(
				"expected response of type instance.ListTemplateUserDataKeysResponse, got %T",
				rawResp,
			)
		}

		return keyList.Keys, nil
	}

	return c
}

func TemplateDeleteUserDataBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true

	return c
}

//
// Cloud Init
//

func TemplateSetCloudInitBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true
	c.ArgSpecs.GetByName("content").CanLoadFile = true

	return c
}

func TemplateGetCloudInitBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("template-id").Positional = true

	return c
}
