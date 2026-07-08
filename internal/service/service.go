package service

import (
	"context"
	"fmt"

	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/models"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/repository"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/pkg/errs"
	"go.uber.org/zap"
)

type MovieService interface {
	CreateMovie(ctx context.Context, m *models.Movie) (int64, error)
	GetMovieByID(ctx context.Context, id int64) (*models.Movie, error)
	UpdateMovie(ctx context.Context, m *models.Movie) error
	DeleteMovie(ctx context.Context, id int64) error
}

type service struct {
	repo repository.MovieRepository
	log  *zap.Logger
}

func New(repo repository.MovieRepository, log *zap.Logger) MovieService {
	return &service{
		repo: repo,
		log:  log,
	}
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
		s.log.Error("failed to create movie in repository", zap.Error(err))
		return 0, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}

	s.log.Info("movie successfully created", zap.Int64("id", id))
	return id, nil
}

func (s *service) GetMovieByID(ctx context.Context, id int64) (*models.Movie, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("failed to get movie from repository", zap.Int64("id", id), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	return m, nil
}

func (s *service) UpdateMovie(ctx context.Context, m *models.Movie) error {
	err := s.repo.Update(ctx, m)
	if err != nil {
		s.log.Error("failed to update movie in repository", zap.Int64("id", m.ID), zap.Error(err))
		return fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	s.log.Info("movie updated", zap.Int64("id", m.ID))
	return nil
}

func (s *service) DeleteMovie(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Error("failed to delete movie from repository", zap.Int64("id", id), zap.Error(err))
		return fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	s.log.Info("movie deleted", zap.Int64("id", id))
	return nil
}
