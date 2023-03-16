package tasks

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/containerd/console"
	buildkit "github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/opencontainers/go-digest"
	"github.com/pkg/errors"
)

type Logger struct {
	status   chan *buildkit.SolveStatus
	done     <-chan struct{}
	err      error
	warnings []buildkit.VertexWarning
}

type LoggerMode string

const (
	PrinterModeQuiet LoggerMode = "quiet"
	PrinterModeAuto  LoggerMode = "auto"
	PrinterModeTty   LoggerMode = "tty"
)

func NewTasksLogger(ctx context.Context, mode LoggerMode) (*Logger, error) {
	// new temp file for logging
	out := os.Stdout
	var writer io.Writer = out

	var cons console.Console
	switch mode {
	case PrinterModeQuiet:
		writer = io.Discard
	case PrinterModeAuto, PrinterModeTty:
		var err error
		cons, err = console.ConsoleFromFile(out)

		if err != nil {
			if mode == PrinterModeTty {
				return nil, errors.Wrap(err, "failed to get console")
			}
		}
	}

	doneChan := make(chan struct{})
	logger := &Logger{
		status: make(chan *buildkit.SolveStatus),
		done:   doneChan,
	}

	go func() {
		// resumeLogs := logutil.Pause(logrus.StandardLogger())
		logger.warnings, logger.err = progressui.DisplaySolveStatus(ctx, "Running workflow", cons, writer, logger.status)
		// resumeLogs()
		close(doneChan)
	}()
	return logger, nil
}

func (l *Logger) Write(s *buildkit.SolveStatus) {
	l.status <- s
}

func (l *Logger) CloseAndWait() error {
	close(l.status)
	<-l.done
	return l.err
}

type LoggerEntry struct {
	Logs io.Writer

	Start    func()
	Complete func(err error)
}

func (l *Logger) AddEntry(name string) *LoggerEntry {
	id := digest.FromString(name)

	r, w := io.Pipe()
	go func() {
		data := make([]byte, 1024)

		for {
			n, err := r.Read(data)
			if err != nil {
				return
			}

			if n == 0 {
				continue
			}

			l.Write(&buildkit.SolveStatus{
				Logs: []*buildkit.VertexLog{
					{
						Vertex: id,
						Data:   data,
					},
				},
			})
		}
	}()

	var started time.Time

	return &LoggerEntry{
		Logs: w,
		Start: func() {
			started = time.Now()
			l.Write(&buildkit.SolveStatus{
				Vertexes: []*buildkit.Vertex{
					{
						Digest:  id,
						Name:    name,
						Started: &started,
					},
				},
			})
		},
		Complete: func(err error) {
			r.Close()
			w.Close()

			completed := time.Now()

			var errStr string
			if err != nil {
				errStr = err.Error()
			}

			l.Write(&buildkit.SolveStatus{
				Vertexes: []*buildkit.Vertex{
					{
						Digest:    id,
						Name:      name,
						Started:   &started,
						Completed: &completed,
						Error:     errStr,
					},
				},
			})
		},
	}
}
