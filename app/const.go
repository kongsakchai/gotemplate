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

	// Business Code

	InvalidAgeCode          = "1001"
	InvalidAgeMsg           = "age invalid; age >= 15 and age <= 60"
	UsernameUnavailableCode = "1002"
	UsernameUnavailableMsg  = "username unavaliable"
	MemberNotFoundCode      = "1003"
	MemberNotFoundMsg       = "member not found"
)
