package config

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

func validateRawArgsForConfigSet(rawArgs args.RawArgs) (profile string, key string, value string, err error) {
	if len(rawArgs) == 0 {
		return "", "", "", notEnoughArgsForConfigSetError()
	}
	if len(rawArgs) == 1 {
		return "", "", "", missingValueForConfigSetError(rawArgs[0])
	}
	if len(rawArgs) > 2 {
		return "", "", "", tooManyArgsForConfigSetError()
	}

	profileAndKey := rawArgs[0]
	value = rawArgs[1]
	profile, key, err = splitProfileKey(profileAndKey)
	if err != nil {
		return
	}
	err = validateProfileKey(key)
	if err != nil {
		return
	}
	err = validateProfileValue(key, value)
	if err != nil {
		return
	}

	return
}

func validateProfileKey(fieldName string) error {
	field := reflect.ValueOf(&scw.Profile{}).Elem().FieldByName(strcase.ToPublicGoName(fieldName))
	if !field.IsValid() {
		return invalidProfileKeyError(fieldName)
	}
	return nil
}

func validateProfileValue(fieldName string, value string) error {
	switch fieldName {
	case "default_organization_id":
		if !validation.IsOrganizationID(value) {
			return invalidDefaultOrganizationIDError(value)
		}
	case "secret_key":
		if !validation.IsSecretKey(value) {
			return core.InvalidSecretKeyError(value)
		}
	case "default_region":
		if !validation.IsRegion(value) {
			return invalidRegionError(value)
		}
	case "default_zone":
		if !validation.IsZone(value) {
			return invalidZoneError(value)
		}
	}
	return nil
}
