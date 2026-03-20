package cache

import (
	"testing"

	"github.com/kongsakchai/gotemplate/config"
	"github.com/stretchr/testify/assert"
)

func TestNewRedis(t *testing.T) {
	t.Run("should ping no error when mysql connection success", func(t *testing.T) {
		rd := NewRedis(config.Redis{
			Host: "localhost",
			Port: "63799",
		})

		result := rd.Ping(t.Context())
		assert.Error(t, result.Err())
	})
}
