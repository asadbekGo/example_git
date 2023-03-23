package storage

import (
	"app/api/models"
)

type StorageI interface {
	CloseDB()
	Book() BookRepoI
}

type BookRepoI interface {
	Create(*models.CreateBook) (string, error)
	GetByID(*models.BookPrimaryKey) (*models.Book, error)
	GetList(*models.GetListBookRequest) (*models.GetListBookResponse, error)
	Update(*models.UpdateBook) (int64, error)
	Delete(*models.BookPrimaryKey) error
}
