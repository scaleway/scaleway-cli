package tutorial

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

type TutorialProgress struct {
	TutorialID     string    `json:"tutorial_id"`
	CurrentStep    int       `json:"current_step"`
	Completed      bool      `json:"completed"`
	StartedAt      time.Time `json:"started_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
	CLIVersion     string    `json:"cli_version"`
}

type ProgressMap map[string]TutorialProgress

func GetProgress(progress ProgressMap, tutorialID string) (TutorialProgress, bool) {
	p, ok := progress[tutorialID]

	return p, ok
}

func SetProgress(progress ProgressMap, p TutorialProgress) {
	progress[p.TutorialID] = p
}

func ResetProgress(progress ProgressMap, tutorialID string) {
	delete(progress, tutorialID)
}

func progressFilePath(_ context.Context) string {
	configDir, err := scw.GetScwConfigDir()
	if err != nil {
		configDir = filepath.Join(os.Getenv("HOME"), ".config", "scw")
	}

	return filepath.Join(configDir, "tutorial-progress.json")
}

func LoadProgress(ctx context.Context) (ProgressMap, error) {
	path := progressFilePath(ctx)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ProgressMap{}, nil
		}

		return nil, err
	}

	var progress ProgressMap
	err = json.Unmarshal(data, &progress)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func SaveProgress(ctx context.Context, progress ProgressMap) error {
	path := progressFilePath(ctx)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0o700)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(progress, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := path + ".tmp"
	err = os.WriteFile(tmpPath, data, 0o600)
	if err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}
