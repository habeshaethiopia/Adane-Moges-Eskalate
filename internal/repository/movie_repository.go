package repository

import (
	"eskalate-movie-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MovieRepository interface {
	Create(movie *models.Movie) error
	Update(movie *models.Movie) error
	Delete(movie *models.Movie) error
	FindByID(id uuid.UUID) (*models.Movie, error)
	FindAll(offset, limit int) ([]models.Movie, int64, error)
	SearchByTitle(title string, offset, limit int) ([]models.Movie, int64, error)
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db}
}

func (r *movieRepository) Create(movie *models.Movie) error {
	return r.db.Create(movie).Error
}

func (r *movieRepository) Update(movie *models.Movie) error {
	return r.db.Save(movie).Error
}

func (r *movieRepository) Delete(movie *models.Movie) error {
	return r.db.Delete(movie).Error
}

func (r *movieRepository) FindByID(id uuid.UUID) (*models.Movie, error) {
	var movie models.Movie
	if err := r.db.First(&movie, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *movieRepository) FindAll(offset, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	r.db.Model(&models.Movie{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&movies).Error
	return movies, total, err
}

func (r *movieRepository) SearchByTitle(title string, offset, limit int) ([]models.Movie, int64, error) {
	var movies []models.Movie
	var total int64
	q := r.db.Model(&models.Movie{})
	if title != "" {
		q = q.Where("LOWER(title) LIKE ?", "%"+title+"%")
	}
	q.Count(&total)
	err := q.Offset(offset).Limit(limit).Find(&movies).Error
	return movies, total, err
}
