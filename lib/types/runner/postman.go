package runner

import (
	"io"

	"github.com/tpyle/testamint/lib/types/result"
)

type PostmanCollection struct {
	CollectionFile       string            `json:"collectionFile"`
	EnvironmentFile      string            `json:"environmentFile"`
	EnvironmentVariables map[string]string `json:"environmentOverrides"`
}

type PostmanRunner struct {
	Collections          []*PostmanCollection `json:"collections"`
	EnvironmentFile      string               `json:"environmentFile"`
	EnvironmentVariables map[string]string    `json:"environmentOverrides"`

	NewmanPath  string `json:"newmanPath"`  // If empty, use the newman in the PATH
	UseDocker   bool   `json:"useDocker"`   // If true, use the dockerImage to run the tests, instead of a native newman
	DockerImage string `json:"dockerImage"` // If useDocker is true, this is the image to use, otherwise it is ignored
}

func (p *PostmanRunner) Run(context interface{}, logFile io.Writer) (*result.ResultGroup, error) {
	return &result.ResultGroup{}, nil
}

func (p *PostmanRunner) Validate(context interface{}, logFile io.Writer) error {
	return nil
}
