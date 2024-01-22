package types

type Test struct {
	Name        string        `json:"name"`
	Setup       *Setup        `json:"setup"`
	ReadyChecks []*ReadyCheck `json:"readyChecks"`
	Runner      *Runner       `json:"runner"`
}
