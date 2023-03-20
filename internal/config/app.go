package config

import "file-storage/internal/clients/minio"

type App struct {
	MigrationsPath string `required:"true"`

	Server Server
	Pg     Pg
	Minio  minio.Config
}

type Server struct {
	Port int `required:"true"`
}

type Pg struct {
	Port     int    `required:"true"`
	Host     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true"`
}
