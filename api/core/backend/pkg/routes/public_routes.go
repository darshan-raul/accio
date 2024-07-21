package routes

import (
	"github.com/darshan-raul/accio/core/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for GET method:
	route.Get("/books", controllers.GetBooks)   // get list of all books
	route.Get("/book/:id", controllers.GetBook) // get one book by ID

	route.Get("/cloudprovs", controllers.GetCloudProvs)
	route.Get("/resourcetypes", controllers.GetResourceTypes)
	route.Get("/projects", controllers.GetProjects)
	route.Get("/project/:id", controllers.CreateProject)

	route.Get("/stacks", controllers.GetStacks)
	route.Get("/stack/:id", controllers.GetStack)
	

	route.Get("/resources", controllers.GetResources)
	route.Get("/resource/:id", controllers.GetResource)

	// Routes for POST method:
	route.Post("/user/sign/up", controllers.UserSignUp) // register a new user
	route.Post("/user/sign/in", controllers.UserSignIn) // auth, return Access & Refresh tokens
}
