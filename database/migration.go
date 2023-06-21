package database

import (
	"fmt"
	"golang/models"
	"golang/pkg/mysql"
)

func RunMigration() {
	// database.DB.AutoMigrate(&entity.User{}, &next-entity)
	err := mysql.DB.AutoMigrate(&models.User{})

	if err != nil {
		fmt.Println(err)
		panic("Migration Failed")
	}

	fmt.Println("Migration Success")
}
