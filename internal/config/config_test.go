// file: internal/config/config_test.go
package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"credit-service/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig_FromEnv_Success(t *testing.T) {
	t.Setenv("APP_NAME", "env-service")
	t.Setenv("APP_VERSION", "2.0.0")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("CONFIG_PATH", "")

	cfg, err := config.NewConfig()
	assert.NoError(t, err)
	assert.Equal(t, "env-service", cfg.App.Name)
	assert.Equal(t, "2.0.0", cfg.App.Version)
	assert.Equal(t, "9090", cfg.HTTP.Port)
}

func TestNewConfig_MissingEnv_Error(t *testing.T) {
	t.Setenv("CONFIG_PATH", filepath.Join(t.TempDir(), "noop.yml"))
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_VERSION")
	os.Unsetenv("HTTP_PORT")

	_, err := config.NewConfig()
	assert.Error(t, err)

	assert.Contains(t, err.Error(), "configuration error")
}

func TestNewConfig_InvalidFile_Error(t *testing.T) {

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "bad_config.yml")
	if err := os.WriteFile(filePath, []byte("not: valid: yaml: :::"), 0o644); err != nil {
		t.Fatalf("unable to write bad config file: %v", err)
	}

	t.Setenv("CONFIG_PATH", filePath)
	t.Setenv("APP_NAME", "irrelevant")
	t.Setenv("APP_VERSION", "1.0.0")
	t.Setenv("HTTP_PORT", ":8080")

	_, err := config.NewConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}
