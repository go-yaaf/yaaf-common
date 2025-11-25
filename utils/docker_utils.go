package utils

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
)

// region Docker configuration object ----------------------------------------------------------------------------------

// DockerContainer provides a fluent API for configuring and managing a Docker container.
// It simplifies the process of building and executing `docker` commands.
type DockerContainer struct {
	image      string            // The Docker image to use for the container.
	name       string            // The name to assign to the container.
	ports      map[string]string // A map of port mappings from host to container (e.g., "8080:80").
	vars       map[string]string // A map of environment variables to set in the container.
	labels     map[string]string // A map of labels to apply to the container.
	entryPoint []string          // The entrypoint command and arguments for the container.
	autoRemove bool              // If true, the container will be automatically removed on exit.
}

// Name sets the name of the container.
func (c *DockerContainer) Name(value string) *DockerContainer {
	c.name = value
	return c
}

// Port adds a single port mapping to the container configuration.
func (c *DockerContainer) Port(external, internal string) *DockerContainer {
	c.ports[external] = internal
	return c
}

// Ports adds multiple port mappings to the container configuration.
func (c *DockerContainer) Ports(ports map[string]string) *DockerContainer {
	for k, v := range ports {
		c.ports[k] = v
	}
	return c
}

// Var adds a single environment variable to the container configuration.
func (c *DockerContainer) Var(key, value string) *DockerContainer {
	c.vars[key] = value
	return c
}

// Vars adds multiple environment variables to the container configuration.
func (c *DockerContainer) Vars(vars map[string]string) *DockerContainer {
	for k, v := range vars {
		c.vars[k] = v
	}
	return c
}

// Label adds a single label to the container configuration.
func (c *DockerContainer) Label(label, value string) *DockerContainer {
	c.labels[label] = value
	return c
}

// Labels adds multiple labels to the container configuration.
func (c *DockerContainer) Labels(labels map[string]string) *DockerContainer {
	for k, v := range labels {
		c.labels[k] = v
	}
	return c
}

// EntryPoint sets the entrypoint for the container.
func (c *DockerContainer) EntryPoint(args ...string) *DockerContainer {
	c.entryPoint = args
	return c
}

// AutoRemove sets whether the container should be automatically removed when it stops.
func (c *DockerContainer) AutoRemove(value bool) *DockerContainer {
	c.autoRemove = value
	return c
}

// Run executes the `docker run` command based on the configured container settings.
// If the container is already running, it does nothing.
func (c *DockerContainer) Run() error {
	if c.IsRunning() {
		return nil
	}

	args := []string{"run", "-d"} // Run in detached mode.

	if c.autoRemove {
		args = append(args, "--rm")
	}
	if c.name != "" {
		args = append(args, "--name", c.name)
	}

	for ext, in := range c.ports {
		args = append(args, "-p", fmt.Sprintf("%s:%s", ext, in))
	}

	if len(c.ports) == 0 {
		args = append(args, "-h", getMachineIP())
	}

	for k, v := range c.vars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	for k, v := range c.labels {
		args = append(args, "-l", fmt.Sprintf("%s=%s", k, v))
	}

	if c.image == "" {
		return fmt.Errorf("docker image is not specified")
	}
	args = append(args, c.image)
	args = append(args, c.entryPoint...)

	cmd := exec.Command("docker", args...)
	return cmd.Run()
}

// Stop stops and removes the container.
func (c *DockerContainer) Stop() error {
	if c.name == "" {
		return fmt.Errorf("container name is not specified")
	}

	// Stop the container.
	stopCmd := exec.Command("docker", "stop", c.name)
	if err := stopCmd.Run(); err != nil {
		// Ignore "No such container" errors, as the container might already be stopped.
		if !strings.Contains(err.Error(), "No such container") {
			return fmt.Errorf("failed to stop container %s: %w", c.name, err)
		}
	}

	// If auto-remove is not enabled, manually remove the container.
	if !c.autoRemove {
		rmCmd := exec.Command("docker", "rm", c.name)
		if err := rmCmd.Run(); err != nil {
			if !strings.Contains(err.Error(), "No such container") {
				return fmt.Errorf("failed to remove container %s: %w", c.name, err)
			}
		}
	}
	return nil
}

// Exists checks if a container with the configured name exists (either running or stopped).
func (c *DockerContainer) Exists() bool {
	if c.name == "" {
		return false
	}
	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", c.name))
	out, err := cmd.Output()
	return err == nil && strings.Contains(string(out), c.name)
}

// IsRunning checks if a container with the configured name is currently running.
func (c *DockerContainer) IsRunning() bool {
	if c.name == "" {
		return false
	}
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", c.name), "--filter", "status=running")
	out, err := cmd.Output()
	return err == nil && strings.Contains(string(out), c.name)
}

// endregion

// region Singleton Pattern --------------------------------------------------------------------------------------------

// dockerUtils is a singleton struct providing Docker utility functions.
type dockerUtils struct{}

var (
	dockerUtilsSingleton *dockerUtils
	once                 sync.Once
)

// DockerUtils returns a singleton instance of the dockerUtils.
// This utility uses shell commands to interact with Docker, avoiding direct dependency on the Docker client library.
func DockerUtils() *dockerUtils {
	once.Do(func() {
		dockerUtilsSingleton = &dockerUtils{}
	})
	return dockerUtilsSingleton
}

// endregion

// region Docker utilities ---------------------------------------------------------------------------------------------

// CreateContainer creates a new DockerContainer configuration with the specified image.
func (c *dockerUtils) CreateContainer(image string) *DockerContainer {
	return &DockerContainer{
		image:      image,
		ports:      make(map[string]string),
		vars:       make(map[string]string),
		labels:     make(map[string]string),
		autoRemove: true,
	}
}

// StopContainer stops and removes a container by its name.
func (c *dockerUtils) StopContainer(containerName string) error {
	if containerName == "" {
		return fmt.Errorf("container name is required")
	}

	// Stop the container.
	stopCmd := exec.Command("docker", "stop", containerName)
	_ = stopCmd.Run() // Ignore error, as container might not be running.

	// Remove the container.
	rmCmd := exec.Command("docker", "rm", containerName)
	_ = rmCmd.Run() // Ignore error, as container might already be removed.

	return nil
}

// ContainerExists checks if a container with the given name exists.
func (c *dockerUtils) ContainerExists(containerName string) bool {
	if containerName == "" {
		return false
	}
	cmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", containerName))
	out, err := cmd.Output()
	return err == nil && strings.Contains(string(out), containerName)
}

// endregion

// region PRIVATE SECTION ----------------------------------------------------------------------------------------------

// getMachineIP returns a non-localhost IPv4 address of the machine.
func getMachineIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}

	for _, addr := range addresses {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

// endregion
