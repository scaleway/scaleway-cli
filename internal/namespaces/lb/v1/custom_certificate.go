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
		Zone   scw.Zone `json:"-"`
		Region scw.Region
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
		tmpRequest := argsI.(*lbCreateCertificateRequestCustom)
		var polished *lb.CreateCertificateRequest
		if tmpRequest.CustomCertificateChain != "" {
			polished = &lb.CreateCertificateRequest{
				Region: tmpRequest.Region,
				LBID:   tmpRequest.LBID,
				Name:   tmpRequest.Name,
				CustomCertificate: &lb.CreateCertificateRequestCustomCertificate{
					CertificateChain: tmpRequest.CustomCertificateChain,
				},
			}
			return runner(ctx, polished)
		}

		if tmpRequest.LetsencryptCommonName != "" {
			polished = &lb.CreateCertificateRequest{
				Region: tmpRequest.Region,
				LBID:   tmpRequest.LBID,
				Name:   tmpRequest.Name,
				Letsencrypt: &lb.CreateCertificateRequestLetsencryptConfig{
					CommonName:             tmpRequest.LetsencryptCommonName,
					SubjectAlternativeName: tmpRequest.LetsencryptAlternativeName,
				},
			}
			return runner(ctx, polished)
		}

		return nil, &core.CliError{
			Err:  fmt.Errorf("missing required argument"),
			Hint: fmt.Sprintf("You need to specify %s or %s", leCommonNameArgSpecs.Name, customeCertificateArgSpecs.Name),
			Code: 1,
		}
	}

	return c
}
