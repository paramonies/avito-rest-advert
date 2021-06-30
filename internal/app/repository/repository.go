package repository

import "github.com/paramonies/avito-rest-advert/internal/app/model"

type Repository interface {
	CreateAdvert(model.Advert) (int, error)
	GetAdvertById(int) (model.Advert, error)
	GetAdvertList(int, string, string) ([]model.Advert, error)
}
