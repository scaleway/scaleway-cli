package interactive

type InterruptError struct{}

func (e *InterruptError) MarshalHuman() string {
	return ""
}

func (e *InterruptError) Error() string {
	return "Readline Interrupt"
}

type EOFError struct{}

func (e *EOFError) Error() string {
	return "Readline EOF"
}
