package marketplace

import (
	"time"

	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v2"
)

func init() {
	// register external types custom human marshaler func
	// for internal types use the human.Marshaler interface (see code at the end of file)
	humanv2.RegisterFormatFunc(func(v marketplace.Image) (string, error) {
		// projection keeps a curated subset/order in the human output
		type humanImage struct {
			ID          string
			Label       string
			Name        string
			UpdatedAt   *time.Time
			CreatedAt   *time.Time
			ValidUntil  *time.Time
			Description string
		}
		str, err := humanv2.Render(humanImage{
			ID:          v.ID,
			Label:       v.Label,
			Name:        v.Name,
			Description: v.Description,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
			ValidUntil:  v.ValidUntil,
		}, humanv2.Layout{})
		if err != nil {
			return "", err
		}

		return terminal.Style("Image:", color.Bold) + "\n" + str, nil
	})
}
