package service

import (
	"context"
	"fmt"

	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/models"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/repository"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/pkg/errs"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateMovie(ctx context.Context, m *models.Movie) (int64, error) {
	allowedRatings := map[int32]bool{
		0:  true,
		6:  true,
		12: true,
		16: true,
		18: true,
		21: true,
	}

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

	id, err := s.repo.Movie.Create(ctx, m)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrInternalServer, err)
	}

	return id, nil
}
