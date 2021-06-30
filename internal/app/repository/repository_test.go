package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/paramonies/avito-rest-advert/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateAdvert(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' while opening a stub database connection", err)
	}

	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	r := NewAdvertRepository(db)

	type args struct {
		advert model.Advert
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO adverts").
					WithArgs("name-test", "desc-test", 1000, "avito/files/ad1,avito/files/ad2,avito/files/ad3").WillReturnRows(rows)
			},
			input: args{
				advert: model.Advert{
					Name:        "name-test",
					Description: "desc-test",
					Price:       1000,
					Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
				},
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Empty Fields",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"})
				mock.ExpectQuery("INSERT INTO adverts").
					WithArgs("", "desc-test", 1000, "avito/files/ad1,avito/files/ad2,avito/files/ad3").WillReturnRows(rows)
			},
			input: args{
				advert: model.Advert{
					Name:        "",
					Description: "desc-test",
					Price:       1000,
					Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
				},
			},
			// want:    1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.CreateAdvert(test.input.advert)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAdvertById(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' while opening a stub database connection", err)
	}

	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	r := NewAdvertRepository(db)

	type args struct {
		advertId int
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    model.Advert
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name", "price", "description", "pictures"}).
					AddRow("name-test", "desc-test", 1000, "avito/files/ad1,avito/files/ad2,avito/files/ad3")

				mock.ExpectQuery("SELECT name, description, price, pictures FROM adverts WHERE (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				advertId: 1,
			},
			want: model.Advert{
				Name:        "name-test",
				Description: "desc-test",
				Price:       1000,
				Pictures:    "avito/files/ad1,avito/files/ad2,avito/files/ad3",
			},
			wantErr: false,
		},
		{
			name: "Not Found - wit `advertisement not found` error",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name", "price", "description", "pictures"})

				mock.ExpectQuery("SELECT name, description, price, pictures FROM adverts WHERE (.+)").
					WithArgs(666).WillReturnRows(rows)
			},
			input: args{
				advertId: 666,
			},
			wantErr: true,
		},
		//make test "Not Found - wit other error"
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetAdvertById(test.input.advertId)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAdvertList(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' while opening a stub database connection", err)
	}

	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	r := NewAdvertRepository(db)

	type args struct {
		page        int
		orderField  string
		orderDirect string
	}

	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []model.Advert
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"name", "price", "pictures"}).
					AddRow("name-test1", 1000, "avito/files/ad1,avito/files/ad2,avito/files/ad3").
					AddRow("name-test2", 100, "avito/files/ad1,avito/files/ad2,avito/files/ad3").
					AddRow("name-test3", 10, "avito/files/ad1,avito/files/ad2,avito/files/ad3")

				mock.ExpectQuery("SELECT (.+) FROM adverts ORDER BY (.+)  LIMIT (.+) OFFSET (.+)").
					WithArgs(1).WillReturnRows(rows)
			},
			input: args{
				page:        1,
				orderField:  "price",
				orderDirect: "desc",
			},
			want: []model.Advert{
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
			},
			wantErr: false,
		},
		//make test "Not Found - wit error"
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetAdvertList(test.input.page, test.input.orderField, test.input.orderDirect)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
