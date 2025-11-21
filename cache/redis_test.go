package cache

import (
	"strings"
	"testing"

	"github.com/kongsakchai/gotemplate/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewRedis(t *testing.T) {
	t.Run("should ping no error when mysql connection success", func(t *testing.T) {
		t.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
		// init container
		ct, err := testcontainers.Run(
			t.Context(),
			"redis:alpine",

			testcontainers.WithProvider(testcontainers.ProviderPodman),
			testcontainers.WithExposedPorts("6379/tcp"),
			testcontainers.WithWaitStrategy(
				wait.ForListeningPort("6379/tcp"),
			),
		)
		require.NoError(t, err)

		endpoint, err := ct.Endpoint(t.Context(), "")
		require.NoError(t, err)

		// check panic
		defer func() {
			p := recover()
			assert.Nil(t, p)
		}()

		p := strings.Split(endpoint, ":")
		rd := NewRedis(config.Redis{
			Host: p[0],
			Port: p[1],
		})

		result := rd.Ping(t.Context())
		assert.NoError(t, result.Err())

		testcontainers.CleanupContainer(t, ct)
	})
}
