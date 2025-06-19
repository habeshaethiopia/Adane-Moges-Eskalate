package services

import (
	"errors"
	"strings"

	"eskalate-movie-api/internal/models"
	"eskalate-movie-api/internal/repository"

	"github.com/google/uuid"
)

type MovieService interface {
	Create(movie *models.Movie) error
	Update(movie *models.Movie, userID uuid.UUID) error
	Delete(movieID uuid.UUID, userID uuid.UUID) error
	GetByID(movieID uuid.UUID) (*models.Movie, error)
	GetAll(pageNumber, pageSize int) ([]models.Movie, int64, error)
	Search(title string, pageNumber, pageSize int) ([]models.Movie, int64, error)
}

type movieService struct {
	repo repository.MovieRepository
}

func NewMovieService(repo repository.MovieRepository) MovieService {
	return &movieService{repo}
}

func (s *movieService) Create(movie *models.Movie) error {
	return s.repo.Create(movie)
}

func (s *movieService) Update(movie *models.Movie, userID uuid.UUID) error {
	m, err := s.repo.FindByID(movie.ID)
	if err != nil {
		return err
	}
	if m.UserID != userID {
		return errors.New("forbidden")
	}
	movie.UserID = userID
	return s.repo.Update(movie)
}

func (s *movieService) Delete(movieID uuid.UUID, userID uuid.UUID) error {
	m, err := s.repo.FindByID(movieID)
	if err != nil {
		return err
	}
	if m.UserID != userID {
		return errors.New("forbidden")
	}
	return s.repo.Delete(m)
}

func (s *movieService) GetByID(movieID uuid.UUID) (*models.Movie, error) {
	return s.repo.FindByID(movieID)
}

func (s *movieService) GetAll(pageNumber, pageSize int) ([]models.Movie, int64, error) {
	offset := (pageNumber - 1) * pageSize
	return s.repo.FindAll(offset, pageSize)
}

func (s *movieService) Search(title string, pageNumber, pageSize int) ([]models.Movie, int64, error) {
	offset := (pageNumber - 1) * pageSize
	return s.repo.SearchByTitle(strings.ToLower(title), offset, pageSize)
}
