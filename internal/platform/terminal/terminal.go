package terminal

import (
	"github.com/scaleway/scaleway-cli/v2/internal/platform"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Platform struct {
	UserAgent string

	cfg *scw.Config
}

func (p *Platform) ScwConfig(path string) (*scw.Config, error) {
	if p == nil {
		return nil, nil
	}

	if p.cfg == nil && path != "" {
		config, err := scw.LoadConfigFromPath(path)
		if err != nil {
			return nil, err
		}

		p.cfg = config
	}

	return p.cfg, nil
}

func (p *Platform) SetScwConfig(cfg *scw.Config) {
	p.cfg = cfg
}

func NewPlatform(useragent string) platform.Platform {
	return &Platform{
		UserAgent: useragent,
	}
}
