package commands

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"github.com/scaleway/scaleway-cli/pkg/config"
)

func shouldBeAnUUID(actual interface{}, expected ...interface{}) string {
	input := actual.(string)
	input = strings.TrimSpace(input)
	uuid := uuid.Parse(input)
	if uuid == nil {
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

func RealAPIContext() *CommandContext {
	config, err := config.GetConfig()
	if err != nil {
		logrus.Warnf("RealAPIContext: failed to call config.GetConfig(): %v", err)
		return nil
	}

	apiClient, err := api.NewScalewayAPI(config.ComputeAPI, config.AccountAPI, config.Organization, config.Token)
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
