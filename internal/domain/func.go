package domain

import (
	"github.com/guackamolly/insta-archiver/internal/model"
)

// Calls a function that leads to an output and error. If the output error is present,
// it's wrapped with a custom error reason.
func WrapResult[I any, O any](
	input I,
	action func(I) (O, error),
	ifError model.ErrorReason,
) (O, error) {
	return WrapResult0(func() (O, error) { return action(input) }, ifError)
}

// Calls a function that leads to an output and error. If the output error is present,
// it's wrapped with a custom error reason.
func WrapResult0[O any](
	action func() (O, error),
	ifError model.ErrorReason,
) (O, error) {
	output, err := action()

	if err != nil {
		return output, model.Wrap(err, ifError)
	}

	return output, nil
}

// Invokes an use case. The [err] parameter
// is used to control whether the use case can be invoked or not.
func Invoke[I any, O any](
	input I,
	usecase Usecase[I, O],
	err error,
) (O, error) {
	var o O

	if err == nil {
		o, err = usecase.Invoke(input)
	}

	return o, err
}
