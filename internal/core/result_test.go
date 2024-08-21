package core_test

import (
	"testing"

	"github.com/stretchr/testify/require"

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
	require.NoError(t, err)
	assert.Equal(t, "", humanOutput)
	jsonOutput, err := result.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, []byte("{}"), jsonOutput)
}
