package xerror

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

const (
	ErrCodeUnauthorized                = "UNAUTHORIZED"
	ErrCodePermissionDenied            = "PERMISSION_DENIED"
	ErrCodeInvalidRequest              = "INVALID_REQUEST"
	ErrCodeNotFound                    = "NOT_FOUND"
	ErrCodeNotImplement                = "NOT_IMPLEMENT"
	ErrCodeInternalError               = "INTERNAL_ERROR"
	ErrCodePartnerInternalError        = "PARTNER_INTERNAL_ERROR"
	ErrCodePartnerBadResponseError     = "PARTNER_BAD_RESPONSE_ERROR"
	ErrCodeServiceInternalError        = "SERVICE_INTERNAL_ERROR"
	ErrCodeServerPanic                 = "SERVER_PANIC"
	ErrCodeInvalidEnum                 = "INVALID_ENUM"
	ErrCodeClientRequestCanceled       = "CLIENT_REQUEST_CANCELED"
	ErrCodeClientRequestDeadlineExceed = "CLIENT_REQUEST_DEADLINE_EXCEED"
	ErrCodeAlreadyExists               = "ALREADY_EXISTS"
	ErrCodeRateLimitExceeded           = "RATE_LIMIT_EXCEEDED"
)

type stack []uintptr

// Error is an xerror information when api returns non-2xx.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	// Detail contains a more human-friendly message which may include long instructions to fix the error
	Caller    string         `json:"caller,omitempty"`
	Detail    string         `json:"detail,omitempty"`
	AppName   string         `json:"app_name,omitempty"`
	Info      map[string]any `json:"info,omitempty"`
	Timestamp Timestamp      `json:"timestamp"`
	Callers   *stack         `json:"callers,omitempty"`
}

// Error implements xerror.
func (err *Error) Error() string {
	return fmt.Sprintf("%s: %s", err.Code, err.Message)
}

func (err *Error) StackTrace() errors.StackTrace {
	if err.Callers == nil {
		return make([]errors.Frame, 0)
	}
	f := make([]errors.Frame, len(*err.Callers))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*err.Callers)[i])
	}
	return f
}

type Timestamp time.Time

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	s := fmt.Sprint(time.Time(*t).UnixMilli())
	return []byte(s), nil
}

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	n, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	// Convert milliseconds to seconds and nanoseconds
	*t = Timestamp(time.Unix(n/1000, (n%1000)*1_000_000))
	return nil
}

// TimestampNow is a shortcut for Timestamp(timeutils.Now()).
func TimestampNow() Timestamp {
	return Timestamp(time.Now().UTC())
}
