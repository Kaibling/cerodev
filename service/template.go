package service

import (
	"fmt"
	"strings"

	"github.com/kaibling/cerodev/model"
	"github.com/kaibling/cerodev/pkg/utils"
)

type templaterepo interface {
	GetByID(id string) (*model.Template, error)
	GetAll() ([]*model.Template, error)
	Create(template *model.Template) (*model.Template, error)
	Delete(id string) error
	Update(template *model.Template) (*model.Template, error)
}

type TemplateService struct {
	dbrepo templaterepo
}

func NewTemplateService(dbrepo templaterepo) *TemplateService {
	return &TemplateService{
		dbrepo: dbrepo,
	}
}

func (s *TemplateService) GetByID(id string) (*model.Template, error) {
	val, err := s.dbrepo.GetByID(id)

	return HandleError[*model.Template](val, err, "failed to GetByID")
}

func (s *TemplateService) GetAll() ([]*model.Template, error) {
	val, err := s.dbrepo.GetAll()

	return HandleError[[]*model.Template](val, err, "failed to GetAll")
}

func (s *TemplateService) Create(template *model.Template) (*model.Template, error) {
	template.ID = utils.GenerateULID()
	template.Dockerfile = baseTemplate
	template.RepoName = strings.ToLower(template.RepoName)
	val, err := s.dbrepo.Create(template)

	return HandleError[*model.Template](val, err, "failed to Create")
}

func (s *TemplateService) Delete(id string) error {
	if err := s.dbrepo.Delete(id); err != nil {
		return fmt.Errorf("failed to db template Delete: %w", err)
	}

	return nil
}

func (s *TemplateService) Update(template *model.Template) (*model.Template, error) {
	val, err := s.dbrepo.Update(template)

	return HandleError[*model.Template](val, err, "failed to Update")
}

const baseTemplate = `
FROM codercom/code-server:latest
ARG ARCHITECTURE

ENV DEBIAN_FRONTEND=noninteractive
# Set Go version — change as needed
ENV GO_VERSION=1.24.2
ENV ARCHITECTURE=${ARCHITECTURE}

# Install git, Go, etc. if needed
USER root
RUN apt-get update && apt-get install -y git make

# Download and install Go
RUN curl -LO https://golang.org/dl/go${GO_VERSION}.linux-${ARCHITECTURE}.tar.gz && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-${ARCHITECTURE}.tar.gz && \
    rm go${GO_VERSION}.linux-${ARCHITECTURE}.tar.gz

# Set Go paths
ENV PATH="/usr/local/go/bin:${PATH}"
ENV GOPATH="/home/coder/go"
ENV PATH="${GOPATH}/bin:${PATH}"

# Fix ownership if running as coder later
RUN mkdir -p /home/coder/go && chown -R coder:coder /home/coder/go

# Ensure extension dir exists
RUN mkdir -p /home/coder/.local/share/code-server/extensions \
    && chown -R coder:coder /home/coder/.local

RUN mkdir -p /home/coder/workspace
RUN chown coder:coder /home/coder/workspace
WORKDIR /home/coder/workspace
# Copy entrypoint script
COPY ./entrypoint.sh /usr/bin/entrypoint.sh
RUN chmod +x /usr/bin/entrypoint.sh
# Switch to coder for installing extensions
USER coder

# Install Go extension as coder
RUN code-server --install-extension golang.go

# Set working directory
WORKDIR /home/coder/workspace
# Expose Code Server port

RUN echo 'export PATH="$PATH:/go/bin:/usr/local/go/bin"' >> ~/.bashrc
ENTRYPOINT ["/bin/bash", "/usr/bin/entrypoint.sh"]
# Start code-server as default CMD
`
