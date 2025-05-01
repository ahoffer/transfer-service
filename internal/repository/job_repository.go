package repository

import (
	"errors"

	"github.com/transfer-service/internal/models"
	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found")

type JobRepository interface {
	Create(job *models.Job) error
	GetByID(id uint) (*models.Job, error)
	Update(job *models.Job) error
	Delete(id uint) error
	List() ([]*models.Job, error)
	Migrate() error
}

type GormJobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &GormJobRepository{db: db}
}

func (r *GormJobRepository) Create(job *models.Job) error {
	return r.db.Create(job).Error
}

func (r *GormJobRepository) GetByID(id uint) (*models.Job, error) {
	var job models.Job
	err := r.db.First(&job, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &job, nil
}

func (r *GormJobRepository) Update(job *models.Job) error {
	return r.db.Save(job).Error
}

func (r *GormJobRepository) Delete(id uint) error {
	return r.db.Delete(&models.Job{}, id).Error
}

func (r *GormJobRepository) List() ([]*models.Job, error) {
	var jobs []*models.Job
	err := r.db.Find(&jobs).Error
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *GormJobRepository) Migrate() error {
	return r.db.AutoMigrate(&models.Job{})
}
