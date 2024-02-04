package types

import (
	"encoding/json"
	"fmt"

	"github.com/tpyle/testamint/lib/types/setup"
)

type TestConfig struct {
	Tests []*Test `json:"tests"`
}

type Test struct {
	Name  string      `json:"name"`
	Setup setup.Setup `json:"setup"`
	// ReadyChecks []ReadyCheck `json:"readyChecks"`
	// Runner      Runner       `json:"runner"`
}

func (t *Test) UnmarshalJSON(data []byte) error {
	type Aux struct {
		Name  string          `json:"name"`
		Setup json.RawMessage `json:"setup"`
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

	return nil
}
