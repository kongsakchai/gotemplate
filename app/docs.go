//go:build docs

package app

type SwaggerSuccessResponse struct {
	// example: 0000
	Code string `json:"code"`
	// example: true
	Success bool `json:"success"`
}

type SwaggerErrorResponse struct {
	// example: 9999
	Code string `json:"code"`
	// example: false
	Success bool `json:"success"`
}
