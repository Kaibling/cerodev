package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
	"github.com/kaibling/cerodev/model"
)

type Repo struct {
	cli         *client.Client
	ctx         context.Context //nolint:containedctx
	volumesPath string
}

func NewRepo(ctx context.Context, volumesPath string) *Repo {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}

	return &Repo{cli, ctx, volumesPath}
}

func (r *Repo) StartContainer(containerID string) error {
	return containerStart(r.ctx, r.cli, containerID)
}

func (r *Repo) CreateContainer(mc *model.Container) (string, error) {
	c := unmarshalContainer(*mc)

	return containerCreate(r.ctx, r.cli, c, r.volumesPath+"/"+mc.ID)
}

func (r *Repo) StopContainer(containerID string) error {
	return containerStop(r.ctx, r.cli, containerID)
}

func (r *Repo) DeleteContainer(containerID string) error {
	return containerDelete(r.ctx, r.cli, containerID)
}

func (r *Repo) GetContainerStatuses(containerID []string) ([]model.ContainerStatus, error) {
	return getAllContainerStatuses(r.ctx, r.cli, containerID)
}

func (r *Repo) Build(t model.Template, tag string, env map[string]*string) error {
	return build(r.ctx, r.cli, t, tag, env)
}

func (r *Repo) GetImages() ([]model.Image, error) {
	return getImages(r.ctx, r.cli)
}

func unmarshalContainer(c model.Container) Container {
	ports := make([]Port, len(c.Ports))

	for i, p := range c.Ports {
		p := strings.Split(p, ":")
		ports[i] = Port{
			HostPort:      p[0],
			ContainerPort: p[1],
		}
	}

	return Container{
		ContainerID:   c.ID,
		ContainerName: c.ContainerName,
		ImageName:     c.ImageName,
		GitRepo:       c.GitRepo,
		UserID:        c.UserID,
		Environment:   c.EnvVars,
		Ports:         ports,
	}
}
