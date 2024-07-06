package routes

import (
	"github.com/darshan-raul/accio/api/db"
	"github.com/darshan-raul/accio/api/models"
	"github.com/darshan-raul/accio/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func GetResourceTypes(c *fiber.Ctx) error {
	
	dbConn := db.DbConn()
	cloudname := &models.Cloud{}
	dbConn.Where("code = ?", "AWS").First(cloudname)
    return c.SendString(cloudname.Name)
}

func CreateResourceType(c *fiber.Ctx) error {
	// dbConn := db.DbConn()
	// cloudname := &models.Cloud{}
	// dbConn.Where("code = ?", "AWS").First(cloudname)
	// dbConn.Create(models.Resouretype{Name: "Instance", CloudID: cloudname.ID})
	controllers.Sendmsg()
	return c.SendString("Created resource type")
}