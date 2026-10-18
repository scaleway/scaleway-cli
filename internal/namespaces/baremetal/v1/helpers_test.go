package baremetal_test

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	baremetalSDK "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}

type OfferFilter func(*baremetalSDK.Offer) bool

func offerFilterAvailable(offer *baremetalSDK.Offer) bool {
	return offer.Stock == baremetalSDK.OfferStockAvailable
}

func offerFilterHDD(offer *baremetalSDK.Offer) bool {
	for _, d := range offer.Disks {
		if d.Type == "HDD" {
			return true
		}
	}

	return false
}

func offerFilterNVMe(offer *baremetalSDK.Offer) bool {
	for _, d := range offer.Disks {
		if d.Type == "NVMe" {
			return true
		}
	}

	return false
}

func offerFilterCompatibleOS(osID string) OfferFilter {
	return func(offer *baremetalSDK.Offer) bool {
		return !slices.Contains(offer.IncompatibleOsIDs, osID)
	}
}

// selectOffer returns the lowest-costed offer from the filtered offers
// in the given zone.
func selectOffer(zone scw.Zone, filters ...OfferFilter) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		api := baremetalSDK.NewAPI(ctx.Client)
		offers, err := api.ListOffers(&baremetalSDK.ListOffersRequest{
			Zone: zone,
		}, scw.WithAllPages())
		if err != nil {
			return fmt.Errorf("checking offers: %w", err)
		}
		var candidates []*baremetalSDK.Offer
	offerLoop:
		for _, offer := range offers.Offers {
			if offer.PricePerHour == nil {
				continue offerLoop
			}
			for _, filter := range filters {
				if filter(offer) == false {
					continue offerLoop
				}
			}
			candidates = append(candidates, offer)
		}
		if len(candidates) == 0 {
			return errors.New("no accepted offer")
		}

		slices.SortStableFunc(candidates, func(a, b *baremetalSDK.Offer) int {
			if units := a.PricePerHour.Units - b.PricePerHour.Units; units != 0 {
				return int(units)
			}

			return int(a.PricePerHour.Nanos - b.PricePerHour.Nanos)
		})
		ctx.Meta["Offer"] = candidates[0]

		return nil
	}
}

// createServerAndWait creates a baremetal instance
// register it in the context Meta at metaKey.
func createServerAndWait() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Server",
		"scw baremetal server create type={{ .Offer.Name }} zone={{ .Offer.Zone }} -w",
	)
}

func waitForServer() core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		api := baremetalSDK.NewAPI(ctx.Client)
		server := ctx.Meta["Server"].(*baremetalSDK.Server)
		_, err := api.WaitForServer(&baremetalSDK.WaitForServerRequest{
			ServerID:      server.ID,
			Zone:          server.Zone,
			Timeout:       new(baremetal.ServerActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})

		return err
	}
}

func createServer() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Server",
		"scw baremetal server create type={{ .Offer.Name }} zone={{ .Offer.Zone }}",
	)
}

func deleteServer() core.AfterFunc {
	return core.ExecAfterCmd(
		"scw baremetal server delete {{ .Server.ID }} zone={{ .Server.Zone }}",
	)
}

// add an ssh key with a given meta key
func addSSH(key string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd("Key", "scw iam ssh-key create public-key='"+key+"'")
}

// delete an ssh key with a given meta key
func deleteSSH() core.AfterFunc {
	return core.ExecAfterCmd("scw iam ssh-key delete {{ .Key.ID }}")
}
