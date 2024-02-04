package setup

import "encoding/json"

type SetupType string

const (
	// DockerCompose represents a docker-compose setup.
	DockerCompose SetupType = "docker-compose"
	None          SetupType = "none"
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
	}

	err := json.Unmarshal(data, setup)
	if err != nil {
		return nil, err
	}

	return setup, nil
}
