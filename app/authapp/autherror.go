package authapp

import (
	"github.com/kongsakchai/gotemplate/app"
	"github.com/kongsakchai/gotemplate/internal/auth"
	"github.com/kongsakchai/gotemplate/pkg/serror"
)

func handleError(err error) app.Error {
	if e, ok := serror.As(err); ok {
		switch e.Code() {
		case auth.ErrUserNotFound.Code:
			return app.NotFound(e.Code(), e.Msg(), e, e.Data...)
		case auth.ErrInvalidPassword.Code, auth.ErrUserAlreadyExists.Code:
			return app.BadRequest(e.Code(), e.Msg(), e, e.Data...)
		}
	}
	return app.InternalError(app.InternalErrorCode, app.InternalErrorMsg, err)
}
