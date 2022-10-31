package service

import (
	"github.com/africarealty/server/src/kit/auth"
	kitConfig "github.com/africarealty/server/src/kit/config"
	kitHttp "github.com/africarealty/server/src/kit/http"
	"github.com/africarealty/server/src/kit/log"
	"github.com/africarealty/server/src/kit/queue"
	kitAero "github.com/africarealty/server/src/kit/storages/aerospike"
	"github.com/africarealty/server/src/kit/storages/pg"
	"os"
	"path/filepath"
)

// Here are defined all types for your configuration
// You can remove not needed types or add your own

type CfgStorages struct {
	Aero *kitAero.Config
	Pg   *pg.DbClusterConfig
}

type CfgAddress struct {
	Host string
	Port string
}

type CfgEmail struct {
	SmtpServer     string `config:"smtp-server"`
	SmtpServerPort string `config:"smtp-port"`
	SmtpUser       string `config:"smtp-user"`
	SmtpPassword   string `config:"smtp-password"`
	SmtpFrom       string `config:"smtp-from"`
}

type CfgCommunications struct {
	Email *CfgEmail
}

type CfgAuth struct {
	Session  *auth.Config
	Password struct {
		MinLen uint `config:"min-len"`
	}
	Activation struct {
		Url string
		Ttl uint32
	}
}

type CfgSdk struct {
	Url string
	Log *log.Config
}

type CfgTests struct {
	User     string
	Password string
}

type Config struct {
	Log            *log.Config
	Http           *kitHttp.Config
	Storages       *CfgStorages
	Auth           *CfgAuth
	Communications *CfgCommunications
	Sdk            *CfgSdk
	Tests          *CfgTests
	Nats           *queue.Config
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
