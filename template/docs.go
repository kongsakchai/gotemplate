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

import "github.com/kongsakchai/gotemplate/template/app"

// swagger:route GET /health common none
// Health check endpoint.
// responses:
//   200: healthResponse
//   500: errorDatabaseResponse

// swagger:response healthResponse
type healthResponseWrapper struct {
	// in:body
	Body struct {
		app.SwaggerSuccessResponse
		// example: health
		Message string `json:"message"`
	}
}

// swagger:route GET /metrics common none
// Metrics endpoint.
// responses:
//   200: metricsResponse

// swagger:response metricsResponse
type metricsResponseWrapper struct {
	// in:body
	Body struct {
		app.SwaggerSuccessResponse
		Data struct {
			/*
							{
				  "code": "0000",
				  "success": true,
				  "data": {
				    "alloc": "0.77 MB",
				    "heapIdle": "2.34 MB",
				    "heapInuse": "1.38 MB",
				    "heapReleased": "2.28 MB",
				    "stackInuse": "0.28 MB",
				    "stackSys": "0.28 MB",
				    "sysAlloc": "8.02 MB",
				    "totalAlloc": "0.77 MB"
				  }
				}
			*/
			// example: 0.77 MB
			Alloc string `json:"alloc"`
			// example: 2.34 MB
			HeapIdle string `json:"heapIdle"`
			// example: 1.38 MB
			HeapInuse string `json:"heapInuse"`
			// example: 2.28 MB
			HeapReleased string `json:"heapReleased"`
			// example: 0.28 MB
			StackInuse string `json:"stackInuse"`
			// example: 0.28 MB
			StackSys string `json:"stackSys"`
			// example: 0.28 MB
			SysAlloc string `json:"sysAlloc"`
			// example: 0.28 MB
			TotalAlloc string `json:"totalAlloc"`
			// example: 0.28 MB

		} `json:"data"`
	}
}
