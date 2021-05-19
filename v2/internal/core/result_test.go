package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	result := &SuccessResult{
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
