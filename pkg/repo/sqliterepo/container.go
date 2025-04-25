package sqliterepo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/model"
)

type Container struct {
	ID            string
	DockerID      string
	ImageName     string
	ContainerName string
	GitRepo       string
	UserID        string
	Environment   []string
	Ports         []string
	UIPort        string
}

type ContainerRepo struct {
	db *sql.DB
	l  log.Writer
}

func NewContainerRepo(db *sql.DB, l log.Writer) *ContainerRepo {
	return &ContainerRepo{db: db, l: l.Named("repo_container")}
}

func (r *ContainerRepo) Create(container *model.Container) error {
	_, err := r.db.Exec(`
		INSERT INTO containers (id,dockerid, imagename, containername, gitrepo, user_id, environment, ports) 
		VALUES (?,?, ?, ?, ?, ?, ?, ?)`,
		container.ID, container.DockerID, container.ImageName, container.ContainerName, container.GitRepo, container.UserID,
		joinStrings(container.EnvVars), joinStrings(container.Ports))
	if err != nil {
		return fmt.Errorf("failed to create container: %w", err)
	}

	return nil
}

func (r *ContainerRepo) Delete(containerID string) error {
	_, err := r.db.Exec("DELETE FROM containers WHERE id = ?", containerID)
	if err != nil {
		return fmt.Errorf("failed to delete container: %w", err)
	}

	return nil
}

func (r *ContainerRepo) GetByID(containerID string) (*model.Container, error) {
	row := r.db.QueryRow("SELECT c.id, c.dockerid, c.imagename, c.containername, c.gitrepo, c.user_id, c.environment, c.ports,p.port FROM containers c JOIN ports p on p.container_id = c.id WHERE id = ?", containerID) //nolint: lll

	var container Container

	var environment, ports string

	err := row.Scan(&container.ID, &container.DockerID, &container.ImageName, &container.ContainerName, &container.GitRepo, &container.UserID, &environment, &ports, &container.UIPort) //nolint: lll
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no container found with ID %s: %w", containerID, err)
		}

		return nil, fmt.Errorf("failed to get container: %w", err)
	}

	container.Environment = splitString(environment)
	container.Ports = splitString(ports)

	return marshalContainer(&container), nil
}

func (r *ContainerRepo) GetAll() ([]model.Container, error) {
	rows, err := r.db.Query("SELECT c.id, c.dockerid, c.imagename, c.containername, c.gitrepo, c.user_id, c.environment, c.ports,p.port FROM containers c JOIN ports p on p.container_id = c.id") //nolint: lll
	if err != nil {
		return nil, fmt.Errorf("failed to fetch containers: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to fetch containers: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.l.Error("failed to close rows: %w", err)
		}
	}()

	containers := []Container{}

	for rows.Next() {
		var container Container

		var environment, ports string

		if err := rows.Scan(&container.ID, &container.DockerID, &container.ImageName, &container.ContainerName, &container.GitRepo, &container.UserID, &environment, &ports, &container.UIPort); err != nil { //nolint: lll
			return nil, fmt.Errorf("failed to scan container: %w", err)
		}

		container.Environment = splitString(environment)
		container.Ports = splitString(ports)
		containers = append(containers, container)
	}

	return marshalContainers(containers), nil
}

func (r *ContainerRepo) GetByUser(userID string) ([]model.Container, error) {
	rows, err := r.db.Query("SELECT id, dockerid, imagename, containername, gitrepo, user_id, environment, ports FROM containers WHERE user_id = ?", userID) //nolint: lll
	if err != nil {
		return nil, fmt.Errorf("failed to fetch containers: %w", err)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to fetch containers: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.l.Error("failed to close rows: %v", err)
		}
	}()

	var containers []Container

	for rows.Next() {
		var container Container

		var environment, ports string

		if err := rows.Scan(&container.ID, &container.DockerID, &container.ImageName, &container.ContainerName, &container.GitRepo, &container.UserID, &environment, &ports); err != nil { //nolint: lll
			return nil, fmt.Errorf("failed to scan container: %w", err)
		}

		container.Environment = splitString(environment)
		container.Ports = splitString(ports)
		containers = append(containers, container)
	}

	return marshalContainers(containers), nil
}

// func (r *ContainerRepo) GetAllocatedPort(containerID string) (int, error) {
// 	var port int
// 	err := r.db.QueryRow("SELECT port FROM ports WHERE container_id = ?", containerID).Scan(&port)
// 	if err != nil {
// 		return 0, fmt.Errorf("cannot fetch port: %w", err)
// 	}
// 	return port, nil
// }

func (r *ContainerRepo) AllocatePort(containerID string) (int, error) {
	var port int

	err := r.db.QueryRow("SELECT port FROM ports WHERE in_use = 0 ORDER BY port LIMIT 1").Scan(&port)
	if err != nil {
		return 0, fmt.Errorf("no available ports: %w", err)
	}

	_, err = r.db.Exec("UPDATE ports SET in_use = 1,container_id = ? WHERE port = ?", containerID, port)
	if err != nil {
		return 0, fmt.Errorf("failed to mark port as used: %w", err)
	}

	return port, nil
}

func (r *ContainerRepo) ReleasePort(containerID string) error {
	stmt := "UPDATE ports SET in_use = 0, container_id = NULL WHERE container_id = ?"
	r.l.Debug(stmt + " " + containerID)

	_, err := r.db.Exec("UPDATE ports SET in_use = 0, container_id = NULL WHERE container_id = ?", containerID)
	if err != nil {
		return fmt.Errorf("failed to release port: %w", err)
	}

	return nil
}

func (r *ContainerRepo) Update(container *model.Container) (*model.Container, error) {
	_, err := r.db.Exec(`
		UPDATE containers 
		SET dockerid = ?, imagename = ?, containername = ?, gitrepo = ?, user_id = ?, environment = ?, ports = ? 
		WHERE id = ?`,
		container.DockerID, container.ImageName, container.ContainerName, container.GitRepo, container.UserID,
		joinStrings(container.EnvVars), joinStrings(container.Ports), container.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update container: %w", err)
	}

	updatedContainer, err := r.GetByID(container.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated container: %w", err)
	}

	return updatedContainer, nil
}

// Helper function to join string slices into a single string (for storing arrays like Environment, Ports).
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

// Helper function to split a single string into a slice of strings (for retrieving arrays like Environment, Ports).
func splitString(s string) []string {
	if len(s) <= 2 { //nolint:mnd
		return []string{}
	}

	return strings.Split(s[1:len(s)-1], ",")
}

func marshalContainer(container *Container) *model.Container {
	return &model.Container{ //nolint:exhaustruct
		ID:            container.ID,
		DockerID:      container.DockerID,
		ImageName:     container.ImageName,
		ContainerName: container.ContainerName,
		GitRepo:       container.GitRepo,
		UserID:        container.UserID,
		EnvVars:       container.Environment,
		Ports:         container.Ports,
		UIPort:        container.UIPort,
	}
}

func marshalContainers(containers []Container) []model.Container {
	result := []model.Container{}
	for _, container := range containers {
		result = append(result, *marshalContainer(&container))
	}

	return result
}
