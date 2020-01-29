package matomo

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	// full reference API https://developer.matomo.org/api-reference/tracking-api
	endpoint     = "stats.cloud.online.net/matomo.php"
	timeout      = 5 * time.Second
	siteID       = "32"
	enableRecord = "1"
	apiVersion   = "1"
)

// ForceTelemetry is used to send telemetry even from a non-released CLI.
// This WILL NOT bypass user policy set in send_telemetry attribute.
var ForceTelemetry = os.Getenv("SCW_FORCE_TELEMETRY") == "true"

type SendCommandTelemetryRequest struct {
	Command       string
	Version       string
	ExecutionTime time.Duration
}

// SendCommandTelemetry will send the telemetry report or return an error on failure.
func SendCommandTelemetry(request *SendCommandTelemetryRequest) error {
	// compute or retrieve telemetry parameters
	terminalResolution := fmt.Sprintf("%dx%d", terminal.GetWidth(), terminal.GetHeight())
	commandDurationInMs := fmt.Sprintf("%d", request.ExecutionTime/time.Millisecond)
	randNumber := generateRandNumber()
	action, actionURL := commandToAction(request.Command)
	userAgent := fakeUserAgent(request.Version)
	organizationID := ""
	config, err := scw.LoadConfig()
	if err == nil {
		profile, err := config.GetActiveProfile()
		if err == nil && profile.DefaultOrganizationID != nil {
			organizationID = *profile.DefaultOrganizationID
		}
	}

	// build the query parameters in the URL
	query := url.Values{}

	// required
	query.Add("idsite", siteID)
	query.Add("rec", enableRecord)

	// recommended
	query.Add("action_name", action)
	query.Add("url", actionURL)
	query.Add("rand", randNumber)
	query.Add("apiv", apiVersion)

	// optional
	query.Add("res", terminalResolution)
	query.Add("ua", userAgent)
	query.Add("uid", organizationID)
	query.Add("gt_ms", commandDurationInMs)

	matomoURL := url.URL{
		Path:     endpoint,
		Scheme:   "https",
		RawQuery: query.Encode(),
	}

	// send the report
	resp, err := (&http.Client{
		Timeout: timeout,
	}).Get(matomoURL.String())
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("non-success status code %d: %s", resp.StatusCode, matomoURL.String())
	}

	return nil
}

// fakeUserAgent creates a fake user agent that follows Matomo requirements.
// We don't use the actual SDK user agent because it is not exposed publicly.
func fakeUserAgent(version string) string {
	return fmt.Sprintf("scaleway-cli/%s (%s; %s; %s)", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// commandToAction will convert the command path to an action (space-separated) and an action URL (required by Matomo).
func commandToAction(command string) (action string, url string) {
	command = "scw " + strings.Replace(command, ".", " ", -1)
	return command, "https://" + strings.Replace(command, " ", "/", -1)
}

// generateRandNumber will generate true random number in order to matomo requests to be cached by a proxy.
func generateRandNumber() string {
	bigRand, err := rand.Int(rand.Reader, big.NewInt(int64(1<<uint64(32))-1))
	if err != nil {
		return ""
	}
	return bigRand.String()
}

// IsTelemetryEnabled returns true when the Opt-In send_telemetry attribute in the config is set.
func IsTelemetryEnabled() bool {
	config, err := scw.LoadConfig()
	if err != nil {
		return false
	}
	return config.SendUsage
}
