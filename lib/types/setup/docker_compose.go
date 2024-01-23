package setup

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/google/uuid"
)

type DockerComposeContainerContext struct {
	Networks       []string          `json:"networks"`
	ID             string            `json:"id"`
	Image          string            `json:"image"`
	Name           string            `json:"name"`
	PublishedPorts map[string]string `json:"published_ports"`
}

type DockerComposeContext struct {
	Networks []string `json:"networks"`
}

type DockerComposeSetup struct {
	projectName string            `json:"-"`
	Files       []string          `json:"files"`
	Services    []string          `json:"services"`
	EnvFile     string            `json:"env_file"`
	ExtraEnv    map[string]string `json:"extra_env"`
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

	if _, err := os.Stat(d.EnvFile); err != nil {
		return fmt.Errorf("env file %s not found: %w", d.EnvFile, err)
	}

	return nil
}

func (d *DockerComposeSetup) generateBaseCommand() []string {
	args_slice := []string{"docker", "compose", "--project-name", d.GetProjectName()}
	for _, v := range d.Files {
		args_slice = append(args_slice, "--file", v)
	}
	if d.EnvFile != "" {
		args_slice = append(args_slice, "--env-file", d.EnvFile)
	}
	return args_slice
}

func (d *DockerComposeSetup) generateUpCommand() []string {
	args_slice := d.generateBaseCommand()
	args_slice = append(args_slice, "up", "-d")
	for _, v := range d.Services {
		args_slice = append(args_slice, v)
	}
	return args_slice
}

func (d *DockerComposeSetup) generateDownCommand() []string {
	args_slice := d.generateBaseCommand()
	args_slice = append(args_slice, "down")
	return args_slice
}

func (d *DockerComposeSetup) Setup(context map[string]string, logFile *io.Writer) (interface{}, error) {
	args_slice := d.generateUpCommand()
	cmd := exec.Command(args_slice[0], args_slice[1:]...)

	cmd.Env = os.Environ()
	for k, v := range d.ExtraEnv {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdout = *logFile
	cmd.Stderr = *logFile

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running docker compose: %w", err)
	}
	return nil, nil
}

func (d *DockerComposeSetup) Teardown() error {
	return nil
}
