package billing

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2alpha1"
)

type billingDownloadRequest struct {
	billing.DownloadInvoiceRequest
	// extra arguments
	FileName string
	FilePath string
}

func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func buildDownloadCommand(command *core.Command) *core.Command {
	command.ArgsType = reflect.TypeOf(billingDownloadRequest{})
	command.ArgSpecs = core.ArgSpecs{
		{
			Name:       "invoice-id",
			Short:      `Invoice ID`,
			Required:   true,
			Deprecated: false,
			Positional: false,
		},
		{
			Name:       "file-type",
			Short:      `Wanted file type`,
			Required:   false,
			Deprecated: false,
			Positional: false,
			EnumValues: []string{"pdf"},
		},
		{
			Name:       "file-name",
			Short:      `Wanted file name`,
			Required:   false,
			Deprecated: false,
			Positional: false,
		},
		{
			Name:       "file-path",
			Short:      `Wanted file locality`,
			Required:   false,
			Deprecated: false,
			Positional: false,
		},
	}
	command.Run = billingDownloadRun

	return command
}

func billingDownloadRun(ctx context.Context, argsI interface{}) (interface{}, error) {
	argsDownload := argsI.(*billingDownloadRequest)

	request := &billing.DownloadInvoiceRequest{
		InvoiceID: argsDownload.InvoiceID,
		FileType:  argsDownload.FileType,
	}

	client := core.ExtractClient(ctx)
	billingAPI := billing.NewAPI(client)
	resp, err := billingAPI.DownloadInvoice(request)
	if err != nil {
		return nil, err
	}

	// set path
	filePath, err := filepath.Abs(argsDownload.FilePath)
	if err != nil {
		return nil, err
	}

	// set file name
	fileName := resp.Name
	if len(argsDownload.FileName) > 0 {
		fileName = fileNameWithoutExtTrimSuffix(argsDownload.FileName)
	}

	// check content-type
	switch resp.ContentType {
	case "application/pdf":
		fileName = fmt.Sprintf("%s.pdf", fileName)
	}

	fileOutput, err := os.Create(filepath.Join(filePath, fileName))
	if err != nil {
		return nil, err
	}
	defer fileOutput.Close()

	_, err = io.Copy(fileOutput, resp.Content)
	if err != nil {
		return nil, err
	}

	return &core.SuccessResult{Empty: true}, nil
}
