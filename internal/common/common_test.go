package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testcase to check if TestGenerateUuid() is able to generate uuid
func TestGenerateUuid(t *testing.T) {
	_, err := GenerateUuid()

	assert.NoError(t, err)
}
