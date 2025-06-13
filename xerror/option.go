package xerror

import "time"

type ErrorOptionFunc func(err *Error) *Error

func WithCaller(skip int) ErrorOptionFunc {
	return func(err *Error) *Error {
		err.Caller = caller(skip)
		return err
	}
}

func WithDetail(detail string) ErrorOptionFunc {
	return func(err *Error) *Error {
		err.Detail = detail
		return err
	}
}

func WithMessage(msg string) ErrorOptionFunc {
	return func(err *Error) *Error {
		err.Message = msg
		return err
	}
}

func WithTimestamp(t time.Time) ErrorOptionFunc {
	return func(err *Error) *Error {
		err.Timestamp = Timestamp(t.UTC())
		return err
	}
}

func WithJSONInfo(key string, value any) ErrorOptionFunc {
	return func(err *Error) *Error {
		if err.Info == nil {
			err.Info = make(map[string]any)
		}
		err.Info[key] = value
		return err
	}
}
