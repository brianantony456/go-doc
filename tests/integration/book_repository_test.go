package integrationtests

import (
	"testing"

	"github.com/brianantony456/go-doc/internal/domain/model"
	"github.com/brianantony456/go-doc/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to the in-memory database: %v", err)
	}
	return db
}

func TestRepositoryCreateBook(t *testing.T) {
	db := setupInMemoryDB(t)
	repo := persistence.NewGormBookRepository(db)

	newBook := model.Book{
		ID:       "1",
		Title:    "A New Book",
		Author:   "John Doe",
		Quantity: 3,
	}

	err := repo.Create(newBook)
	assert.NoError(t, err)

	createdBook, err := repo.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, newBook.Title, createdBook.Title)
	assert.Equal(t, newBook.Author, createdBook.Author)
	assert.Equal(t, newBook.Quantity, createdBook.Quantity)
}

func TestRepositoryGetAllBooks(t *testing.T) {
	db := setupInMemoryDB(t)
	repo := persistence.NewGormBookRepository(db)

	books := []model.Book{
		{ID: "1", Title: "Book One", Author: "Author One", Quantity: 5},
		{ID: "2", Title: "Book Two", Author: "Author Two", Quantity: 2},
	}

	for _, book := range books {
		_ = repo.Create(book)
	}

	allBooks, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, allBooks, 2)
}

func TestRepositoryGetBookByID(t *testing.T) {
	db := setupInMemoryDB(t)
	repo := persistence.NewGormBookRepository(db)

	book := model.Book{ID: "1", Title: "A New Book", Author: "John Doe", Quantity: 3}
	_ = repo.Create(book)

	foundBook, err := repo.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, book.Title, foundBook.Title)
	assert.Equal(t, book.Author, foundBook.Author)
	assert.Equal(t, book.Quantity, foundBook.Quantity)
}

func TestRepositoryUpdateBook(t *testing.T) {
	db := setupInMemoryDB(t)
	repo := persistence.NewGormBookRepository(db)

	book := model.Book{ID: "1", Title: "A New Book", Author: "John Doe", Quantity: 3}
	_ = repo.Create(book)

	book.Quantity = 5
	err := repo.Update(book)
	assert.NoError(t, err)

	updatedBook, err := repo.GetByID("1")
	assert.NoError(t, err)
	assert.Equal(t, 5, updatedBook.Quantity)
}

func TestRepositoryGetNonExistentBookByID(t *testing.T) {
	db := setupInMemoryDB(t)
	repo := persistence.NewGormBookRepository(db)

	_, err := repo.GetByID("non-existent")
	assert.Error(t, err)
	assert.Equal(t, "book not found", err.Error())
}
