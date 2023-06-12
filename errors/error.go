package errors

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	NoType              ErrorType = ""
	BadRequest          ErrorType = "BAD_REQUEST"
	NotFound            ErrorType = "NOT_FOUND"
	Unexpected          ErrorType = "UNEXPECTED"
	ExternalAPIError    ErrorType = "EXTERNAL_API_ERROR"
	DataProcessingError ErrorType = "DATA_PROCESSING_ERROR"
)

type ErrorType string
type Operation string

type Error struct {
	operations []Operation
	errorType  ErrorType
	err        error
	severity   logrus.Level
}

func (e Error) Error() string {
	return e.err.Error()
}

func (e Error) ErrorType() ErrorType {
	return e.errorType
}

func (e Error) Severity() logrus.Level {
	return e.severity
}

func (e *Error) Operations() []Operation {
	return e.operations
}

func (e *Error) WithOperation(operation Operation) *Error {
	e.operations = append(e.operations, operation)
	return e
}

// New creates a new customError
func (errorType ErrorType) New(msg string, severity logrus.Level) error {
	return Error{errorType: errorType, err: errors.New(msg), severity: severity}
}

// New creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, severity logrus.Level, args ...interface{}) error {
	return Error{errorType: errorType, err: fmt.Errorf(msg, args...), severity: severity}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string, severity logrus.Level) error {
	return errorType.Wrapf(err, msg, severity)
}

// Wrap creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, severity logrus.Level, args ...interface{}) error {
	return Error{errorType: errorType, err: errors.Wrapf(err, msg, args...), severity: severity}
}

// New creates a no type error
func New(msg string) error {
	return Error{errorType: NoType, err: errors.New(msg), severity: logrus.InfoLevel}
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	return Error{errorType: NoType, err: errors.New(fmt.Sprintf(msg, args...)), severity: logrus.InfoLevel}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// GetType returns the error type
func GetTypeAndLogLevel(err error) (ErrorType, logrus.Level) {
	if customErr, ok := err.(Error); ok {
		return customErr.errorType, customErr.severity
	}

	return NoType, logrus.InfoLevel
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(Error); ok {
		return Error{
			errorType: customErr.errorType,
			err:       wrappedError,
		}
	}

	return Error{errorType: NoType, err: wrappedError}
}
