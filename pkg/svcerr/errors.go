package svcerr

type apiErr string

const (
	ErrNotAuthorized apiErr = "not authorized"
	ErrBadRequest    apiErr = "bad request"
	ErrNotFound      apiErr = "not found"
	ErrConflict      apiErr = "conflict"
	ErrInternal      apiErr = "internal error"
)

type Error struct {
	Message string
	RawErr  error
	AppErr  apiErr
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.RawErr
}

func NewError(msg string, rawErr error, appErr apiErr) *Error {
	return &Error{
		Message: msg,
		RawErr:  rawErr,
		AppErr:  appErr,
	}
}
