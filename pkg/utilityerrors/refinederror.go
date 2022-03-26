package utilityerrors

func NewRefined(message string) *RefinedError {
	return &RefinedError{
		message: message,
	}
}

type RefinedError struct {
	message     string
	clarifications []string
}

func (e RefinedError) Clarify(clarification string) *RefinedError {
	e.clarifications = append(e.clarifications, clarification)
	return &e
}

func (e RefinedError) Error() string {
	if len(e.clarifications) == 0 {
		return e.message
	}

	finalClarification := ""

	for i, clarification := range e.clarifications {
		finalClarification += indent(clarification, "  ", i + 1)

		if len(e.clarifications) != i + 1 {
			finalClarification += ":\n"
		}
	}

	return e.message + ":\n" + finalClarification
}

func (e RefinedError) Is(targetError error) bool {
	v, ok := targetError.(*RefinedError)
	if !ok {
		return false
	}

	return e.message == v.message
}
