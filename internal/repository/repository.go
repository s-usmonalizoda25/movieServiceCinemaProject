package repository

import (
	"context"

	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/models"
)

type MovieRepository interface {
	Create(ctx context.Context, m *models.Movie) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Movie, error)
	Update(ctx context.Context, m *models.Movie) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]*models.Movie, error)
}
