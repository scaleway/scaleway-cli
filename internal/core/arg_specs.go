package core

import "github.com/scaleway/scaleway-sdk-go/scw"

type ArgSpecs []*ArgSpec

func (s ArgSpecs) GetByName(name string) *ArgSpec {
	for _, spec := range s {
		if spec.Name == name {
			return spec
		}
	}
	return nil
}

type ArgSpec struct {
	// Name of the argument.
	Name string

	// Short description.
	Short string

	// Required defines whether the argument is required.
	Required bool

	// Default is the argument default value.
	Default DefaultFunc

	// EnumValues contains all possible values of an enum.
	EnumValues []string

	// AutoCompleteFunc is used to autocomplete possible values for a given argument.
	AutoCompleteFunc AutoCompleteArgFunc

	// ValidateFunc validates an argument.
	ValidateFunc ArgSpecValidateFunc
}

type DefaultFunc func() (value string, doc string)

func ZoneArgSpec(zones ...scw.Zone) *ArgSpec {
	enumValues := []string(nil)
	for _, zone := range zones {
		enumValues = append(enumValues, zone.String())
	}
	return &ArgSpec{
		Name:       "zone",
		Short:      "Zone to target. If none is passed will use default zone from the config",
		EnumValues: enumValues,
	}
}

func RegionArgSpec(regions ...scw.Region) *ArgSpec {
	enumValues := []string(nil)
	for _, region := range regions {
		enumValues = append(enumValues, region.String())
	}
	return &ArgSpec{
		Name:       "region",
		Short:      "Region to target. If none is passed will use default region from the config",
		EnumValues: enumValues,
	}
}

func OrganizationIDArgSpec() *ArgSpec {
	return &ArgSpec{
		Name:         "organization-id",
		Short:        "Organization ID to use. If none is passed will use default organization ID from the config",
		ValidateFunc: ValidateOrganisationID(),
	}
}

func OrganizationArgSpec() *ArgSpec {
	return &ArgSpec{
		Name:         "organization",
		Short:        "Organization ID to use. If none is passed will use default organization ID from the config",
		ValidateFunc: ValidateOrganisationID(),
	}
}
