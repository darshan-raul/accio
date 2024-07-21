package controllers

import (

	"strconv"
	"time"

	"github.com/darshan-raul/accio/core/app/models"
	"github.com/darshan-raul/accio/core/pkg/utils"
	"github.com/darshan-raul/accio/core/platform/database"
	"github.com/gofiber/fiber/v2"
)

// GetResources func gets all exists resources.
// @Description Get all exists resources.
// @Summary get all exists resources
// @Tags Resources
// @Accept json
// @Produce json
// @Success 200 {array} models.Resource
// @Router /v1/resources [get]
func GetResources(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all resources.
	resources, err := db.GetResources()
	if err != nil {
		// Return, if resources not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      "resources were not found",
			"count":    0,
			"resources": nil,
			"err":      err,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"count":    len(resources),
		"resources": resources,
	})
}

// GetResource func gets resource by given ID or 404 error.
// @Description Get resource by given ID.
// @Summary get resource by given ID
// @Tags Resource
// @Accept json
// @Produce json
// @Param id path string true "Resource ID"
// @Success 200 {object} models.Resource
// @Router /v1/resource/{id} [get]
func GetResource(c *fiber.Ctx) error {
	// Catch resource ID from URL.
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

	// Get resource by ID.
	resource, err := db.GetResource(id)
	if err != nil {
		// Return, if resource not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"msg":     "resource with the given ID is not found",
			"resource": nil,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error":   false,
		"msg":     nil,
		"resource": resource,
	})
}

// CreateResource func for creates a new resource.
// @Description Create a new resource.
// @Summary create a new resource
// @Tags Resource
// @Accept json
// @Produce json
// @Param name body string true "Name"
// @Success 200 {object} models.CreateResourceRequest
// @Security ApiKeyAuth
// @Router /v1/resource [post]
func CreateResource(c *fiber.Ctx) error {
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

	// Set expiration time from JWT data of current resource.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new Resource struct
	resource := &models.Resource{}

	resourcereq := &models.CreateResourceRequest{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(resourcereq); err != nil {
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


	// Create a new validator for a Resource model.
	validate := utils.NewValidator()


	// Validate resource fields.
	if err := validate.Struct(resourcereq); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}


	var stack models.Stack
	var cloudprov models.CloudProv
	var resourcetype models.ResourceType

    
	stack, err = db.GetStackByName(resourcereq.StackName)

	if err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
	})
	}


	cloudprov, err = db.GetCloudProvBySlug(resourcereq.CloudProviderSlug)

	if err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
	})
	}

	resourcetype, err = db.GetResourceTypeByName(resourcereq.Type)

	if err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
	})
	}

	resource.Name = resourcereq.ResourceName
	resource.CloudProvId = cloudprov.ID
	resource.StackId = stack.ID
	resource.TypeId = resourcetype.ID


	// Set initialized default data for project:
	resource.CreatedAt = time.Now()
	
	// // Create resource by given model.
	if err := db.CreateResource(resource); err != nil {
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
		"resource": resource,
	})
}


// DeleteResource func for deletes resource by given ID.
// @Description Delete resource by given ID.
// @Summary delete resource by given ID
// @Tags Resource
// @Accept json
// @Produce json
// @Param id body string true "Resource ID"
// @Success 204 {string} status "ok"
// @Security ApiKeyAuth
// @Router /v1/resource [delete]
func DeleteResource(c *fiber.Ctx) error {
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

	// Set expiration time from JWT data of current resource.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create new Resource struct
	resource := &models.Resource{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(resource); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a Resource model.
	validate := utils.NewValidator()

	// Validate resource fields.
	if err := validate.StructPartial(resource, "id"); err != nil {
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

	// Checking, if resource with given ID is exists.
	foundedResource, err := db.GetResource(resource.ID)
	if err != nil {
		// Return status 404 and resource not found error.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "msg":   "resource with this ID not found"})
	}

	if err := db.DeleteResource(foundedResource.ID); err != nil {
		// Return status 500 and error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)

}
