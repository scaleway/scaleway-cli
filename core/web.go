package core

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/skratchdot/open-golang/open"
)

func commandHasWeb(cmd *Command) bool {
	return cmd.WebURL != ""
}

func runWeb(cmd *Command, respI any) (any, error) {
	url := cmd.WebURL

	if respI != nil {
		tmpl, err := template.New("url").Parse(url)
		if err != nil {
			return nil, err
		}
		buf := bytes.NewBuffer(nil)
		err = tmpl.Execute(buf, respI)
		if err != nil {
			return nil, err
		}
		url = buf.String()
	}

	err := open.Start(url)
	if err != nil {
		return nil, &CliError{
			Err:     err,
			Message: "Failed to open web url",
			Details: "You can open it: " + url,
			Hint:    "You may not have a default browser configured",
			Code:    1,
		}
	}

	return fmt.Sprintf("Opening %s\n", url), nil
}
