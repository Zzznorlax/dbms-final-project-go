package util

type ErrNotFound struct {
	Err error
}

func (e ErrNotFound) Error() string {
	return "not found"
}

func (e ErrNotFound) Unwrap() error { return e.Err }

type ErrUnauthorized struct {
	Reason string
}

func (e ErrUnauthorized) Error() string {
	return e.Reason
}
