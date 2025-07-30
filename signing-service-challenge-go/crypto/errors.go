package crypto

import "fmt"

// SignOperationError definition
type SignOperationError struct {
	Algorithm string
	Err       error
}

func (e *SignOperationError) Error() string {
	return fmt.Sprintf("error while signing the data with algorithm %s: %s", e.Algorithm, e.Err)
}

func (e *SignOperationError) Unwrap() error { //Unwrap included for logging purposes
	return e.Err
}

func NewSignOperationError(algorithm string, err error) *SignOperationError {
	return &SignOperationError{Algorithm: algorithm, Err: err}
}

// MarshalError definition
type MarshalError struct {
	Algorithm string
	Err       error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("error while marshaling the data with the algorithm %s: %s", e.Algorithm, e.Err)
}

func (e *MarshalError) Unwrap() error { //Unwrap included for logging purposes
	return e.Err
}

func NewMarshalError(algorithm string, err error) *MarshalError {
	return &MarshalError{Algorithm: algorithm, Err: err}
}
