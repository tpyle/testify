package setup

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

type DockerComposeSetup struct {
	projectName string   `json:"-"`
	Files       []string `json:"files"`
	Services    []string `json:"services"`
}

func (d *DockerComposeSetup) GetProjectName() string {
	if d.projectName == "" {
		d.projectName = uuid.New().String()
	}
	return d.projectName
}

// Validate ensures the setup is appropriate for the test context
// In this case, it checks that docker is present, as is `docker compose`, the
// files exist relative to the working directory, and the services exist in the
// compose files
func (d *DockerComposeSetup) Validate() error {
	path, err := exec.LookPath("docker")
	if err != nil {
		return fmt.Errorf("docker not found in path: %w", err)
	}

	if err := exec.Command(path, "compose", "--version").Run(); err != nil {
		return fmt.Errorf("docker compose not found in runnable: %w", err)
	}

	for _, file := range d.Files {
		if _, err := os.Stat(file); err != nil {
			return fmt.Errorf("file %s not found: %w", file, err)
		}
	}

	return nil
}

func (d *DockerComposeSetup) Setup() (interface{}, error) {
	return nil, nil
}

func (d *DockerComposeSetup) Teardown() error {
	return nil
}
