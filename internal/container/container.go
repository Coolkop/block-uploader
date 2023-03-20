package container

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"

	"file-storage/internal/clients/minio"
	"file-storage/internal/config"
	getHandler "file-storage/internal/handlers/get"
	putHandler "file-storage/internal/handlers/put"
	"file-storage/internal/repository"
	getProcessor "file-storage/internal/services/get"
	putProcessor "file-storage/internal/services/put"
	"file-storage/internal/services/storage"
)

type Container struct {
	Config        config.App
	WeightsLoader *storage.WeightLoader
	HttpHandler   http.Handler
}

func New() (*Container, error) {
	var cfg config.App
	err := envconfig.Process("FILE_STORAGE", &cfg)
	if err != nil {
		return nil, err
	}

	db, err := initDbWithRetries(cfg)
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", cfg.MigrationsPath), "postgres", driver)
	if err != nil {
		return nil, err
	}

	m.Up()

	externalClients, err := minio.NewExternalClients(cfg.Minio) // TODO add health checks
	if err != nil {
		return nil, err
	}
	client := minio.New(externalClients)

	weightsCache := storage.NewCache()
	weightsLoader := storage.NewWeightLoader(len(cfg.Minio.Hosts), weightsCache)
	err = weightsLoader.Load()
	if err != nil {
		return nil, err
	}

	repo := repository.New(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{fileID}", getHandler.New(getProcessor.New(client, repo)).Handle)
	r.Put(
		"/",
		putHandler.New(putProcessor.New(
			client,
			weightsCache,
			repo,
			putProcessor.NewChunker(),
		)).Handle,
	)

	return &Container{
		Config:        cfg,
		WeightsLoader: weightsLoader,
		HttpHandler:   r,
	}, nil
}

func initDbWithRetries(cfg config.App) (*sqlx.DB, error) {
	retries := 3
	attempt := 0

	for {
		db, err := initDB(cfg)
		if err != nil {
			if attempt == retries {
				return nil, err
			}

			log.Println(err)

			attempt++
			time.Sleep(10 * time.Second)

			continue
		}

		return db, nil
	}
}

func initDB(cfg config.App) (*sqlx.DB, error) {
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=public sslmode=disable",
			cfg.Pg.Host,
			cfg.Pg.Port,
			cfg.Pg.User,
			cfg.Pg.Password,
		),
	)
	if err != nil {
		return nil, err
	}
	if db.Ping() != nil { // TODO move to health check
		return nil, fmt.Errorf("db not responding")
	}

	return db, nil
}
