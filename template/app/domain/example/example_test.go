package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDomain(t *testing.T) {
	assert.NotNil(t, New())
}
