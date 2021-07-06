package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/paramonies/avito-rest-advert/internal/app/mock"
	"github.com/paramonies/avito-rest-advert/internal/app/model"
)

func TestHandler_createAdvert(t *testing.T) {
	type mockBehaviorType func(s *mock.MockService, advert model.Advert)

	tests := []struct {
		name                 string
		inputBody            string
		inputAdvert          model.Advert
		mockBehavior         mockBehaviorType
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name":"name-test", "description":"desc-test", "price":1000, "pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"}`,
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			mockBehavior: func(s *mock.MockService, advert model.Advert) {
				s.EXPECT().CreateAdvert(advert).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Bad input",
			inputBody:            `{"name":"", "description":"desc-test", "price":1000, "pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"}`,
			inputAdvert:          model.Advert{},
			mockBehavior:         func(s *mock.MockService, advert model.Advert) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "Server error",
			inputBody: `{"name":"name-test", "description":"desc-test", "price":1000, "pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"}`,
			inputAdvert: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			mockBehavior: func(s *mock.MockService, advert model.Advert) {
				s.EXPECT().CreateAdvert(advert).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock.NewMockService(c)
			test.mockBehavior(mockService, test.inputAdvert)

			handler := NewHandler(mockService)
			router := gin.New()
			router.POST("/create", handler.createAdvert)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create", bytes.NewBufferString(test.inputBody))
			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getAdvertById(t *testing.T) {
	type mockBehaviorType func(*mock.MockService, int, []string)

	tests := []struct {
		name                 string
		inputURL             string
		inputId              int
		inputFields          []string
		mockBehavior         mockBehaviorType
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			inputURL:    "/get/1",
			inputId:     1,
			inputFields: []string{},
			mockBehavior: func(s *mock.MockService, advertId int, fields []string) {
				s.EXPECT().GetAdvertById(1, []string{}).Return(model.Advert{
					Name:        "name-test",
					Description: "desc-test",
					Price:       1000,
					Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"name":"name-test","description":"desc-test","price":1000,"pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"}`,
		},
		{
			name:        "Ok with fields params",
			inputURL:    "/get/1?fields=description,pictures",
			inputId:     1,
			inputFields: []string{"description", "pictures"},
			mockBehavior: func(s *mock.MockService, advertId int, fields []string) {
				s.EXPECT().GetAdvertById(1, []string{"description", "pictures"}).Return(model.Advert{
					Name:        "name-test",
					Description: "desc-test",
					Price:       1000,
					Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"name":"name-test","description":"desc-test","price":1000,"pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"}`,
		},
		{
			name:                 "Bad input",
			inputURL:             "/get/1a",
			inputId:              1,
			inputFields:          []string{},
			mockBehavior:         func(s *mock.MockService, advertId int, fields []string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"advertisement id must be integer"}`,
		},
		{
			name:        "Server error",
			inputURL:    "/get/1",
			inputId:     1,
			inputFields: []string{},
			mockBehavior: func(s *mock.MockService, advertId int, fields []string) {
				s.EXPECT().GetAdvertById(1, []string{}).Return(model.Advert{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock.NewMockService(c)
			test.mockBehavior(mockService, test.inputId, test.inputFields)

			handler := NewHandler(mockService)
			router := gin.New()
			router.GET("/get/:id", handler.getAdvertById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", test.inputURL, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getList(t *testing.T) {
	type mockBehaviorType func(*mock.MockService, int, string)
	tests := []struct {
		name                 string
		inputURL             string
		inputPage            int
		inputOrderBy         string
		mockBehavior         mockBehaviorType
		expectedResponseCode int
		expectedResponseBody string
	}{
		{
			name:         "Ok",
			inputURL:     "/list",
			inputPage:    1,
			inputOrderBy: "",
			mockBehavior: func(s *mock.MockService, page int, orderBy string) {
				s.EXPECT().GetAdvertList(1, "createdat_desc").Return([]model.Advert{
					{
						Name:     "name-test1",
						Price:    1000,
						Pictures: "avito/files/ad1,avito/files/ad2,avito/files/ad3",
					},
					{
						Name:     "name-test2",
						Price:    100,
						Pictures: "avito/files/ad1,avito/files/ad2,avito/files/ad3",
					},
					{
						Name:     "name-test3",
						Price:    10,
						Pictures: "avito/files/ad1,avito/files/ad2,avito/files/ad3",
					},
				}, nil)
			},
			expectedResponseCode: 200,
			expectedResponseBody: `[{"name":"name-test1","price":1000,"pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"},{"name":"name-test2","price":100,"pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"},{"name":"name-test3","price":10,"pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"}]`,
		},
		{
			name:         "Ok with params",
			inputURL:     "/list?page=1&order_by=price_desc",
			inputPage:    1,
			inputOrderBy: "price_desc",
			mockBehavior: func(s *mock.MockService, page int, orderBy string) {
				s.EXPECT().GetAdvertList(1, "price_desc").Return([]model.Advert{
					{
						Name:     "name-test1",
						Price:    1000,
						Pictures: "avito/files/ad1,avito/files/ad2,avito/files/ad3",
					},
					{
						Name:     "name-test2",
						Price:    100,
						Pictures: "avito/files/ad1,avito/files/ad2,avito/files/ad3",
					},
					{
						Name:     "name-test3",
						Price:    10,
						Pictures: "avito/files/ad1,avito/files/ad2,avito/files/ad3",
					},
				}, nil)
			},
			expectedResponseCode: 200,
			expectedResponseBody: `[{"name":"name-test1","price":1000,"pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"},{"name":"name-test2","price":100,"pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"},{"name":"name-test3","price":10,"pictures":"avito/files/ad1,avito/files/ad2,avito/files/ad3"}]`,
		},
		{
			name:         "Server error",
			inputURL:     "/list",
			inputPage:    1,
			inputOrderBy: "",
			mockBehavior: func(s *mock.MockService, page int, orderBy string) {
				s.EXPECT().GetAdvertList(1, "createdat_desc").Return([]model.Advert{}, errors.New("something went wrong"))
			},
			expectedResponseCode: 500,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
		{
			name:         "Records not found",
			inputURL:     "/list?page=2",
			inputPage:    2,
			inputOrderBy: "createdat_desc",
			mockBehavior: func(s *mock.MockService, page int, orderBy string) {
				s.EXPECT().GetAdvertList(2, "createdat_desc").Return([]model.Advert{}, nil)
			},
			expectedResponseCode: 404,
			expectedResponseBody: `{"error":"advertisements not found"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockService := mock.NewMockService(c)
			test.mockBehavior(mockService, test.inputPage, test.inputOrderBy)

			handler := NewHandler(mockService)
			router := gin.New()
			router.GET("/list", handler.getList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", test.inputURL, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedResponseCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}
