package errorsx

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrFileTooLarge   = errors.New("file too large")
	ErrSaveFailed     = errors.New("save file failed")
	ErrInternal       = errors.New("internal error")
)

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) Error() string {
	return e.Msg
}

func NewCodeError(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}
