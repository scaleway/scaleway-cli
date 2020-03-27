package k8s

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	poolActionTimeout = 10 * time.Minute
)

//
// Marshalers
//

// poolStatusMarshalerFunc marshals a k8s.PoolStatus.
var (
	poolStatusAttributes = human.Attributes{
		k8s.PoolStatusScaling:   color.FgBlue,
		k8s.PoolStatusReady:     color.FgGreen,
		k8s.PoolStatusLocked:    color.FgRed,
		k8s.PoolStatusUpgrading: color.FgBlue,
		k8s.PoolStatusWarning:   color.FgHiYellow,
	}
)

const (
	poolActionCreate = iota
	poolActionUpdate
	poolActionUpgrade
	poolActionDelete
)

func poolCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionCreate)
	return c
}

func poolDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionDelete)
	return c
}

func poolUpgradeBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionUpgrade)
	return c
}

func poolUpdateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForPoolFunc(poolActionUpdate)
	return c
}

func waitForPoolFunc(action int) core.WaitFunc {
	return func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		pool, err := k8s.NewAPI(core.ExtractClient(ctx)).WaitForPool(&k8s.WaitForPoolRequest{
			Region:  respI.(*k8s.Pool).Region,
			PoolID:  respI.(*k8s.Pool).ID,
			Timeout: scw.TimeDurationPtr(poolActionTimeout),
		})
		switch action {
		case poolActionCreate:
			return pool, err
		case poolActionUpdate:
			return pool, err
		case poolActionUpgrade:
			return pool, err
		case poolActionDelete:
			if err != nil {
				// if we get a 404 here, it means the resource was successfully deleted
				notFoundError := &scw.ResourceNotFoundError{}
				responseError := &scw.ResponseError{}
				if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound || errors.As(err, &notFoundError) {
					return fmt.Sprintf("Pool %s successfully deleted.", respI.(*k8s.Pool).ID), nil
				}
			}
		}
		return nil, err
	}
}
