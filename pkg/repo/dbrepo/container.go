package dbrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/repo/sqlcrepo"
	_ "modernc.org/sqlite" // Import the SQLite driver
)

type ContainerRepo struct {
	ctx context.Context
	db  *sql.DB
	l   log.Writer
}

func NewContainerRepo(ctx context.Context, db *sql.DB, l log.Writer) *ContainerRepo {
	return &ContainerRepo{ctx: ctx, db: db, l: l.Named("repo_container")}
}

func (r *ContainerRepo) GetByID(id string) (*model.Container, error) {
	container, err := sqlcrepo.New(r.db).GetContainerByID(r.ctx, id)
	if err != nil {
		return nil, ToAppError(fmt.Errorf("GetContainerByID failed: %w", err))
	}

	env := container.EnvVars.String
	ports := container.Ports.String

	return &model.Container{
		ID:            container.ID,
		DockerID:      container.DockerID,
		ContainerName: container.ContainerName,
		ImageName:     container.ImageName,
		EnvVars:       splitString(env),
		Ports:         splitString(ports),
		UIPort:        strconv.FormatInt(container.UiPort, 10),
	}, nil
}

func (r *ContainerRepo) GetAll() ([]model.Container, error) {
	containers, err := sqlcrepo.New(r.db).GetAllContainers(r.ctx)
	if err != nil {
		return nil, ToAppError(fmt.Errorf("GetAllContainers failed: %w", err))
	}

	result := []model.Container{}

	for _, container := range containers {
		env := container.EnvVars.String
		ports := container.Ports.String

		result = append(result, model.Container{
			ID:            container.ID,
			DockerID:      container.DockerID,
			ContainerName: container.ContainerName,
			ImageName:     container.ImageName,
			EnvVars:       splitString(env),
			Ports:         splitString(ports),
			UIPort:        strconv.FormatInt(container.UiPort, 10),
		})
	}

	return result, nil
}

func (r *ContainerRepo) Create(container *model.Container) (*model.Container, error) {
	containerID, err := sqlcrepo.New(r.db).CreateContainer(r.ctx, sqlcrepo.CreateContainerParams{ //nolint:exhaustruct
		ID:            container.ID,
		DockerID:      container.DockerID,
		ContainerName: container.ContainerName,
		ImageName:     container.ImageName,
		EnvVars:       sql.NullString{String: joinStrings(container.EnvVars), Valid: true},
		Ports:         sql.NullString{String: joinStrings(container.Ports), Valid: true},
	})
	if err != nil {
		return nil, ToAppError(err)
	}

	return r.GetByID(containerID)
}

func (r *ContainerRepo) Delete(id string) error {
	return sqlcrepo.New(r.db).DeleteContainer(r.ctx, id)
}

func (r *ContainerRepo) Update(container *model.Container) (*model.Container, error) {
	err := sqlcrepo.New(r.db).UpdateContainer(r.ctx, sqlcrepo.UpdateContainerParams{ //nolint:exhaustruct
		ID:            container.ID,
		DockerID:      container.DockerID,
		ContainerName: container.ContainerName,
		ImageName:     container.ImageName,
		EnvVars:       sql.NullString{String: joinStrings(container.EnvVars), Valid: true},
		Ports:         sql.NullString{String: joinStrings(container.Ports), Valid: true},
	})
	if err != nil {
		return nil, ToAppError(err)
	}

	return r.GetByID(container.ID)
}

func (r *ContainerRepo) ReleasePort(containerID string) error {
	return sqlcrepo.New(r.db).ReleasePortByContainer(r.ctx, sql.NullString{String: containerID, Valid: true})
}

func (r *ContainerRepo) GetFreePort() (int, error) {
	port, err := sqlcrepo.New(r.db).GetFreePort(r.ctx)
	if err != nil {
		return 0, ToAppError(err)
	}

	return int(port), nil
}

func (r *ContainerRepo) AllocatePort(containerID string, port int) error {
	return sqlcrepo.New(r.db).AllocatePort(r.ctx, sqlcrepo.AllocatePortParams{
		ContainerID: sql.NullString{String: containerID, Valid: true},
		Port:        int64(port),
	},
	)
}

func (r *ContainerRepo) GetPortCount() (int, error) {
	count, err := sqlcrepo.New(r.db).GetPortCount(r.ctx)

	return int(count), ToAppError(err)
}

func (r *ContainerRepo) FillPorts(minPort, maxPort int) error {
	tx, err := r.db.BeginTx(r.ctx, nil)
	if err != nil {
		return ToAppError(err)
	}

	qtx := sqlcrepo.New(tx)

	for port := minPort; port <= maxPort; port++ {
		err := qtx.CreatePort(r.ctx, int64(port))
		if err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				r.l.Error("failed to rollback transaction", rerr)

				return rerr
			}

			return ToAppError(err)
		}
	}

	return tx.Commit()
}

func joinStrings(s []string) string {
	if len(s) == 0 {
		return ""
	}

	return fmt.Sprintf("{%s}", stringJoin(s, ","))
}

// Helper function to join string slices with separator.
func stringJoin(s []string, sep string) string {
	result := ""

	for i, str := range s {
		if i > 0 {
			result += sep
		}

		result += str
	}

	return result
}

func splitString(s string) []string {
	if len(s) <= 2 { //nolint:mnd
		return []string{}
	}

	return strings.Split(s[1:len(s)-1], ",")
}
