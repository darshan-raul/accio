package controllers

import (
	// "fmt"
	"fmt"
	"strconv"
	"time"

	"github.com/darshan-raul/accio/core/app/models"
	"github.com/darshan-raul/accio/core/pkg/utils"
	"github.com/darshan-raul/accio/core/platform/database"
	"github.com/gofiber/fiber/v2"
)

// GetStacks func gets all exists stacks.
// @Description Get all exists stacks.
// @Summary get all exists stacks
// @Tags Stacks
// @Accept json
// @Produce json
// @Success 200 {array} models.Stack
// @Router /v1/stacks [get]
func GetStacks(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all stacks.
	stacks, err := db.GetStacks()
	if err != nil {
		// Return, if stacks not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "stacks were not found",
			"count":    0,
			"stacks": nil,
			"err":      err,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"count":    len(stacks),
		"stacks": stacks,
	})
}

// GetStack func gets stack by given ID or 404 error.
// @Description Get stack by given ID.
// @Summary get stack by given ID
// @Tags Stack
// @Accept json
// @Produce json
// @Param id path string true "Stack ID"
// @Success 200 {object} models.Stack
// @Router /v1/stack/{id} [get]
func GetStack(c *fiber.Ctx) error {
	// Catch stack ID from URL.
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

	// Get stack by ID.
	stack, err := db.GetStack(id)
	if err != nil {
		// Return, if stack not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"msg":     "stack with the given ID is not found",
			"stack": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"stack": stack,
	})
}

// CreateStack func for creates a new stack.
// @Description Create a new stack.
// @Summary create a new stack
// @Tags Stack
// @Accept json
// @Produce json
// @Param name body string true "Name"
// @Success 200 {object} models.CreateStackRequest
// @Security ApiKeyAuth
// @Router /v1/stack [post]
func CreateStack(c *fiber.Ctx) error {
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

	// Set expiration time from JWT data of current stack.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new Stack struct
	stack := &models.Stack{}

	stackreq := &models.CreateStackRequest{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(stackreq); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	fmt.Println(stackreq.ProjectName)
	fmt.Println(stackreq.StackName)

	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}


	// Create a new validator for a Stack model.
	validate := utils.NewValidator()


	// Validate stack fields.
	if err := validate.Struct(stackreq); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}


	var project models.Project
    
	project, err = db.GetProjectByName(stackreq.ProjectName)

	if err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
	})
	}

	stack.Name = stackreq.StackName
	stack.ProjectId = project.ID

	// // Create stack by given model.
	if err := db.CreateStack(stack); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"stack": stack,
	})
}

// DeleteStack func for deletes stack by given ID.
// @Description Delete stack by given ID.
// @Summary delete stack by given ID
// @Tags Stack
// @Accept json
// @Produce json
// @Param id body string true "Stack ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/stack [delete]
func DeleteStack(c *fiber.Ctx) error {
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

	// Set expiration time from JWT data of current stack.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new Stack struct
	stack := &models.Stack{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(stack); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Stack model.
	validate := utils.NewValidator()

	// Validate stack fields.
	if err := validate.StructPartial(stack, "id"); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
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
	fmt.Println(stack.ID)
	// Checking, if stack with given ID is exists.
	foundedStack, err := db.GetStack(stack.ID)
	if err != nil {
		// Return status 404 and stack not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "msg":   "stack with this ID not found"})
	}
	fmt.Println(foundedStack)
	if err := db.DeleteStack(foundedStack.ID); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)

}
