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



// // GetBook func gets book by given ID or 404 error.
// // @Description Get book by given ID.
// // @Summary get book by given ID
// // @Tags Book
// // @Accept json
// // @Produce json
// // @Param id path string true "Book ID"
// // @Success 200 {object} models.Book
// // @Router /v1/book/{id} [get]
// func GetBook(c *fiber.Ctx) error {
// 	// Catch book ID from URL.
// 	id, err := uuid.Parse(c.Params("id"))
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   err.Error(),
// 		})
// 	}

// 	// Create database connection.
// 	db, err := database.OpenDBConnection()
// 	if err != nil {
// 		// Return status 500 and database connection error.
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   err.Error(),
// 		})
// 	}

// 	// Get book by ID.
// 	book, err := db.GetBook(id)
// 	if err != nil {
// 		// Return, if book not found.
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"error": true,
// 			"msg":   "book with the given ID is not found",
// 			"book":  nil,
// 		})
// 	}

// 	// Return status 200 OK.
// 	return c.JSON(fiber.Map{
// 		"error": false,
// 		"msg":   nil,
// 		"book":  book,
// 	})
// }