package docker

import (
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

type ContainerConfig struct {
	Name        string
	Image       string
	Env         []string
	Labels      map[string]string
	NetworkName string
}

type ContainerStatus struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	State   string `json:"state"`
	Status  string `json:"status"`
	Running bool   `json:"running"`
}

func NewClient() (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &Client{cli: cli}, nil
}

func (d *Client) Close() error {
	return d.cli.Close()
}

// CreateAndStartContainer creates a container and starts it.
func (d *Client) CreateAndStartContainer(ctx context.Context, cfg ContainerConfig) (string, error) {
	containerCfg := &container.Config{
		Image:  cfg.Image,
		Env:    cfg.Env,
		Labels: cfg.Labels,
	}

	hostCfg := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: container.RestartPolicyUnlessStopped,
		},
	}

	networkCfg := &network.NetworkingConfig{}
	if cfg.NetworkName != "" {
		networkCfg.EndpointsConfig = map[string]*network.EndpointSettings{
			cfg.NetworkName: {},
		}
	}

	resp, err := d.cli.ContainerCreate(ctx, containerCfg, hostCfg, networkCfg, nil, cfg.Name)
	if err != nil {
		return "", fmt.Errorf("failed to create container %s: %w", cfg.Name, err)
	}

	if err := d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container %s: %w", cfg.Name, err)
	}

	return resp.ID, nil
}

// GetContainerStatus returns the status of a container by name.
func (d *Client) GetContainerStatus(ctx context.Context, containerName string) (*ContainerStatus, error) {
	args := filters.NewArgs()
	args.Add("name", containerName)

	containers, err := d.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: args,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	// Find exact match (docker name filter is a prefix match)
	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName {
				return &ContainerStatus{
					Name:    containerName,
					ID:      c.ID[:12],
					State:   c.State,
					Status:  c.Status,
					Running: c.State == "running",
				}, nil
			}
		}
	}

	return &ContainerStatus{
		Name:    containerName,
		State:   "not_found",
		Running: false,
	}, nil
}

// PullImage pulls a Docker image. Useful to ensure images are available before creating containers.
func (d *Client) PullImage(ctx context.Context, imageName string) error {
	reader, err := d.cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageName, err)
	}
	defer reader.Close()

	// Drain the reader to complete the pull
	_, _ = io.Copy(io.Discard, reader)

	return nil
}
