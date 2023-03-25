package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"app/api/models"
	"app/pkg/helper"
)

type bookRepo struct {
	db *pgxpool.Pool
}

func NewBookRepo(db *pgxpool.Pool) *bookRepo {
	return &bookRepo{
		db: db,
	}
}

func (r *bookRepo) Create(ctx context.Context, req *models.CreateBook) (string, error) {

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

	_, err := r.db.Exec(ctx, query,
		id.String(),
		req.Name,
		req.Price,
	)

	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (r *bookRepo) GetByID(ctx context.Context, req *models.BookPrimaryKey) (*models.Book, error) {

	var (
		query     string
		id        sql.NullString
		name      sql.NullString
		price     sql.NullFloat64
		status    sql.NullString
		createdAt sql.NullString
		updatedAt sql.NullString
	)

	query = `
		SELECT
			id,
			name,
			price,
			status,
			created_at,
			updated_at
		FROM book
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&name,
		&price,
		&status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &models.Book{
		Id:        id.String,
		Name:      name.String,
		Price:     price.Float64,
		Status:    status.String,
		CreatedAt: createdAt.String,
		UpdatedAt: updatedAt.String,
	}, nil
}

func (r *bookRepo) GetList(ctx context.Context, req *models.GetListBookRequest) (resp *models.GetListBookResponse, err error) {

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
			COALESCE(status, ''),
			TO_CHAR(created_at, 'YYYY-MM-DD HH24-MI-SS'),
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

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var book models.Book
		err = rows.Scan(
			&resp.Count,
			&book.Id,
			&book.Name,
			&book.Price,
			&book.Status,
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

func (r *bookRepo) Update(ctx context.Context, req *models.UpdateBook) (int64, error) {

	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE
			book
		SET
			name = :name,
			price = :price,
			updated_at = now()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":    req.Id,
		"name":  req.Name,
		"price": req.Price,
	}

	query, args := helper.ReplaceQueryParams(query, params)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *bookRepo) Patch(ctx context.Context, req *models.PatchRequest) (int64, error) {

	var (
		query string
		set   string
	)

	if len(req.Fields) <= 0 {
		return 0, errors.New("no fields")
	}

	for key := range req.Fields {
		set += fmt.Sprintf(" %s = :%s, ", key, key)
	}

	query = `
		UPDATE
			book
		SET
	` + set + ` updated_at = now()
		WHERE id = :id
	`

	req.Fields["id"] = req.ID

	fmt.Println(req.Fields)

	query, args := helper.ReplaceQueryParams(query, req.Fields)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (r *bookRepo) Delete(ctx context.Context, req *models.BookPrimaryKey) error {

	_, err := r.db.Exec(ctx,
		"DELETE FROM book WHERE id = $1", req.Id,
	)

	if err != nil {
		return err
	}

	return nil
}
