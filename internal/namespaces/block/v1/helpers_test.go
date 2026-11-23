package block_test

import "github.com/scaleway/scaleway-cli/v2/core"

// deleteVolume deletes a volume.
// metaKey must be the key in meta where the volume state is stored.
func deleteVolume(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw block volume delete {{." + metaKey + ".ID}}")
}

// deleteSnapshot deletes a snapshot.
// metaKey must be the key in meta where the volume state is stored.
func deleteSnapshot(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw block snapshot delete {{." + metaKey + ".ID}}")
}
