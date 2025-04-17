package parallel

type PanicError struct {
	error
	ErrMsg string
}

func (e PanicError) Error() string {
	return e.ErrMsg
}
