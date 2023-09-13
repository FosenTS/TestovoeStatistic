package config

import (
	"Testovoe1/pkg/gconfig"
	"Testovoe1/pkg/mysync"
	log "github.com/sirupsen/logrus"
	"path"
)

const (
	envLocalFilename = ".env.local"
	envDevFilename   = ".env.dev"
	envProdFilename  = ".env.prod"
)

// EnvConfig Env
type EnvConfig struct {
	ProjectAbsPath string `env:"PROJECT_ABS_PATH" env-required:"true"`
	ConfigPath     string `env:"CONFIG_PATH" env-required:"true"`

	// Postgres
	PostgresHost         string `env:"POSTGRES_HOST" env-required:"true"`
	PostgresPort         string `env:"POSTGRES_PORT" env-required:"true"`
	PostgresUsername     string `env:"POSTGRES_USERNAME" env-required:"true"`
	PostgresPassword     string `env:"POSTGRES_PASSWORD" env-required:"true"`
	PostgresDatabaseName string `env:"POSTGRES_DATABASE_NAME" env-required:"true"`
	PostgresUseCA        bool   `env:"POSTGRES_USE_CA" env-required:"true"`
	PostgresCaPath       string `env:"POSTGRES_CA_PATH" env_required:"true"`
	ConfigAbsPath        string
}

var (
	envCfgInst  = &EnvConfig{}
	loadEnvOnce = mysync.NewOnce()
)

func LoadEnv() {
	envFilename := ".env"
	loadEnvOnce.Do(func() {
		err := gconfig.ReadEnv(envFilename, envCfgInst)
		if err != nil {
			log.Fatalf("fatal reading env: %s\n", err)
		}

		envCfgInst.ConfigAbsPath = path.Join(envCfgInst.ProjectAbsPath, envCfgInst.ConfigPath)

		log.Infoln("Env successfully read")
	})
}

func Env() EnvConfig {
	if !loadEnvOnce.IsDone() {
		log.Fatalln("env must be loaded before use!")
	}
	return *envCfgInst
}
