package marketplace

import (
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v2"
)

func init() {
	// register external types custom human marshaler func
	// for internal types use the human.Marshaler interface (see code at the end of file)
	human.RegisterMarshalerFunc(marketplace.Image{}, imageMarshalerFunc)
}

// imagesMarshalerFunc marshals marketplace.Image.
func imageMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
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
		ID:          image.ID,
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
