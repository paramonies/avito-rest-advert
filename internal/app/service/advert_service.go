package service

import (
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/paramonies/avito-rest-advert/internal/app/model"
	"github.com/paramonies/avito-rest-advert/internal/app/repository"
)

type AdvertService struct {
	repo repository.Repository
}

func NewAdvertService(repo repository.Repository) *AdvertService {
	return &AdvertService{repo: repo}
}

func (s *AdvertService) CreateAdvert(advert model.Advert) (int, error) {
	if err := validate(advert); err != nil {
		return 0, err
	}

	return s.repo.CreateAdvert(advert)
}

func (s *AdvertService) GetAdvertById(advertId int, fields []string) (model.Advert, error) {
	advert, err := s.repo.GetAdvertById(advertId)
	if err != nil {
		return advert, err
	}

	return checkFields(advert, fields), nil
}

func (s *AdvertService) GetAdvertList(page int, orderBy string) ([]model.Advert, error) {
	order := strings.Split(orderBy, "_")
	orderField, orderDirect := order[0], order[1]
	return s.repo.GetAdvertList(page, orderField, orderDirect)
}

func validate(advert model.Advert) error {
	var messageErrors []string
	if utf8.RuneCountInString(advert.Name) > 200 {
		messageErrors = append(messageErrors, `length of the field "name" should not exceed 200`)
	}

	if utf8.RuneCountInString(advert.Description) > 1000 {
		messageErrors = append(messageErrors, `length of the field "description" should not exceed 1000`)
	}

	if advert.Price < 0 {
		messageErrors = append(messageErrors, `the field "price" must have a value greater than 0`)
	}

	if advert.Pictures != "" {
		pictures := strings.Split(advert.Pictures, ",")
		if len(pictures) > 3 {
			messageErrors = append(messageErrors, `the field "pictures" must contain no more than 3 photos`)
		}
	}

	if len(messageErrors) > 0 {
		err := strings.Join(messageErrors, ", ")
		return errors.New(err)
	}

	return nil
}

func checkFields(advert model.Advert, fields []string) model.Advert {
	if advert.Pictures != "" {
		advert.MainPicture = strings.Split(advert.Pictures, ",")[0]
	} else {
		advert.MainPicture = ""
	}

	if len(fields) == 0 {
		advert.Description = ""
		advert.Pictures = ""
		return advert
	}

	if contains(fields, "description") && contains(fields, "pictures") {
		return advert
	}

	if contains(fields, "pictures") {
		advert.Description = ""
		return advert
	}

	if contains(fields, "description") {
		advert.Pictures = ""
		return advert
	}

	return advert
}

func contains(s []string, el string) bool {
	for _, v := range s {
		if v == el {
			return true
		}
	}
	return false
}
