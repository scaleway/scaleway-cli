package commands

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/moul/anonuuid"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/config"
	"github.com/scaleway/scaleway-cli/pkg/scwversion"
)

func shouldBeAnUUID(actual interface{}, expected ...interface{}) string {
	input := actual.(string)
	input = strings.TrimSpace(input)
	if err := anonuuid.IsUUID(input); err != nil {
		return fmt.Sprintf("%q should be an UUID", actual)
	}
	return ""
}

func getScopedCtx(sessionCtx *CommandContext) (*CommandContext, *bytes.Buffer, *bytes.Buffer) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	var newCtx CommandContext
	newCtx = *sessionCtx
	newCtx.Stdout = &stdout
	newCtx.Stderr = &stderr

	return &newCtx, &stdout, &stderr
}

// RealAPIContext returns a CommandContext with a configured API
func RealAPIContext() *CommandContext {
	if os.Getenv("TEST_WITH_REAL_API") == "0" {
		return nil
	}
	config, err := config.GetConfig()
	if err != nil {
		logrus.Warnf("RealAPIContext: failed to call config.GetConfig(): %v", err)
		return nil
	}

	apiClient, err := api.NewScalewayAPI(config.Organization, config.Token, scwversion.UserAgent(), "par1")
	if err != nil {
		logrus.Warnf("RealAPIContext: failed to call api.NewScalewayAPI(): %v", err)
		return nil
	}

	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	ctx := CommandContext{
		Streams: Streams{
			Stdin:  os.Stdin,
			Stdout: &stdout,
			Stderr: &stderr,
		},
		Env: []string{
			"HOME" + os.Getenv("HOME"),
		},
		RawArgs: []string{},
		API:     apiClient,
	}
	return &ctx
}
