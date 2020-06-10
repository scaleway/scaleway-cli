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

func InjectHttpClient(ctx context.Context, httpClient *http.Client) context.Context {
	return context.WithValue(ctx, contextKey, httpClient)
}

func extractHttpClient(ctx context.Context) *http.Client {
	httpClient, isHttpClient := ctx.Value(contextKey).(*http.Client)
	if httpClient != nil && isHttpClient {
		return httpClient
	}
	return http.DefaultClient
}
