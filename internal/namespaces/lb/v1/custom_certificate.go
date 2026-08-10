package lb

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var certificateStatusMarshalSpecs = human.EnumMarshalSpecs{
	lb.CertificateStatusError:   &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
	lb.CertificateStatusPending: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "pending"},
	lb.CertificateStatusReady:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
}

func lbCertificateMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp lb.Certificate
	certificate := tmp(i.(lb.Certificate))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "LB",
		},
	}

	str, err := human.Marshal(certificate, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

func certificateCreateBuilder(c *core.Command) *core.Command {
	leCommonNameArgSpecs := c.ArgSpecs.GetByName("letsencrypt.common-name")
	leAlternativeNames := c.ArgSpecs.GetByName("letsencrypt.subject-alternative-name.{index}")
	customCertificateArgSpecs := c.ArgSpecs.GetByName("custom-certificate.certificate-chain")

	leCommonNameArgSpecs.Required = false
	leCommonNameArgSpecs.Name = "letsencrypt-common-name"
	leCommonNameArgSpecs.ConflictWith(customCertificateArgSpecs)

	leAlternativeNames.Name = "letsencrypt-alternative-name.{index}"
	leCommonNameArgSpecs.ConflictWith(customCertificateArgSpecs)

	customCertificateArgSpecs.Name = "custom-certificate-chain"
	customCertificateArgSpecs.Required = false

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

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
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
			res, err := runner(ctx, createCertificateRequest)
			if err != nil {
				return nil, err
			}

			if len(res.(*lb.Certificate).LB.Tags) != 0 &&
				res.(*lb.Certificate).LB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}

			return res, nil
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
			res, err := runner(ctx, createCertificateRequest)
			if err != nil {
				return nil, err
			}

			if len(res.(*lb.Certificate).LB.Tags) != 0 &&
				res.(*lb.Certificate).LB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}

			return res, nil
		}

		return nil, &core.CliError{
			Err: errors.New("missing required argument"),
			Hint: fmt.Sprintf(
				"You need to specify %s or %s",
				leCommonNameArgSpecs.Name,
				customCertificateArgSpecs.Name,
			),
			Code: 1,
		}
	}

	return c
}

func certificateGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptCertificate()

	return c
}

func certificateUpdateBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptCertificate()

	return c
}

func certificateDeleteBuilder(c *core.Command) *core.Command {
	c.Interceptor = interceptCertificate()

	return c
}

func interceptCertificate() core.CommandInterceptor {
	return func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		var getCertificate *lb.Certificate
		var err error

		client := core.ExtractClient(ctx)
		api := lb.NewZonedAPI(client)

		if _, ok := argsI.(*lb.ZonedAPIDeleteCertificateRequest); ok {
			getCertificate, err = api.GetCertificate(&lb.ZonedAPIGetCertificateRequest{
				Zone:          argsI.(*lb.ZonedAPIDeleteCertificateRequest).Zone,
				CertificateID: argsI.(*lb.ZonedAPIDeleteCertificateRequest).CertificateID,
			})
			if err != nil {
				return nil, err
			}
		}

		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		if _, ok := res.(*core.SuccessResult); ok {
			if len(getCertificate.LB.Tags) != 0 && getCertificate.LB.Tags[0] == kapsuleTag {
				return warningKapsuleTaggedMessageView(), nil
			}
		}

		return res, nil
	}
}
