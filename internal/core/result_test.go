package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	result := &SuccessResult{Empty: true}
	humanOutput, err := result.MarshalHuman()
	assert.NoError(t, err)
	assert.Equal(t, "", humanOutput)
}
