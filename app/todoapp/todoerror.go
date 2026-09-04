package todoapp

import (
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/todo"
	"github.com/kongsakchai/gotemplate/pkg/serror"
)

func handleError(err error) app.Error {
	if e, ok := serror.As(err); ok {
		switch e.Code() {
		case todo.ErrUserNotFound.Code, todo.ErrTodoNotFound.Code:
			return app.NotFound(e.Code(), e.Msg(), e, e.Data...)
		}
	}
	return app.InternalError(app.InternalErrorCode, app.InternalErrorMsg, err)
}
