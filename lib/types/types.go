package types

import (
	"encoding/json"
	"fmt"

	"github.com/tpyle/testamint/lib/types/check"
	"github.com/tpyle/testamint/lib/types/setup"
)

type TestConfig struct {
	Tests []*Test `json:"tests"`
}

type Test struct {
	Name        string        `json:"name"`
	Setup       setup.Setup   `json:"setup"`
	ReadyChecks []check.Check `json:"readyChecks"`
	// Runner      Runner       `json:"runner"`
}

func (t *Test) UnmarshalJSON(data []byte) error {
	type Aux struct {
		Name        string            `json:"name"`
		Setup       json.RawMessage   `json:"setup"`
		ReadyChecks []json.RawMessage `json:"readyChecks"`
	}

	var aux Aux
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Name = aux.Name
	if aux.Setup == nil {
		return fmt.Errorf("error: setup is required")
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

	return nil
}
