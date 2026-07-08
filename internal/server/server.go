package server

import (
	"context"

	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/models"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/internal/service"
	pb "github.com/s-usmonalizoda25/movieServiceCinemaProject/moviepb/v1"
	"github.com/s-usmonalizoda25/movieServiceCinemaProject/pkg/errs"
	"go.uber.org/zap"
)

type Server struct {
	pb.UnimplementedMovieServiceServer
	log *zap.Logger
	svc *service.Service
}

func New(log *zap.Logger, svc *service.Service) *Server {
	return &Server{
		log: log,
		svc: svc,
	}
}

func (s *Server) Create(ctx context.Context, req *pb.CreateMovieRequest) (*pb.CreateMovieResponse, error) {
	id, err := s.svc.CreateMovie(ctx, &models.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		AgeLimit:    req.AgeLimit,
	})
	if err != nil {
		return nil, s.handleError(errs.MsgFailedCreate, err)
	}
	return &pb.CreateMovieResponse{Id: id}, nil
}

func (s *Server) GetByID(ctx context.Context, req *pb.GetMovieRequest) (*pb.GetMovieResponse, error) {
	m, err := s.svc.GetMovieByID(ctx, req.Id)
	if err != nil {
		return nil, s.handleError(errs.MsgFailedGet, err)
	}
	return &pb.GetMovieResponse{
		Id:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Duration:    m.Duration,
		AgeLimit:    m.AgeLimit,
	}, nil
}

func (s *Server) Update(ctx context.Context, req *pb.UpdateMovieRequest) (*pb.UpdateMovieResponse, error) {
	err := s.svc.UpdateMovie(ctx, &models.Movie{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		AgeLimit:    req.AgeLimit,
	})
	if err != nil {
		return nil, s.handleError(errs.MsgFailedUpdate, err)
	}
	return &pb.UpdateMovieResponse{Code: 200, Message: "success"}, nil
}

func (s *Server) Delete(ctx context.Context, req *pb.DeleteMovieRequest) (*pb.DeleteMovieResponse, error) {
	err := s.svc.DeleteMovie(ctx, req.Id)
	if err != nil {
		return nil, s.handleError(errs.MsgFailedDelete, err)
	}
	return &pb.DeleteMovieResponse{Code: 200, Message: "success"}, nil
}

func (s *Server) List(ctx context.Context, req *pb.ListMovieRequest) (*pb.ListMovieResponse, error) {
	movies, err := s.svc.GetAllMovies(ctx)
	if err != nil {
		return nil, s.handleError(errs.MsgFailedGet, err)
	}

	var pbMovies []*pb.GetMovieResponse
	for _, m := range movies {
		pbMovies = append(pbMovies, &pb.GetMovieResponse{
			Id:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			Duration:    m.Duration,
			AgeLimit:    m.AgeLimit,
		})
	}
	return &pb.ListMovieResponse{Movies: pbMovies}, nil
}

func (s *Server) handleError(msg string, err error) error {
	s.log.Error(msg, zap.Error(err))
	return errs.MapToGRPC(err)
}
