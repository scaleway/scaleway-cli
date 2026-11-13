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
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/karrick/tparse/v2"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	backupActionTimeout = 20 * time.Minute
)

var backupStatusMarshalSpecs = human.EnumMarshalSpecs{
	rdb.DatabaseBackupStatusUnknown: &human.EnumMarshalSpec{
		Attribute: color.Faint,
		Value:     "unknown",
	},
	rdb.DatabaseBackupStatusCreating: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "creating",
	},
	rdb.DatabaseBackupStatusReady: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "ready",
	},
	rdb.DatabaseBackupStatusRestoring: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "restoring",
	},
	rdb.DatabaseBackupStatusDeleting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "deleting",
	},
	rdb.DatabaseBackupStatusError: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "error",
	},
	rdb.DatabaseBackupStatusExporting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
		Value:     "exporting",
	},
}

type backupWaitRequest struct {
	BackupID string
	Region   scw.Region
	Timeout  time.Duration
}

func backupWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a backup to reach a stable state`,
		Long:      `Wait for a backup to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "rdb",
		Resource:  "backup",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeFor[backupWaitRequest](),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := rdb.NewAPI(core.ExtractClient(ctx))

			return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
				DatabaseBackupID: argsI.(*backupWaitRequest).BackupID,
				Region:           argsI.(*backupWaitRequest).Region,
				Timeout:          new(argsI.(*backupWaitRequest).Timeout),
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
			core.WaitTimeoutArgSpec(backupActionTimeout),
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
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))

		return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
			DatabaseBackupID: respI.(*rdb.DatabaseBackup).ID,
			Region:           respI.(*rdb.DatabaseBackup).Region,
			Timeout:          &timeout,
			RetryInterval:    core.DefaultRetryInterval,
		})
	}
	c.Examples = []*core.Example{
		{
			Short: "Create a backup and wait until it is ready",
			Raw:   "scw rdb backup create name=foobar expires-at=2032-01-02T15:04:05Z instance-id=11111111-1111-1111-1111-111111111111 database-name=rdb --wait",
		},
	}

	return c
}

