package tutorial_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/tutorial"
	"github.com/stretchr/testify/assert"
)

func TestListTutorials(t *testing.T) {
	tutorials := tutorial.ListTutorials()

	assert.NotEmpty(t, tutorials)

	var found bool
	for _, tut := range tutorials {
		if tut.ID == "getting-started" {
			found = true
			assert.True(t, tut.Recommended, "getting-started should be recommended")
			assert.Equal(t, tutorial.DifficultyBeginner, tut.Difficulty)
		}
	}
	assert.True(t, found, "should have getting-started tutorial")
}

func TestGetTutorial(t *testing.T) {
	tut, ok := tutorial.GetTutorial("getting-started")
	assert.True(t, ok)
	assert.Equal(t, "Getting Started", tut.Title)
	assert.NotEmpty(t, tut.Steps)
}

func TestGetTutorial_NotFound(t *testing.T) {
	_, ok := tutorial.GetTutorial("nonexistent")
	assert.False(t, ok)
}

func Test_TutorialList(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(),
		Cmd:      "scw tutorial list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_TutorialResume_NoProgress(t *testing.T) {
	t.Run("no-progress", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(),
		Cmd:      "scw tutorial resume",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}
