package controllers

import (
	"strconv"
	"time"

	"github.com/darshan-raul/accio/core/app/models"
	"github.com/darshan-raul/accio/core/pkg/utils"
	"github.com/darshan-raul/accio/core/platform/database"
	"github.com/gofiber/fiber/v2"
)

// GetProjects func gets all exists projects.
// @Description Get all exists projects.
// @Summary get all exists projects
// @Tags Projects
// @Accept json
// @Produce json
// @Success 200 {array} models.Project
// @Router /v1/projects [get]
func GetProjects(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all projects.
	projects, err := db.GetProjects()
	if err != nil {
		// Return, if projects not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "projects were not found",
			"count": 0,
			"projects": nil,
			"err": err,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(projects),
		"projects": projects,
	})
}




// GetProject func gets project by given ID or 404 error.
// @Description Get project by given ID.
// @Summary get project by given ID
// @Tags Project
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} models.Project
// @Router /v1/project/{id} [get]
func GetProject(c *fiber.Ctx) error {
	// Catch project ID from URL.
	strid := c.Params("id")
	id, err := strconv.Atoi(strid)
    if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
    }
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get project by ID.
	project, err := db.GetProject(id)
	if err != nil {
		// Return, if project not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "project with the given ID is not found",
			"project":  nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"project":  project,
	})
}

// CreateProject func for creates a new project.
// @Description Create a new project.
// @Summary create a new project
// @Tags Project
// @Accept json
// @Produce json
// @Param name body string true "Name"
// @Success 200 {object} models.Project
// @Security ApiKeyAuth
// @Router /v1/project [post]
func CreateProject(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current project.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}


	// Create new Project struct
	project := &models.Project{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(project); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Project model.
	validate := utils.NewValidator()

	// Set initialized default data for project:
	project.CreatedAt = time.Now()

	// Validate project fields.
	if err := validate.Struct(project); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create project by given model.
	if err := db.CreateProject(project); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"project":  project,
	})
}