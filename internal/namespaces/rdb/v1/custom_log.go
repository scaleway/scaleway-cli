package rdb

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

var logStatusMarshalSpecs = human.EnumMarshalSpecs{
	rdb.InstanceLogStatusUnknown: &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
	rdb.InstanceLogStatusReady:   &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
	rdb.InstanceLogStatusCreating: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "creating",
	},
	rdb.InstanceLogStatusError: &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
}

func logPrepareBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		getResp := respI.(*rdb.PrepareInstanceLogsResponse)
		api := rdb.NewAPI(core.ExtractClient(ctx))
		readyLogs := make([]*rdb.InstanceLog, len(getResp.InstanceLogs))
		for i := range getResp.InstanceLogs {
			logs, err := api.WaitForInstanceLog(&rdb.WaitForInstanceLogRequest{
				InstanceLogID: getResp.InstanceLogs[i].ID,
				Region:        getResp.InstanceLogs[i].Region,
				Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
			if err != nil {
				return nil, err
			}
			readyLogs[i] = logs
		}
		respI.(*rdb.PrepareInstanceLogsResponse).InstanceLogs = readyLogs

		return respI.(*rdb.PrepareInstanceLogsResponse), nil
	}

	return c
}

type logDownloadArgs struct {
	InstanceID string
	From       *time.Time
	To         *time.Time
	Output     string
	Region     scw.Region
}

type logDownloadResult struct {
	Size     scw.Size `json:"size"`
	FileName string   `json:"file_name"`
}

func logDownloadResultMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	result := i.(logDownloadResult)
	sizeStr, err := human.Marshal(result.Size, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Log downloaded to %s successfully (%s written)",
		result.FileName,
		sizeStr,
	), nil
}

func getLogDefaultFileName(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	splitURL := strings.Split(u.Path, "/")
	filename := splitURL[len(splitURL)-1]

	return filename, nil
}

func logDownloadCommand() *core.Command {
	return &core.Command{
		Short:     `Download logs from a database instance`,
		Long:      `Prepare, wait for, and download logs from a Database Instance. This command automatically prepares the logs, waits for them to be ready, and downloads them to a file.`,
		Namespace: "rdb",
		Resource:  "log",
		Verb:      "download",
		ArgsType:  reflect.TypeOf(logDownloadArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "instance-id",
				Short:      `UUID of the Database Instance you want logs of`,
				Required:   true,
				Positional: true,
			},
			{
				Name:     "from",
				Short:    `Start datetime of your log. Supports absolute RFC3339 timestamps and relative times (see ` + "`scw help date`" + `).`,
				Required: false,
			},
			{
				Name:     "to",
				Short:    `End datetime of your log. Supports absolute RFC3339 timestamps and relative times (see ` + "`scw help date`" + `).`,
				Required: false,
			},
			{
				Name:  "output",
				Short: "Filename to write to",
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*logDownloadArgs)
			api := rdb.NewAPI(core.ExtractClient(ctx))

			prepareRequest := &rdb.PrepareInstanceLogsRequest{
				InstanceID: args.InstanceID,
				Region:     args.Region,
			}
			if args.From != nil {
				prepareRequest.StartDate = args.From
			}
			if args.To != nil {
				prepareRequest.EndDate = args.To
			}

			_, err = interactive.Print("Preparing logs... ")
			if err != nil {
				return nil, err
			}

			prepareResp, err := api.PrepareInstanceLogs(prepareRequest)
			if err != nil {
				return nil, err
			}

			if len(prepareResp.InstanceLogs) == 0 {
				return nil, errors.New("no logs found for the specified time range")
			}

			_, err = interactive.Println("OK")
			if err != nil {
				return nil, err
			}

			_, err = interactive.Print("Waiting for logs to be ready... ")
			if err != nil {
				return nil, err
			}

			readyLogs := make([]*rdb.InstanceLog, len(prepareResp.InstanceLogs))
			for i := range prepareResp.InstanceLogs {
				logs, err := api.WaitForInstanceLog(&rdb.WaitForInstanceLogRequest{
					InstanceLogID: prepareResp.InstanceLogs[i].ID,
					Region:        prepareResp.InstanceLogs[i].Region,
					Timeout:       scw.TimeDurationPtr(instanceActionTimeout),
					RetryInterval: core.DefaultRetryInterval,
				})
				if err != nil {
					return nil, err
				}
				readyLogs[i] = logs
			}

			_, err = interactive.Println("OK")
			if err != nil {
				return nil, err
			}

			if len(readyLogs) == 0 {
				return nil, errors.New("no logs ready after waiting")
			}

			logToDownload := readyLogs[0]
			if logToDownload.DownloadURL == nil {
				return nil, errors.New("download URL is not available")
			}

			httpClient := core.ExtractHTTPClient(ctx)

			_, err = interactive.Print("Downloading logs... ")
			if err != nil {
				return nil, err
			}

			res, err := httpClient.Get(*logToDownload.DownloadURL)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			defaultFilename, err := getLogDefaultFileName(*logToDownload.DownloadURL)
			if err != nil {
				return nil, err
			}
			filename := defaultFilename
			if args.Output != "" {
				fi, err := os.Stat(args.Output)
				if err != nil {
					if !os.IsNotExist(err) {
						return nil, err
					}
					filename = args.Output
				} else {
					switch mode := fi.Mode(); {
					case mode.IsDir():
						filename = path.Join(args.Output, defaultFilename)
					case mode.IsRegular():
						filename = args.Output
					}
				}
			}

			out, err := os.Create(filename)
			if err != nil {
				return nil, err
			}
			defer out.Close()

			size, err := io.Copy(out, res.Body)
			if err != nil {
				return nil, err
			}

			_, err = interactive.Println("OK")
			if err != nil {
				return nil, err
			}

			return logDownloadResult{
				Size:     scw.Size(size),
				FileName: filename,
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Download logs from a database instance",
				ArgsJSON: `{"instance_id": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Download logs with a time range",
				ArgsJSON: `{"instance_id": "11111111-1111-1111-1111-111111111111", "from": "2023-01-01T00:00:00Z", "to": "2023-01-02T00:00:00Z"}`,
			},
			{
				Short:    "Download logs to a specific file",
				ArgsJSON: `{"instance_id": "11111111-1111-1111-1111-111111111111", "output": "myLogs.txt"}`,
			},
		},
	}
}
