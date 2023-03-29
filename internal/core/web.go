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

func runWeb(cmd *Command, respI interface{}) error {
	url := cmd.WebURL

	if respI != nil {
		tmpl, err := template.New("url").Parse(url)
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(nil)
		err = tmpl.Execute(buf, respI)
		if err != nil {
			return err
		}
		url = buf.String()
	}

	err := open.Start(url)
	if err != nil {
		return &CliError{
			Err:     nil,
			Message: "Failed to open web url",
			Details: fmt.Sprintf("You can open it: %s", url),
			Hint:    "You may not have a default browser configured",
			Code:    1,
		}
	}

	return nil
}
