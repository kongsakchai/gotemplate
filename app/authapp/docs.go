//go:build docs

package authapp

import (
	"github.com/kongsakchai/gotemplate/app"
)

// swagger:route POST /api/v1/auth/register auth registerUser
//
// Register a new user account.
//
// Creates a user with the given username and password.
// The password is hashed before storing.
//
// responses:
//   201: authRegisterResponse
//   400: authErrorResponse
//   500: errorInternalResponse

// swagger:route POST /api/v1/auth/login auth loginUser
//
// Login with username and password.
//
// Returns a JWT token to be sent in the Authorization header
// (format: `Bearer <token>`) of protected endpoints.
//
// responses:
//   200: authLoginResponse
//   400: authErrorResponse
//   500: errorInternalResponse

// Parameters for registering a new user.
//
// swagger:parameters registerUser
type RegisterUserParams struct {
	// User credentials.
	// in: body
	Body registerRequest
}

// Parameters for logging in.
//
// swagger:parameters loginUser
type LoginUserParams struct {
	// User credentials.
	// in: body
	Body loginRequest
}

// User account created successfully (no data payload).
//
// swagger:response authRegisterResponse
type AuthRegisterResponse struct {
	// in: body
	Body struct {
		app.SwaggerSuccessResponse
	}
}

// Login result containing the JWT token.
//
// swagger:response authLoginResponse
type AuthLoginResponse struct {
	// in: body
	Body struct {
		app.SwaggerSuccessResponse
		Data struct {
			// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxIiwidXNlcm5hbWUiOiJhZG1pbiJ9.signature
			Token string `json:"token"`
		} `json:"data"`
	}
}

// Error response for auth endpoints.
//
// Possible codes:
//   - 1000: bad request (cannot read request body)
//   - 1001: invalid request (validation failed)
//   - 2000: user already exists
//   - 2001: user not found
//   - 2002: invalid password
//
// swagger:response authErrorResponse
type AuthErrorResponse struct {
	// in: body
	Body struct {
		// example: 2002
		Code string `json:"code"`
		// example: false
		Success bool `json:"success"`
		// example: invalid password
		Message string `json:"message"`
	}
}
