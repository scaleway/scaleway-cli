package rdb

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	backupActionTimeout = 20 * time.Minute
)

var (
	backupStatusMarshalSpecs = human.EnumMarshalSpecs{
		rdb.DatabaseBackupStatusUnknown:   &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		rdb.DatabaseBackupStatusCreating:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "creating"},
		rdb.DatabaseBackupStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		rdb.DatabaseBackupStatusRestoring: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "restoring"},
		rdb.DatabaseBackupStatusDeleting:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "deleting"},
		rdb.DatabaseBackupStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		rdb.DatabaseBackupStatusExporting: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "exporting"},
	}
)

type backupWaitRequest struct {
	BackupID string
	Region   scw.Region
}

func backupWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a backup to reach a stable state`,
		Long:      `Wait for a backup to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "rdb",
		Resource:  "backup",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(backupWaitRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := rdb.NewAPI(core.ExtractClient(ctx))
			return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
				DatabaseBackupID: argsI.(*backupWaitRequest).BackupID,
				Region:           argsI.(*backupWaitRequest).Region,
				Timeout:          scw.TimeDurationPtr(backupActionTimeout),
				RetryInterval:    core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backup-id",
				Short:      `ID of the backup you want to wait for.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a backup to reach a stable state",
				ArgsJSON: `{"backup_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func backupCreateBuilder(c *core.Command) *core.Command {
	timeout := backupActionTimeout
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
			DatabaseBackupID: respI.(*rdb.DatabaseBackup).ID,
			Region:           respI.(*rdb.DatabaseBackup).Region,
			Timeout:          &timeout,
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}

func backupExportBuilder(c *core.Command) *core.Command {
	timeout := backupActionTimeout
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
			DatabaseBackupID: respI.(*rdb.DatabaseBackup).ID,
			Region:           respI.(*rdb.DatabaseBackup).Region,
			Timeout:          &timeout,
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}

func backupRestoreBuilder(c *core.Command) *core.Command {
	timeout := backupActionTimeout
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
			DatabaseBackupID: respI.(*rdb.DatabaseBackup).ID,
			Region:           respI.(*rdb.DatabaseBackup).Region,
			Timeout:          &timeout,
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}

func getDefaultFileName(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	splitURL := strings.Split(u.Path, "/")
	filename := splitURL[len(splitURL)-1]
	return filename, nil
}

type backupDownloadResult struct {
	Size     scw.Size `json:"size"`
	FileName string   `json:"file_name"`
}

func backupResultMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	backupResult := i.(backupDownloadResult)
	sizeStr, err := human.Marshal(backupResult.Size, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Backup downloaded to %s successfully (%s written)", backupResult.FileName, sizeStr), nil
}

func backupDownloadCommand() *core.Command {
	type backupDownloadArgs struct {
		BackupID string
		Region   scw.Region
		Output   string
	}

	return &core.Command{
		Short:     `Download a backup locally`,
		Long:      `Download a backup locally.`,
		Namespace: "rdb",
		Resource:  "backup",
		Verb:      "download",
		ArgsType:  reflect.TypeOf(backupDownloadArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			args := argsI.(*backupDownloadArgs)
			api := rdb.NewAPI(core.ExtractClient(ctx))
			backup, err := api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
				DatabaseBackupID: args.BackupID,
				Region:           args.Region,
				Timeout:          scw.TimeDurationPtr(backupActionTimeout),
				RetryInterval:    core.DefaultRetryInterval,
			})
			if err != nil {
				return nil, err
			}
			if backup.DownloadURL == nil {
				return nil, fmt.Errorf("no download URL found")
			}

			httpClient := core.ExtractHTTPClient(ctx)
			res, err := httpClient.Get(*backup.DownloadURL)
			if err != nil {
				return nil, err
			}
			defer res.Body.Close()

			// Find the filename for the dump
			defaultFilename, err := getDefaultFileName(*backup.DownloadURL)
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
						// do directory stuff
						filename = path.Join(args.Output, defaultFilename)
					case mode.IsRegular():
						// do file stuff
						filename = args.Output
					}
				}
			}

			// Create the file
			out, err := os.Create(filename)
			if err != nil {
				return nil, err
			}
			defer out.Close()

			// Write the body to file
			size, err := io.Copy(out, res.Body)
			if err != nil {
				return nil, err
			}
			return backupDownloadResult{
				Size:     scw.Size(size),
				FileName: filename,
			}, nil
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backup-id",
				Short:      `ID of the backup you want to download.`,
				Required:   true,
				Positional: true,
			},
			{
				Name:  "output",
				Short: "Filename to write to",
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:    "Download a backup",
				ArgsJSON: `{"backup_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
