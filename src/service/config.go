package service

import (
	"github.com/africarealty/server/src/kit/auth"
	kitConfig "github.com/africarealty/server/src/kit/config"
	kitHttp "github.com/africarealty/server/src/kit/http"
	"github.com/africarealty/server/src/kit/log"
	kitAero "github.com/africarealty/server/src/kit/storages/aerospike"
	"github.com/africarealty/server/src/kit/storages/pg"
	"os"
	"path/filepath"
)

// Here are defined all types for your configuration
// You can remove not needed types or add your own

type Storages struct {
	Aero *kitAero.Config
	Pg   *pg.DbClusterConfig
}

type Address struct {
	Host string
	Port string
}

type Config struct {
	Log      *log.Config
	Http     *kitHttp.Config
	Storages *Storages
	Auth     *auth.Config
}

func LoadConfig() (*Config, error) {

	// get root folder from env
	rootPath := os.Getenv("ARROOT")
	if rootPath == "" {
		return nil, kitConfig.ErrEnvRootPathNotSet("ARROOT")
	}

	// config path
	configPath := filepath.Join(rootPath, "africarealty", "config.yml")

	// .env path
	envPath := filepath.Join(rootPath, "africarealty", ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		envPath = ""
	}

	// load config
	config := &Config{}
	err := kitConfig.NewConfigLoader(LF()).
		WithConfigPath(configPath).
		WithEnvPath(envPath).
		Load(config)

	if err != nil {
		return nil, err
	}
	return config, nil
}
