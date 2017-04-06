package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ImagesAPI interface {
	GetImages() (*[]MarketImage, error)
	GetImage(id string) (*ScalewayImage, error)
	PostImage(volumeID string, name string, bootscript string, arch string) (string, error)

	GetMarketPlaceImages(uuidImage string) (*MarketImages, error)
	PostMarketPlaceImage(images MarketImage) error
	PutMarketPlaceImage(uudiImage string, images MarketImage) error
	DeleteMarketPlaceImage(uudImage string) error

	GetMarketPlaceImageVersions(uuidImage, uuidVersion string) (*MarketVersions, error)
	PostMarketPlaceImageVersion(uuidImage string, version MarketVersion) error
	PutMarketPlaceImageVersion(uuidImage, uuidVersion string, version MarketVersion) error
	DeleteMarketPlaceImageVersion(uuidImage, uuidVersion string) error

	GetMarketPlaceLocalImages(uuidImage, uuidVersion, uuidLocalImage string) (*MarketLocalImages, error)
	PostMarketPlaceLocalImage(uuidImage, uuidVersion, uuidLocalImage string, local MarketLocalImage) error
	PutMarketPlaceLocalImage(uuidImage, uuidVersion, uuidLocalImage string, local MarketLocalImage) error
	DeleteMarketPlaceLocalImage(uuidImage, uuidVersion, uuidLocalImage string) error
}

// ResolveImage attempts to find a matching Identifier for the input string
func (s *ScalewayAPI) ResolveImage(needle string) (ScalewayResolverResults, error) {
	images, err := s.Cache.LookUpImages(needle, true)
	if err != nil {
		return images, err
	}
	if len(images) == 0 {
		if _, err = s.GetImages(); err != nil {
			return nil, err
		}
		images, err = s.Cache.LookUpImages(needle, true)
	}
	return images, err
}

// GetImages gets the list of images from the ScalewayAPI
func (s *ScalewayAPI) GetImages() (*[]MarketImage, error) {
	images, err := s.GetMarketPlaceImages("")
	if err != nil {
		return nil, err
	}
	s.Cache.ClearImages()
	for i, image := range images.Images {
		if image.CurrentPublicVersion != "" {
			for _, version := range image.Versions {
				if version.ID == image.CurrentPublicVersion {
					for _, localImage := range version.LocalImages {
						images.Images[i].Public = true
						s.Cache.InsertImage(localImage.ID, localImage.Zone, localImage.Arch, image.Organization.ID, image.Name, image.CurrentPublicVersion)
					}
				}
			}
		}
	}
	values := url.Values{}
	values.Set("organization", s.Organization)
	resp, err := s.GetResponsePaginate(s.computeAPI, "images", values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var OrgaImages ScalewayImages

	if err = json.Unmarshal(body, &OrgaImages); err != nil {
		return nil, err
	}

	for _, orgaImage := range OrgaImages.Images {
		images.Images = append(images.Images, MarketImage{
			Categories:           []string{"MyImages"},
			CreationDate:         orgaImage.CreationDate,
			CurrentPublicVersion: orgaImage.Identifier,
			ModificationDate:     orgaImage.ModificationDate,
			Name:                 orgaImage.Name,
			Public:               false,
			MarketVersions: MarketVersions{
				Versions: []MarketVersionDefinition{
					{
						CreationDate:     orgaImage.CreationDate,
						ID:               orgaImage.Identifier,
						ModificationDate: orgaImage.ModificationDate,
						MarketLocalImages: MarketLocalImages{
							LocalImages: []MarketLocalImageDefinition{
								{
									Arch: orgaImage.Arch,
									ID:   orgaImage.Identifier,
									// TODO: fecth images from ams1 and par1
									Zone: s.Region,
								},
							},
						},
					},
				},
			},
		})
		s.Cache.InsertImage(orgaImage.Identifier, s.Region, orgaImage.Arch, orgaImage.Organization, orgaImage.Name, "")
	}
	return &images.Images, nil
}

// GetImage gets an image from the ScalewayAPI
func (s *ScalewayAPI) GetImage(imageID string) (*ScalewayImage, error) {
	resp, err := s.GetResponsePaginate(s.computeAPI, "images/"+imageID, url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var oneImage ScalewayOneImage

	if err = json.Unmarshal(body, &oneImage); err != nil {
		return nil, err
	}
	// FIXME owner, title
	s.Cache.InsertImage(oneImage.Image.Identifier, s.Region, oneImage.Image.Arch, oneImage.Image.Organization, oneImage.Image.Name, "")
	return &oneImage.Image, nil
}

// DeleteImage deletes a image
func (s *ScalewayAPI) DeleteImage(imageID string) error {
	defer s.Cache.RemoveImage(imageID)
	resp, err := s.DeleteResponse(s.computeAPI, fmt.Sprintf("images/%s", imageID))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := s.handleHTTPError([]int{http.StatusNoContent}, resp); err != nil {
		return err
	}
	return nil
}

// PostImage creates a new image
func (s *ScalewayAPI) PostImage(volumeID string, name string, bootscript string, arch string) (string, error) {
	definition := ScalewayImageDefinition{
		SnapshotIDentifier: volumeID,
		Name:               name,
		Organization:       s.Organization,
		Arch:               arch,
	}
	if bootscript != "" {
		definition.DefaultBootscript = &bootscript
	}

	resp, err := s.PostResponse(s.computeAPI, "images", definition)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusCreated}, resp)
	if err != nil {
		return "", err
	}
	var image ScalewayOneImage

	if err = json.Unmarshal(body, &image); err != nil {
		return "", err
	}
	// FIXME region, arch, owner, title
	s.Cache.InsertImage(image.Image.Identifier, "", image.Image.Arch, image.Image.Organization, image.Image.Name, "")
	return image.Image.Identifier, nil
}

