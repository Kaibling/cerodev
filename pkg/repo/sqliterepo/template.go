package sqliterepo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/model"
)

type Template struct {
	ID         string
	Name       string
	ImageName  string
	Dockerfile string
}

type TemplateRepository struct {
	db *sql.DB
	l  log.Writer
}

func NewTemplateRepo(db *sql.DB, l log.Writer) *TemplateRepository {
	return &TemplateRepository{db: db, l: l.Named("repo_template")}
}

func (r *TemplateRepository) Create(template *model.Template) error {
	_, err := r.db.Exec("INSERT INTO templates (id, name, image_name, dockerfile) VALUES (?, ?, ?, ?)", template.ID, template.Name, template.RepoName, template.Dockerfile) //nolint: lll
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}

func (r *TemplateRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM templates WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	return nil
}

func (r *TemplateRepository) Update(template *model.Template) (*model.Template, error) {
	_, err := r.db.Exec("UPDATE templates SET name = ?, image_name = ?, dockerfile = ? WHERE id = ?", template.Name, template.RepoName, template.Dockerfile, template.ID) //nolint: lll
	if err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	return r.GetByID(template.ID)
}

func (r *TemplateRepository) GetByID(id string) (*model.Template, error) {
	row := r.db.QueryRow("SELECT id, name, image_name, dockerfile FROM templates WHERE id = ?", id)

	var template Template

	err := row.Scan(&template.ID, &template.Name, &template.ImageName, &template.Dockerfile)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no template found with ID %s: %w", id, err)
		}

		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return marshalTemplate(&template), nil
}

func (r *TemplateRepository) GetAll() ([]*model.Template, error) {
	rows, err := r.db.Query("SELECT id, name, image_name, dockerfile FROM templates")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch templates: %w", err)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to fetch templates: %w", rows.Err())
	}

	defer func() {
		if err := rows.Close(); err != nil {
			r.l.Error("failed to close rows", err)
		}
	}()

	var templates []Template

	for rows.Next() {
		var template Template
		if err := rows.Scan(&template.ID, &template.Name, &template.ImageName, &template.Dockerfile); err != nil {
			return nil, fmt.Errorf("failed to scan template: %w", err)
		}

		templates = append(templates, template)
	}

	return marshalTemplates(templates), nil
}

func marshalTemplate(template *Template) *model.Template {
	return &model.Template{
		ID:         template.ID,
		Name:       template.Name,
		RepoName:   template.ImageName,
		Dockerfile: template.Dockerfile,
	}
}

func marshalTemplates(templates []Template) []*model.Template {
	marshaledTemplates := make([]*model.Template, len(templates))
	for i, template := range templates {
		marshaledTemplates[i] = marshalTemplate(&template)
	}

	return marshaledTemplates
}
