package app

type Response struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type RequestContext interface {
	Bind(i any) error
	Validate(i any) error
}

func Request[Req any](ctx RequestContext) (Req, error) {
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return req, BadRequest(BadRequestCode, BadRequestMsg, err)
	}
	if err := ctx.Validate(req); err != nil {
		return req, BadRequest(InValidCode, InValidMsg, err, err)
	}
	return req, nil
}

type Context interface {
	JSON(code int, i any) error
}

func firstMsg(msg []string) string {
	if len(msg) > 0 {
		return msg[0]
	}
	return ""
}

func Ok(ctx Context, data any, msg ...string) error {
	return ctx.JSON(200, Response{
		Code:    SuccessCode,
		Success: true,
		Data:    data,
		Message: firstMsg(msg),
	})
}

func Created(ctx Context, data any, msg ...string) error {
	return ctx.JSON(201, Response{
		Code:    SuccessCode,
		Success: true,
		Data:    data,
		Message: firstMsg(msg),
	})
}

func Fail(ctx Context, err Error) error {
	return ctx.JSON(err.HTTPCode, Response{
		Code:    err.Code,
		Success: false,
		Data:    err.Data,
		Message: err.Message,
	})
}
