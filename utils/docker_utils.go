package utils

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
)

// region Docker configuration object ----------------------------------------------------------------------------------

// DockerContainer is used to construct docker container spec for the docker engine.
type DockerContainer struct {
	image      string            // Docker image
	name       string            // Container name
	ports      map[string]string // Container ports mapping
	vars       map[string]string // Environment variables
	labels     map[string]string // Container labels
	entryPoint []string          // Entry point
	autoRemove bool              // Automatically remove container when stopped (default: true)
}

// Name sets the container name.
func (c *DockerContainer) Name(value string) *DockerContainer {
	c.name = value
	return c
}

// Port adds a port mapping
func (c *DockerContainer) Port(external, internal string) *DockerContainer {
	c.ports[external] = internal
	return c
}

// Ports adds multiple port mappings
func (c *DockerContainer) Ports(ports map[string]string) *DockerContainer {
	for k, v := range ports {
		c.ports[k] = v
	}
	return c
}

// Var adds an environment variable
func (c *DockerContainer) Var(key, value string) *DockerContainer {
	c.vars[key] = value
	return c
}

// Vars adds multiple environment variables
func (c *DockerContainer) Vars(vars map[string]string) *DockerContainer {
	for k, v := range vars {
		c.vars[k] = v
	}
	return c
}

// Label adds custom label
func (c *DockerContainer) Label(label, value string) *DockerContainer {
	c.labels[label] = value
	return c
}

// Labels adds multiple labels
func (c *DockerContainer) Labels(label map[string]string) *DockerContainer {
	for k, v := range label {
		c.labels[k] = v
	}
	return c
}

// EntryPoint sets the entrypoint arguments of the container.
func (c *DockerContainer) EntryPoint(args ...string) *DockerContainer {
	c.entryPoint = append(c.entryPoint, args...)
	return c
}

// AutoRemove determines whether to automatically remove the container when it has stopped
func (c *DockerContainer) AutoRemove(value bool) *DockerContainer {
	c.autoRemove = value
	return c
}

// Run builds and run command
func (c *DockerContainer) Run() error {

	// First, stop container if exists
	if c.Exists() {
		_ = c.Stop()
	}

	// construct the docker shell command
	command := "docker"
	args := make([]string, 0)
	args = append(args, "run")

	if len(c.name) > 0 {
		args = append(args, "--name")
		args = append(args, c.name)
	}

	// Expose ports if defined, otherwise, use host network
	if len(c.ports) > 0 {
		for k, v := range c.ports {
			args = append(args, "-p")
			args = append(args, fmt.Sprintf("%s:%s", k, v))
		}
	} else {
		args = append(args, "-h")
		args = append(args, getMachineIP())
	}

	// Add environment variables
	if len(c.vars) > 0 {
		for k, v := range c.vars {
			args = append(args, "-e")
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Add metadata (labels)
	if len(c.labels) > 0 {
		for k, v := range c.labels {
			args = append(args, "-l")
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Add docker image
	if len(c.image) > 0 {
		args = append(args, c.image)
	} else {
		return fmt.Errorf("missing image field")
	}

	// Add entry point
	if len(c.entryPoint) > 0 {
		for _, v := range c.entryPoint {
			args = append(args, v)
		}
	}

	cmd := exec.Command(command, args...)
	go cmd.Run()
	return nil
}

// Stop and kill container
func (c *DockerContainer) Stop() error {
	// construct the docker shell command
	command := "docker"
	args := make([]string, 0)
	args = append(args, "stop")

	if len(c.name) > 0 {
		args = append(args, c.name)
	} else {
		return fmt.Errorf("missing container name")
	}

	cmd := exec.Command(command, args...)
	if er := cmd.Run(); er != nil {
		return er
	}

	args2 := make([]string, 0)
	args2 = append(args2, "rm")
	args2 = append(args2, c.name)

	cmd2 := exec.Command(command, args2...)
	return cmd2.Run()
}

// Exists return true if the container exists
func (c *DockerContainer) Exists() bool {
	cmd := fmt.Sprintf("docker ps -a | grep '%s'", c.name)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return false
	} else {
		return strings.Contains(string(out), c.name)
	}
}

// endregion

// region Singleton Pattern --------------------------------------------------------------------------------------------

type dockerUtils struct {
}

var onlyOnce sync.Once
var dockerUtilsSingleton *dockerUtils = nil

// DockerUtils is a simple utility to execute docker commands using shell
// This utility does not us the docker client library (due to compatability issues) but instead utilizes
// the shell command executor (from os/exec package).
// The pre-requisite
func DockerUtils() (du *dockerUtils) {
	onlyOnce.Do(func() {
		dockerUtilsSingleton = &dockerUtils{}
	})
	return dockerUtilsSingleton
}

// endregion

// region Docker utilities ---------------------------------------------------------------------------------------------

// CreateContainer create a docker container configuration via a fluent interface.
func (c *dockerUtils) CreateContainer(image string) *DockerContainer {
	return &DockerContainer{
		image:      image,
		ports:      make(map[string]string),
		vars:       make(map[string]string),
		labels:     make(map[string]string),
		entryPoint: make([]string, 0),
		autoRemove: true,
	}
}

// StopContainer stops and kill container
func (c *dockerUtils) StopContainer(container string) error {

	// Verify container
	if len(container) == 0 {
		return fmt.Errorf("missing container name")
	}

	// construct the docker stop shell command
	command := "docker"
	args := make([]string, 0)
	args = append(args, "stop", container)
	_ = exec.Command(command, args...).Run()

	// construct the docker rm shell command
	args = nil
	args = make([]string, 0)
	args = append(args, "rm", container)
	_ = exec.Command(command, args...).Run()

	return nil
}

// ContainerExists return true if the container exists
func (c *dockerUtils) ContainerExists(container string) bool {
	cmd := fmt.Sprintf("docker ps -a | grep '%s'", container)
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return false
	} else {
		return strings.Contains(string(out), container)
	}
}

// endregion

// region PRIVATE SECTION ----------------------------------------------------------------------------------------------

// getMachineIP returns one of the IPv4 (which is not the localhost)
func getMachineIP() (result string) {

	// Initialize with localhost
	result = "127.0.0.1"
	results := make([]string, 0)
	results = append(results, result)

	if interfaces, err := net.Interfaces(); err != nil {
		return
	} else {
		for _, i := range interfaces {
			if addresses, er := i.Addrs(); er == nil {
				for _, addr := range addresses {
					ip := addr.String()
					if strings.Count(ip, ".") == 3 {
						if idx := strings.Index(ip, "/"); idx > 0 {
							results = append(results, ip[0:idx])
						}
					}
				}
			}
		}
	}
	return results[len(results)-1]
}

// endregion
