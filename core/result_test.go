package core_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResult(t *testing.T) {
	result := &core.SuccessResult{
		Empty:    true,
		Details:  "dummy",
		Message:  "dummy",
		Resource: "dummy",
		Verb:     "dummy",
	}

	humanOutput, err := humanv2.Format(result)
	require.NoError(t, err)
	assert.Empty(t, humanOutput)
	jsonOutput, err := result.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, "{}", string(jsonOutput))
}
