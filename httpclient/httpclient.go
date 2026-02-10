package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"time"

	"github.com/kongsakchai/gotemplate/config"
)

const (
	maxIdleConns        = 100
	maxConnsPerHost     = 100
	maxIdleConnsPerHost = 100
	timeout             = 10 * time.Second
)

type Client struct {
	*http.Client
	refIDKey string
	options  []OptionFunc

	logEnable bool

	_forceResponseNil bool //this use for only unit test
}

func New(cfg config.Config, options ...OptionFunc) *Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	c := &http.Client{
		Transport: t,
		Timeout:   timeout,
	}
	return &Client{
		Client:    c,
		options:   options,
		refIDKey:  cfg.Header.RefIDKey,
		logEnable: cfg.Log.HttpEnable,
	}
}

type Response[T any] struct {
	Code    int
	Data    T
	RawData []byte
}

func newRequest(ctx context.Context, client *Client, method, url string, payload any, headers ...http.Header) (*http.Request, error) {
	var buf bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&buf).Encode(payload); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, err
	}

	for _, header := range headers {
		maps.Copy(req.Header, header)
	}

	for _, option := range client.options {
		ctx = option(req, ctx)
	}

	return req.WithContext(ctx), nil
}

func doRequest[Resp any](client *Client, req *http.Request) (Response[Resp], error) {
	traceID, _ := req.Context().Value(client.refIDKey).(string)
	if client.logEnable {
		logHTTPRequest(traceID, req)
	}

	response := Response[Resp]{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if client._forceResponseNil {
		resp.Body.Close()
	}

	bytesResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	response.Code = resp.StatusCode
	response.RawData = bytesResponse

	if client.logEnable {
		logHTTPResponse(traceID, string(bytesResponse), resp.StatusCode, req)
	}

	if err = json.Unmarshal(bytesResponse, &response.Data); err == nil {
		response.Code = resp.StatusCode
		return response, nil
	}

	if _, ok := any(response.Data).(string); ok {
		response.Data = any(string(bytesResponse)).(Resp)
		return response, nil
	}

	return response, err
}

func logHTTPRequest(traceID string, req *http.Request) {
	body := []byte{}
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	slog.Info(
		"HTTP Client Request",
		"url", req.URL.Path,
		"method", req.Method,
		"body", string(body),
		"trace_id", traceID,
	)
}

func logHTTPResponse(traceID, data string, status int, req *http.Request) {
	if status >= http.StatusBadRequest {
		slog.Info(
			"HTTP Client Response",
			"url", req.URL.Path,
			"method", req.Method,
			"status", status,
			"body", data,
			"trace_id", traceID,
		)
	} else {
		slog.Error(
			"HTTP Client Response",
			"url", req.URL.Path,
			"method", req.Method,
			"status", status,
			"body", data,
			"trace_id", traceID,
		)
	}
}

func callRequest[Resp any](ctx context.Context, client *Client, mathod, url string, payload any, headers ...http.Header) (response Response[Resp], err error) {
	req, err := newRequest(ctx, client, http.MethodPost, url, payload, headers...)
	if err != nil {
		return response, err
	}

	return doRequest[Resp](client, req)
}

func Get[Resp any](ctx context.Context, client *Client, url string, headers ...http.Header) (Response[Resp], error) {
	return callRequest[Resp](ctx, client, http.MethodGet, url, nil, headers...)
}

func Post[Resp any](ctx context.Context, client *Client, url string, payload any, headers ...http.Header) (Response[Resp], error) {
	return callRequest[Resp](ctx, client, http.MethodPost, url, payload, headers...)
}

func Put[Resp any](ctx context.Context, client *Client, url string, payload any, headers ...http.Header) (Response[Resp], error) {
	return callRequest[Resp](ctx, client, http.MethodPut, url, payload, headers...)
}

func Delete[Resp any](ctx context.Context, client *Client, url string, payload any, headers ...http.Header) (Response[Resp], error) {
	return callRequest[Resp](ctx, client, http.MethodDelete, url, payload, headers...)
}
