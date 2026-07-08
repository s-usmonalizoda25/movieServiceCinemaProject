package service

import (
	"context"
	"fmt"

	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/models"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/repository"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/pkg/errs"
)

type MovieService interface {
	CreateMovie(ctx context.Context, m *models.Movie) (int64, error)
	GetMovieByID(ctx context.Context, id int64) (*models.Movie, error)
	UpdateMovie(ctx context.Context, m *models.Movie) error
	DeleteMovie(ctx context.Context, id int64) error
}

type service struct {
	repo repository.MovieRepository
}

func New(repo repository.MovieRepository) MovieService {
	return &service{repo: repo}
}

func (s *service) CreateMovie(ctx context.Context, m *models.Movie) (int64, error) {
	allowedRatings := map[int32]bool{0: true, 6: true, 12: true, 16: true, 18: true, 21: true}

	if !allowedRatings[m.AgeLimit] {
		return 0, errs.ErrInvalidAgeLimit
	}
	if len(m.Title) < 1 {
		return 0, errs.ErrTitleEmpty
	}
	if m.Duration <= 0 {
		return 0, errs.ErrDurationInvalid
	}
	if len(m.Description) < 10 {
		return 0, errs.ErrDescriptionShort
	}

	id, err := s.repo.Create(ctx, m)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}

	return id, nil
}

func (s *service) GetMovieByID(ctx context.Context, id int64) (*models.Movie, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	return m, nil
}

func (s *service) UpdateMovie(ctx context.Context, m *models.Movie) error {
	err := s.repo.Update(ctx, m)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	return nil
}

func (s *service) DeleteMovie(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	return nil
}
