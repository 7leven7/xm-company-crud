package configs

import (
	"fmt"

	"github.com/7leven7/xm-company-crud/app/database"
	"github.com/7leven7/xm-company-crud/app/models"
)

// PGConnection connecting to a database using the configuration loaded from a local file. It returns an error if the config file cannot be loaded or if the database connection fails.
func PGConnection() error {
	config, err := database.LoadConfig(".")

	if err != nil {
		return fmt.Errorf("could not load environment variables: %v", err)
	}

	err = database.ConnectDB(&config)
	if err != nil {
		return fmt.Errorf("could not connect to database: %v", err)
	}

	return nil
}

// Migrate database migrations for the given model
func Migrate() {
	err := database.DB.AutoMigrate(&models.Company{}, &models.User{})
	if err != nil {
		panic(err)
	}

	fmt.Println("? Migration complete")
}
