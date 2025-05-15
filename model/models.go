package model

type Container struct {
	ID            string   `json:"id"`
	DockerID      string   `json:"docker_id"`
	ImageName     string   `json:"image_name"`
	Status        string   `json:"status"` // "running"
	State         string   `json:"state"`  // "Up 4 hours"
	ContainerName string   `json:"container_name"`
	GitRepo       string   `json:"git_repo"`
	UserID        string   `json:"user_id"`
	EnvVars       []string `json:"env_vars"` // ["ENV=prod"]
	Ports         []string `json:"ports"`    // ["8080:8098/tcp"]
	UIPort        string   `json:"ui_port"`  // "32102"
}

type ContainerStatus struct {
	DockerID string `json:"docker_id"`
	Status   string `json:"status"` // "running"
	State    string `json:"state"`  // "Up 4 hours"
}

type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Password string   `json:"password,omitempty"`
	Tokens   []string `json:"tokens"`
}

type Token struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

type Template struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	RepoName   string `json:"repo_name"`
	Dockerfile string `json:"dockerfile"`
}

type Image struct {
	RepoName string `json:"repo_name"` // "gocode"
	ImageID  string `json:"image_id"`  // "sha256:abc123"
	Tag      string `json:"tag"`       // "latest"
}

type BuildParams struct {
	RepoName   string             `json:"repo_name"` // "gocode"
	TemplateID string             `json:"template_id"`
	Tag        string             `json:"tag"`        // "latest"
	BuildArgs  map[string]*string `json:"build_args"` // ["ENV=prod"]
}

func (bp *BuildParams) Validate() {
	if bp.BuildArgs == nil {
		bp.BuildArgs = map[string]*string{}
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Token    string `json:"token"`
}

type WebSocketMessage struct {
	Timestamp   string `json:"timestamp"`
	MessageType string `json:"message_type"`
	Message     string `json:"message"`
}
