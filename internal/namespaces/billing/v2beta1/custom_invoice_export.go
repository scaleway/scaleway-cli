package billing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
)

type billingExportRequest struct {
	billing.ExportInvoicesRequest
	FilePath     string
	ForceReplace bool
}

func invoiceExportBuilder(command *core.Command) *core.Command {
	command.ArgsType = reflect.TypeOf(billingExportRequest{})
	command.ArgSpecs = core.ArgSpecs{
		{
			Name:       "organization-id",
			Short:      `Organization ID. If specified, only invoices from this Organization will be returned`,
			Required:   false,
			Positional: false,
		},
		{
			Name:       "billing-period-start-after",
			Short:      `Return only invoice with start date greater than billing_period_start`,
			Required:   false,
			Positional: false,
		},
		{
			Name:       "billing-period-start-before",
			Short:      `Return only invoice with start date less than billing_period_start`,
			Required:   false,
			Positional: false,
		},
		{
			Name:       "invoice-type",
			Short:      `Invoice type. It can either be ` + "`" + `periodic` + "`" + ` or ` + "`" + `purchase` + "`" + ``,
			Required:   false,
			Positional: false,
		},
		{
			Name:       "file-path",
			Short:      `Wanted file path`,
			Required:   false,
			Positional: false,
			Default:    core.DefaultValueSetter("./"),
		},
		{
			Name:       "file-type",
			Short:      `Wanted file extension`,
			Required:   false,
			Positional: false,
			Default:    core.DefaultValueSetter(billing.ExportInvoicesRequestFileTypeCsv.String()),
		},
		{
			Name:       "force-replace",
			Short:      `Force file replacement`,
			Required:   false,
			Positional: false,
			Default:    core.DefaultValueSetter("false"),
		},
	}
	command.Run = billingExportRun
	command.PreValidateFunc = func(ctx context.Context, argsI interface{}) error {
		args := argsI.(*billingExportRequest)
		askPrompt := false

		dir, file := getDirFile(args.FilePath)
		if len(file) > 0 {
			fileExtension := filepath.Ext(file)
			if extensionOnFile := checkExportInvoiceExt(fileExtension); !extensionOnFile {
				return errors.New("file has not supported extension")
			}
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}

		if args.ForceReplace {
			return nil
		}

		defaultFileName := fmt.Sprintf(
			"%s-%s.%s",
			invoiceDefaultPrefix,
			time.Now().Format("2006-01"),
			args.FileType,
		)
		// read only in the parent path
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			// case default name on directory
			if len(file) == 0 && e.Name() == defaultFileName {
				file = defaultFileName
				askPrompt = true
			}

			if file == e.Name() {
				askPrompt = true
			}
		}

		if askPrompt {
			_, _ = interactive.PrintlnWithoutIndent(`
					Current file exist is located at ` + terminal.Style(args.FilePath, color.Faint))
			overrideFile, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
				Prompt:       fmt.Sprintf("Do you want to override the current file: %s ?", file),
				DefaultValue: false,
				Ctx:          ctx,
			})
			if err != nil {
				return err
			}
			if !overrideFile {
				return errors.New("export file canceled")
			}
		}

		return nil
	}

	return command
}

func billingExportRun(ctx context.Context, argsI interface{}) (interface{}, error) {
	argsExport := argsI.(*billingExportRequest)

	request := &billing.ExportInvoicesRequest{
		OrganizationID:           argsExport.OrganizationID,
		BillingPeriodStartAfter:  argsExport.BillingPeriodStartAfter,
		BillingPeriodStartBefore: argsExport.BillingPeriodStartBefore,
		InvoiceType:              argsExport.InvoiceType,
		FileType:                 argsExport.FileType,
	}

	client := core.ExtractClient(ctx)
	billingAPI := billing.NewAPI(client)
	resp, err := billingAPI.ExportInvoices(request)
	if err != nil {
		return nil, err
	}

	fileName := fileNameWithoutExtTrimSuffix(argsExport.FilePath)
	// check file path argument
	fInfo, err := os.Stat(argsExport.FilePath)
	if err == nil {
		if fInfo.IsDir() {
			// case when filepath is a directory: join default name with custom path
			defaultFileName := fmt.Sprintf(
				"%s-%s",
				invoiceDefaultPrefix,
				time.Now().Format("2006-01"),
			)
			pathAbs, err := filepath.Abs(argsExport.FilePath)
			if err != nil {
				return nil, err
			}
			fileName = filepath.Join(pathAbs, defaultFileName)
		}
	}
	// add supported extension
	fileName = addExportExt(fileName, resp.ContentType)

	fileOutput, err := os.Create(fileName)
	if err != nil {
		dir, file := getDirFile(fileName)

		return nil, fmt.Errorf("unavailable to create file %s on directory %s", file, dir)
	}
	defer fileOutput.Close()

	_, err = io.Copy(fileOutput, resp.Content)
	if err != nil {
		return nil, err
	}

	return &core.SuccessResult{Message: "Invoice exported successfully"}, nil
}

func addExportExt(fileName, contentType string) string {
	if contentType == "text/csv" {
		fileName += ".csv"
	}

	return fileName
}

func checkExportInvoiceExt(ext string) bool {
	return ext == ".csv"
}
