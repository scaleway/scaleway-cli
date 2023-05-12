package platform

import (
	"net/http"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Web struct {
	Config
	JWT string
}

func (p *Web) CreateClient(client *http.Client, _ string, _ string) (*scw.Client, error) {
	opts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithUserAgent(p.UserAgent),
		scw.WithProfile(scw.LoadEnvProfile()),
		scw.WithHTTPClient(client),
		scw.WithJWT(p.JWT),
	}

	return scw.NewClient(opts...)
}

func (p *Web) ScwConfig() *scw.Config {
	return nil
}

func (p *Web) SetScwConfig(_ *scw.Config) {}

var _ Platform = (*Web)(nil)
