package web

import (
	"net/http"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Platform struct {
	UserAgent        string
	JWT              string
	DefaultProjectID string
}

func (p *Platform) CreateClient(client *http.Client, _ string, _ string) (*scw.Client, error) {
	opts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithUserAgent(p.UserAgent),
		scw.WithUserAgent("cli/web"),
		scw.WithProfile(scw.LoadEnvProfile()),
		scw.WithHTTPClient(client),
		scw.WithJWT(p.JWT),
	}

	if p.DefaultProjectID != "" {
		opts = append(opts, scw.WithDefaultProjectID(p.DefaultProjectID))
	}

	return scw.NewClient(opts...)
}

func (p *Platform) ScwConfig() *scw.Config {
	return nil
}

func (p *Platform) SetScwConfig(_ *scw.Config) {}
