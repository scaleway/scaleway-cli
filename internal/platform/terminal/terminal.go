package terminal

import (
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type Platform struct {
	UserAgent string

	cfg *scw.Config
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
