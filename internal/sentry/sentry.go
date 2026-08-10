package sentry

import (
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

const (
	dsn = "https://e2d045fcf9a549199b7906c83493e2b6@sentry-par.scaleway.com/55"
)

const ErrorBanner = `---------------------------------------------------------------------------------------
An error occurred, we are sorry, please consider opening a ticket on github using: 'scw feedback bug'
Give us as many details as possible so we can reproduce the error and fix it.
---------------------------------------------------------------------------------------`

// RecoverPanicAndSendReport is to be called after recover.
// It expects tags returned by core.BuildInfo and the recovered error
func RecoverPanicAndSendReport(tags map[string]string, version string, e any) {
	sentryHub, err := sentryHub(tags, version)
	if err != nil {
		logger.Debugf("cannot get sentry hub: %s", err)
	}

	err, isError := e.(error)
	if isError {
		logAndSentry(sentryHub, err)
	} else {
		logAndSentry(sentryHub, fmt.Errorf("unknown error: %v", e))
	}
}

func logAndSentry(sentryHub *sentry.Hub, err error) {
	logger.Errorf("%s", err)
	if sentryHub != nil {
		event := sentryHub.Recover(err)
		if event == nil {
			logger.Debugf("failed to capture exception with sentry")

			return
		}
		logger.Debugf("sending sentry report: %s", *event)
		if !sentry.Flush(time.Second * 2) {
			logger.Debugf("failed to send report")
		}
	}
}

// newSentryClient create a sentry client with build info tags.
func newSentryClient(version string) (*sentry.Client, error) {
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		BeforeSend: func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
			filterStackFrames(event)

			return event
		},
		Release: version,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func configureSentryScope(tags map[string]string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTags(tags)
	})
}

// AddCommandContext is used to pass executed command
func AddCommandContext(line string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("command", map[string]any{
			"line": line,
		})
	})
}

func AddArgumentsContext(args [][2]string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		argMap := map[string]any{}

		for _, arg := range args {
			argMap[arg[0]] = len(arg[1])
		}

		scope.SetContext("arguments", argMap)
	})
}

func sentryHub(tags map[string]string, version string) (*sentry.Hub, error) {
	hub := sentry.CurrentHub()

	if hub.Client() == nil {
		client, err := newSentryClient(version)
		if err != nil {
			return nil, fmt.Errorf("cannot create sentry client: %w", err)
		}
		hub.BindClient(client)
		configureSentryScope(tags)
	}

	return hub, nil
}

// Filter the stack frames so that the top frame is the one causing panic.
// On top of the culprit one there should be
// - the deferred recover function
// - The two functions called in this package
func filterStackFrames(event *sentry.Event) {
	for _, e := range event.Exception {
		if e.Stacktrace == nil {
			continue
		}
		frames := e.Stacktrace.Frames[:0]
		for _, frame := range e.Stacktrace.Frames {
			if frame.Module == "main" && strings.HasPrefix(frame.Function, "cleanup") {
				continue
			}
			if frame.Module == "github.com/scaleway/scaleway-cli/v2/internal/sentry" &&
				strings.HasPrefix(frame.Function, "RecoverPanicAndSendReport") {
				continue
			}
			if frame.Module == "github.com/scaleway/scaleway-cli/v2/internal/sentry" &&
				strings.HasPrefix(frame.Function, "logAndSentry") {
				continue
			}
			frames = append(frames, frame)
		}
		e.Stacktrace.Frames = frames
	}
}
