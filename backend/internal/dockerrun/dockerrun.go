// Package dockerrun provides functionality for running Docker containers.
// It handles container creation, starting, stopping, and removal, with Traefik label configuration.
package dockerrun

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// Runner handles running Docker containers with Traefik routing labels.
// It wraps the Docker client and provides methods for container lifecycle management.
type Runner struct {
	// client is the Docker API client used to communicate with the Docker daemon
	client *client.Client
}

// NewRunner creates a new Runner instance connected to the Docker daemon.
//
// Parameters:
//   - dockerHost: The Docker daemon address (e.g., "unix:///var/run/docker.sock" or "tcp://host:port")
//
// Returns:
//   - *Runner: A new Runner instance ready to run containers, or nil on error
//   - error: Error if Docker client creation fails (connection issue, invalid host, etc.)
func NewRunner(dockerHost string) (*Runner, error) {
	// Create Docker client with host configuration and API version negotiation
	cli, err := client.NewClientWithOpts(
		client.WithHost(dockerHost),
		client.WithAPIVersionNegotiation(), // Automatically negotiate API version with daemon
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &Runner{client: cli}, nil
}

// Run creates and starts a Docker container with Traefik routing labels.
// The container will be accessible via Traefik at {subdomain}.{baseDomain}.
//
// Traefik labels configured:
//   - traefik.enable: Enables Traefik to route to this container
//   - traefik.http.routers.{subdomain}.rule: Host routing rule
//   - traefik.http.services.{subdomain}.loadbalancer.server.port: Container port (80)
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - imageName: The Docker image name to run (e.g., "mvp-myapp:123")
//   - subdomain: The subdomain for this deployment (e.g., "myapp-123")
//   - baseDomain: The base domain for routing (e.g., "localhost" or "example.com")
//
// Returns:
//   - string: The Docker container ID, or empty string on error
//   - error: Error if container creation fails, image not found, or container cannot start
func (r *Runner) Run(ctx context.Context, imageName, subdomain, baseDomain string) (string, error) {
	// Configure container settings
	containerConfig := &container.Config{
		Image: imageName, // The Docker image to run
		// Traefik labels for automatic routing
		Labels: map[string]string{
			// Enable Traefik for this container
			"traefik.enable": "true",
			// Configure routing rule: requests to {subdomain}.{baseDomain} go to this container
			"traefik.http.routers." + subdomain + ".rule": fmt.Sprintf("Host(`%s.%s`)", subdomain, baseDomain),
			// Configure the port Traefik should forward to (container's port 8080)
			"traefik.http.services." + subdomain + ".loadbalancer.server.port": "8080",
		},
	}

	// Configure host-specific settings
	hostConfig := &container.HostConfig{
		AutoRemove: false, // Don't auto-remove container when it stops
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped", // Restart container unless explicitly stopped
		},
	}

	// Configure networking (using default bridge network)
	networkConfig := &network.NetworkingConfig{}

	// Create the container (but don't start it yet)
	// The container name is set to the subdomain for easy identification
	resp, err := r.client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, nil, subdomain)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	if err := r.client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	// Return the container ID for tracking
	return resp.ID, nil
}

// Stop stops a running container gracefully.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - containerID: The Docker container ID to stop
//
// Returns:
//   - error: Error if container cannot be stopped or doesn't exist
func (r *Runner) Stop(ctx context.Context, containerID string) error {
	return r.client.ContainerStop(ctx, containerID, container.StopOptions{})
}

// Remove removes a container from the system.
// The Force option ensures removal even if the container is running.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - containerID: The Docker container ID to remove
//
// Returns:
//   - error: Error if container cannot be removed
func (r *Runner) Remove(ctx context.Context, containerID string) error {
	return r.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}
