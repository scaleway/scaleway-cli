package interactive

import (
	"context"
	"io"
)

//
// Code in this file is useful if you are trying to test code that use interactive package
//

type contextKeyType struct{}

var contextKey = contextKeyType{}

func InjectMockResponseToContext(ctx context.Context, mockValues []string) context.Context {
	return context.WithValue(ctx, contextKey, &mockValues)
}

type mockResponseReader struct {
	mockResponses []string
	defaultReader io.ReadCloser
}

func (m *mockResponseReader) Read(p []byte) (n int, err error) {
	if len(m.mockResponses) > 0 {
		mockResponse := m.mockResponses[0]
		m.mockResponses = m.mockResponses[1:]
		buff := []byte(mockResponse + "\n")
		copy(p, buff)

		return len(buff), nil
	}

	return m.defaultReader.Read(p)
}

func (m *mockResponseReader) Close() error {
	return m.defaultReader.Close()
}
