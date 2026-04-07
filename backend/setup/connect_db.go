package setup

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func FormatDBConnectionString(c Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.DBHost,
		c.DBUserName,
		c.DBUserPassword,
		c.DBName,
		c.DBPort,
	)
}

func ConnectDB(connectionString string) (*gorm.DB, error) {
	DB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Println("? Connected Successfully to the Database")

	return DB, nil
}
