package persistence

import (
	"errors"

	"github.com/brianantony456/go-doc/internal/domain/model"
	"gorm.io/gorm"
)

type GormBookRepository struct {
	db *gorm.DB
}

func NewGormBookRepository(db *gorm.DB) *GormBookRepository {
	db.AutoMigrate(&model.Book{})
	return &GormBookRepository{db: db}
}

func (r *GormBookRepository) GetAll() ([]model.Book, error) {
	var books []model.Book
	result := r.db.Find(&books)
	return books, result.Error
}

func (r *GormBookRepository) GetByID(id string) (*model.Book, error) {
	var book model.Book
	result := r.db.First(&book, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("book not found")
	}
	return &book, result.Error
}

func (r *GormBookRepository) Create(book model.Book) error {
	return r.db.Create(&book).Error
}

func (r *GormBookRepository) Update(updatedBook model.Book) error {
	return r.db.Save(&updatedBook).Error
}
