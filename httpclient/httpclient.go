package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
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
		Client:   c,
		options:  options,
		refIDKey: cfg.Header.RefIDKey,
	}
}

type Response[T any] struct {
	Code    int
	Data    T
	RawData []byte
}

func (c *Client) newRequest(ctx context.Context, method, url string, payload any, headers ...http.Header) (*http.Request, error) {
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

	for _, option := range c.options {
		ctx = option(req, ctx)
	}

	return req.WithContext(ctx), nil
}

func doRequest[Resp any](client *Client, req *http.Request) (Response[Resp], error) {
	traceID, _ := req.Context().Value(client.refIDKey).(string)

	body := []byte{}
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
	}
	req.Body.Close()
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	slog.Info(fmt.Sprintf("request %s", req.URL),
		"method", req.Method,
		"body", string(body),
		"trace_id", traceID,
	)

	response := Response[Resp]{}
	resp, err := client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	bytesResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		slog.Error("response", "status", resp.Status, "url", req.URL.Path, "data", bytesResponse, "trace_id", traceID)
	} else {
		slog.Info("response", "status", resp.Status, "url", req.URL.Path, "data", bytesResponse, "trace_id", traceID)
	}

	response.Code = resp.StatusCode
	response.RawData = bytesResponse

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

func Get[Resp any](ctx context.Context, client *Client, url string, headers ...http.Header) (response Response[Resp], err error) {
	req, err := client.newRequest(ctx, http.MethodGet, url, nil, headers...)
	if err != nil {
		return response, err
	}

	return doRequest[Resp](client, req)
}

func Post[Resp any](ctx context.Context, client *Client, url string, payload any, headers ...http.Header) (response Response[Resp], err error) {
	req, err := client.newRequest(ctx, http.MethodPost, url, payload, headers...)
	if err != nil {
		return response, err
	}

	return doRequest[Resp](client, req)
}
