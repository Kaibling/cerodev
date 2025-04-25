package docker

type Port struct {
	HostPort      string `json:"host_port"`      //  "8089"
	ContainerPort string `json:"container_port"` //  "8080/tcp"
}

type ImageTemplate struct {
	ImageName  string `json:"image_name"` // "gocode"
	Dockerfile string `json:"dockerfile"` // "Dockerfile"
}

const containerPrefix = "cd"

type ContainerStatus struct {
	ContainerName string `json:"container_name"`
	Status        string `json:"status"` // "running"
	State         string `json:"state"`  // "Up 4 hours"
}

type Container struct {
	ContainerID   string   `json:"container_id"`
	ContainerName string   `json:"container_name"`
	ImageName     string   `json:"image_name"`
	GitRepo       string   `json:"git_repo"`
	UserID        string   `json:"user_id"`
	Environment   []string `json:"environment"` // ["ENV=prod"]
	Ports         []Port   `json:"ports"`
}