// GetImageID returns exactly one image matching
func (s *ScalewayAPI) GetImageID(needle, arch string) (*ScalewayImageIdentifier, error) {
	// Parses optional type prefix, i.e: "image:name" -> "name"
	_, needle = parseNeedle(needle)

	images, err := s.ResolveImage(needle)
	if err != nil {
		return nil, fmt.Errorf("Unable to resolve image %s: %s", needle, err)
	}
	images = FilterImagesByArch(images, arch)
	images = FilterImagesByRegion(images, s.Region)
	if len(images) == 1 {
		return &ScalewayImageIdentifier{
			Identifier: images[0].Identifier,
			Arch:       images[0].Arch,
			// FIXME region, owner hardcoded
			Region: images[0].Region,
			Owner:  "",
		}, nil
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("No such image (zone %s, arch %s) : %s", s.Region, arch, needle)
	}
	return nil, showResolverResults(needle, images)
}

// FilterImagesByArch removes entry that doesn't match with architecture
func FilterImagesByArch(res ScalewayResolverResults, arch string) (ret ScalewayResolverResults) {
	if arch == "*" {
		return res
	}
	for _, result := range res {
		if result.Arch == arch {
			ret = append(ret, result)
		}
	}
	return
}

// FilterImagesByRegion removes entry that doesn't match with region
func FilterImagesByRegion(res ScalewayResolverResults, region string) (ret ScalewayResolverResults) {
	if region == "*" {
		return res
	}
	for _, result := range res {
		if result.Region == region {
			ret = append(ret, result)
		}
	}
	return
}

