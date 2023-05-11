package platform

import "github.com/scaleway/scaleway-sdk-go/scw"

type Default struct {
	Config

	cfg *scw.Config
}

var _ Platform = (*Default)(nil)

func (p *Default) ScwConfig() *scw.Config {
	return p.cfg
}

func (p *Default) SetScwConfig(cfg *scw.Config) {
	p.cfg = cfg
}
