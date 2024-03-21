package main

import (
	"github.com/gofiber/fiber/v2"
	"fmt"

	"github.com/darshan-raul/accio/api/db"
	"github.com/darshan-raul/accio/api/models"
    "github.com/darshan-raul/accio/api/routes"
)

func main() {
    app := fiber.New()

    app.Get("/ping", func(c *fiber.Ctx) error {
        return c.SendString("pong")
    })

    app.Get("/cloud", routes.GetClouds)
    app.Get("/cloud/:cloud", routes.GetCloud)
    app.Get("/type", routes.GetResourceTypes)

    defer app.Listen(":3000")

    handleDbMigrationsSeeding()
    
}

func handleDbMigrationsSeeding(){

    dbConn := db.DbConn()

    dbConn.AutoMigrate(&models.Cloud{},&models.Resource{},&models.Resouretype{})
    
    var count int64 //because Count() returns int64

	dbConn.Model(&models.Cloud{}).Count(&count)
	if count > 0 {
		fmt.Println("seed data already exists, skipping creation of initial data")
	}else {
        cloudData := []models.Cloud{
            {Name: "Amazon Web Services", Code: "AWS"},
            {Name: "Azure", Code: "Azure"},
            {Name: "Google Cloud Platform", Code: "GCP"},
        }
        for _,entry := range cloudData{
            dbConn.Create(&entry)
        }

    }
}
