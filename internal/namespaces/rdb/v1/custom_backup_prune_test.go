package rdb_test

import (
	"strings"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	rdbSDK "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
)

func TestBackupsToDelete(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 11, 12, 12, 0, 0, 0, time.UTC)
	backup := func(id string, age time.Duration) *rdbSDK.DatabaseBackup {
		createdAt := now.Add(-age)

		return &rdbSDK.DatabaseBackup{ID: id, CreatedAt: &createdAt}
	}

	backups := []*rdbSDK.DatabaseBackup{
		backup("new", 1*time.Hour),
		backup("mid", 10*24*time.Hour),
		backup("old", 40*24*time.Hour),
		backup("older", 60*24*time.Hour),
	}

	keep := uint32(2)
	within30d := 30 * 24 * time.Hour
	within5d := 5 * 24 * time.Hour

	t.Run("keep-last", func(t *testing.T) {
		t.Parallel()

		got := rdb.BackupsToDelete(backups, &keep, nil, now)
		assertBackupIDs(t, got, "old", "older")
	})

	t.Run("keep-within", func(t *testing.T) {
		t.Parallel()

		got := rdb.BackupsToDelete(backups, nil, &within30d, now)
		assertBackupIDs(t, got, "old", "older")
	})

	t.Run("OR keeps mid when beyond keep-last but within window", func(t *testing.T) {
		t.Parallel()

		keepOne := uint32(1)
		got := rdb.BackupsToDelete(backups, &keepOne, &within30d, now)
		assertBackupIDs(t, got, "old", "older")
	})

	t.Run("OR with narrow within still keeps by keep-last", func(t *testing.T) {
		t.Parallel()

		got := rdb.BackupsToDelete(backups, &keep, &within5d, now)
		assertBackupIDs(t, got, "old", "older")
	})

	t.Run("keep-within relative to newest not wall clock", func(t *testing.T) {
		t.Parallel()

		stale := []*rdbSDK.DatabaseBackup{
			backup("newest", 20*24*time.Hour),
			backup("edge", 30*24*time.Hour),
			backup("gone", 40*24*time.Hour),
		}
		within15d := 15 * 24 * time.Hour
		got := rdb.BackupsToDelete(stale, nil, &within15d, now)
		assertBackupIDs(t, got, "gone")
	})

	t.Run("nil CreatedAt sorts last and can match keep-last", func(t *testing.T) {
		t.Parallel()

		withNil := []*rdbSDK.DatabaseBackup{
			backup("a", 1*time.Hour),
			{ID: "nil-date"},
			backup("b", 2*time.Hour),
		}
		keepOne := uint32(1)
		got := rdb.BackupsToDelete(withNil, &keepOne, nil, now)
		assertBackupIDs(t, got, "b", "nil-date")
	})
}

func TestParseKeepWithinDuration(t *testing.T) {
	t.Parallel()

	got, err := rdb.ParseKeepWithinDuration("30d")
	if err != nil {
		t.Fatal(err)
	}
	if got < 29*24*time.Hour || got > 31*24*time.Hour {
		t.Fatalf("got %v, want ~30d", got)
	}

	gotNeg, err := rdb.ParseKeepWithinDuration("-7d")
	if err != nil {
		t.Fatal(err)
	}
	if gotNeg < 6*24*time.Hour || gotNeg > 8*24*time.Hour {
		t.Fatalf("got %v, want ~7d", gotNeg)
	}
}

func TestBackupPruneResultMarshaler(t *testing.T) {
	t.Parallel()

	createdAt := time.Date(2026, 11, 13, 10, 0, 0, 0, time.UTC)
	res := rdb.NewBackupPruneResultForTest(
		[]string{"id-1"},
		1,
		2,
		true,
		[]rdb.BackupPruneItemForTest{
			{ID: "id-1", Name: "ux-test", CreatedAt: &createdAt},
		},
	)

	got, err := rdb.BackupPruneResultMarshaler(res, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "Would delete 1 backup(s), 2 remaining") {
		t.Fatalf("missing summary: %s", got)
	}
	if !strings.Contains(got, "Would delete:") {
		t.Fatalf("expected dry-run label, got: %s", got)
	}
	if strings.Contains(got, "Deleted IDs:") {
		t.Fatalf("unexpected apply label in dry-run: %s", got)
	}
	if !strings.Contains(got, "id-1  ux-test  2026-11-13T10:00:00Z") {
		t.Fatalf("missing backup details: %s", got)
	}
}

func assertBackupIDs(t *testing.T, got []*rdbSDK.DatabaseBackup, want ...string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len=%d, want %d (%v)", len(got), len(want), idsOf(got))
	}
	for i := range want {
		if got[i].ID != want[i] {
			t.Fatalf("got %v, want %v", idsOf(got), want)
		}
	}
}

func idsOf(backups []*rdbSDK.DatabaseBackup) []string {
	ids := make([]string, len(backups))
	for i, backup := range backups {
		ids[i] = backup.ID
	}

	return ids
}
