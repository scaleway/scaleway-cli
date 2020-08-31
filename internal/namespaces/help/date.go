package help

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func helpDate() *core.Command {
	longText := `Date parsing

You have two ways for managing date in the CLI: Absolute and Relative

- Absolute time

  Absolute time refers to a specific and absolute point in time.
  CLI uses RFC3339 to parse those time and pass a time.Time go structure to the underlying functions.

	Example: "2006-01-02T15:04:05Z07:00"

- Relative time

  Relative time refers to a time calculated from adding a given duration to the time when a command is launched.

	Example: now+1d4m => current time plus 1 day and 4 minutes

- Units of time

	Nanosecond: ns
	Microsecond: us, µs (U+00B5 = micro symbol), μs (U+03BC = Greek letter mu)
	Millisecond: ms
	Second: s, sec, second, seconds
	Minute: m, min, minute, minutes
	Hour: h, hr, hour, hours
	Day: d, day, days
	Week: w, wk, week, weeks
	Month: mo, mon, month, months
	Year: y, yr, year, years
`
	return &core.Command{
		Short:                "Get help about how date parsing works in the CLI",
		Long:                 longText,
		Namespace:            "help",
		Resource:             "date",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(struct{}{}),
		ArgSpecs:             core.ArgSpecs{},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			return longText, nil
		},
	}
}
