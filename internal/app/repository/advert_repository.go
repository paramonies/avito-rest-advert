package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/paramonies/avito-rest-advert/internal/app/model"
)

const (
	ADVERTSTABLE = "adverts"
)

type AdvertRepository struct {
	DB *sqlx.DB
}

func NewAdvertRepository(db *sqlx.DB) *AdvertRepository {
	return &AdvertRepository{DB: db}
}

func (r *AdvertRepository) CreateAdvert(advert model.Advert) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, description, price, pictures) VALUES ($1, $2, $3, $4) RETURNING id", ADVERTSTABLE)
	row := r.DB.QueryRow(query, advert.Name, advert.Description, advert.Price, advert.Pictures)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AdvertRepository) GetAdvertById(advertId int) (model.Advert, error) {
	query := fmt.Sprintf("SELECT name, description, price, pictures FROM %s WHERE id = $1", ADVERTSTABLE)
	row := r.DB.QueryRow(query, advertId)
	var advert model.Advert
	if err := row.Scan(&advert.Name, &advert.Description, &advert.Price, &advert.Pictures); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return advert, errors.New("advertisement not found")
		default:
			return advert, err
		}
	}
	return advert, nil

}

func (r *AdvertRepository) GetAdvertList(page int, orderField string, orderDirect string) ([]model.Advert, error) {
	var adverts []model.Advert
	query := fmt.Sprintf("SELECT name, price, pictures FROM %s  ORDER BY %s %s  LIMIT 10 OFFSET ($1-1)*10", ADVERTSTABLE, orderField, orderDirect)
	if err := r.DB.Select(&adverts, query, page); err != nil {
		return nil, err
	}
	return adverts, nil
}
