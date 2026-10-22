//go:build wasm

package tutorial

import (
	"context"
	"fmt"
)

func RunTutorial(ctx context.Context, tutorial Tutorial, startStep int) error {
	return fmt.Errorf("tutorial system is not available in WASM builds")
}

func CheckPrerequisites(ctx context.Context, tutorial Tutorial) error {
	return nil
}

func CheckVersionMismatch(ctx context.Context, progress ProgressMap, tutorialID string) (bool, error) {
	return false, nil
}
