package models

import (
	"net/http"
	"os"
	"regexp"

	uuid "github.com/satori/go.uuid"
)

// Config struct
type Config struct {
	Port          int               `yaml:"port"`
	Storage       *Storage          `yaml:"storage"`
	Sessions      map[string]string `yaml:"sessions"`
	MultiUserMode bool              `yaml:"multiUser"`
	RootPath      string            `yaml:"rootPath"`
}

// Storage struct
type Storage struct {
	Driver     string `yaml:"driver"`
	Connection string `yaml:"connection"`
	Name       string `yaml:"name"`
}

type contextKey int

const (
	//ContextLoggerKey key
	ContextLoggerKey contextKey = iota
	//ContextReqIDKey key
	ContextReqIDKey contextKey = iota
	//ContextAuthMockKey key
	ContextAuthMockKey contextKey = iota
	//XRequestID key
	XRequestID = "X-Request-Id"
	//EnvPort key
	EnvPort = "PORT"
	//EnvSessKey key
	EnvSessKey = "SESSIONS_KEY"
	//EnvDbName key
	EnvDbName = "DB_NAME"
	//EnvDbConnection key
	EnvDbConnection = "DB_CONNECTION"
	//EnvConfigPath key
	EnvConfigPath = "APP_CONFIG_PATH"
)

const (
	// CtxAPIVersion key
	CtxAPIVersion contextKey = iota
	// CtxValueOwner key
	CtxValueOwner
	// CtxValueEnvID key
	CtxValueEnvID
	// CtxValueAuth key
	CtxValueAuth
)

// OwnerFromContext returns context value for project owner
func OwnerFromContext(r *http.Request) string {
	owner := r.Context().Value(CtxValueOwner)
	return owner.(string)
}

// IsCodeValid returns true if code has valid value
func IsCodeValid(code string) bool {
	res, _ := regexp.Match(`^([a-z0-9\-]+)$`, []byte(code))
	return res
}

// SetConfigDefaults sets defaults for config env variables
func SetConfigDefaults() {
	if os.Getenv(EnvConfigPath) == "" {
		os.Setenv(EnvConfigPath, "./configs/application.yml")
	}
	if os.Getenv(EnvPort) == "" {
		os.Setenv(EnvPort, "8080")
	}
	if os.Getenv(EnvSessKey) == "" {
		os.Setenv(EnvSessKey, uuid.Must(uuid.NewV4(), nil).String())
	}
	if os.Getenv(EnvDbName) == "" {
		os.Setenv(EnvDbName, "toggly")
	}
	if os.Getenv(EnvDbConnection) == "" {
		os.Setenv(EnvDbConnection, "mongodb://localhost")
	}
}
