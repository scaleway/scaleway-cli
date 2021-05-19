package interactive

import (
	"context"
	"fmt"
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

func popMockResponseFromContext(ctx context.Context) (string, bool) {
	contextValue := ctx.Value(contextKey)
	if contextValue == nil {
		return "", false
	}

	mockValues := contextValue.(*[]string)
	if mockValues == nil || len(*mockValues) == 0 {
		return "", false
	}
	str := (*mockValues)[0]
	*mockValues = (*mockValues)[1:]
	return str, true
}

type mockResponseReader struct {
	ctx           context.Context
	defaultReader io.ReadCloser
}

func (m *mockResponseReader) Read(p []byte) (n int, err error) {
	if m.ctx != nil {
		if mockResponse, exist := popMockResponseFromContext(m.ctx); exist {
			buff := []byte(fmt.Sprintf("%s\n", mockResponse))
			copy(p, buff)
			return len(buff), nil
		}
	}
	return m.defaultReader.Read(p)
}

func (m *mockResponseReader) Close() error {
	return m.defaultReader.Close()
}
