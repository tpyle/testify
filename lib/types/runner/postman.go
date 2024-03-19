package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/tpyle/testamint/lib/types/result"
)

type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostmanCollection struct {
	CollectionFile       string                 `json:"collectionFile"`
	EnvironmentFile      string                 `json:"environmentFile"`
	EnvironmentVariables []*EnvironmentVariable `json:"environmentOverrides"`
}

type PostmanRunner struct {
	Collections          []*PostmanCollection   `json:"collections"`
	EnvironmentFile      string                 `json:"environmentFile"`
	EnvironmentVariables []*EnvironmentVariable `json:"environmentOverrides"`

	NewmanPath  string `json:"newmanPath"`  // If empty, use the newman in the PATH
	UseDocker   bool   `json:"useDocker"`   // If true, use the dockerImage to run the tests, instead of a native newman
	DockerImage string `json:"dockerImage"` // If useDocker is true, this is the image to use, otherwise it is ignored
}

func (p *PostmanRunner) Run(context interface{}, logFile io.Writer) (*result.ResultGroup, error) {
	return &result.ResultGroup{}, nil
}

func (p *PostmanRunner) Validate(context interface{}, logFile io.Writer) error {
	if p.UseDocker && p.DockerImage == "" {
		return fmt.Errorf("docker image must be specified if useDocker is true")
	}
	for _, collection := range p.Collections {
		if collection.CollectionFile == "" {
			return fmt.Errorf("collectionFile must be specified for each collection")
		}
		info, err := os.Stat(collection.CollectionFile)
		if err != nil {
			return fmt.Errorf("collectionFile %s does not exist", collection.CollectionFile)
		}
		if info.IsDir() {
			return fmt.Errorf("collectionFile %s is a directory", collection.CollectionFile)
		}
	}
	return nil
}
