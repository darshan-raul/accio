package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
  
func DbConn() *gorm.DB {
	dsn := "host=localhost user=postgres password=asdf123 dbname=accio port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {

		fmt.Println(err)
	}

	return db
}