func backupExportBuilder(c *core.Command) *core.Command {
	timeout := backupActionTimeout
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
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
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
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

func backupListBuilder(c *core.Command) *core.Command {
	type backupListArgs struct {
		Name           *string
		OrderBy        rdb.ListDatabaseBackupsRequestOrderBy
		InstanceID     *string
		ProjectID      *string
		OrganizationID *string
		Region         scw.Region
		CreatedAfter   *time.Time
		CreatedBefore  *time.Time
	}

	c.ArgsType = reflect.TypeFor[backupListArgs]()

	// Insert created-after / created-before date filters before region, keeping
	// existing arguments (name, order-by, instance-id, project-id, organization-id)
	// unchanged and still passed to the API request.
	c.ArgSpecs.AddBefore("region", &core.ArgSpec{
		Name: "created-before",
		Short: "Only list backups created before this date. " +
			"Supports absolute RFC3339 timestamps and relative times (see `scw help date`).",
		Required: false,
	})
	c.ArgSpecs.AddBefore("region", &core.ArgSpec{
		Name: "created-after",
		Short: "Only list backups created after this date. " +
			"Supports absolute RFC3339 timestamps and relative times (see `scw help date`).",
		Required: false,
	})

	c.Run = func(ctx context.Context, argsI any) (any, error) {
		args := argsI.(*backupListArgs)

		client := core.ExtractClient(ctx)
		api := rdb.NewAPI(client)

		req := &rdb.ListDatabaseBackupsRequest{
			Name:           args.Name,
			OrderBy:        args.OrderBy,
			InstanceID:     args.InstanceID,
			ProjectID:      args.ProjectID,
			OrganizationID: args.OrganizationID,
			Region:         args.Region,
		}

		opts := []scw.RequestOption{scw.WithAllPages()}
		if req.Region == scw.Region(core.AllLocalities) {
			opts = append(opts, scw.WithRegions(api.Regions()...))
			req.Region = ""
		}

		resp, err := api.ListDatabaseBackups(req, opts...)
		if err != nil {
			return nil, err
		}

		backups := resp.DatabaseBackups

		// Client-side filtering on CreatedAt using parsed absolute or relative dates.
		if args.CreatedAfter != nil || args.CreatedBefore != nil {
			filtered := make([]*rdb.DatabaseBackup, 0, len(backups))
			for _, b := range backups {
				createdAt := b.CreatedAt

				if args.CreatedAfter != nil {
					if createdAt == nil || createdAt.Before(*args.CreatedAfter) {
						continue
					}
				}

				if args.CreatedBefore != nil {
					if createdAt == nil || createdAt.After(*args.CreatedBefore) {
						continue
					}
				}

				filtered = append(filtered, b)
			}
			backups = filtered
		}

		return backups, nil
	}

	type customBackup struct {
		ID           string                   `json:"ID"`
		InstanceID   string                   `json:"instance_ID"`
		DatabaseName string                   `json:"database_name"`
		Name         string                   `json:"name"`
		Status       rdb.DatabaseBackupStatus `json:"status"`
		Size         *scw.Size                `json:"size"`
		ExpiresAt    *time.Time               `json:"expires_at"`
		CreatedAt    *time.Time               `json:"created_at"`
		UpdatedAt    *time.Time               `json:"updated_at"`
		InstanceName string                   `json:"instance_name"`
		IsExported   bool                     `json:"is_exported"`
		URLExpired   bool                     `json:"url_expired"`
		Region       scw.Region               `json:"region"`
		SameRegion   bool                     `json:"same_region"`
	}

	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "Name",
				FieldName: "Name",
			},
			{
				Label:     "Created At",
				FieldName: "CreatedAt",
			},
			{
				Label:     "Database Name",
				FieldName: "DatabaseName",
			},
			{
				Label:     "Size",
				FieldName: "Size",
			},
			{
				Label:     "Status",
				FieldName: "Status",
			},
			{
				Label:     "Instance ID",
				FieldName: "InstanceID",
			},
			{
				Label:     "URL Expired",
				FieldName: "URLExpired",
			},
			{
				Label:     "Is Exported",
				FieldName: "IsExported",
			},
			{
				Label:     "Expires At",
				FieldName: "ExpiresAt",
			},
			{
				Label:     "Updated At",
				FieldName: "UpdatedAt",
			},
			{
				Label:     "Region",
				FieldName: "Region",
			},
			{
				Label:     "Same Region",
				FieldName: "SameRegion",
			},
		},
	}

	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
			listBackupResp, err := runner(ctx, argsI)
			if err != nil {
				return listBackupResp, err
			}
			backupList := listBackupResp.([]*rdb.DatabaseBackup)
			var res []customBackup
			for _, backup := range backupList {
				isExported := false
				if backup.DownloadURL != nil {
					isExported = true
				}
				res = append(res, customBackup{
					ID:           backup.ID,
					InstanceID:   backup.InstanceID,
					DatabaseName: backup.DatabaseName,
					Name:         backup.Name,
					Status:       backup.Status,
					Size:         backup.Size,
					ExpiresAt:    backup.ExpiresAt,
					CreatedAt:    backup.CreatedAt,
					UpdatedAt:    backup.UpdatedAt,
					InstanceName: backup.InstanceName,
					IsExported:   isExported,
					URLExpired:   urlExpired(backup.DownloadURLExpiresAt),
					Region:       backup.Region,
					SameRegion:   backup.SameRegion,
				})
			}

			return res, nil
		},
	)

	return c
}

// urlExpired: indicates if the backup url is still valid after the indicated date.
func urlExpired(expirationDate *time.Time) bool {
	if expirationDate == nil {
		return true
	}

	return time.Now().After(*expirationDate)
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

func backupResultMarshallerFunc(i any, _ *human.MarshalOpt) (string, error) {
	backupResult := i.(backupDownloadResult)
	sizeStr, err := human.Marshal(backupResult.Size, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Backup downloaded to %s successfully (%s written)",
		backupResult.FileName,
		sizeStr,
	), nil
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
		ArgsType:  reflect.TypeFor[backupDownloadArgs](),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*backupDownloadArgs)
			api := rdb.NewAPI(core.ExtractClient(ctx))
			backupRequest := &rdb.WaitForDatabaseBackupRequest{
				DatabaseBackupID: args.BackupID,
				Region:           args.Region,
				Timeout:          new(backupActionTimeout),
				RetryInterval:    core.DefaultRetryInterval,
			}

			backup, err := api.WaitForDatabaseBackup(backupRequest)
			if err != nil {
				return nil, err
			}
			needExport := backup.DownloadURL == nil
			if needExport {
				_, err = interactive.Print("Exporting backup... ")
				if err != nil {
					return nil, err
				}
				exportRequest := rdb.ExportDatabaseBackupRequest{
					DatabaseBackupID: args.BackupID,
					Region:           args.Region,
				}
				_, err = api.ExportDatabaseBackup(&exportRequest)
				if err != nil {
					return nil, err
				}
			}

			backup, err = api.WaitForDatabaseBackup(backupRequest)
			if err != nil {
				return nil, err
			}
			if needExport {
				_, err = interactive.Println("OK")
			}
			if err != nil {
				return nil, err
			}

			httpClient := core.ExtractHTTPClient(ctx)

			if backup.DownloadURL == nil {
				return nil, errors.New("download URL is still nil after export")
			}

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

type backupPruneArgs struct {
	InstanceID string
	Name       *string
	KeepLast   *uint32
	KeepWithin *string
	DryRun     bool
	Region     scw.Region
}

type backupPruneItem struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	CreatedAt *time.Time `json:"created_at"`
}

