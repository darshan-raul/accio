package routes

import (
	"github.com/darshan-raul/accio/api/db"
	"github.com/darshan-raul/accio/api/models"
	"github.com/gofiber/fiber/v2"
)

func GetResourceTypes(c *fiber.Ctx) error {
	
	dbConn := db.DbConn()
	cloudname := &models.Cloud{}
	dbConn.Where("code = ?", "AWS").First(cloudname)
    return c.SendString(cloudname.Name)
}