package platform

import (
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var _ Platform = (*Default)(nil)

type Default struct {
	Config

	cfg *scw.Config
}

func NewDefault(useragent string) *Default {
	return &Default{
		Config: Config{
			UserAgent: useragent,
		},
		cfg: nil,
	}
}

func (p *Default) ScwConfig() *scw.Config {
	return p.cfg
}

func (p *Default) SetScwConfig(cfg *scw.Config) {
	p.cfg = cfg
}
