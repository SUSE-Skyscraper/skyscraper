package apikeys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_Generate(t *testing.T) {
	generator := NewGenerator(64*1024, 1, 4)
	apiKey, hash, err := generator.Generate()
	assert.Nil(t, err)
	assert.NotEqualf(t, "", apiKey, "apiKey should not be empty")
	assert.NotEqualf(t, "", hash, "hash should not be empty")
}
