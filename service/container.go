package service

import (
	"os"
	"strconv"
	"strings"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/config"
	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/utils"
)

type dbrepo interface {
	GetByID(id string) (*model.Container, error)
	GetAll() ([]model.Container, error)
	Create(container *model.Container) (*model.Container, error)
	Delete(id string) error
	Update(container *model.Container) (*model.Container, error)
	ReleasePort(containerID string) error
	GetFreePort() (int, error)
	AllocatePort(containerID string, port int) error
	GetPortCount() (int, error)
	FillPorts(minPort, maxPort int) error
}

type dockerrepo interface {
	StartContainer(containerID string) error
	CreateContainer(container *model.Container) (string, error)
	StopContainer(containerID string) error
	DeleteContainer(containerID string) error
	GetContainerStatuses(containerID []string) ([]model.ContainerStatus, error)
	Build(t model.Template, tag string, env map[string]*string) error
	GetImages() ([]model.Image, error)
}

type ContainerService struct {
	dbrepo       dbrepo
	dockerrepo   dockerrepo
	templaterepo templaterepo
	l            log.Writer
	cfg          config.Configuration
}

func NewContainerService(dbrepo dbrepo,
	dockerrepo dockerrepo,
	templaterepo templaterepo,
	l log.Writer,
	cfg config.Configuration,
) *ContainerService {
	return &ContainerService{
		dbrepo:       dbrepo,
		dockerrepo:   dockerrepo,
		templaterepo: templaterepo,
		l:            l.Named("container_service"),
		cfg:          cfg,
	}
}

func (s *ContainerService) GetByID(id string) (*model.Container, error) {
	container, err := s.dbrepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	status, err := s.dockerrepo.GetContainerStatuses([]string{container.DockerID})
	if err != nil {
		return nil, err
	}

	container.Status = status[0].Status
	container.State = status[0].State

	return container, nil
}

func (s *ContainerService) GetAll() ([]model.Container, error) {
	containers, err := s.dbrepo.GetAll()
	if err != nil {
		return nil, err
	}

	containerIDs := make([]string, len(containers))
	for i, c := range containers {
		containerIDs[i] = c.DockerID
	}

	statuses, err := s.dockerrepo.GetContainerStatuses(containerIDs)
	if err != nil {
		return nil, err
	}

	for i, c := range containers {
		for _, s := range statuses {
			if c.DockerID == s.DockerID {
				containers[i].Status = s.Status
				containers[i].State = s.State
			}
		}
	}

	return containers, nil
}

func (s *ContainerService) Create(container *model.Container) (*model.Container, error) {
	container.ID = utils.GenerateULID()

	freePort, err := s.dbrepo.GetFreePort()
	if err != nil {
		s.l.Error("Failed to get free port", err)

		return nil, err
	}

	// add vscode port to container
	if err := s.dbrepo.AllocatePort(container.ID, freePort); err != nil {
		return nil, err
	}

	container.Ports = append(container.Ports, strconv.Itoa(freePort)+":8765/tcp")

	// container data validation
	container.ContainerName = utils.ContainerName(container.UserID, container.GitRepo)
	container.EnvVars = utils.RemoveEmptyStrings(container.EnvVars)
	container.Ports = utils.RemoveEmptyStrings(container.Ports)
	container.EnvVars = append(container.EnvVars, "GIT_REPO="+container.GitRepo)

	ctrID, err := s.dockerrepo.CreateContainer(container)
	if err != nil {
		return nil, err
	}

	container.DockerID = ctrID

	return s.dbrepo.Create(container)
}

func (s *ContainerService) Update(container *model.Container) (*model.Container, error) {
	return s.dbrepo.Update(container)
}

func (s *ContainerService) StartContainer(containerID string) error {
	m, err := s.GetByID(containerID)
	if err != nil {
		return err
	}

	return s.dockerrepo.StartContainer(m.DockerID)
}

func (s *ContainerService) StopContainer(containerID string) error {
	m, err := s.GetByID(containerID)
	if err != nil {
		return err
	}

	return s.dockerrepo.StopContainer(m.DockerID)
}

func (s *ContainerService) DeleteContainer(containerID string) error {
	m, err := s.GetByID(containerID)
	if err != nil {
		return err
	}

	if err := s.dockerrepo.DeleteContainer(m.DockerID); err != nil {
		if !strings.Contains(err.Error(), "No such container") {
			// Container already deleted, no need to return an error
			return err
		}
	}

	// delete volumes directory
	volumeDir := s.cfg.VolumesPath + "/" + m.ID
	if err := os.RemoveAll(volumeDir); err != nil {
		s.l.Warn("Failed to remove volumes directory: %s", err.Error())

		return err
	}

	s.l.Debug("Removed volumes directory: %s", volumeDir)

	if err := s.dbrepo.ReleasePort(containerID); err != nil {
		return err
	}

	s.l.Debug("Released UI port: %s", m.UIPort)

	return s.dbrepo.Delete(containerID)
}

func (s *ContainerService) BuildTemplate(templateID string, tag string, env map[string]*string) error {
	t, err := s.templaterepo.GetByID(templateID)
	if err != nil {
		return err
	}

	return s.dockerrepo.Build(*t, tag, env)
}

func (s *ContainerService) GetImages() ([]model.Image, error) {
	return s.dockerrepo.GetImages()
}

func (s *ContainerService) GetPortCount() (int, error) {
	return s.dbrepo.GetPortCount()
}

func (s *ContainerService) FillPorts(minPort, maxPort int) error {
	return s.dbrepo.FillPorts(minPort, maxPort)
}
