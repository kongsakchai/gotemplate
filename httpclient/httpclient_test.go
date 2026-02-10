package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kongsakchai/gotemplate/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequest(t *testing.T) {
	type testcase struct {
		title    string
		method   string
		url      string
		payload  any
		headers  []http.Header
		validate func(req *http.Request, err error)
	}

	testcases := []testcase{
		{
			title:   "should return error when invalid url",
			payload: "",
			url:     ":",
			method:  http.MethodPost,
			headers: []http.Header{
				{"Header": {"header pass"}},
			},
			validate: func(req *http.Request, err error) {
				assert.Nil(t, req)
				assert.Error(t, err)
			},
		},
		{
			title:  "should return error when invalid method",
			method: "mock method",
			validate: func(req *http.Request, err error) {
				assert.Nil(t, req)
				assert.Error(t, err)
			},
		},
		{
			title:   "should return error when invalid payload",
			payload: new(chan struct{}),
			validate: func(req *http.Request, err error) {
				assert.Nil(t, req)
				assert.Error(t, err)
			},
		},
		{
			title:   "should return http request when success",
			payload: "",
			method:  http.MethodPost,
			headers: []http.Header{
				{"Header": {"header pass"}},
			},
			validate: func(req *http.Request, err error) {
				assert.Equal(t, "option pass", req.Header.Get("Option"))
				assert.Equal(t, "header pass", req.Header.Get("Header"))
				assert.NoError(t, err)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			c := New(config.Config{}, func(r *http.Request, ctx context.Context) context.Context {
				r.Header.Set("Option", "option pass")
				return ctx
			})

			req, err := newRequest(context.Background(), c, tc.method, tc.url, tc.payload, tc.headers...)

			tc.validate(req, err)
		})
	}
}

func TestDoRequest(t *testing.T) {
	t.Run("should return success when client success and response is string", func(t *testing.T) {
		serve := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			b, _ := io.ReadAll(r.Body)

			assert.Equal(t, "\"some payload\"\n", string(b))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}))
		defer serve.Close()

		c := New(config.Config{})
		c.logEnable = true

		req, err := newRequest(context.Background(), c, http.MethodPost, serve.URL, "some payload")
		require.NoError(t, err)

		// act
		resp, err := doRequest[string](c, req)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "success", string(resp.Data))
	})

	t.Run("should return json when client success and response is json", func(t *testing.T) {
		serve := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			b, _ := io.ReadAll(r.Body)
			assert.Equal(t, "\"some payload\"\n", string(b))

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{\"message\":\"success\"}"))
		}))
		defer serve.Close()

		c := New(config.Config{})
		c.logEnable = true

		req, err := newRequest(context.Background(), c, http.MethodPost, serve.URL, "some payload")
		require.NoError(t, err)

		type responseStruct struct {
			Message string `json:"message"`
		}

		// act
		resp, err := doRequest[responseStruct](c, req)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "success", resp.Data.Message)
	})

	t.Run("should return error when client success but response is bad request", func(t *testing.T) {
		serve := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			b, _ := io.ReadAll(r.Body)
			assert.Equal(t, "\"some payload\"\n", string(b))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\":\"error\"}"))
		}))
		defer serve.Close()

		c := New(config.Config{})
		c.logEnable = true

		req, err := newRequest(context.Background(), c, http.MethodPost, serve.URL, "some payload")
		require.NoError(t, err)

		type responseStruct struct {
			Message string `json:"message"`
		}

		// act
		resp, err := doRequest[responseStruct](c, req)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, "error", resp.Data.Message)
	})

	t.Run("should return error when missing response", func(t *testing.T) {
		serve := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			b, _ := io.ReadAll(r.Body)
			assert.Equal(t, "\"some payload\"\n", string(b))
		}))
		defer serve.Close()

		c := New(config.Config{})
		c.logEnable = true

		req, err := newRequest(context.Background(), c, http.MethodPost, serve.URL, "some payload")
		require.NoError(t, err)

		type responseStruct struct {
			Message string `json:"message"`
		}

		// act
		_, err = doRequest[responseStruct](c, req)
		t.Log(err.Error())

		// assert
		assert.Error(t, err)
	})

	t.Run("should return error when invalid endpoint", func(t *testing.T) {
		c := New(config.Config{})
		req, err := newRequest(context.Background(), c, http.MethodPost, "error", "some payload")
		require.NoError(t, err)

		type responseStruct struct {
			Message string `json:"message"`
		}

		// act
		_, err = doRequest[responseStruct](c, req)
		t.Log(err.Error())

		// assert
		assert.Error(t, err)
	})

	t.Run("should return error when client read all error", func(t *testing.T) {
		serve := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			b, _ := io.ReadAll(r.Body)
			assert.Equal(t, "\"some payload\"\n", string(b))

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\":\"error\"}"))
		}))
		defer serve.Close()

		c := New(config.Config{})
		c.logEnable = true
		c._forceResponseNil = true

		req, err := newRequest(context.Background(), c, http.MethodPost, serve.URL, "some payload")
		require.NoError(t, err)

		type responseStruct struct {
			Message string `json:"message"`
		}

		// act
		_, err = doRequest[responseStruct](c, req)
		t.Log(err.Error())

		// assert
		assert.Error(t, err)
	})
}

func TestCallRequest(t *testing.T) {
	t.Run("should return success when client success and response is string", func(t *testing.T) {
		serve := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}))
		defer serve.Close()

		c := New(config.Config{})

		// act
		resp, err := callRequest[string](t.Context(), c, http.MethodGet, serve.URL, nil)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		assert.Equal(t, "success", string(resp.Data))
	})

	t.Run("should return error when invalid url", func(t *testing.T) {
		c := New(config.Config{})

		// act
		_, err := callRequest[string](t.Context(), c, http.MethodGet, ":", nil)

		// assert
		assert.Error(t, err)
	})
}

func TestMethod(t *testing.T) {
	type testcase struct {
		title string
		call  func(ctx context.Context, client *Client, url string) (Response[string], error)
	}

	testcases := []testcase{
		{
			title: "call method get",
			call: func(ctx context.Context, client *Client, url string) (Response[string], error) {
				return Get[string](ctx, client, url)
			},
		},
		{
			title: "call method post",
			call: func(ctx context.Context, client *Client, url string) (Response[string], error) {
				return Post[string](ctx, client, url, nil)
			},
		},
		{
			title: "call method put",
			call: func(ctx context.Context, client *Client, url string) (Response[string], error) {
				return Put[string](ctx, client, url, nil)
			},
		},
		{
			title: "call method delete",
			call: func(ctx context.Context, client *Client, url string) (Response[string], error) {
				return Delete[string](ctx, client, url, nil)
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.title, func(t *testing.T) {
			serve := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))
			defer serve.Close()

			c := New(config.Config{})

			// act
			resp, err := tc.call(t.Context(), c, serve.URL)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.Code)
			assert.Equal(t, "success", resp.Data)
		})
	}
}
