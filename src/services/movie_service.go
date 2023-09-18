package services

import (
	"clean_api/src/api/dto"
	"clean_api/src/api/filters"
	"clean_api/src/data/db"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type MovieService struct {
	DB *sql.DB
}

func NewMovieService() *MovieService {
	return &MovieService{
		DB: db.GetDB(),
	}
}

/*

	Id
	Title
	Year
	Genres
	Runtime
	CreatedAt
*/

func (m *MovieService) GetById(id int32) (*dto.MovieResponse, error) {
	q := `
	SELECT id, title, year, genres, runtime, created_at, version
	FROM movies 
	WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	response := dto.MovieResponse{}
	err := m.DB.QueryRowContext(ctx, q, id).Scan(
		&response.Id,
		&response.Title,
		&response.Year,
		pq.Array(&response.Genres),
		&response.Runtime,
		&response.CreatedAt,
		&response.Version,
	)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (m *MovieService) Create(req *dto.CreateMovie) (*dto.MovieResponse, error) {
	q := `
	INSERT INTO movies (title, year, runtime, genres)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	args := []interface{}{req.Title, req.Year, req.Runtime, pq.Array(req.Genres)}
	var id int32
	err := m.DB.QueryRow(q, args...).Scan(&id)
	if err != nil {
		return nil, err
	}
	return m.GetById(id)
}


func (m *MovieService) Update(req *dto.UpdateMovie, id int32,version int32) (*dto.MovieResponse, error) {
	q := `
	UPDATE movies SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1 
	WHERE ID = $5 AND version = $6 
	RETURNING id`

	args := []interface{}{
		req.Title,
		req.Year,
		req.Runtime,
		pq.Array(req.Genres),
		id, 
		version,
	}
	var uid int32
	err := m.DB.QueryRow(q, args...).Scan(&uid)
	if err != nil {
		return nil, err
	}
	return m.GetById(uid)
}

func (m *MovieService) Delete(id int64) (error) {
	q := `
	DELETE FROM movies WHERE id = $1
	`
	res ,err := m.DB.Exec(q,id)
	if err != nil {
		return err
	}
	rowsDeleted, err := res.RowsAffected()
	if err != nil {
		return err
	} 
	if rowsDeleted == 0 {
		return errors.New("no rows deleted")
	}
	return nil
}

func (m *MovieService) GetAll(title string, genres []string, filter filters.Filter) ([]*dto.MovieResponse,*filters.MetaData ,error) {
	q := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND (genres @> $2 OR $2 = '{}')
	ORDER BY %s %s ,id ASC
	LIMIT $3 OFFSET $4`, filter.SortCol(), filter.SortDirection())

	args := []interface{}{
		title, 
		pq.Array(genres),
		filter.Limit(),
		filter.Offset(),
	}

	rows, err := m.DB.Query(q, args...)
	if err != nil {
		return nil,nil, err
	}
	defer rows.Close()

	totalRecord := 0
	res :=[]*dto.MovieResponse{}
	for rows.Next() {
		var mv dto.MovieResponse
		
		err := rows.Scan(
			&totalRecord,
			&mv.Id,
			&mv.CreatedAt,
			&mv.Title,
			&mv.Year,
			&mv.Runtime,
			pq.Array(&mv.Genres),
			&mv.Version,
		)
		if err != nil{
			return nil, nil, err
		}
		res = append(res, &mv)
	}

	mData := filters.CalculateMetaData(totalRecord, filter.Page, filter.PageSize)
	return res, &mData, nil
}