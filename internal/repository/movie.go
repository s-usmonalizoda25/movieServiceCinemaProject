package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/models"
)

type movieRepository struct {
	pool *pgxpool.Pool
}

func NewMovieRepository(pool *pgxpool.Pool) MovieRepository {
	return &movieRepository{pool: pool}
}

func (r *movieRepository) Create(ctx context.Context, m *models.Movie) (int64, error) {
	var id int64
	const query = `INSERT INTO movies (title, description, duration, age_limit, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	err := r.pool.QueryRow(ctx, query, m.Title, m.Description, m.Duration, m.AgeLimit).Scan(&id)
	return id, err
}

func (r *movieRepository) GetByID(ctx context.Context, id int64) (*models.Movie, error) {
	m := &models.Movie{}
	query := `SELECT id, title, description, duration, age_limit, created_at, updated_at FROM movies WHERE id = $1`
	err := r.pool.QueryRow(ctx, query, id).Scan(&m.ID, &m.Title, &m.Description, &m.Duration, &m.AgeLimit, &m.CreatedAt, &m.UpdatedAt)
	return m, err
}

func (r *movieRepository) Update(ctx context.Context, m *models.Movie) error {
	const query = `UPDATE movies SET title=$1, description=$2, duration=$3, age_limit=$4, updated_at=NOW() WHERE id=$5`
	_, err := r.pool.Exec(ctx, query, m.Title, m.Description, m.Duration, m.AgeLimit, m.ID)
	return err
}

func (r *movieRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM movies WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}
