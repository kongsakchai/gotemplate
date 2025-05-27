package app

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEchoApp(t *testing.T) {
	t.Run("should start and stop the app without error", func(t *testing.T) {
		w := sync.WaitGroup{}
		w.Add(1)

		app := NewEchoApp()
		go func() {
			defer w.Done()
			err := app.Start(":8080")
			assert.Error(t, http.ErrServerClosed, err)
		}()

		ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancle()
		assert.NoError(t, app.Shutdown(ctx))

		w.Wait()
	})
}
