package storage

import (
	"app/api/models"
	"context"
)

type StorageI interface {
	CloseDB()
	Book() BookRepoI
}

type BookRepoI interface {
	Create(context.Context, *models.CreateBook) (string, error)
	GetByID(context.Context, *models.BookPrimaryKey) (*models.Book, error)
	GetList(context.Context, *models.GetListBookRequest) (*models.GetListBookResponse, error)
	Update(context.Context, *models.UpdateBook) (int64, error)
	Patch(ctx context.Context, req *models.PatchRequest) (int64, error)
	Delete(context.Context, *models.BookPrimaryKey) error
}
