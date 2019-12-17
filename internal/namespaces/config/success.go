package config

//
// Functions in this file can only return non-nil *core.SuccessResult.
//

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func configSetSuccess(profileAndKey string, value string) *core.SuccessResult {
	return &core.SuccessResult{
		Message: fmt.Sprintf("set %v %v successfully", profileAndKey, value),
	}
}

func configUnsetSuccess(profileAndKey string) *core.SuccessResult {
	return &core.SuccessResult{
		Message: fmt.Sprintf("unset %v successfully", profileAndKey),
	}
}

func configDeleteProfileSuccess(profileName string) *core.SuccessResult {
	return &core.SuccessResult{
		Message: fmt.Sprintf("profile '%s' deleted successfully", profileName),
	}
}

func configResetSuccess() *core.SuccessResult {
	return &core.SuccessResult{
		Message: "reset config successfully",
	}
}
