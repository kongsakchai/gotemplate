package app

type Context interface {
	Param(key string) string
	Bind(obj any) error
	JSON(code int, obj any) error
	OK(obj any) error
	Created(obj any) error
	NotFound(err Error) error
	InternalServer(err Error) error
	BadRequest(err Error) error
}

type Router interface {
	GET(path string, handler func(c Context))
	POST(path string, handler func(c Context))
	PUT(path string, handler func(c Context))
	DELETE(path string, handler func(c Context))
	PATCH(path string, handler func(c Context))
}

type Consumer interface {
	Consume(path string, handler func(c Context))
}

type response struct {
	Code    string `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
