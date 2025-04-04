package terminal

import (
	"github.com/scaleway/scaleway-cli/v2/internal/platform"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Platform struct {
	cfg       *scw.Config
	UserAgent string
}

func (p *Platform) ScwConfig() *scw.Config {
	if p == nil {
		return nil
	}

	return p.cfg
}

func (p *Platform) SetScwConfig(cfg *scw.Config) {
	p.cfg = cfg
}

func NewPlatform(useragent string) platform.Platform {
	return &Platform{
		UserAgent: useragent,
	}
}
