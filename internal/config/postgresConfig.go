package config

import (
	"path"
	"sync"
	"time"
)

const postgresConfigFilename = "postgres.config.yaml"

type PostgresConfig struct {
	MaxConnections          int  `yaml:"maxConnections" env-required:"true"`
	MaxConnectionAttempts   int  `yaml:"maxConnectionAttempts" env-required:"true"`
	WaitTimeoutMilliseconds int  `yaml:"waitTimeoutMilliseconds" env-required:"true"`
	SimpleQueryMode         bool `yaml:"simpleQueryMode" env-required:"true"`

	Host         string
	Port         string
	Username     string
	Password     string
	DatabaseName string
	UseCA        bool
	CAAbsPath    string
	WaitTimeout  time.Duration
}

var (
	postgresCfgInst  = &PostgresConfig{}
	loadPostgresOnce = sync.Once{}
)

func Postgres() PostgresConfig {
	loadPostgresOnce.Do(func() {
		env := Env()
		readConfig(path.Join(env.ConfigAbsPath, postgresConfigFilename), postgresCfgInst)

		postgresCfgInst.Host = env.PostgresHost
		postgresCfgInst.Port = env.PostgresPort
		postgresCfgInst.Username = env.PostgresUsername
		postgresCfgInst.Password = env.PostgresPassword
		postgresCfgInst.DatabaseName = env.PostgresDatabaseName
		postgresCfgInst.CAAbsPath = path.Join(env.ProjectAbsPath, env.PostgresCaPath)
		postgresCfgInst.UseCA = env.PostgresUseCA
		postgresCfgInst.WaitTimeout = time.Millisecond * time.Duration(postgresCfgInst.WaitTimeoutMilliseconds)
	})

	return *postgresCfgInst
}
