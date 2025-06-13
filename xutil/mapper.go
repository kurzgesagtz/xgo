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
	if valuesSize == 0 {
		return []O{}, nil
	}

	// Use a smaller number of workers if there are fewer values than requested workers
	if n > valuesSize {
		n = valuesSize
	}

	inputs := make([]chan argc[I], n)
	for i := range inputs {
		inputs[i] = make(chan argc[I])
	}

	outs := make([]chan *argc[O], n)
	for i := range outs {
		outs[i] = make(chan *argc[O])
	}

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Send input values to worker channels
	go func() {
		for i, value := range values {
			select {
			case <-ctx.Done():
				return
			case inputs[i%n] <- argc[I]{index: i, v: &value}:
			}
		}
		// Close input channels after all values are sent
		for _, input := range inputs {
			close(input)
		}
	}()

	// Start worker goroutines
	for i, input := range inputs {
		go func(i int, input chan argc[I]) {
			for arg := range input {
				select {
				case <-ctx.Done():
					return
				default:
					result, err := process(ctx, *arg.v)
					if err != nil {
						outs[arg.index%n] <- &argc[O]{
							err: err,
						}
					} else {
						outs[arg.index%n] <- &argc[O]{
							index: arg.index,
							v:     &result,
						}
					}
				}
			}
			// Close output channel when input channel is closed
			close(outs[i])
		}(i, input)
	}

	count := make(chan error, valuesSize)
	result := make([]O, valuesSize)

	// Start collector goroutines
	for _, out := range outs {
		go func(out chan *argc[O]) {
			for o := range out {
				if o != nil {
					if o.err != nil {
						select {
						case <-ctx.Done():
							return
						case count <- o.err:
						}
					} else {
						result[o.index] = *o.v
						select {
						case <-ctx.Done():
							return
						case count <- nil:
						}
					}
				}
			}
		}(out)
	}

	// Collect results
	size := 0
	for size < valuesSize {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-count:
			if err != nil {
				cancel() // Cancel all goroutines
				return nil, err
			}
			size++
		}
	}

	close(count)
	return result, nil
}
