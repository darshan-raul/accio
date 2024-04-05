package routes

import (

	"github.com/darshan-raul/accio/api/db"
	"github.com/darshan-raul/accio/api/models"
	"github.com/gofiber/fiber/v2"
)

//TODO: move the cloud grab logic to diff controller package

func GetClouds(c *fiber.Ctx) error {
	clouds := []models.Cloud{}
	dbConn := db.DbConn()
	dbConn.Find(&clouds)
    return c.JSON(clouds)
}

//TODO: move the cloud grab logic to diff controller package

func GetCloud(c *fiber.Ctx) error {
	cloud := c.Params("cloud")
	dbConn := db.DbConn()
	cloudname := &models.Cloud{}
	result := dbConn.Where("code = ?", cloud).First(cloudname)
	
	if int(result.RowsAffected) == 0 {
		return c.SendString("No entry exists")
	}
    return c.SendString(cloudname.Name)
}