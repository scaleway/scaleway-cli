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

func extractHTTPClient(ctx context.Context) *http.Client {
	httpClient, isHTTPClient := ctx.Value(contextKey).(*http.Client)
	if httpClient != nil && isHTTPClient {
		return httpClient
	}
	return http.DefaultClient
}
