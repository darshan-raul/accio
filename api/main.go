package main

import (
    //"github.com/gofiber/fiber/v2"
    "github.com/darshan-raul/accio/api/db"
    "github.com/darshan-raul/accio/api/models"
)
func main() {
    // app := fiber.New()

    // app.Get("/ping", func(c *fiber.Ctx) error {
    //     return c.SendString("pong")
    // })

    // app.Get("/resource", func(c *fiber.Ctx) error {

    //     return c.SendString("resource")
    // })

    // app.Listen(":3000")

    dbConn := db.DbConn()

    dbConn.AutoMigrate(&models.Cloud{},&models.Resource{},&models.Type{})

    dbConn.Create(&models.Cloud{Name: "AWS"})
    dbConn.Create(&models.Cloud{Name: "Azure"})
    dbConn.Create(&models.Cloud{Name: "GCP"})
}