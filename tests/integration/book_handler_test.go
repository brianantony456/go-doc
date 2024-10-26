package integrationtests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brianantony456/go-doc/internal/domain/model"
	"github.com/brianantony456/go-doc/internal/infrastructure/gin_handler"
	"github.com/brianantony456/go-doc/internal/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	bookRepo := persistence.NewGormBookRepository(db)
	bookHandler := gin_handler.NewBookHandler(bookRepo)

	router := gin.Default()
	apiRoutes := router.Group("/api")
	{
		apiRoutes.GET("/books", bookHandler.GetBooks)
		apiRoutes.POST("/books", bookHandler.CreateBook)
		apiRoutes.GET("/books/:id", bookHandler.BookById)
		apiRoutes.PATCH("/checkout", bookHandler.CheckoutBook)
		apiRoutes.PATCH("/return", bookHandler.ReturnBook)
	}

	return router
}

func createOneRow(router *gin.Engine) {
	// Insert a row
	newBook := model.Book{
		ID:       "1",
		Title:    "New Book",
		Author:   "New Author",
		Quantity: 10,
	}

	jsonValue, _ := json.Marshal(newBook)
	post_req, _ := http.NewRequest("POST", "/api/books", strings.NewReader(string(jsonValue)))
	post_req.Header.Set("Content-Type", "application/json")

	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, post_req)
}

func TestGetBooks(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/books", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateBook(t *testing.T) {
	router := setupRouter()

	newBook := model.Book{
		ID:       "4",
		Title:    "New Book",
		Author:   "New Author",
		Quantity: 10,
	}

	jsonValue, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/api/books", strings.NewReader(string(jsonValue)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// TODO: Fix tests these fail
func TestBookById(t *testing.T) {
	router := setupRouter()

	createOneRow(router)	

	req, _ := http.NewRequest("GET", "/api/books/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCheckoutBook(t *testing.T) {
	router := setupRouter()
	createOneRow(router)
	
	req, _ := http.NewRequest("PATCH", "/api/checkout?id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReturnBook(t *testing.T) {
	router := setupRouter()
	createOneRow(router)

	req, _ := http.NewRequest("PATCH", "/api/return?id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
