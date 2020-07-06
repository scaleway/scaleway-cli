package rdb

import (
	"context"
	"io/ioutil"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func certificateGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		getCertificateResp, err := runner(ctx, argsI)
		if err != nil {
			return getCertificateResp, err
		}
		getCertificate := getCertificateResp.(*scw.File)
		if b, err := ioutil.ReadAll(getCertificate.Content); err == nil {
			return string(b), nil
		}
		return nil, err
	}

	return c
}
