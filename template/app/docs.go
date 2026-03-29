//go:build docs

package app

type SwaggerSuccessResponse struct {
	// example: 0000
	Code string `json:"code"`
	// example: true
	Success bool `json:"success"`
}

// swagger:response errorInternalResponse
type errorInternalResponse struct {
	// in:body
	Body struct {
		// example: 9999
		Code string `json:"code"`
		// example: false
		Success bool `json:"success"`
		// example: internal error
		Message string `json:"message"`
	}
}

// swagger:response errorDatabaseResponse
type errorDatabaseResponse struct {
	// in:body
	Body struct {
		// example: 9998
		Code string `json:"code"`
		// example: false
		Success bool `json:"success"`
		// example: database is not ready
		Message string `json:"message"`
	}
}
