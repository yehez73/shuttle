package services

import (
	"context"
	"errors"
	"log"
	"time"

	"shuttle/databases"
	"shuttle/models"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Login(email, password string) (models.User, error) {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return models.User{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, errors.New("user not found")
	}

	if !validatePassword(password, user.Password) {
		return user, errors.New("password does not match")
	}

	return user, nil
}

func GetMyProfile(userID string) (models.UserResponse, error) {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return models.UserResponse{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.UserResponse
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return user, err
	}

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func validatePassword(providedPassword, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
	return err == nil
}

func GetStoredRefreshToken(userID string) (string, error) {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
		return "", err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("refresh_tokens")
	
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", err
	}

	var storedToken models.RefreshToken
	err = collection.FindOne(context.Background(), bson.M{"user_id": objectID}).Decode(&storedToken)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			log.Println("refresh token not found for the user")
			return "", errors.New("refresh token not found for the user")
		}
		log.Print(err)
		return "", err
	}

	if time.Now().After(storedToken.ExpiredAt) {
		return "", errors.New("refresh token has expired")
	}

	return storedToken.RefreshToken, nil
}