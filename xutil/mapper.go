package xutil

import (
	"context"
	"fmt"
	"github.com/kurzgesagtz/xgo/xerror"
)

func EnumToValue[K comparable, V any](m map[K]V, in K) (V, error) {
	v, ok := m[in]
	if !ok {
		return v, xerror.NewError(xerror.ErrCodeInvalidEnum, xerror.WithMessage(fmt.Sprintf("invalid value %v", in)))
	}
	return v, nil
}

// ValueToEnum revert value, find by value in source return key instead
func ValueToEnum[E comparable, V comparable](source map[E]V, in V) (E, error) {
	for k, v := range source {
		if v == in {
			return k, nil
		}
	}
	var defaultValue E
	return defaultValue, xerror.NewError(xerror.ErrCodeInvalidEnum, xerror.WithMessage(fmt.Sprintf("invalid enum value %v", in)))
}

func MapToSlice[I, O any](mapper func(I) (O, error), input []I) ([]O, error) {
	output := make([]O, 0)
	for _, i := range input {
		o, err := mapper(i)
		if err != nil {
			return nil, err
		}
		output = append(output, o)
	}
	return output, nil
}

func MapToSliceWithOption[I, O, P any](mapper func(I, ...P) (O, error), input []I, option ...P) ([]O, error) {
	output := make([]O, 0)
	for _, i := range input {
		o, err := mapper(i, option...)
		if err != nil {
			return nil, err
		}
		output = append(output, o)
	}
	return output, nil
}

type argc[V any] struct {
	index int
	v     *V
	err   error
}

func MapToSliceAsync[I any, O any](ctx context.Context, n int, process func(ctx context.Context, value I) (O, error), values []I) ([]O, error) {
	valuesSize := len(values)

	inputs := make([]chan argc[I], n)
	for i := range inputs {
		inputs[i] = make(chan argc[I])
	}

	outs := make([]chan *argc[O], n)
	for i := range outs {
		outs[i] = make(chan *argc[O])
	}

	for i, value := range values {
		go func() {
			inputs[i%n] <- argc[I]{index: i, v: &value}
		}()
	}

	for _, input := range inputs {
		go func() {
			for arg := range input {
				result, err := process(ctx, *arg.v)
				if err != nil {
					outs[arg.index%n] <- &argc[O]{
						err: err,
					}
					return
				}
				outs[arg.index%n] <- &argc[O]{
					index: arg.index,
					v:     &result,
				}
			}
		}()
	}

	count := make(chan error, n)

	result := make([]O, len(values))
	for _, out := range outs {
		go func() {
			for o := range out {
				if o != nil {
					if o.err != nil {
						count <- o.err
					} else {
						result[o.index] = *o.v
						count <- nil
					}
				}
			}
		}()
	}

	size := 0
	for err := range count {
		if err != nil {
			return nil, err
		}
		size++
		if size == valuesSize {
			break
		}
	}

	for _, input := range inputs {
		close(input)
	}
	for _, out := range outs {
		close(out)
	}
	close(count)

	return result, nil
}
