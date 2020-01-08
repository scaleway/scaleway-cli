package core

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

var (
	ZoneArgSpec = &ArgSpec{
		Name:     "zone",
		Short:    "Zone to target. If none is passed will use default zone from the config",
		Required: false,
		// TODO: Default:          nil,
		// TODO: EnumValues:       nil,
		// TODO: AutoCompleteFunc: nil,
		// TODO: ValidateFunc:
	}
	RegionArgSpec = &ArgSpec{
		Name:     "region",
		Short:    "Region to target. If none is passed will use default region from the config",
		Required: false,
		// TODO: Default:          nil,
		// TODO: EnumValues:       nil,
		// TODO: AutoCompleteFunc: nil,
		// TODO: ValidateFunc:
	}
	OrganizationIDArgSpec = &ArgSpec{
		Name:     "organization_id",
		Short:    "Organization ID to use. If none is passed will use default organization ID from the config",
		Required: false,
		// TODO: Default:          nil,
		// TODO: EnumValues:       nil,
		// TODO: AutoCompleteFunc: nil,
		// TODO: ValidateFunc:
	}
	OrganizationArgSpec = &ArgSpec{
		Name:     "organization",
		Short:    "Organization ID to use. If none is passed will use default organization ID from the config",
		Required: false,
		// TODO: Default:          nil,
		// TODO: EnumValues:       nil,
		// TODO: AutoCompleteFunc: nil,
		// TODO: ValidateFunc:
	}
)
