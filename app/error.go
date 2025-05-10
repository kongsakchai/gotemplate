package app

type Error struct {
	Code    string
	Message string
	Err     error
}

func NewError(code string, msg string, err ...error) Error {
	var e error
	if len(err) > 0 {
		e = err[0]
	}

	return Error{
		Code:    code,
		Message: msg,
		Err:     e,
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}
