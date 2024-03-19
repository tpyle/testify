package types

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tpyle/testamint/lib/types/check"
	"github.com/tpyle/testamint/lib/types/result"
	"github.com/tpyle/testamint/lib/types/runner"
	"github.com/tpyle/testamint/lib/types/setup"
)

type TestConfig struct {
	Parallelization uint16  `json:"parallelization"`
	Tests           []*Test `json:"tests"`
}

type Test struct {
	Name        string        `json:"name"`
	Setup       setup.Setup   `json:"setup"`
	ReadyChecks []check.Check `json:"readyChecks"`
	Runner      runner.Runner `json:"runner"`
	// Runner      Runner       `json:"runner"`
}

type TestContext struct {
	Setup       interface{}
	ReadyChecks []interface{}
}

func (t *Test) Validate(context interface{}, logFile io.Writer) error {
	err := t.Setup.Validate(context, logFile)
	if err != nil {
		return fmt.Errorf("error in setup: %w", err)
	}
	for _, rc := range t.ReadyChecks {
		err = rc.Validate(context, logFile)
		if err != nil {
			return fmt.Errorf("error in readyCheck: %w", err)
		}
	}
	err = t.Runner.Validate(context, logFile)
	if err != nil {
		return fmt.Errorf("error in runner: %w", err)
	}
	return nil
}

func (t *Test) Run(logFile io.Writer) (*result.ResultGroup, error) {
	var context TestContext
	err := t.Validate(context, logFile)
	if err != nil {
		return nil, err
	}

	setupContext, err := t.Setup.Setup(context, logFile)
	if err != nil {
		return nil, err
	}
	defer t.Setup.Teardown(logFile)
	context.Setup = setupContext

	// for k := range context.Setup.(setup.DockerComposeContext).Containers {
	// 	fmt.Printf("Container: %s\n", k)
	// 	fmt.Printf("Container: %+v\n", context.Setup.(setup.DockerComposeContext).Containers[k])
	// }

	// fmt.Printf("Total Context %+v\n", context)

	for _, rc := range t.ReadyChecks {
		err = rc.WaitForReady(context, logFile)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (t *Test) UnmarshalJSON(data []byte) error {
	type Aux struct {
		Name        string            `json:"name"`
		Setup       json.RawMessage   `json:"setup"`
		ReadyChecks []json.RawMessage `json:"readyChecks"`
		Runner      json.RawMessage   `json:"runner"`
	}

	var aux Aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Name = aux.Name
	if aux.Setup == nil {
		return fmt.Errorf("setup is required")
	}

	var err error
	t.Setup, err = setup.UnmarshalSetup(aux.Setup)
	if err != nil {
		return err
	}

	var checks []check.Check
	for _, rc := range aux.ReadyChecks {
		c, err := check.UnmarshalCheck(rc)
		if err != nil {
			return err
		}
		checks = append(checks, c)
	}
	t.ReadyChecks = checks

	if aux.Runner == nil {
		return fmt.Errorf("runner is required")
	}

	t.Runner, err = runner.UnmarshalRunner(aux.Runner)
	if err != nil {
		return err
	}

	return nil
}