type backupPruneResult struct {
	DeletedIDs []string          `json:"deleted_ids"`
	Deleted    int               `json:"deleted"`
	Remaining  int               `json:"remaining"`
	DryRun     bool              `json:"dry_run"`
	Backups    []backupPruneItem `json:"backups"`
}

func backupPruneResultMarshalerFunc(i any, _ *human.MarshalOpt) (string, error) {
	res := i.(backupPruneResult)
	action := "Deleted"
	idsLabel := "Deleted"
	if res.DryRun {
		action = "Would delete"
		idsLabel = "Would delete"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s %d backup(s), %d remaining\n", action, res.Deleted, res.Remaining)
	if len(res.Backups) == 0 {
		fmt.Fprintf(&b, "%s:\n  (none)", idsLabel)

		return b.String(), nil
	}

	fmt.Fprintf(&b, "%s:\n", idsLabel)
	for _, item := range res.Backups {
		createdAt := "-"
		if item.CreatedAt != nil {
			createdAt = item.CreatedAt.Format(time.RFC3339)
		}
		fmt.Fprintf(&b, "  %s  %s  %s\n", item.ID, item.Name, createdAt)
	}

	return strings.TrimRight(b.String(), "\n"), nil
}

// parseKeepWithinDuration parses a positive duration (e.g. 30d, 1w2d) using the same
// relative time units as `scw help date`. A leading sign is ignored; 30d and -30d both
// mean a 30-day window.
func parseKeepWithinDuration(value string) (time.Duration, error) {
	if value == "" {
		return 0, errors.New("empty duration")
	}

	signed := value
	if value[0] != '+' && value[0] != '-' {
		signed = "+" + value
	}

	now := time.Now()
	if args.TestForceNow != nil {
		now = *args.TestForceNow
	}

	parsed, err := tparse.ParseWithMap(time.RFC3339, "t"+signed, map[string]time.Time{"t": now})
	if err != nil {
		return 0, fmt.Errorf("invalid keep-within duration %q: %w", value, err)
	}

	duration := parsed.Sub(now)
	if duration < 0 {
		duration = -duration
	}
	if duration == 0 {
		return 0, errors.New("keep-within duration must be greater than 0")
	}

	return duration, nil
}

func latestBackupTime(backups []*rdb.DatabaseBackup, now time.Time) *time.Time {
	var latest *time.Time
	for _, backup := range backups {
		if backup.CreatedAt == nil || !backup.CreatedAt.Before(now) {
			continue
		}
		if latest == nil || backup.CreatedAt.After(*latest) {
			latest = backup.CreatedAt
		}
	}

	return latest
}

// backupsToDelete applies restic-like retention: keep if any policy matches (OR),
// delete otherwise. keep-within is relative to the newest backup (not wall clock).
func backupsToDelete(
	backups []*rdb.DatabaseBackup,
	keepLast *uint32,
	keepWithin *time.Duration,
	now time.Time,
) []*rdb.DatabaseBackup {
	sorted := slices.Clone(backups)
	slices.SortFunc(sorted, func(a, b *rdb.DatabaseBackup) int {
		switch {
		case a.CreatedAt == nil && b.CreatedAt == nil:
			return 0
		case a.CreatedAt == nil:
			return 1
		case b.CreatedAt == nil:
			return -1
		default:
			return b.CreatedAt.Compare(*a.CreatedAt)
		}
	})

	var withinCutoff *time.Time
	if keepWithin != nil {
		if latest := latestBackupTime(sorted, now); latest != nil {
			cutoff := latest.Add(-*keepWithin)
			withinCutoff = &cutoff
		}
	}

	out := make([]*rdb.DatabaseBackup, 0)
	for i, backup := range sorted {
		keep := false
		if keepLast != nil && uint32(i) < *keepLast {
			keep = true
		}
		if withinCutoff != nil && backup.CreatedAt != nil && backup.CreatedAt.After(*withinCutoff) {
			keep = true
		}
		if !keep {
			out = append(out, backup)
		}
	}

	return out
}

func backupPruneCommand() *core.Command {
	return &core.Command{
		Short: `Delete old database backups according to a retention policy`,
		Long: `Delete old database backups according to a retention policy.

Retention matches restic forget semantics: a backup is kept if it matches any
policy (OR). Use keep-last to retain the N most recent backups, and/or
keep-within to retain backups newer than a duration relative to the newest
backup (e.g. 30d). See also scw help date for duration units.

dry-run=true previews deletions without calling the API delete endpoint.`,
		Namespace: "rdb",
		Resource:  "backup",
		Verb:      "prune",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeFor[backupPruneArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "instance-id",
				Short:    "UUID of the Database Instance",
				Required: true,
			},
			{
				Name:  "name",
				Short: "Only consider backups with this name",
			},
			{
				Name:  "keep-last",
				Short: "Number of most recent backups to keep",
			},
			{
				Name: "keep-within",
				Short: "Keep backups newer than this duration relative to the newest backup " +
					"(e.g. 30d). See `scw help date` for units.",
			},
			{
				Name:    "dry-run",
				Short:   "Preview which backups would be deleted without deleting them",
				Default: core.DefaultValueSetter("false"),
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*backupPruneArgs)

			if args.KeepLast == nil && args.KeepWithin == nil {
				return nil, &core.CliError{
					Message: "missing retention policy",
					Hint:    "Set keep-last and/or keep-within (e.g. keep-last=5 or keep-within=30d)",
				}
			}

			var keepWithin *time.Duration
			if args.KeepWithin != nil {
				duration, err := parseKeepWithinDuration(*args.KeepWithin)
				if err != nil {
					return nil, &core.CliError{
						Err:     err,
						Message: "invalid keep-within value",
						Hint:    "Use a duration like 30d or 1w2d (see `scw help date`)",
					}
				}
				keepWithin = &duration
			}

			api := rdb.NewAPI(core.ExtractClient(ctx))
			listReq := &rdb.ListDatabaseBackupsRequest{
				Region:     args.Region,
				InstanceID: &args.InstanceID,
				Name:       args.Name,
				OrderBy:    rdb.ListDatabaseBackupsRequestOrderByCreatedAtDesc,
			}

			resp, err := api.ListDatabaseBackups(listReq, scw.WithAllPages(), scw.WithContext(ctx))
			if err != nil {
				return nil, fmt.Errorf("could not list backups: %w", err)
			}

			toDelete := backupsToDelete(resp.DatabaseBackups, args.KeepLast, keepWithin, time.Now())
			deletedIDs := make([]string, 0, len(toDelete))
			items := make([]backupPruneItem, 0, len(toDelete))

			for _, backup := range toDelete {
				deletedIDs = append(deletedIDs, backup.ID)
				items = append(items, backupPruneItem{
					ID:        backup.ID,
					Name:      backup.Name,
					CreatedAt: backup.CreatedAt,
				})
				if args.DryRun {
					continue
				}

				_, err = api.DeleteDatabaseBackup(&rdb.DeleteDatabaseBackupRequest{
					Region:           args.Region,
					DatabaseBackupID: backup.ID,
				}, scw.WithContext(ctx))
				if err != nil {
					return nil, fmt.Errorf("could not delete backup %s: %w", backup.ID, err)
				}
			}

			return backupPruneResult{
				DeletedIDs: deletedIDs,
				Deleted:    len(deletedIDs),
				Remaining:  len(resp.DatabaseBackups) - len(deletedIDs),
				DryRun:     args.DryRun,
				Backups:    items,
			}, nil
		},
		Examples: []*core.Example{
			{
				Short:    "Keep only the 5 most recent backups named monthly",
				ArgsJSON: `{"instance_id": "11111111-1111-1111-1111-111111111111", "name": "monthly", "keep_last": 5}`,
			},
			{
				Short:    "Keep backups within 30 days of the newest, dry-run",
				ArgsJSON: `{"instance_id": "11111111-1111-1111-1111-111111111111", "keep_within": "30d", "dry_run": true}`,
			},
		},
	}
}
