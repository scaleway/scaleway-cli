package marketplace

import (
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func init() {
	// register external types custom human marshaler func
	// for internal types use the human.Marshaler interface (see code at the end of file)
	human.RegisterMarshalerFunc([]*marketplace.Image{}, imagesMarshalerFunc)
	human.RegisterMarshalerFunc(marketplace.Image{}, imageMarshalerFunc)
}

// imagesMarshalerFunc marshals []*marketplace.Image.
func imagesMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// humanServerInList is the custom Server type used for list view.
	type humanImageInList struct {
		Label            string
		Name             string
		Zones            []scw.Zone
		Archs            []string
		ModificationDate *time.Time
		CreationDate     *time.Time
	}

	images := i.([]*marketplace.Image)
	humanImages := make([]*humanImageInList, 0)
	for _, image := range images {
		zones := []scw.Zone(nil)
		archs := []string(nil)
		for _, version := range image.Versions {
			for _, localImage := range version.LocalImages {
				zones = append(zones, localImage.Zone)
				archs = append(archs, localImage.Arch)
			}
		}
		zones = uniqueZones(zones)
		archs = uniqueStrings(archs)
		sort.Strings(archs)
		sort.Slice(zones, func(i, j int) bool {
			return zones[i] < zones[j]
		})
		humanImages = append(humanImages, &humanImageInList{
			Label:            image.Label,
			Name:             image.Name,
			Zones:            zones,
			Archs:            archs,
			CreationDate:     image.CreationDate,
			ModificationDate: image.ModificationDate,
		})
	}
	return human.Marshal(humanImages, opt)
}

// imagesMarshalerFunc marshals marketplace.Image.
func imageMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	// Image
	image := i.(marketplace.Image)
	humanImage := struct {
		Label            string
		Name             string
		ModificationDate *time.Time
		CreationDate     *time.Time
		ValidUntil       *time.Time
		Description      string
	}{
		Label:            image.Label,
		Name:             image.Name,
		Description:      image.Description,
		CreationDate:     image.CreationDate,
		ModificationDate: image.ModificationDate,
		ValidUntil:       image.ValidUntil,
	}
	imageContent, err := human.Marshal(humanImage, opt)
	if err != nil {
		return "", err
	}

	// Local Images
	type humanLocalImage struct {
		ID                        string
		Zone                      scw.Zone
		Arch                      string
		CompatibleCommercialTypes []string
	}
	humanLocalImages := []humanLocalImage(nil)
	for _, version := range image.Versions {
		for _, localImage := range version.LocalImages {
			types := localImage.CompatibleCommercialTypes
			sort.Strings(localImage.CompatibleCommercialTypes)
			humanLocalImages = append(humanLocalImages, humanLocalImage{
				ID:                        localImage.ID,
				Zone:                      localImage.Zone,
				Arch:                      localImage.Arch,
				CompatibleCommercialTypes: types,
			})
		}
	}
	localImagesContent, err := human.Marshal(humanLocalImages, opt)
	if err != nil {
		return "", err
	}

	// Concatenate
	return terminal.Style("Image:", color.Bold) + "\n" +
		imageContent + "\n\n" +
		terminal.Style("Local Images:", color.Bold) + "\n" +
		localImagesContent, nil
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
