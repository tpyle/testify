package runner

import (
	"github.com/tpyle/testamint/lib/types/result"
)

type PostmanCollection struct {
	CollectionFile       string            `json:"collectionFile"`
	EnvironmentFile      string            `json:"environmentFile"`
	EnvironmentVariables map[string]string `json:"environmentOverrides"`
}

type PostmanRunner struct {
	Collections []*PostmanCollection `json:"collections"`
}

func (p *PostmanRunner) Run() (result.ResultGroup, error) {
	return result.ResultGroup{}, nil
}
