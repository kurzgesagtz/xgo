package xerror

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"runtime"
)

const xErrorAppNameKey = "APP_NAME" // default is localhost
var appName = "local"

func init() {
	if aName, ok := os.LookupEnv(xErrorAppNameKey); ok {
		appName = aName
	}
}

func NewError(code string, fn ...ErrorOptionFunc) *Error {
	c := caller(2)
	cs := callers(2)
	err := &Error{
		Code:      code,
		AppName:   appName,
		Caller:    c,
		Callers:   cs,
		Detail:    c,
		Timestamp: TimestampNow(),
	}
	for _, optionFunc := range fn {
		err = optionFunc(err)
	}
	return err
}

func newError(code string, fn ...ErrorOptionFunc) *Error {
	c := caller(3)
	cs := callers(3)
	err := &Error{
		Code:      code,
		AppName:   appName,
		Caller:    c,
		Callers:   cs,
		Detail:    c,
		Timestamp: TimestampNow(),
	}
	for _, optionFunc := range fn {
		err = optionFunc(err)
	}
	return err
}

func callers(skip int) *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

func caller(skip int) string {
	_, file, no, ok := runtime.Caller(skip)
	if ok {
		return fmt.Sprintf("%s:%d", file, no)
	}
	return ""
}

// IsErrorCode Check xerror has same code
func IsErrorCode(err error, code string) bool {
	if err == nil {
		return false
	}
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.Code == code
	} else {
		// if it isn't api.Error then make it as internal xerror
		return ErrCodeInternalError == code
	}
}

// GetErrorMessage - Get xerror message
func GetErrorMessage(err error) string {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr.Message
	}

	return ""
}
