package toml

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCfg struct {
	Path string `toml:"path"`
}

func TestLoadConfig(t *testing.T) {
	cfgFile := filepath.Join(t.TempDir(), "cfg")
	assert.NotNil(t, LoadConfig(cfgFile, cfgFile, &TestCfg{}))

	f, err := os.Create(cfgFile)
	assert.NoError(t, err)
	assert.NotNil(t, f)
	_, _ = f.WriteString("Hello World")
	assert.NotNil(t, LoadConfig(cfgFile, cfgFile, &TestCfg{}))
	_ = f.Close()

	_ = EncodeToml(cfgFile, &TestCfg{Path: "/data/path"})
	cfg := TestCfg{}
	err = LoadConfig(cfgFile, cfgFile, &cfg)
	assert.NoError(t, err)
	assert.Equal(t, TestCfg{Path: "/data/path"}, cfg)

	err = LoadConfig("", cfgFile, &cfg)
	assert.NoError(t, err)
	assert.Equal(t, TestCfg{Path: "/data/path"}, cfg)
}
