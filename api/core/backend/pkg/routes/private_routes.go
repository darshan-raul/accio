package routes

import (
	"github.com/darshan-raul/accio/core/app/controllers"
	"github.com/darshan-raul/accio/core/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)           // create a new book
	route.Post("/project", middleware.JWTProtected(), controllers.CreateProject)
	route.Post("/stack", middleware.JWTProtected(), controllers.CreateStack)
	route.Post("/resource", middleware.JWTProtected(), controllers.CreateResource)
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)   // renew Access & Refresh tokens

	// Routes for PUT method:
	route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID

	// Routes for DELETE method:
	route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID
	route.Delete("/project", middleware.JWTProtected(), controllers.DeleteProject)
	route.Delete("/stack", middleware.JWTProtected(), controllers.DeleteStack)
	route.Delete("/resource", middleware.JWTProtected(), controllers.DeleteResource)

}
