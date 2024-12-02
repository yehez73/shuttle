package main

import (
	"context"
	"log"
	"shuttle/databases"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	SuperAdmin Role = "superadmin"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Role     Role               `json:"role" bson:"role"`
	RoleCode string             `json:"role_code" bson:"role_code"`
}

func main() {
	client, err := database.MongoConnection()
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	count, err := collection.CountDocuments(context.Background(), bson.M{"role": "superadmin"})
	if err != nil {
		log.Fatal("Failed to count users:", err)
		return
	}

	if count > 0 {
		color.Yellow("Superadmin already exists, no need to seed.")
		return
	}

	hashedPassword, err := hashPassword("12345678")
	if err != nil {
		log.Fatal("Error hashing password:", err)
		return
	}

	userID, err := primitive.ObjectIDFromHex("000000000000000000000000")
	if err != nil {
		log.Fatal("Failed to create ObjectID:", err)
		return
	}

	users := []interface{}{
		User{
			ID:       userID,
			Email:    "faker@gmail.com",
			Password: hashedPassword,
			Role:     SuperAdmin,
			RoleCode: "SA",
		},
	}

	_, err = collection.InsertMany(context.Background(), users)
	if err != nil {
		log.Fatal("Failed to insert users:", err)
		return
	}

	color.Green("Users seeded successfully!")
	color.Yellow("Please login and create a new user with superadmin role immediately and delete this user.")

	defer client.Disconnect(context.Background())
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}