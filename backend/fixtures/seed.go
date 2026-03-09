package main

import (
	"context"
	"crypto/rand"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/caaldrid/mindtracer/backend/models"
	"github.com/caaldrid/mindtracer/backend/setup"
)

const testUserEmail = "seed@test.local"

func generatePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	bytePassword := make([]byte, length)
	_, err := rand.Read(bytePassword)
	if err != nil {
		return "", err
	}
	for i := range length {
		bytePassword[i] = charset[int(bytePassword[i])%len(charset)]
	}
	return string(bytePassword), nil
}

func main() {
	config, err := setup.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	DB, err := setup.ConnectDB(config)
	if err != nil {
		log.Fatal("Could connect to database instance", err)
	}
	cntx := context.Background()

	usr := &models.User{}

	result := DB.WithContext(cntx).Where("email = ? ", testUserEmail).First(usr)

	if result.Error == nil {
		log.Println("Seed data is already present in the DB")
	} else if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Println("Seeding Db...")

		usr.Email = testUserEmail
		usr.UserName = "Seed User"

		randPass, err := generatePassword(32)
		if err != nil {
			log.Fatalf("Failed to generate a password: %s", err.Error())
		}
		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(randPass),
			bcrypt.DefaultCost,
		)
		if err != nil {
			log.Fatalf("Could not encrypt seed user pswd: %s", err.Error())
		}
		usr.Password = string(passwordHash)

		result = DB.WithContext(cntx).Create(usr)

		if result.Error != nil {
			log.Fatalf("Failed to create a user in the DB: %s", result.Error.Error())
		}
	} else {
		log.Fatalf("%s", result.Error.Error())
	}
}
