package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig_ValidEnvs(t *testing.T) {
	err := os.Setenv("FOXYLOCK_REDIS_PASS", "expected-password")
	assert.NoError(t, err)
	defer func() {
		err := os.Unsetenv("FOXYLOCK_REDIS_PASS")
		assert.NoError(t, err)
	}()

	c, err := LoadConfig()

	assert.NoError(t, err)
	assert.Equal(t, *c, Config{":6379", "expected-password", ""})
}

func TestLoadConfig_EmptyEnvs(t *testing.T) {
	_, err := LoadConfig()

	assert.Error(t, err)
}
