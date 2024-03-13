package core_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	result := &core.SuccessResult{
		Empty:    true,
		Details:  "dummy",
		Message:  "dummy",
		Resource: "dummy",
		Verb:     "dummy",
	}

	humanOutput, err := result.MarshalHuman()
	assert.NoError(t, err)
	assert.Equal(t, "", humanOutput)
	jsonOutput, err := result.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("{}"), jsonOutput)
}
