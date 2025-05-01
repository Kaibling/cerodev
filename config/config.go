package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	osPrefix                  = "CD"
	AppName                   = "cerodev"
	defaultPasswordCost       = 11
	defaultTokenLength        = 32
	defaultContainerPortRange = "30000-40000"
	defaultVolumesPath        = "/var/lib/cerodev/volumes"
)

var (
	Version      string //nolint:gochecknoglobals
	BuildTime    string //nolint:gochecknoglobals
	Architecture string //nolint:gochecknoglobals
)

type Configuration struct {
	APIBindingIP      string
	APIBindingPort    string
	APIEnableTLS      bool
	APITLSCertPath    string
	APITLSCertKeyPath string
	AdminUser         string
	AdminPassword     string
	AdminToken        string
	TokenLength       int
	PasswordCost      int
	ContainerMinPort  int
	ContainerMaxPort  int
	DBConfig          DBConfiguration
	VolumesPath       string
	PublicURL         string
}
type DBConfiguration struct {
	FilePath string
}

func Load() Configuration {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file") //nolint: forbidigo
	}

	minPort, maxPort := splitPortRange(getEnv("CONTAINER_PORT_RANGE", defaultContainerPortRange))
	if minPort == 0 || maxPort == 0 {
		minPort, maxPort = splitPortRange(defaultContainerPortRange)
	}

	return Configuration{
		APIBindingIP:      getEnv("API_BINDING_IP", "0.0.0.0"),
		APIBindingPort:    getEnv("API_BINDING_PORT", "8088"),
		APIEnableTLS:      toBool(getEnv("API_ENABLE_TLS", "true")),
		APITLSCertPath:    getEnv("API_TLS_CERT_PATH", "./server.crt"),
		APITLSCertKeyPath: getEnv("API_TLS_CERT_KEY_PATH", "./server.key"),
		AdminUser:         getEnv("ADMIN_USER", "admin"),
		AdminPassword:     getEnv("ADMIN_PASSWORD", "abc123"),
		AdminToken:        getEnv("ADMIN_TOKEN", ""),
		TokenLength:       getEnvAsInt("TOKEN_LENGTH", defaultTokenLength),
		PasswordCost:      defaultPasswordCost,
		ContainerMinPort:  minPort,
		ContainerMaxPort:  maxPort,
		VolumesPath:       getEnv("VOLUMES_PATH", defaultVolumesPath),
		DBConfig: DBConfiguration{
			FilePath: getEnv("DB_FILE_PATH", "cerodev.db"),
		},
		PublicURL: getEnv("PUBLIC_URL", "http://localhost"),
	}
}

func getEnv(key string, defaultValue string) string {
	fullKey := osPrefix + "_" + key

	val := os.Getenv(fullKey)
	if val == "" {
		if defaultValue != "" {
			return defaultValue
		}
	}

	return val
}

func getEnvAsInt(key string, defaultValue int) int {
	fullKey := osPrefix + "_" + key

	val := os.Getenv(fullKey)
	if val == "" {
		return defaultValue
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return intVal
}

func toBool(s string) bool {
	return strings.ToLower(s) == "true"
}

func splitPortRange(portRange string) (int, int) {
	parts := strings.Split(portRange, "-")
	if len(parts) != 2 { //nolint: mnd
		return 0, 0
	}

	minPort, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0
	}

	maxPort, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0
	}

	return minPort, maxPort
}
