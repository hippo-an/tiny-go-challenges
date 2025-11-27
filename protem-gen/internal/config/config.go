package config

// Database represents the database choice
type Database string

const (
	DatabasePostgres Database = "postgres"
	DatabaseMySQL    Database = "mysql"
	DatabaseSQLite   Database = "sqlite"
	DatabaseNone     Database = "none"
)

// ProjectConfig holds all configuration for project generation
type ProjectConfig struct {
	// Basic project info
	Name       string `json:"name"`
	ModulePath string `json:"module_path"`
	OutputDir  string `json:"output_dir"`

	// Database configuration
	Database Database `json:"database"`

	// Optional features
	IncludeGRPC bool `json:"include_grpc"`
	IncludeAuth bool `json:"include_auth"`
	IncludeAI   bool `json:"include_ai"`

	// Development options
	IncludeDocker bool `json:"include_docker"`
	IncludeMake   bool `json:"include_make"`
}

// NewDefaultConfig returns a ProjectConfig with sensible defaults
func NewDefaultConfig() *ProjectConfig {
	return &ProjectConfig{
		Database:      DatabasePostgres,
		IncludeGRPC:   false,
		IncludeAuth:   false,
		IncludeAI:     false,
		IncludeDocker: true,
		IncludeMake:   true,
	}
}

// Validate checks if the configuration is valid
func (c *ProjectConfig) Validate() error {
	if c.Name == "" {
		return ErrNameRequired
	}
	if c.ModulePath == "" {
		return ErrModulePathRequired
	}
	return nil
}

// DatabaseOptions returns available database options
func DatabaseOptions() []Database {
	return []Database{
		DatabasePostgres,
		DatabaseMySQL,
		DatabaseSQLite,
		DatabaseNone,
	}
}
