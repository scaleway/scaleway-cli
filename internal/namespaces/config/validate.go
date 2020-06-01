package config

import (
	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

func validateRawArgsForConfigSet(rawArgs args.RawArgs) (key string, value string, err error) {
	if len(rawArgs) == 0 {
		return "", "", notEnoughArgsForConfigSetError()
	}
	if len(rawArgs) == 1 {
		return "", "", missingValueForConfigSetError(rawArgs[0])
	}
	if len(rawArgs) > 2 {
		return "", "", tooManyArgsForConfigSetError()
	}

	key = rawArgs[0]
	value = rawArgs[1]
	err = validateProfileKey(key)
	if err != nil {
		return "", "", err
	}
	err = validateProfileValue(key, value)
	if err != nil {
		return "", "", err
	}

	return key, value, nil
}

func validateProfileKey(fieldName string) error {
	if fieldName == sendTelemetryKey {
		return nil
	}

	_, err := getProfileField(&scw.Profile{}, fieldName)
	return err
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
