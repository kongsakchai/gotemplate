package generate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {
	u := NewUUID()
	assert.NotNil(t, u)
}

func TestGenUUID(t *testing.T) {
	u := NewUUID()
	id := u.GenUUID()
	assert.NotEmpty(t, id)
	assert.Len(t, id, 36)

	id2 := u.GenUUID()
	assert.NotEqual(t, id, id2)
}
