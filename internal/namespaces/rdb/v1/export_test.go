package rdb

import (
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

// Exported for black-box tests of backup prune helpers.

func BackupsToDelete(
	backups []*rdbSDK.DatabaseBackup,
	keepLast *uint32,
	keepWithin *time.Duration,
	now time.Time,
) []*rdbSDK.DatabaseBackup {
	return backupsToDelete(backups, keepLast, keepWithin, now)
}

func ParseKeepWithinDuration(value string) (time.Duration, error) {
	return parseKeepWithinDuration(value)
}

func BackupPruneResultMarshaler(i any, opt *human.MarshalOpt) (string, error) {
	return backupPruneResultMarshalerFunc(i, opt)
}

func NewBackupPruneResultForTest(
	deletedIDs []string,
	deleted int,
	remaining int,
	dryRun bool,
	backups []BackupPruneItemForTest,
) any {
	items := make([]backupPruneItem, 0, len(backups))
	for _, backup := range backups {
		items = append(items, backupPruneItem{
			ID:        backup.ID,
			Name:      backup.Name,
			CreatedAt: backup.CreatedAt,
		})
	}

	return backupPruneResult{
		DeletedIDs: deletedIDs,
		Deleted:    deleted,
		Remaining:  remaining,
		DryRun:     dryRun,
		Backups:    items,
	}
}

type BackupPruneItemForTest struct {
	ID        string
	Name      string
	CreatedAt *time.Time
}

// Exported for black-box tests of config snippet rendering.

func RenderPHPConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypePHP, info)
}

func RenderNodeConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypeNode, info)
}

func RenderTypeScriptConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypeTypeScript, info)
}

func RenderPythonConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypePython, info)
}

func RenderGoConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypeGo, info)
}

func RenderRustConfig(info *ConnectionInfo) (core.RawResult, error) {
	return renderConfigTemplate(rdbConfigTypeRust, info)
}
