package human

type Formatter interface {
	FormatHuman() (string, error)
}

func Format(data any) (string, error) {
	return "", nil
}
