package models

// Config struct
type Config struct {
	Port     int               `yaml:"port"`
	Storage  *Storage          `yaml:"storage"`
	Sessions map[string]string `yaml:"sessions"`
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
)
