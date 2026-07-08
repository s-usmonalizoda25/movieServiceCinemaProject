package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct{
	Movie *MovieRepository
}

func New(pool *pgxpool.Pool) *Repository{
	return &Repository{
		Movie: NewMovieRepository(pool),
	}
}

