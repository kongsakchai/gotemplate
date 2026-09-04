//go:build docs

package todoapp

import (
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/todo"
)

// swagger:route GET /api/v1/todos todo listTodos
//
// List all todos of the authenticated user.
//
// Security:
//   key: []
//
// responses:
//   200: todoListResponse
//   400: todoErrorResponse
//   401: unauthorizedResponse
//   500: errorInternalResponse

// swagger:route POST /api/v1/todos todo createTodo
//
// Create a new todo for the authenticated user.
//
// Security:
//   key: []
//
// responses:
//   201: todoCreatedResponse
//   400: todoErrorResponse
//   401: unauthorizedResponse
//   500: errorInternalResponse

// swagger:route PUT /api/v1/todos/{id} todo updateTodo
//
// Update an existing todo of the authenticated user.
//
// Security:
//   key: []
//
// responses:
//   200: todoUpdatedResponse
//   400: todoErrorResponse
//   401: unauthorizedResponse
//   500: errorInternalResponse

// swagger:route DELETE /api/v1/todos/{id} todo deleteTodo
//
// Delete an existing todo of the authenticated user.
//
// Security:
//   key: []
//
// responses:
//   200: todoDeletedResponse
//   400: todoErrorResponse
//   401: unauthorizedResponse
//   500: errorInternalResponse

// Parameters for creating a todo.
//
// swagger:parameters createTodo
type CreateTodoParams struct {
	// Todo properties.
	// in: body
	Body createTodoRequest
}

// Parameters for updating a todo.
//
// swagger:parameters updateTodo
type UpdateTodoParams struct {
	// Todo ID.
	// in: path
	// required: true
	ID string `json:"id"`
	// Todo properties to update (all fields optional).
	// in: body
	Body struct {
		// example: Buy groceries
		Name string `json:"name"`
		// example: Milk and eggs
		Description string `json:"description"`
		// example: pending
		Status string `json:"status"`
	}
}

// Parameters for deleting a todo.
//
// swagger:parameters deleteTodo
type DeleteTodoParams struct {
	// Todo ID.
	// in: path
	// required: true
	ID string `json:"id"`
}

// List of todos belonging to the authenticated user.
//
// swagger:response todoListResponse
type TodoListResponse struct {
	// in: body
	Body struct {
		app.SwaggerSuccessResponse
		Data []todo.Todo `json:"data"`
	}
}

// Todo created successfully (no data payload).
//
// swagger:response todoCreatedResponse
type TodoCreatedResponse struct {
	// in: body
	Body struct {
		app.SwaggerSuccessResponse
	}
}

// Todo updated successfully (no data payload).
//
// swagger:response todoUpdatedResponse
type TodoUpdatedResponse struct {
	// in: body
	Body struct {
		app.SwaggerSuccessResponse
	}
}

// Todo deleted successfully (no data payload).
//
// swagger:response todoDeletedResponse
type TodoDeletedResponse struct {
	// in: body
	Body struct {
		app.SwaggerSuccessResponse
	}
}

// Error response for todo endpoints.
//
// Possible codes:
//   - 1000: bad request (cannot read request body)
//   - 1001: invalid request (validation failed)
//   - 3000: user not found
//   - 3001: todo not found
//
// swagger:response todoErrorResponse
type TodoErrorResponse struct {
	// in: body
	Body struct {
		// example: 3001
		Code string `json:"code"`
		// example: false
		Success bool `json:"success"`
		// example: todo not found
		Message string `json:"message"`
	}
}

// Unauthorized response.
//
// The Authorization header is missing or the token is invalid/expired.
//
// Possible codes:
//   - 10002: missing token
//   - 10003: unauthorized
//   - 10004: token expired
//   - 10005: invalid token
//
// swagger:response unauthorizedResponse
type UnauthorizedResponse struct {
	// in: body
	Body struct {
		// example: 10002
		Code string `json:"code"`
		// example: false
		Success bool `json:"success"`
		// example: missing token
		Message string `json:"message"`
	}
}
