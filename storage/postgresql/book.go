package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	"app/api/models"
)

type bookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *bookRepo {
	return &bookRepo{
		db: db,
	}
}

func (r *bookRepo) Create(req *models.CreateBook) (string, error) {

	var (
		query string
		id    = uuid.New()
	)

	query = `
		INSERT INTO book(
			id, 
			name, 
			price, 
			updated_at
		)
		VALUES ($1, $2, $3, now())
	`

	_, err := r.db.Exec(query,
		id.String(),
		req.Name,
		req.Price,
	)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (r *bookRepo) GetByID(req *models.BookPrimaryKey) (*models.Book, error) {

	var (
		query string
		book  models.Book
	)

	query = `
		SELECT
			id,
			name,
			price,
			created_at,
			updated_at
		FROM book
		WHERE id = $1
	`

	err := r.db.QueryRow(query, req.Id).Scan(
		&book.Id,
		&book.Name,
		&book.Price,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (r *bookRepo) GetList(req *models.GetListBookRequest) (resp *models.GetListBookResponse, err error) {

	resp = &models.GetListBookResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			name,
			price,
			created_at,
			updated_at
		FROM book
	`

	if len(req.Search) > 0 {
		filter += " AND name ILIKE '%' || '" + req.Search + "' || '%' "
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + offset + limit

	rows, err := r.db.Query(query)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var book models.Book
		err = rows.Scan(
			&resp.Count,
			&book.Id,
			&book.Name,
			&book.Price,
			&book.CreatedAt,
			&book.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Books = append(resp.Books, &book)
	}

	return resp, nil
}
