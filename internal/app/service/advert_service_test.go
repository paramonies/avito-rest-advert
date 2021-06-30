package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/paramonies/avito-rest-advert/internal/app/mock"
	"github.com/paramonies/avito-rest-advert/internal/app/model"
)

func TestService_CreateAdvert(t *testing.T) {
	type mockBehaviortype func(*mock.MockRepository, model.Advert)
	tests := []struct {
		name           string
		inputAdvert    model.Advert
		mockBehavior   mockBehaviortype
		expectedResult int
		expectedError  error
	}{
		{
			name: "OK",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			mockBehavior: func(r *mock.MockRepository, advert model.Advert) {
				r.EXPECT().CreateAdvert(advert).Return(1, nil)
			},
			expectedResult: 1,
			expectedError:  nil,
		},
		{
			name: "Input model.Advert with invalid Name field",
			inputAdvert: model.Advert{
				Name:        strings.Repeat("t", 201),
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			mockBehavior: func(r *mock.MockRepository, advert model.Advert) {
			},
			expectedResult: 0,
			expectedError:  errors.New(`length of the field "name" should not exceed 200`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockRepository := mock.NewMockRepository(c)
			test.mockBehavior(mockRepository, test.inputAdvert)

			service := NewAdvertService(mockRepository)

			resultId, resultError := service.CreateAdvert(test.inputAdvert)

			// t.Logf("!!! %s\nexpected: %v %v\ngot: %v %v", test.name, test.expectedResult, test.expectedError, resultId, resultError)

			assert.Equal(t, resultError, test.expectedError)
			assert.Equal(t, resultId, test.expectedResult)

		})
	}
}

func TestService_validate(t *testing.T) {
	tests := []struct {
		name           string
		inputAdvert    model.Advert
		expectedResult error
	}{
		{
			name: "Ok",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			expectedResult: nil,
		},
		{
			name: "Advert.Name error",
			inputAdvert: model.Advert{
				Name:        strings.Repeat("t", 201),
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			expectedResult: errors.New(`length of the field "name" should not exceed 200`),
		},
		{
			name: "Advert.Description error",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: strings.Repeat("t", 1001),
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			expectedResult: errors.New(`length of the field "description" should not exceed 1000`),
		},
		{
			name: "Advert.Price error",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       -1,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			expectedResult: errors.New(`the field "price" must have a value greater than 0`),
		},
		{
			name: "Advert.Pictures error",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3,avito/files/ad4",
			},
			expectedResult: errors.New(`the field "pictures" must contain no more than 3 photos`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := validate(test.inputAdvert)
			t.Logf("!!! %s expected: %v\ngot: %v", test.name, test.expectedResult, result)
			assert.Equal(t, result, test.expectedResult)
		})
	}
}

func TestService_checkFields(t *testing.T) {
	tests := []struct {
		name           string
		inputAdvert    model.Advert
		inputFields    []string
		expectedResult model.Advert
	}{
		{
			name: "Fields empty",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			inputFields: []string{},
			expectedResult: model.Advert{
				Name:        "name-test",
				Description: "",
				Price:       1000,
				Pictures:    "",
				MainPicture: "avito/files/ad1",
			},
		},
		{
			name: "Fields empty with Advert.Pictures empty ",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "",
			},
			inputFields: []string{},
			expectedResult: model.Advert{
				Name:        "name-test",
				Description: "",
				Price:       1000,
				Pictures:    "",
				MainPicture: "",
			},
		},
		{
			name: "Fields contains description and pictures",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			inputFields: []string{"description", "pictures"},
			expectedResult: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
				MainPicture: "avito/files/ad1",
			},
		},
		{
			name: "Fields contains pictures",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			inputFields: []string{"pictures"},
			expectedResult: model.Advert{
				Name:        "name-test",
				Description: "",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
				MainPicture: "avito/files/ad1",
			},
		},
		{
			name: "Fields contains description",
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			inputFields: []string{"description"},
			expectedResult: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "",
				MainPicture: "avito/files/ad1",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := checkFields(test.inputAdvert, test.inputFields)
			// t.Logf("!!! %s expected: %v\ngot: %v", test.name, test.expectedResult, result)
			assert.Equal(t, result, test.expectedResult)
		})
	}
}
