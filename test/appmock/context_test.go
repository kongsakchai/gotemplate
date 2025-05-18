package appmock

import (
	"context"
	"log/slog"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContext(t *testing.T) {
	t.Run("TestContext", func(t *testing.T) {
		ctx := &Context{}
		ctx.On("Query", "key").Return("value")
		ctx.On("Param", "key").Return("value")
		ctx.On("Bind", mock.Anything).Return(nil)
		ctx.On("JSON", 200, mock.Anything).Return(nil)
		ctx.On("OK", mock.Anything).Return(nil)
		ctx.On("OKWithMessage", "message", mock.Anything).Return(nil)
		ctx.On("Created", mock.Anything).Return(nil)
		ctx.On("CreatedWithMessage", "message", mock.Anything).Return(nil)
		ctx.On("Error", mock.Anything).Return(nil)
		ctx.On("Ctx").Return(context.Background())
		ctx.On("Set", "key", "value")
		ctx.On("Get", "key").Return("value")
		ctx.On("Logger").Return(slog.Default())

		assert.Equal(t, "value", ctx.Query("key"))
		assert.Equal(t, "value", ctx.Param("key"))
		assert.NoError(t, ctx.Bind(nil))
		assert.NoError(t, ctx.JSON(200, nil))
		assert.NoError(t, ctx.OK(nil))
		assert.NoError(t, ctx.OKWithMessage("message", nil))
		assert.NoError(t, ctx.Created(nil))
		assert.NoError(t, ctx.CreatedWithMessage("message", nil))
		assert.NoError(t, ctx.Error(app.NewError(404, "0000", "1111")))
		assert.Equal(t, context.Background(), ctx.Ctx())
		ctx.Set("key", "value")
		assert.Equal(t, "value", ctx.Get("key"))
		assert.Equal(t, slog.Default(), ctx.Logger())

		ctx.AssertExpectations(t)
	})
}
