package core

import "github.com/guackamolly/insta-archiver/internal/domain"

// This function simplifies the error handling of function results that contain errors.
func Process[I any, O any](
	input I,
	action func(I) (O, error),
	onError func(error),
) O {
	output, err := action(input)

	if err != nil {
		onError(err)
		return output
	}

	return output
}

// Applies [Process] function to an usecase.
func Invoke[I any, O any](
	input I,
	usecase domain.Usecase[I, O],
	onError func(error),
) O {
	return Process(input, usecase.Invoke, onError)
}
