package marketplace

import (
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	// register external types custom human marshaler func
	// for internal types use the human.Marshaler interface (see code at the end of file)
	human.RegisterMarshalerFunc(marketplace.Image{}, imageMarshalerFunc)
}

// imagesMarshalerFunc marshals marketplace.Image.
func imageMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// Image
	image := i.(marketplace.Image)
	humanImage := struct {
		ID          string
		Label       string
		Name        string
		UpdatedAt   *time.Time
		CreatedAt   *time.Time
		ValidUntil  *time.Time
		Description string
	}{
		Label:       image.Label,
		Name:        image.Name,
		Description: image.Description,
		CreatedAt:   image.CreatedAt,
		UpdatedAt:   image.UpdatedAt,
		ValidUntil:  image.ValidUntil,
	}
	imageContent, err := human.Marshal(humanImage, opt)
	if err != nil {
		return "", err
	}

	// Concatenate
	return terminal.Style("Image:", color.Bold) + "\n" +
		imageContent, nil
}

func uniqueZones(zones []scw.Zone) []scw.Zone {
	u := make([]scw.Zone, 0, len(zones))
	m := make(map[scw.Zone]bool)
	for _, val := range zones {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

func uniqueStrings(strs []string) []string {
	u := make([]string, 0, len(strs))
	m := make(map[string]bool)
	for _, val := range strs {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
