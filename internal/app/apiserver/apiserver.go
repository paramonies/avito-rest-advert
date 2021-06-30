package apiserver

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/paramonies/avito-rest-advert/internal/app/handler"
	"github.com/paramonies/avito-rest-advert/internal/app/repository"
	"github.com/paramonies/avito-rest-advert/internal/app/service"
)

func Start(config *Config) error {

	db, err := newDB(config)
	if err != nil {
		return err
	}
	defer db.Close()

	repo := repository.NewAdvertRepository(db)
	service := service.NewAdvertService(repo)
	handler := handler.NewHandler(service)

	srv := new(Server)
	if err := srv.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

	return nil
}

func newDB(config *Config) (*sqlx.DB, error) {
	dbURL := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
