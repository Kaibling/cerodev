package errs

import (
	"errors"

	"github.com/kaibling/cerodev/errs/msg"
)

func ServiceBuildError(serviceName string, err error) string {
	return "could not build service " + serviceName + ": " + err.Error()
}

func ErrMsg(msg string, err error) string {
	return msg + ": " + err.Error()
}

var (
	ErrWrongCredentials = errors.New(msg.WrongCredentials)
	ErrInvalidToken     = errors.New(msg.InvalidToken)

	ErrContainerNotInProvider = errors.New(msg.ContainerNotInProvider)
)
