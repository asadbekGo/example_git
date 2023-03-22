package models

type Book struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type BookPrimaryKey struct {
	Id string `json:"id"`
}

type CreateBook struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type UpdateBook struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type GetListBookRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetListBookResponse struct {
	Count int     `json:"count"`
	Books []*Book `json:"books"`
}
