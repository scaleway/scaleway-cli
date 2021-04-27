package lb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var (
	certificateStatusMarshalSpecs = human.EnumMarshalSpecs{
		lb.CertificateStatusError:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		lb.CertificateStatusPending: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "pending"},
		lb.CertificateStatusReady:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
	}
)

func certificateCreateBuilder(c *core.Command) *core.Command {
	leCommonNameArgSpecs := c.ArgSpecs.GetByName("letsencrypt.common-name")
	leAlternativeNames := c.ArgSpecs.GetByName("letsencrypt.subject-alternative-name.{index}")
	customeCertificateArgSpecs := c.ArgSpecs.GetByName("custom-certificate.certificate-chain")

	leCommonNameArgSpecs.Required = false
	leCommonNameArgSpecs.Name = "letsencrypt-common-name"
	leCommonNameArgSpecs.ConflictWith(customeCertificateArgSpecs)

	leAlternativeNames.Name = "letsencrypt-alternative-name.{index}"
	leCommonNameArgSpecs.ConflictWith(customeCertificateArgSpecs)

	customeCertificateArgSpecs.Name = "custom-certificate-chain"
	customeCertificateArgSpecs.Required = false

	type lbCreateCertificateRequestCustom struct {
		Zone scw.Zone `json:"-"`
		// OrganizationID with which the server will be created
		OrganizationID string `json:"organization_id"`
		// Name of the server (≠hostname)
		Name                       string `json:"name"`
		LBID                       string
		CustomCertificateChain     string
		LetsencryptCommonName      string
		LetsencryptAlternativeName []string
	}

	c.ArgsType = reflect.TypeOf(lbCreateCertificateRequestCustom{})

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*lbCreateCertificateRequestCustom)
		var createCertificateRequest *lb.ZonedAPICreateCertificateRequest
		if args.CustomCertificateChain != "" {
			createCertificateRequest = &lb.ZonedAPICreateCertificateRequest{
				Zone: args.Zone,
				LBID: args.LBID,
				Name: args.Name,
				CustomCertificate: &lb.CreateCertificateRequestCustomCertificate{
					CertificateChain: args.CustomCertificateChain,
				},
			}
			return runner(ctx, createCertificateRequest)
		}

		if args.LetsencryptCommonName != "" {
			createCertificateRequest = &lb.ZonedAPICreateCertificateRequest{
				Zone: args.Zone,
				LBID: args.LBID,
				Name: args.Name,
				Letsencrypt: &lb.CreateCertificateRequestLetsencryptConfig{
					CommonName:             args.LetsencryptCommonName,
					SubjectAlternativeName: args.LetsencryptAlternativeName,
				},
			}
			return runner(ctx, createCertificateRequest)
		}

		return nil, &core.CliError{
			Err:  fmt.Errorf("missing required argument"),
			Hint: fmt.Sprintf("You need to specify %s or %s", leCommonNameArgSpecs.Name, customeCertificateArgSpecs.Name),
			Code: 1,
		}
	}

	return c
}
