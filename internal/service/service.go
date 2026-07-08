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
	GetAllMovies(ctx context.Context) ([]*models.Movie, error)
}

type Service struct {
	repo repository.MovieRepository
	log  *zap.Logger
}

func New(repo repository.MovieRepository, log *zap.Logger) MovieService {
	return &Service{
		repo: repo,
		log:  log,
	}
}

func (s *Service) CreateMovie(ctx context.Context, m *models.Movie) (int64, error) {
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
		s.log.Error(errs.MsgFailedCreate, zap.Error(err))
		return 0, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}

	s.log.Info("movie successfully created", zap.Int64("id", id))
	return id, nil
}

func (s *Service) GetMovieByID(ctx context.Context, id int64) (*models.Movie, error) {
	m, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.log.Error(errs.MsgFailedGet, zap.Int64("id", id), zap.Error(err))
		return nil, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	return m, nil
}

func (s *Service) UpdateMovie(ctx context.Context, m *models.Movie) error {
	err := s.repo.Update(ctx, m)
	if err != nil {
		s.log.Error(errs.MsgFailedUpdate, zap.Int64("id", m.ID), zap.Error(err))
		return fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	s.log.Info("movie updated", zap.Int64("id", m.ID))
	return nil
}

func (s *Service) DeleteMovie(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Error(errs.MsgFailedDelete, zap.Int64("id", id), zap.Error(err))
		return fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}
	s.log.Info("movie deleted", zap.Int64("id", id))
	return nil
}

func (s *Service) GetAllMovies(ctx context.Context) ([]*models.Movie, error) {
	movies, err := s.repo.GetAll(ctx)
	if err != nil {
		s.log.Error(errs.MsgFailedGet, zap.Error(err))
		return nil, errs.ErrInternalServer
	}
	return movies, nil
}
