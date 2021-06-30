package service

import "github.com/paramonies/avito-rest-advert/internal/app/model"

type Service interface {
	CreateAdvert(model.Advert) (int, error)
	GetAdvertById(int, []string) (model.Advert, error)
	GetAdvertList(int, string) ([]model.Advert, error)
}
