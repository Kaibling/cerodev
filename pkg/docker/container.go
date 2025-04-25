package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/kaibling/cerodev/model"
)

const entrypoint = `
#!/bin/bash

set -e

# Require GIT_REPO
if [ -z "$GIT_REPO" ]; then
  echo "❌ ERROR: GIT_REPO environment variable is not set."
  echo "👉 Set it with: -e GIT_REPO=https://github.com/your/repo.git"
  exit 1
fi

WORKDIR="/home/coder/workspace"

# Clean workspace before clone (optional: use caution)
if [ -n "$(ls -A $WORKDIR 2>/dev/null)" ]; then
  echo "⚠️ Workspace is not empty. Skipping clone."
else
  echo "📥 Cloning $GIT_REPO into $WORKDIR..."
  git clone "$GIT_REPO" "$WORKDIR"
fi

# Start Code Server
echo "🚀 Starting Code Server..."
exec code-server --bind-addr 0.0.0.0:8765 \
  --auth none \
  --user-data-dir /home/coder/.local/share/code-server  \
  "$WORKDIR"
`

func createTar(dockerfileContent string, files map[string]string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Add Dockerfile
	if err := addFileToTar(tw, "Dockerfile", dockerfileContent); err != nil {
		return nil, err
	}

	// Add additional files
	for name, content := range files {
		if err := addFileToTar(tw, name, content); err != nil {
			return nil, err
		}
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

func addFileToTar(tw *tar.Writer, name string, content string) error {
	hdr := &tar.Header{ //nolint:exhaustruct
		Name:    name,
		Mode:    0o644, //nolint:mnd
		Size:    int64(len(content)),
		ModTime: time.Now(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	_, err := tw.Write([]byte(content))

	return err
}

func build(ctx context.Context, cli *client.Client, t model.Template, tag string, buildArgs map[string]*string) error {
	tarBuffer, err := createTar(t.Dockerfile, map[string]string{
		"entrypoint.sh": entrypoint,
	})
	if err != nil {
		return err
	}

	// Use it in ImageBuild
	res, err := cli.ImageBuild(ctx, tarBuffer, types.ImageBuildOptions{ //nolint:exhaustruct
		Tags:       []string{containerPrefix + "-" + t.RepoName + ":" + tag},
		Dockerfile: "Dockerfile",
		Remove:     true,
		BuildArgs:  buildArgs,
	})
	if err != nil {
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Printf("failed to close rows: %s", err.Error()) //nolint:forbidigo
		}
	}()
	// Stream build logs
	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil {
		return err
	}

	return nil
}

func containerCreate(ctx context.Context, cli *client.Client, c Container, volumesPath string) (string, error) {
	exposedPorts := nat.PortSet{}
	for _, port := range c.Ports {
		exposedPorts[nat.Port(port.ContainerPort)] = struct{}{}
	}

	portbindings := nat.PortMap{}
	for _, port := range c.Ports {
		portbindings[nat.Port(port.ContainerPort)] = []nat.PortBinding{ //nolint:exhaustruct,nolintlint
			{HostPort: port.HostPort}, //nolint:exhaustruct
		}
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{ //nolint:exhaustruct
		Image:        c.ImageName,
		Env:          c.Environment,
		ExposedPorts: exposedPorts,
		Volumes:      map[string]struct{}{volumesPath: {}},
	}, &container.HostConfig{ //nolint:exhaustruct
		PortBindings: portbindings,
		Binds:        []string{volumesPath + ":/home/coder/workspace"},
	}, nil, nil, c.ContainerName)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func containerStart(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerStart(ctx, containerID, container.StartOptions{}) //nolint:exhaustruct
}

func containerStop(ctx context.Context, cli *client.Client, containerID string) error {
	timeout := 10

	err := cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}) //nolint:exhaustruct
	if err != nil {
		return err
	}

	return nil
}

func containerDelete(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerRemove(ctx, containerID, container.RemoveOptions{ //nolint:exhaustruct
		Force: true,
	})
}

func (c Container) Status(ctx context.Context, cli *client.Client) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{ //nolint:exhaustruct
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", c.ContainerName)),
	})
	if err != nil {
		log.Printf("🚫 Container error: %v", err)

		return
	}

	if len(containers) == 0 {
		log.Printf("🚫 Container not found")

		return
	}
}

func getAllContainerStatuses(ctx context.Context, cli *client.Client, containerIDs []string) ([]model.ContainerStatus, error) { //nolint: lll
	// Create a filter for the provided container IDs
	filtersArgs := filters.NewArgs()
	for _, id := range containerIDs {
		filtersArgs.Add("id", id)
	}

	containers, err := cli.ContainerList(ctx, container.ListOptions{ //nolint:exhaustruct
		All:     true,
		Filters: filtersArgs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	containerStatuses := []model.ContainerStatus{}
	for _, c := range containers {
		containerStatuses = append(containerStatuses, model.ContainerStatus{
			DockerID: c.ID,
			Status:   c.Status,
			State:    c.State,
		})
	}

	return containerStatuses, nil
}

// func getUserContainerStatuses(ctx context.Context, cli *client.Client, username string) {
// 	searchPrefix := containerPrefix + "-" + username + "-"
// 	containers, err := cli.ContainerList(ctx, container.ListOptions{
// 		All:     true,
// 		Filters: filters.NewArgs(filters.Arg("name", searchPrefix)),
// 	})
// 	if err != nil {
// 		log.Fatalf("🚫 : %v", err)
// 	}

// 	//var matchingContainers []container.Summary
// 	var containerStatuses []ContainerStatus
// 	for _, c := range containers {
// 		for _, name := range c.Names {
// 			// Docker prepends "/" to container names
// 			if strings.HasPrefix(name, "/"+searchPrefix) {
// 				containerStatuses = append(containerStatuses, ContainerStatus{
// 					ContainerName: name[1:],
// 					Status:        c.Status,
// 					State:         c.State,
// 				})
// 				break
// 			}
// 		}
// 	}
// 	for _, m := range containerStatuses {
// 		a, _ := json.MarshalIndent(m, "", "  ")
// 		fmt.Println(string(a))
// 	}
// }

func getImages(ctx context.Context, cli *client.Client) ([]model.Image, error) {
	images, err := cli.ImageList(ctx, image.ListOptions{}) //nolint:exhaustruct
	if err != nil {
		return nil, err
	}

	imageList := []model.Image{}

	for _, image := range images {
		for _, repoTag := range image.RepoTags {
			repo := strings.Split(repoTag, ":")
			if strings.HasPrefix(repo[0], "cd-") {
				imageList = append(imageList, model.Image{
					RepoName: repo[0],
					ImageID:  image.ID,
					Tag:      repo[1],
				})
			}
		}
	}

	return imageList, nil
}
