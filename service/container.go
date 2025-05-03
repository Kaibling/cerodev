package service

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/config"
	"github.com/kaibling/cerodev/errs"
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
		return nil, fmt.Errorf("failed to GetByID: %w", err)
	}

	status, err := s.dockerrepo.GetContainerStatuses([]string{container.DockerID})
	if err != nil {
		return nil, fmt.Errorf("failed to provider GetContainerStatuses: %w", err)
	}

	if len(status) == 0 {
		return container, errs.ErrContainerNotInProvider
	}

	container.Status = status[0].Status
	container.State = status[0].State

	return container, nil
}

func (s *ContainerService) GetAll() ([]model.Container, error) {
	containers, err := s.dbrepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to GetAll: %w", err)
	}

	containerIDs := make([]string, len(containers))
	for i, c := range containers {
		containerIDs[i] = c.DockerID
	}

	statuses, err := s.dockerrepo.GetContainerStatuses(containerIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to GetContainerStatuses: %w", err)
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

		return nil, fmt.Errorf("failed to GetFreePort: %w", err)
	}

	// add vscode port to container
	if err := s.dbrepo.AllocatePort(container.ID, freePort); err != nil {
		return nil, fmt.Errorf("failed to AllocatePort: %w", err)
	}

	container.Ports = append(container.Ports, strconv.Itoa(freePort)+":8765/tcp")

	// container data validation
	container.ContainerName = utils.ContainerName(container.UserID, container.GitRepo)
	container.EnvVars = utils.RemoveEmptyStrings(container.EnvVars)
	container.Ports = utils.RemoveEmptyStrings(container.Ports)
	container.EnvVars = append(container.EnvVars, "GIT_REPO="+container.GitRepo)

	ctrID, err := s.dockerrepo.CreateContainer(container)
	if err != nil {
		return nil, fmt.Errorf("failed to CreateContainer: %w", err)
	}

	container.DockerID = ctrID

	val, err := s.dbrepo.Create(container)

	return HandleError[*model.Container](val, err, "failed to db Create")
}

func (s *ContainerService) Update(container *model.Container) (*model.Container, error) {
	return s.dbrepo.Update(container)
}

func (s *ContainerService) StartContainer(containerID string) error {
	m, err := s.GetByID(containerID)
	if err != nil {
		return fmt.Errorf("failed to GetByID: %w", err)
	}

	if err := s.dockerrepo.StartContainer(m.DockerID); err != nil {
		return fmt.Errorf("failed to provider StartContainer: %w", err)
	}

	return nil
}

func (s *ContainerService) StopContainer(containerID string) error {
	m, err := s.GetByID(containerID)
	if err != nil {
		return fmt.Errorf("failed to GetByID: %w", err)
	}

	if err := s.dockerrepo.StopContainer(m.DockerID); err != nil {
		return fmt.Errorf("failed to provider StopContainer: %w", err)
	}

	return nil
}

func (s *ContainerService) DeleteContainer(containerID string) error {
	m, err := s.GetByID(containerID)
	if err != nil { //nolint:nestif
		if errors.Is(err, errs.ErrContainerNotInProvider) {
			s.l.Warn("container is not in provider. skip deletion in provider")
		} else {
			return fmt.Errorf("could not get container: %w", err)
		}
	} else {
		if err := s.dockerrepo.DeleteContainer(m.DockerID); err != nil {
			if !strings.Contains(err.Error(), "No such container") {
				return fmt.Errorf("failed to DeleteContainer: %w", err)
			}
		}
	}

	// delete volumes directory
	volumeDir := s.cfg.VolumesPath + "/" + m.ID
	if err := os.RemoveAll(volumeDir); err != nil {
		s.l.Warn("Failed to remove volumes directory: %s", err.Error())

		return fmt.Errorf("failed to delete volume: %w", err)
	}

	s.l.Debug("Removed volumes directory: %s", volumeDir)

	if err := s.dbrepo.ReleasePort(containerID); err != nil {
		return fmt.Errorf("failed to ReleasePort: %w", err)
	}

	s.l.Debug("Released UI port: %s", m.UIPort)

	if err := s.dbrepo.Delete(containerID); err != nil {
		return fmt.Errorf("failed to db Delete: %w", err)
	}

	return nil
}

func (s *ContainerService) BuildTemplate(templateID string, tag string, env map[string]*string) error {
	t, err := s.templaterepo.GetByID(templateID)
	if err != nil {
		return fmt.Errorf("failed to GetByID: %w", err)
	}

	env["ARCHITECTURE"] = &config.Architecture

	if err := s.dockerrepo.Build(*t, tag, env); err != nil {
		return fmt.Errorf("failed to provider Build: %w", err)
	}

	return nil
}

func (s *ContainerService) GetImages() ([]model.Image, error) {
	val, err := s.dockerrepo.GetImages()

	return HandleError[[]model.Image](val, err, "failed to GetImages")
}

func (s *ContainerService) GetPortCount() (int, error) {
	val, err := s.dbrepo.GetPortCount()

	return HandleError[int](val, err, "failed to GetPortCount")
}

func (s *ContainerService) FillPorts(minPort, maxPort int) error {
	if err := s.dbrepo.FillPorts(minPort, maxPort); err != nil {
		return fmt.Errorf("failed to db FillPorts: %w", err)
	}

	return nil
}
