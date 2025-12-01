package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	require.NotNil(t, cfg)
	assert.Equal(t, DatabasePostgres, cfg.Database, "default database should be postgres")
	assert.False(t, cfg.IncludeGRPC, "gRPC should be disabled by default")
	assert.False(t, cfg.IncludeAuth, "auth should be disabled by default")
	assert.False(t, cfg.IncludeAI, "AI should be disabled by default")
	assert.True(t, cfg.IncludeDocker, "docker should be enabled by default")
	assert.True(t, cfg.IncludeMake, "make should be enabled by default")
	assert.Empty(t, cfg.Name, "name should be empty by default")
	assert.Empty(t, cfg.ModulePath, "module path should be empty by default")
	assert.Empty(t, cfg.OutputDir, "output dir should be empty by default")
}

func TestProjectConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *ProjectConfig
		wantErr error
	}{
		{
			name: "valid config",
			config: &ProjectConfig{
				Name:       "my-app",
				ModulePath: "github.com/user/my-app",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			config: &ProjectConfig{
				Name:       "",
				ModulePath: "github.com/user/my-app",
			},
			wantErr: ErrNameRequired,
		},
		{
			name: "empty module path",
			config: &ProjectConfig{
				Name:       "my-app",
				ModulePath: "",
			},
			wantErr: ErrModulePathRequired,
		},
		{
			name: "both empty",
			config: &ProjectConfig{
				Name:       "",
				ModulePath: "",
			},
			wantErr: ErrNameRequired, // name checked first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDatabaseOptions(t *testing.T) {
	options := DatabaseOptions()

	require.Len(t, options, 3, "should return 3 database options")
	assert.Contains(t, options, DatabasePostgres)
	assert.Contains(t, options, DatabaseSQLite)
	assert.Contains(t, options, DatabaseNone)
}

func TestDatabaseConstants(t *testing.T) {
	assert.Equal(t, Database("postgres"), DatabasePostgres)
	assert.Equal(t, Database("sqlite"), DatabaseSQLite)
	assert.Equal(t, Database("none"), DatabaseNone)
}
