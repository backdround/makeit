package utilityerrors

var newWrappedErrorId int = 0

func NewWrapped(message string) *WrappedError {
	newWrappedErrorId++
	return &WrappedError{message: message, id: newWrappedErrorId}
}

type WrappedError struct {
	message       string
	internalError error
	id            int
}

func (e WrappedError) Error() string {
	if e.internalError != nil {
		internalMessage := indent(e.internalError.Error(), "  ",1)
		return e.message + ":\n" + internalMessage
	}

	return e.message
}

func (e WrappedError) Wrap(internalError error) *WrappedError {
	e.internalError = internalError
	return &e
}

func (e WrappedError) Unwrap() error {
	return e.internalError
}

func (e WrappedError) Is(targetError error) bool {
	v, ok := targetError.(*WrappedError)
	if !ok {
		return false
	}

	return e.id == v.id
}
