package db

import (
	"fmt"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/darshan-raul/accio/api/utils"
)
  
func DbConn() *gorm.DB {
	postgres_host := utils.GetEnv("POSTGRES_HOST")
	postgres_user := utils.GetEnv("POSTGRES_USER")
	postgres_user_password := utils.GetEnv("POSTGRES_PASSWORD")
	postgres_db := utils.GetEnv("POSTGRES_DB")
	postgres_port := utils.GetEnv("POSTGRES_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", postgres_host,postgres_user,postgres_user_password,postgres_db,postgres_port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {

		fmt.Println(err)
	}
	fmt.Println("db connectivity successful")
	return db
}
