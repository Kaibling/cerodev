package dbrepo

import (
	"context"
	"database/sql"

	"github.com/kaibling/apiforge/log"
	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/repo/sqlcrepo"
)

type TemplateRepo struct {
	ctx      context.Context
	sqlcRepo *sqlcrepo.Queries
	l        log.Writer
}

func NewTemplateRepo(ctx context.Context, db *sql.DB, l log.Writer) *TemplateRepo {
	return &TemplateRepo{ctx: ctx, sqlcRepo: sqlcrepo.New(db), l: l.Named("repo_template")}
}

func (r *TemplateRepo) GetByID(id string) (*model.Template, error) {
	r.l.Info("Fetching template by ID", "id", id)

	template, err := r.sqlcRepo.GetTemplate(r.ctx, id)
	if err != nil {
		r.l.Error("Error fetching template by ID", err)

		return nil, err
	}

	return unmarshalTemplate(template), nil
}

func (r *TemplateRepo) GetAll() ([]*model.Template, error) {
	r.l.Info("Fetching all templates")

	templates, err := r.sqlcRepo.GetAllTemplates(r.ctx)
	if err != nil {
		r.l.Error("Error fetching all templates", err)

		return nil, err
	}

	return unmarshalTemplates(templates), nil
}

func (r *TemplateRepo) Create(template *model.Template) (*model.Template, error) {
	r.l.Info("Creating new template", "template", template)

	createdTemplateID, err := r.sqlcRepo.CreateTemplate(r.ctx, sqlcrepo.CreateTemplateParams{
		ID:         template.ID,
		Name:       template.Name,
		RepoName:   template.RepoName,
		Dockerfile: template.Dockerfile,
	})
	if err != nil {
		r.l.Error("Error creating template", err)

		return nil, err
	}

	return r.GetByID(createdTemplateID)
}

func (r *TemplateRepo) Delete(id string) error {
	r.l.Info("Deleting template by ID", "id", id)

	err := r.sqlcRepo.DeleteTemplate(r.ctx, id)
	if err != nil {
		r.l.Error("Error deleting template by ID", err)

		return err
	}

	return nil
}

func (r *TemplateRepo) Update(template *model.Template) (*model.Template, error) {
	r.l.Info("Updating template", template)

	if err := r.sqlcRepo.UpdateTemplate(r.ctx, sqlcrepo.UpdateTemplateParams{
		Name:       template.Name,
		RepoName:   template.RepoName,
		Dockerfile: template.Dockerfile,
		ID:         template.ID,
	}); err != nil {
		r.l.Error("Error updating template", err)

		return nil, err
	}

	return r.GetByID(template.ID)
}

func unmarshalTemplate(template sqlcrepo.Template) *model.Template {
	return &model.Template{
		ID:         template.ID,
		Name:       template.Name,
		RepoName:   template.RepoName,
		Dockerfile: template.Dockerfile,
	}
}

func unmarshalTemplates(templates []sqlcrepo.Template) []*model.Template {
	result := make([]*model.Template, len(templates))
	for i, template := range templates {
		result[i] = unmarshalTemplate(template)
	}

	return result
}
