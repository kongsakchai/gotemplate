//go:build docs

// Package classification Go Template.
//
// Documentation of our Go Template API.
//
//	Schemes: http, https
//	BasePath: /
//	Version: 1.0.0
//	Host: localhost:8080
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- key: []
//
//	SecurityDefinitions:
//	key:
//	  type: apiKey
//	  in: header
//	  name: Authorization
//
// swagger:meta
package main

import "github.com/kongsakchai/gotemplate/app"

// swagger:route GET /health common idOfHealthEndpoint
// Health check endpoint.
// responses:
//   200: healthResponse
//   500: errorHealthResponse

// swagger:response healthResponse
type healthResponseWrapper struct {
	// in:body
	Body struct {
		app.SwaggerSuccessResponse
		// example: health
		Message string `json:"message"`
	}
}

// swagger:response errorHealthResponse
type errorHealthResponseWrapper struct {
	// in:body
	Body struct {
		app.SwaggerErrorResponse
		// example: Internal Server Error
		Message string `json:"message"`
	}
}

// swagger:route GET /metrics common idOfMetricsEndpoint
// Metrics endpoint.
// responses:
//   200: metricsResponse

// swagger:response metricsResponse
type metricsResponseWrapper struct {
	// in:body
	Body struct {
		app.SwaggerSuccessResponse
		Data struct {
			// example: 1.23 MB
			Alloc string `json:"alloc"`
			// example: 4.56 MB
			TotalAlloc string `json:"totalAlloc"`
			// example: 7.89 MB
			SysAlloc string `json:"sysAlloc"`
		} `json:"data"`
	}
}
