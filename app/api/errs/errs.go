package errs

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    ErrCode
	Message string
}

func New(code ErrCode, msg string) Error {
	return Error{Code: code, Message: msg}
}

func Newf(code ErrCode, format string, a ...any) Error {
	return Error{
		Code:    code,
		Message: fmt.Sprintf(format, a...),
	}
}

func (e Error) Error() string {
	return e.Message
}

func IsError(err error) bool {
	var errs Error
	return errors.As(err, &errs)
}

func GetError(err error) Error {
	var er Error
	if errors.As(err, &er) {
		return er
	}
	return Error{}
}