// GetMarketPlaceImages returns images from marketplace
func (s *ScalewayAPI) GetMarketPlaceImages(uuidImage string) (*MarketImages, error) {
	resp, err := s.GetResponsePaginate(MarketplaceAPI, fmt.Sprintf("images/%s", uuidImage), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var ret MarketImages

	if uuidImage != "" {
		ret.Images = make([]MarketImage, 1)

		var img MarketImage

		if err = json.Unmarshal(body, &img); err != nil {
			return nil, err
		}
		ret.Images[0] = img
	} else {
		if err = json.Unmarshal(body, &ret); err != nil {
			return nil, err
		}
	}
	return &ret, nil
}

// GetMarketPlaceImageVersions returns image version
func (s *ScalewayAPI) GetMarketPlaceImageVersions(uuidImage, uuidVersion string) (*MarketVersions, error) {
	resp, err := s.GetResponsePaginate(MarketplaceAPI, fmt.Sprintf("images/%v/versions/%s", uuidImage, uuidVersion), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var ret MarketVersions

	if uuidImage != "" {
		var version MarketVersion
		ret.Versions = make([]MarketVersionDefinition, 1)

		if err = json.Unmarshal(body, &version); err != nil {
			return nil, err
		}
		ret.Versions[0] = version.Version
	} else {
		if err = json.Unmarshal(body, &ret); err != nil {
			return nil, err
		}
	}
	return &ret, nil
}

// GetMarketPlaceImageCurrentVersion return the image current version
func (s *ScalewayAPI) GetMarketPlaceImageCurrentVersion(uuidImage string) (*MarketVersion, error) {
	resp, err := s.GetResponsePaginate(MarketplaceAPI, fmt.Sprintf("images/%v/versions/current", uuidImage), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var ret MarketVersion

	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}
	return &ret, nil
}

// GetMarketPlaceLocalImages returns images from local region
func (s *ScalewayAPI) GetMarketPlaceLocalImages(uuidImage, uuidVersion, uuidLocalImage string) (*MarketLocalImages, error) {
	resp, err := s.GetResponsePaginate(MarketplaceAPI, fmt.Sprintf("images/%v/versions/%s/local_images/%s", uuidImage, uuidVersion, uuidLocalImage), url.Values{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := s.handleHTTPError([]int{http.StatusOK}, resp)
	if err != nil {
		return nil, err
	}
	var ret MarketLocalImages
	if uuidLocalImage != "" {
		var localImage MarketLocalImage
		ret.LocalImages = make([]MarketLocalImageDefinition, 1)

		if err = json.Unmarshal(body, &localImage); err != nil {
			return nil, err
		}
		ret.LocalImages[0] = localImage.LocalImages
	} else {
		if err = json.Unmarshal(body, &ret); err != nil {
			return nil, err
		}
	}
	return &ret, nil
}

// PostMarketPlaceImage adds new image
func (s *ScalewayAPI) PostMarketPlaceImage(images MarketImage) error {
	resp, err := s.PostResponse(MarketplaceAPI, "images/", images)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusAccepted}, resp)
	return err
}

// PostMarketPlaceImageVersion adds new image version
func (s *ScalewayAPI) PostMarketPlaceImageVersion(uuidImage string, version MarketVersion) error {
	resp, err := s.PostResponse(MarketplaceAPI, fmt.Sprintf("images/%v/versions", uuidImage), version)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusAccepted}, resp)
	return err
}

// PostMarketPlaceLocalImage adds new local image
func (s *ScalewayAPI) PostMarketPlaceLocalImage(uuidImage, uuidVersion, uuidLocalImage string, local MarketLocalImage) error {
	resp, err := s.PostResponse(MarketplaceAPI, fmt.Sprintf("images/%v/versions/%s/local_images/%v", uuidImage, uuidVersion, uuidLocalImage), local)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusAccepted}, resp)
	return err
}

// PutMarketPlaceImage updates image
func (s *ScalewayAPI) PutMarketPlaceImage(uudiImage string, images MarketImage) error {
	resp, err := s.PutResponse(MarketplaceAPI, fmt.Sprintf("images/%v", uudiImage), images)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// PutMarketPlaceImageVersion updates image version
func (s *ScalewayAPI) PutMarketPlaceImageVersion(uuidImage, uuidVersion string, version MarketVersion) error {
	resp, err := s.PutResponse(MarketplaceAPI, fmt.Sprintf("images/%v/versions/%v", uuidImage, uuidVersion), version)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// PutMarketPlaceLocalImage updates local image
func (s *ScalewayAPI) PutMarketPlaceLocalImage(uuidImage, uuidVersion, uuidLocalImage string, local MarketLocalImage) error {
	resp, err := s.PostResponse(MarketplaceAPI, fmt.Sprintf("images/%v/versions/%s/local_images/%v", uuidImage, uuidVersion, uuidLocalImage), local)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusOK}, resp)
	return err
}

// DeleteMarketPlaceImage deletes image
func (s *ScalewayAPI) DeleteMarketPlaceImage(uudImage string) error {
	resp, err := s.DeleteResponse(MarketplaceAPI, fmt.Sprintf("images/%v", uudImage))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}

// DeleteMarketPlaceImageVersion delete image version
func (s *ScalewayAPI) DeleteMarketPlaceImageVersion(uuidImage, uuidVersion string) error {
	resp, err := s.DeleteResponse(MarketplaceAPI, fmt.Sprintf("images/%v/versions/%v", uuidImage, uuidVersion))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}

// DeleteMarketPlaceLocalImage deletes local image
func (s *ScalewayAPI) DeleteMarketPlaceLocalImage(uuidImage, uuidVersion, uuidLocalImage string) error {
	resp, err := s.DeleteResponse(MarketplaceAPI, fmt.Sprintf("images/%v/versions/%s/local_images/%v", uuidImage, uuidVersion, uuidLocalImage))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = s.handleHTTPError([]int{http.StatusNoContent}, resp)
	return err
}
