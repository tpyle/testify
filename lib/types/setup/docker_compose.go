package setup

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

// URL represents the URL (address) of the published port.
// TargetPort represents the target port of the published port (in the container).
// PublishedPort represents the published port number (on the host).
// Protocol represents the protocol used for the published port (e.g. tcp/udp).
type DockerComposePublishedPort struct {
	URL           string `json:"URL"`
	TargetPort    int    `json:"TargetPort"`
	PublishedPort int    `json:"publishedPort"`
	Protocol      string `json:"Protocol"`
}

// DockerComposeContainerContext represents the context of a container in a docker compose setup.
// Networks represents the networks the container is connected to.
// ID represents the ID of the container.
// Image represents the image of the container.
// Name represents the name of the container.
// Service represents the service the container belongs to.
// PublishedPorts represents the published ports of the container.
type DockerComposeContainerContext struct {
	Networks       []string                              `json:"-"`
	ID             string                                `json:"ID"`
	Image          string                                `json:"Image"`
	Name           string                                `json:"Name"`
	Service        string                                `json:"Service"`
	PublishedPorts map[string]DockerComposePublishedPort `json:"-"`
}

func (d *DockerComposeContainerContext) UnmarshalJSON(data []byte) error {
	type auxS struct {
		Networks   string                       `json:"Networks"`
		ID         string                       `json:"ID"`
		Image      string                       `json:"Image"`
		Name       string                       `json:"Name"`
		Service    string                       `json:"Service"`
		Publishers []DockerComposePublishedPort `json:"Publishers"`
	}
	var aux auxS
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse and trim the networks
	d.Networks = strings.Split(aux.Networks, ",")
	for i, v := range d.Networks {
		d.Networks[i] = strings.TrimSpace(v)
	}

	// Parse and trim the published ports
	d.PublishedPorts = make(map[string]DockerComposePublishedPort)
	for _, v := range aux.Publishers {
		stringTargetPort := fmt.Sprintf("%d", v.TargetPort)
		d.PublishedPorts[stringTargetPort] = v
	}

	d.ID = aux.ID
	d.Image = aux.Image
	d.Name = aux.Name
	d.Service = aux.Service

	return nil
}

type DockerComposeContext struct {
	Containers map[string]*DockerComposeContainerContext `json:"containers"`
}

func NewDockerComposeContext() DockerComposeContext {
	return DockerComposeContext{
		Containers: make(map[string]*DockerComposeContainerContext),
	}
}

func (d *DockerComposeContext) AddContainer(name string, container *DockerComposeContainerContext) {
	d.Containers[name] = container
}

type EnvEntry struct {
	Key string `json:"key"`
	Val string `json:"value"`
}

type DockerComposeSetup struct {
	projectName string     `json:"-"`
	Files       []string   `json:"files"`
	Services    []string   `json:"services"`
	EnvFile     string     `json:"env_file"`
	ExtraEnv    []EnvEntry `json:"extra_env"`
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
	args_slice = append(args_slice, d.Services...)
	return args_slice
}

func (d *DockerComposeSetup) generateDownCommand() []string {
	args_slice := d.generateBaseCommand()
	args_slice = append(args_slice, "down")
	return args_slice
}

func (d *DockerComposeSetup) generatePsCommand() []string {
	args_slice := d.generateBaseCommand()
	args_slice = append(args_slice, "ps", "--format", "json")
	return args_slice
}

func (d *DockerComposeSetup) GetContext() (interface{}, error) {
	args_slice := d.generatePsCommand()

	cmd := exec.Command(args_slice[0], args_slice[1:]...)

	cmd.Env = os.Environ()

	out, err := cmd.Output()

	if err != nil {
		return nil, fmt.Errorf("error running docker compose: %w", err)
	}

	dockerContext := NewDockerComposeContext()

	//Loop over output lines
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		container := &DockerComposeContainerContext{}
		if err := json.Unmarshal([]byte(line), container); err != nil {
			return nil, fmt.Errorf("error unmarshalling docker compose output: %w", err)
		}
		dockerContext.AddContainer(container.Name, container)
	}

	return dockerContext, nil
}

func (d *DockerComposeSetup) Setup(ctx map[string]string, logFile io.Writer) (interface{}, error) {
	args_slice := d.generateUpCommand()
	cmd := exec.Command(args_slice[0], args_slice[1:]...)

	cmd.Env = os.Environ()
	for _, envEntry := range d.ExtraEnv {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", envEntry.Key, envEntry.Val))
	}
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running docker compose: %w", err)
	}

	return d.GetContext()
}

func (d *DockerComposeSetup) Teardown(logFile io.Writer) error {
	args_slice := d.generateDownCommand()
	cmd := exec.Command(args_slice[0], args_slice[1:]...)

	cmd.Env = os.Environ()
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running docker compose: %w", err)
	}
	return nil
}
