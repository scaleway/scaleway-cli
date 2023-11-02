package account

import (
	"context"
	"net/http"
)

//
// Code in this file is useful if you are trying to test code that use interactive package
//

type contextKeyType struct{}

var contextKey = contextKeyType{}

func InjectHTTPClient(ctx context.Context, httpClient *http.Client) context.Context {
	return context.WithValue(ctx, contextKey, httpClient)
}
