package tasks

import (
	"time"

	"github.com/briandowns/spinner"
)

// Loader will print loading activity to the terminal
type Loader struct {
	spinner *spinner.Spinner
}

func setupLoader() *Loader {
	return &Loader{
		spinner.New(spinner.CharSets[11], 100*time.Millisecond),
	}
}

func (l *Loader) Start() {
	l.spinner.Start()
}

func (l *Loader) Stop() {
	l.spinner.Stop()
}
