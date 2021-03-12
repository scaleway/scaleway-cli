package sentry

import (
	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

const (
	dsn = "https://a3d5d1ef6ae94810952ab245ce61af17@sentry.scaleway.com/186"
)

const ErrorBanner = `---------------------------------------------------------------------------------------
An error occurred, we are sorry, please consider opening a ticket on github using: 'scw feedback bug'
Give us as many details as possible so we can reproduce the error and fix it.
---------------------------------------------------------------------------------------`

// RecoverPanicAndSendReport will recover error if any, log them, and send them to sentry.
// It must be called with the defer built-in.
func RecoverPanicAndSendReport(buildInfo *core.BuildInfo, e interface{}) {
	sentryClient, err := newSentryClient(buildInfo)
	if err != nil {
		logger.Debugf("cannot create sentry client: %s", err)
	}

	err, isError := e.(error)
	if isError {
		logAndSentry(sentryClient, err)
	} else {
		logAndSentry(sentryClient, fmt.Errorf("unknownw error: %v", e))
	}
}

func logAndSentry(sentryClient *raven.Client, err error) {
	logger.Errorf("%s", err)
	if sentryClient != nil {
		logger.Debugf("sending sentry report: %s", sentryClient.CaptureErrorAndWait(err, nil))
	}
}

// newSentryClient create a sentry client with build info tags.
func newSentryClient(buildInfo *core.BuildInfo) (*raven.Client, error) {
	client, err := raven.New(dsn)
	if err != nil {
		return nil, err
	}

	tagsContext := map[string]string{
		"version":    buildInfo.Version.String(),
		"go_arch":    buildInfo.GoArch,
		"go_os":      buildInfo.GoOS,
		"go_version": buildInfo.GoVersion,
	}

	client.SetTagsContext(tagsContext)

	return client, nil
}
