package pretty_errors

import (
	"fmt"
)

type PrettyError struct {
	Message  string
	Path     string
	Level    ErrorLevel
	Details  map[string]any
	InnerErr error
	nested   []PrettyError
}

func (er PrettyError) WithInner(err error) PrettyError {
	er.InnerErr = err
	perr, ok := er.InnerErr.(PrettyError)
	if !ok {
		return er
	}
	er.nested = append(er.nested, perr)
	return er
}

func (er PrettyError) Inners() []PrettyError {
	return er.nested
}

func (e PrettyError) Error() string {
	var detailsStr string
	for k, v := range e.Details {
		detailsStr = fmt.Sprintf("%s,%s=%+v", detailsStr, k, v)
	}

	return fmt.Sprintf(
		"[%s] from [%s] with [%s]",
		e.Message, e.Path, detailsStr[1:])
}

func (e PrettyError) Unwrap() error {
	return e.InnerErr
}

type ErrorLevel string

var (
	UnexpectedOperation ErrorLevel = "UnexpectedOperation"
	UnrecoverableError  ErrorLevel = "UnrecoverableError"
	NoResources         ErrorLevel = "NoResources"
	OutOfBounds         ErrorLevel = "OutOfBounds"
	NoAccess            ErrorLevel = "NoAccess"
)
