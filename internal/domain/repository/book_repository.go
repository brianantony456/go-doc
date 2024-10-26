package repository

import "github.com/brianantony456/go-doc/internal/domain/model"

type BookRepository interface {
	GetAll() ([]model.Book, error)
	GetByID(id string) (*model.Book, error)
	Create(book model.Book) error
	Update(book model.Book) error
}
