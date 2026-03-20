package example

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/kongsakchai/gotemplate/app"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/echotest"
	"github.com/stretchr/testify/assert"
)

func TestHandlerGetUsers(t *testing.T) {
	t.Run("should return internal error when storage users fails", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		storage.EXPECT().Users().Return(nil, errors.New("db error"))

		ctx := echotest.ContextConfig{
			Headers: http.Header{
				echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
			},
		}.ToContext(t)

		// act
		err := h.GetUsers(ctx)

		// assert
		assert.Error(t, err)
		appErr, ok := err.(app.Error)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, appErr.HTTPCode)
		assert.Equal(t, "5001", appErr.Code)
	})

	t.Run("should return ok with users data when success", func(t *testing.T) {
		// arrange
		storage := newMockStorager(t)
		h := NewHandler(storage)

		expectedUsers := []User{
			{FirstName: "john", LastName: "doe", Age: 30},
			{FirstName: "jane", LastName: "doe", Age: 25},
		}
		storage.EXPECT().Users().Return(expectedUsers, nil)

		ctx, rec := echotest.ContextConfig{
			Headers: http.Header{
				echo.HeaderContentType: []string{echo.MIMEApplicationJSON},
			},
		}.ToContextRecorder(t)

		// act
		err := h.GetUsers(ctx)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]any
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, app.SuccessCode, resp["code"])
		assert.Equal(t, true, resp["success"])

		data, ok := resp["data"].([]any)
		assert.True(t, ok)
		assert.Len(t, data, 2)
	})
}
