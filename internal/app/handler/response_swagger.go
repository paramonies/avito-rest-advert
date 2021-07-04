package handler

//types for swagger
type InputAdvert struct {
	Name        string `json:"name" example:"name-test"`
	Description string `json:"description" example:"desc-test"`
	Price       int    `json:"price" example:1000`
	Pictures    string `json:"pictures" example:"avito/files/ad1,avito/files/ad2,avito/files/ad3"`
}

type CreateMessageOk struct {
	Id int `json:"id" example:"1"`
}

type CreateMessage400 struct {
	Message string `json:"error" example:"invalid input body"`
}

type CreateMessage500 struct {
	Message string `json:"error" example:"internal server error"`
}

type GetMessageOk struct {
	Name        string `json:"name" example:"name-test"`
	Description string `json:"description" example:"desc-test"`
	Price       int    `json:"price" example:1000`
	Pictures    string `json:"pictures" example:"avito/files/ad1,avito/files/ad2,avito/files/ad3"`
}

type GetMessage400 struct {
	Message string `json:"error" example:"advertisement id must be integer"`
}

type GetMessage500 struct {
	Message string `json:"error" example:"internal server error"`
}

type ListMessageOk struct {
	Name        string `json:"name" example:"name-test"`
	Description string `json:"description" example:"desc-test"`
	Price       int    `json:"price" example:1000`
	Pictures    string `json:"pictures" example:"avito/files/ad1,avito/files/ad2,avito/files/ad3"`
}

type ListMessageOk1 []ListMessageOk

type ListMessage500 struct {
	Message string `json:"error" example:"internal server error"`
}
