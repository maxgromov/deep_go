package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	Errors []error
}

// вернуть текст ошибки
func (e *MultiError) Error() string {
	if len(e.Errors) == 0 || e == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(strconv.Itoa(len(e.Errors)))
	builder.WriteString(" errors occured:\n")
	for _, err := range e.Errors {
		builder.WriteString("\t* " + err.Error())
	}
	builder.WriteString("\n")

	return builder.String()
}

// добавить ошибку к существующей
func Append(err error, errs ...error) *MultiError {
	if err == nil && len(errs) == 0 {
		return nil
	}

	var mErr *MultiError
	if errors.As(err, &mErr) {
		mErr.Errors = append(mErr.Errors, errs...)
		return mErr
	}

	e := make([]error, 0, len(errs)+1)
	e = append(e, errs...)

	return &MultiError{Errors: e}
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
