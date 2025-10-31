package interactive

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

func Print(a ...any) (int, error) {
	if IsInteractive {
		return fmt.Fprint(outputWriter, a...)
	}

	return 0, nil
}

func Println(a ...any) (int, error) {
	if IsInteractive {
		return fmt.Fprintln(outputWriter, a...)
	}

	return 0, nil
}

func PrintlnWithoutIndent(a string) (int, error) {
	if IsInteractive {
		return fmt.Fprintln(outputWriter, RemoveIndent(a))
	}

	return 0, nil
}

func Printf(format string, a ...any) (int, error) {
	if IsInteractive {
		return fmt.Fprintf(outputWriter, format, a...)
	}

	return 0, nil
}

func Line(char string) string {
	return strings.Repeat(char, terminal.GetWidth())
}

func Center(str string) string {
	longestLine := 0
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		longestLine = int(math.Max(float64(longestLine), float64(len(line))))
	}

	return Indent(str, (terminal.GetWidth()-longestLine)/2)
}

func Indent(str string, indent int) string {
	padding := strings.Repeat(" ", indent)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = padding + line
		}
	}

	return strings.Join(lines, "\n")
}

func RemoveIndent(str string) string {
	return strings.Trim(regexp.MustCompile("\n[ \t]*").ReplaceAllString(str, "\n"), "\n")
}
