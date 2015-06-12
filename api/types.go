package api

// Config is a Scaleway CLI configuration file
type Config struct {
	// APIEndpoint is the endpoint to the Scaleway API
	APIEndPoint string `json:"api_endpoint"`

	// Organization is the identifier of the Scaleway orgnization
	Organization string `json:"organization"`

	// Token is the authentication token for the Scaleway organization
	Token string `json:"token"`
}
