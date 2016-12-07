package api

import (
	"fmt"
	"net/http"
)

type UserAPI interface {
	PatchUserSSHKey(UserID string, definition ScalewayUserPatchSSHKeyDefinition) error
}

// PatchUserSSHKey updates a user
func (s *ScalewayAPI) PatchUserSSHKey(UserID string, definition ScalewayUserPatchSSHKeyDefinition) error {
	resp, err := s.PatchResponse(AccountAPI, fmt.Sprintf("users/%s", UserID), definition)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if _, err := s.handleHTTPError([]int{http.StatusOK}, resp); err != nil {
		return err
	}
	return nil
}
