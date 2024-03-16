package runner

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tpyle/testamint/lib/types/result"
)

type Runner interface {
	Validate(context interface{}, logFile io.Writer) error
	Run(context interface{}, logFile io.Writer) (*result.ResultGroup, error)
}

type RunnerType string

const (
	Postman RunnerType = "postman"
	Junit   RunnerType = "junit"
	Jmeter  RunnerType = "jmeter"
	K6      RunnerType = "k6"
	ReqX    RunnerType = "reqx"
	Go      RunnerType = "go"
	Gradle  RunnerType = "gradle"
	Cargo   RunnerType = "cargo"
)

func UnmarshalRunner(data []byte) (Runner, error) {
	type auxR struct {
		RunnerType RunnerType `json:"type"`
	}

	var aux auxR

	if err := json.Unmarshal(data, &aux); err != nil {
		return nil, err
	}

	var runner Runner
	switch aux.RunnerType {
	case Postman:
		var p PostmanRunner
		runner = &p
	default:
		return nil, fmt.Errorf("unknown runner type: %s", aux.RunnerType)
	}

	err := json.Unmarshal(data, runner)
	if err != nil {
		return nil, err
	}

	return runner, nil
}
