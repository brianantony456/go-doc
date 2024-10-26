package router

import (
	"net/http"

	"github.com/brianantony456/go-doc/internal/infrastructure/gin_handler"
	"github.com/brianantony456/go-doc/internal/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupRouter initializes the router with all the routes and returns it along with the DB connection
func SetupRouter() (*gin.Engine, *gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

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

	router.LoadHTMLGlob("internal/ui/web/templates/*")
	router.GET("/books", gin_handler.ServeBooksPage)
	router.GET("/create", func(c *gin.Context) {
		c.HTML(http.StatusOK, "create_book.html", gin.H{})
	})
	router.GET("/books/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "book_details.html", gin.H{})
	})
	router.GET("/checkout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "checkout_return.html", gin.H{})
	})

	return router, db, nil
}
