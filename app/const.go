package app

const (
	TraceIDKey = "traceID"
	TagKey     = "tag"

	// Common Code

	SuccessCode    = "0000"
	SuccessMessage = "success"

	BadRequestCode = "1000"
	BadRequestMsg  = "bad request"
	InValidCode    = "1001"
	InValidMsg     = "invalid request"

	DatabaseNotReadyCode = "9998"
	DatabaseNotReadyMsg  = "database is not ready"
	InternalErrorCode    = "9999"
	InternalErrorMsg     = "internal error"
)
