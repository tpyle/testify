package setup

import (
	"encoding/json"
	"fmt"
	"io"
)

type Setup interface {
	Validate(context interface{}, logFile io.Writer) error
	Setup(context interface{}, logFile io.Writer) (interface{}, error)
	Teardown(logFile io.Writer) error
}

type SetupType string

const (
	DockerCompose SetupType = "docker-compose"
	None          SetupType = "none"
	Exec          SetupType = "exec"
)

func UnmarshalSetup(data []byte) (Setup, error) {
	type auxS struct {
		SetupType SetupType `json:"type"`
	}

	var aux auxS

	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, err
	}

	var setup Setup
	switch aux.SetupType {
	case DockerCompose:
		var dc DockerComposeSetup
		setup = &dc
	case None:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown setup type: %s", aux.SetupType)
	}

	err := json.Unmarshal(data, setup)
	if err != nil {
		return nil, err
	}

	return setup, nil
}
