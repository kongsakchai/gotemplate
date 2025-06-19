package app

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAppResponse(t *testing.T) {
	t.Run("should return 200 OK when use Ok", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusOK
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"data\":\"Success\"}\n"

		Ok(ctx, "Success")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 200 OK with message when use OkWithMessage", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusOK
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"message\":\"Pong!\",\"data\":\"Success\"}\n"

		OkWithMessage(ctx, "Success", "Pong!")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 201 Created when use Created", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusCreated
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"data\":\"Created\"}\n"

		Created(ctx, "Created")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 201 Created with message when use CreatedWithMessage", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusCreated
		expectedResp := "{\"code\":\"0000\",\"status\":\"SUCCESS\",\"message\":\"Created successfully!\",\"data\":\"Created\"}\n"

		CreatedWithMessage(ctx, "Created", "Created successfully!")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 500 Internal Server Error when use FailWithError", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusInternalServerError
		expectedResp := "{\"code\":\"5000\",\"status\":\"FAIL\",\"message\":\"Internal Server Error\"}\n"

		err := InternalServer("5000", "Internal Server Error", errors.New("unexpected error"))
		Fail(ctx, err)

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})

	t.Run("should return 400 Bad Request with error message when use FailWithData", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		expectedStatus := http.StatusBadRequest
		expectedResp := "{\"code\":\"4000\",\"status\":\"FAIL\",\"message\":\"Bad Request\",\"data\":\"Invalid data\"}\n"

		err := BadRequest("4000", "Bad Request", errors.New("invalid input"))
		FailWithData(ctx, err, "Invalid data")

		assert.Equal(t, expectedStatus, rec.Code)
		assert.JSONEq(t, expectedResp, rec.Body.String())
	})
}

func TestMakeResponse(t *testing.T) {
	t.Run("should create response display with title and message", func(t *testing.T) {
		display := []string{"Title", "Message"}
		respDisplay := makeResponseDisplay(display)

		assert.Equal(t, "Title", respDisplay.Title)
		assert.Equal(t, "Message", respDisplay.Message)
	})

	t.Run("should create response display with only title", func(t *testing.T) {
		display := []string{"Only Title"}
		respDisplay := makeResponseDisplay(display)

		assert.Equal(t, "Only Title", respDisplay.Title)
		assert.Equal(t, "", respDisplay.Message)
	})

	t.Run("should create empty response display", func(t *testing.T) {
		display := []string{}
		respDisplay := makeResponseDisplay(display)

		assert.Equal(t, "", respDisplay.Title)
		assert.Equal(t, "", respDisplay.Message)
	})
}

func TestQuery(t *testing.T) {
	type testcases struct {
		title        string
		req          string
		queryType    string
		defaultValue any
		expected     any
	}
	testCases := []testcases{
		{
			title:        "should return string query",
			req:          "name=John",
			defaultValue: "",
			queryType:    "string",
			expected:     "John",
		},
		{
			title:        "should return default string query",
			req:          "name=",
			defaultValue: "John",
			queryType:    "string",
			expected:     "John",
		},
		{
			title:        "should return int query",
			req:          "age=30",
			defaultValue: int64(0),
			queryType:    "int",
			expected:     int64(30),
		},
		{
			title:        "should return default int query",
			req:          "age=",
			defaultValue: int64(25),
			queryType:    "int",
			expected:     int64(25),
		},
		{
			title:        "should return bool query",
			req:          "active=true",
			defaultValue: false,
			queryType:    "bool",
			expected:     true,
		},
		{
			title:        "should return default bool query",
			req:          "active=",
			defaultValue: false,
			queryType:    "bool",
			expected:     false,
		},
		{
			title:        "should return float query",
			req:          "price=19.99",
			defaultValue: 0.0,
			queryType:    "float",
			expected:     19.99,
		},
		{
			title:        "should return default float query",
			req:          "price=",
			defaultValue: 9.99,
			queryType:    "float",
			expected:     9.99,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/?"+tc.req, nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			var result any
			switch tc.queryType {
			case "string":
				result = QueryString(ctx, "name", tc.defaultValue.(string))
			case "int":
				result = QueryInt(ctx, "age", tc.defaultValue.(int64))
			case "bool":
				result = QueryBool(ctx, "active", tc.defaultValue.(bool))
			case "float":
				result = QueryFloat(ctx, "price", tc.defaultValue.(float64))
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestParam(t *testing.T) {
	type testcases struct {
		title        string
		req          string
		paramType    string
		defaultValue any
		expected     any
	}
	testCases := []testcases{
		{
			title:        "should return string param",
			req:          "John",
			defaultValue: "",
			paramType:    "string",
			expected:     "John",
		},
		{
			title:        "should return default string param",
			req:          "",
			defaultValue: "John",
			paramType:    "string",
			expected:     "John",
		},
		{
			title:        "should return int param",
			req:          "30",
			defaultValue: int64(0),
			paramType:    "int",
			expected:     int64(30),
		},
		{
			title:        "should return default int param",
			req:          "",
			defaultValue: int64(25),
			paramType:    "int",
			expected:     int64(25),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetParamNames("/:value")
			ctx.SetParamNames("value")
			ctx.SetParamValues(tc.req)

			var result any
			switch tc.paramType {
			case "string":
				result = ParamString(ctx, "value", tc.defaultValue.(string))
			case "int":
				result = ParamInt(ctx, "value", tc.defaultValue.(int64))
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
