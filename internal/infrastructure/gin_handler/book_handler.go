package gin_handler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/brianantony456/go-doc/internal/domain/model"
	"github.com/brianantony456/go-doc/internal/domain/repository"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	repo repository.BookRepository
}

func NewBookHandler(repo repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

func generateID() string {
	return uuid.New().String()
}

func (h *BookHandler) GetBooks(c *gin.Context) {
	books, err := h.repo.GetAll()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, books)
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	// Read the raw data from the request body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// Print the raw JSON body for debugging purposes
	fmt.Println("Provided JSON: ", string(body))

	// Reset the request body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Bind the JSON to newBook struct
	var newBook model.Book
	if err := c.BindJSON(&newBook); err != nil {
		fmt.Println("create book, bind error: ", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// Log the values to ensure they are captured correctly
	fmt.Printf("Title: %s, Author: %s, Quantity: %d\n", newBook.Title, newBook.Author, newBook.Quantity)

	// Assign a new ID if not provided
	if newBook.ID == "" {
		newBook.ID = generateID()
	}

	// Create the book entry in the repository
	if err := h.repo.Create(newBook); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not create book"})
		return
	}
	c.IndentedJSON(http.StatusCreated, newBook)
}

func (h *BookHandler) BookById(c *gin.Context) {
	id := c.Param("id")
	book, err := h.repo.GetByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func (h *BookHandler) CheckoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id parameter"})
		return
	}

	book, err := h.repo.GetByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not Found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})
		return
	}

	book.Quantity -= 1
	if err := h.repo.Update(*book); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not update book"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func (h *BookHandler) ReturnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id parameter"})
		return
	}

	book, err := h.repo.GetByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not Found"})
		return
	}

	book.Quantity += 1
	if err := h.repo.Update(*book); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Could not update book"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func ServeBooksPage(c *gin.Context) {
	c.HTML(http.StatusOK, "books.html", gin.H{})
}
