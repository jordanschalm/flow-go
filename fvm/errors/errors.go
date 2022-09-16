package errors

import (
	stdErrors "errors"

	"github.com/onflow/cadence/runtime"
	"github.com/onflow/cadence/runtime/errors"
)

// Error covers all non-fatal errors happening
// while validating and executing a transaction or a script.
type Error interface {
	// Code returns the code for this error
	Code() ErrorCode
	// and anything else that is needed to be an error
	error
}

// Failure captures fatal unexpected virtual machine errors,
// we capture this type of error instead of panicking
// to collect all necessary data before crashing
// if any of these errors occurs we should halt the
// execution.
type Failure interface {
	// FailureCode returns the failure code for this error
	FailureCode() FailureCode
	// and anything else that is needed to be an error
	error
}

type errorWrapper struct {
	err error
}

func (e errorWrapper) Unwrap() error {
	return e.err
}

// Is is a utility function to call std error lib `Is` function for instance equality checks.
func Is(err error, target error) bool {
	return stdErrors.Is(err, target)
}

// As is a utility function to call std error lib `As` function.
// As finds the first error in err's chain that matches target,
// and if so, sets target to that error value and returns true. Otherwise, it returns false.
// The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap.
func As(err error, target interface{}) bool {
	return stdErrors.As(err, target)
}

// SplitErrorTypes splits the error into fatal (failures) and non-fatal errors
func SplitErrorTypes(inp error) (err Error, failure Failure) {
	// failures should get the priority
	// this method will check all the levels for these failures
	if As(inp, &failure) {
		return nil, failure
	}
	// then we should try to match known non-fatal errors
	if As(inp, &err) {
		return err, nil
	}
	// anything else that is left is an unknown failure
	// (except the ones green listed for now to be considered as txErrors)
	if inp != nil {
		return nil, NewUnknownFailure(inp)
	}
	return nil, nil
}

// HandleRuntimeError handles runtime errors and separates
// errors generated by runtime from fvm errors (e.g. environment errors)
func HandleRuntimeError(err error) error {
	if err == nil {
		return nil
	}

	// if is not a runtime error return as vm error
	// this should never happen unless a bug in the code
	runErr, ok := err.(runtime.Error)
	if !ok {
		return NewUnknownFailure(err)
	}

	// External errors are reported by the runtime but originate from the VM.
	// External errors may be fatal or non-fatal, so additional handling by SplitErrorTypes
	if externalErr, ok := errors.GetExternalError(err); ok {
		if recoveredErr, ok := externalErr.Recovered.(error); ok {
			// If the recovered value is an error, pass it to the original
			// error handler to distinguish between fatal and non-fatal errors.
			return recoveredErr
		}
		// if not recovered return
		return NewUnknownFailure(externalErr)
	}

	// All other errors are non-fatal Cadence errors.
	return NewCadenceRuntimeError(runErr)
}
