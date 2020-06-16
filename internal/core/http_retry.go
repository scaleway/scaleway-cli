package core

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-hclog"
)

var (
	// Default retry configuration
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
	defaultRetryMax     = 4
)

// CheckForRetry specifies a policy for handling retries. It is called
// following each request with the response and error values returned by
// the http.Client. If it returns false, the Client stops retrying
// and returns the response to the caller. If it returns an error ,
// that error value is returned in lieu of the error from the request .
type CheckForRetry func(resp *http.Response, err error) (bool, error)

// DefaultRetryPolicy provides a default callback for Client.CheckForRetry,
// will retry on connection errors and server errors .
func DefaultRetryPolicy(resp *http.Response, err error) (bool, error) {
	if err != nil {
		return true, err
	}
	// Check the response code. Here, we retry on 500—range responses to allow
	//the server time to recover
	if resp.StatusCode == 0 || resp.StatusCode >= 500 {
		return true, nil
	}
	return false, nil
}

// Backoff specifies a policy for how long to wait between retries.
// It is called after a failing request to determine the amount of time
// that should pass before trying again.
type Backoff func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration

// DefaultBackoff provides a default callback for Client.Backoff which
// will perform exponential backoff based on the attempt number and limited
// by the provided minimum and maximum durations.
func DefaultBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	mult := math.Pow(2, float64(attemptNum)) * float64(min)
	sleep := time.Duration(mult)
	if float64(sleep) != mult || sleep > max {
		sleep = max
	}
	return sleep
}

// Client is used to make HTTP requests. It adds additional functionality
// like automatic retries to tolerate minor outages.
type Client struct {
	HTTPClient   *http.Client  // Internal HTTP client.
	RetryWaitMin time.Duration // Minimum time to wait
	RetryWaitMax time.Duration // Maximum time to wait
	RetryMax     int           // Maximum number of retries

	// CheckRetry specifies the policy for handling retries, and is called
	// after each request. The default policy is DefaultRetryPolicy.
	CheckForRetry CheckForRetry

	// Backoff specifies the policy for how long to wait between retries
	Backoff Backoff
}

func NewClient() *Client {
	return &Client{
		HTTPClient:    cleanhttp.DefaultClient(),
		RetryWaitMin:  defaultRetryWaitMin,
		RetryWaitMax:  defaultRetryWaitMax,
		RetryMax:      defaultRetryMax,
		CheckForRetry: DefaultRetryPolicy,
		Backoff:       DefaultBackoff,
	}
}

// Request wraps the metadata needed to create HTTP requests.
type Request struct {
	// body is a seekable reader over the request body payload. This is
	// used to rewind the request data in between retries.
	body io.ReadSeeker

	// Embed an HTTP request directly. This makes a *Request act exactly
	// like an *http.Request so that all meta methods are supported.
	*http.Request
}

// Try to read the response body so we can reuse this connection.
func (c *Client) drainBody(body io.ReadCloser) {
	defer body.Close()
	_, err := io.Copy(ioutil.Discard, io.LimitReader(body, respReadLimit))
	if err != nil {
		fmt.Printf("[ERR] error reading response body: %v", err)
	}
}

// Get is a convenience helper for doing simple GET requests.
func (c *Client) Get(url string) (*http.Response, error) {
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post is a convenience method for doing simple POST requests.
func (c *Client) Post(url, bodyType string, body interface{}) (*http.Response, error) {
	req, err := NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	return c.Do(req)
}

// NewRequest creates a new wrapped request.
func NewRequest(method, url string, body io.ReadSeeker) (*Request, error) {
	// Wrap the body in a noop ReadCloser if non-nil. This prevents the
	// reader from being closed by the HTTP client.
	var rcBody io.ReadCloser
	if body != nil {
		rcBody = ioutil.NopCloser(body)
	}

	// Make the request with the noop-closer for the body.
	httpReq, err := http.NewRequest(method, url, rcBody)
	if err != nil {
		return nil, err
	}

	return &Request{body, httpReq}, nil
}

// Do wraps calling an HTTP method with retries.
func (c *Client) Do(req *Request) (*http.Response, error) {
	if c.HTTPClient == nil {
		c.HTTPClient = cleanhttp.DefaultPooledClient()
	}

	logger := c.logger()

	if logger != nil {
		switch v := logger.(type) {
		case Logger:
			v.Printf("[DEBUG] %s %s", req.Method, req.URL)
		case hclog.Logger:
			v.Debug("performing request", "method", req.Method, "url", req.URL)
		}
	}

	var resp *http.Response
	var err error

	for i := 0; ; i++ {
		var code int // HTTP response code

		// Always rewind the request body when non-nil.
		if req.body != nil {
			body, err := req.body()
			if err != nil {
				c.HTTPClient.CloseIdleConnections()
				return resp, err
			}
			if c, ok := body.(io.ReadCloser); ok {
				req.Body = c
			} else {
				req.Body = ioutil.NopCloser(body)
			}
		}

		if c.RequestLogHook != nil && logger != nil {
			switch v := logger.(type) {
			case Logger:
				c.RequestLogHook(v, req.Request, i)
			case hclog.Logger:
				c.RequestLogHook(hookLogger{v}, req.Request, i)
			default:
				c.RequestLogHook(nil, req.Request, i)
			}
		}

		// Attempt the request
		resp, err = c.HTTPClient.Do(req.Request)
		if resp != nil {
			code = resp.StatusCode
		}

		// Check if we should continue with retries.
		checkOK, checkErr := c.CheckRetry(req.Context(), resp, err)

		if logger != nil {
			if err != nil {
				switch v := logger.(type) {
				case Logger:
					v.Printf("[ERR] %s %s request failed: %v", req.Method, req.URL, err)
				case hclog.Logger:
					v.Error("request failed", "error", err, "method", req.Method, "url", req.URL)
				}
			} else {
				// Call this here to maintain the behavior of logging all requests,
				// even if CheckRetry signals to stop.
				if c.ResponseLogHook != nil {
					// Call the response logger function if provided.
					switch v := logger.(type) {
					case Logger:
						c.ResponseLogHook(v, resp)
					case hclog.Logger:
						c.ResponseLogHook(hookLogger{v}, resp)
					default:
						c.ResponseLogHook(nil, resp)
					}
				}
			}
		}

		// Now decide if we should continue.
		if !checkOK {
			if checkErr != nil {
				err = checkErr
			}
			c.HTTPClient.CloseIdleConnections()
			return resp, err
		}

		// We do this before drainBody beause there's no need for the I/O if
		// we're breaking out
		remain := c.RetryMax - i
		if remain <= 0 {
			break
		}

		// We're going to retry, consume any response to reuse the connection.
		if err == nil && resp != nil {
			c.drainBody(resp.Body)
		}

		wait := c.Backoff(c.RetryWaitMin, c.RetryWaitMax, i, resp)
		desc := fmt.Sprintf("%s %s", req.Method, req.URL)
		if code > 0 {
			desc = fmt.Sprintf("%s (status: %d)", desc, code)
		}
		fmt.Printf("[DEBUG] %s: retrying in %s (%d left)", desc, wait, remain)
		time.Sleep(wait)
	}

	// Return an error if we fail out of the retry loop
	return nil, fmt.Errorf("%s %s giving up after %d attempts",
		req.Method, req.URL, c.RetryMax+1)
}
