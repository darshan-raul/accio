package controllers

import (
	"github.com/darshan-raul/accio/core/platform/database"
	"github.com/gofiber/fiber/v2"
)

// GetBooks func gets all exists cloud providers.
// @Description Get all exists cloud providers.
// @Summary get all exists cloud providers
// @Tags CloudProvs
// @Accept json
// @Produce json
// @Success 200 {array} models.CloudProv
// @Router /v1/cloudprovs [get]
func GetCloudProvs(c *fiber.Ctx) error {
	// Create database connection.
	db, err := database.OpenDBConnection()
	if err != nil {
		// Return status 500 and database connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get all cloudprovs.
	cloudprovs, err := db.GetCloudProvs()
	if err != nil {
		// Return, if cloudprovs not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   "cloudprovs were not found",
			"count": 0,
			"cloudprovs": nil,
			"err": err,
		})
	}

	// Return status 200 OK.
	return c.JSON(fiber.Map{
		"error": false,
		"msg":   nil,
		"count": len(cloudprovs),
		"cloudprovs": cloudprovs,
	})
}

