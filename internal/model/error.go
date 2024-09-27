package model

import (
	"fmt"
)

// Declare here but other packages can declare enum values like this
//
// const (
//
//	ExpressiveReason1 model.ErrorReason = iota + 1
//	ExpressionReason2
//	...
//
// )
type ErrorReason int

// Custom error for application level errors.
type Error struct {
	Reason ErrorReason
	Err    error
}

// Call this function to wrap an error with a specific reason.
func Wrap(err error, reason ErrorReason) error {
	if err == nil {
		return err
	}

	return &Error{
		Reason: reason,
		Err:    err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("(Application Error) reason: %d, err:%v", e.Reason, e.Err)
}
