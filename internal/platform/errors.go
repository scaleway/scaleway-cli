package platform

// ClientError define an error when failing to create client
type ClientError struct {
	Err     error
	Details string
}

func (s *ClientError) Error() string {
	return s.Err.Error()
}
