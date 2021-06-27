package model

type Advert struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Pictures    string `json:"pictures,omitempty" binding:"required"`
	MainPicture string `json:"main-picture,omitempty""`
}
