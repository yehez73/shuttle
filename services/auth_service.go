package services

import (
	"context"
	"path/filepath"
	"time"

	"shuttle/databases"
	"shuttle/models"
	"shuttle/errors"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Login(email, password string) (models.User, error) {
	client, err := database.MongoConnection()
	if err != nil {
		return models.User{}, err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("users")

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return user, errors.New("user not found", 0)
	}

	if !validatePassword(password, user.Password) {
		return user, errors.New("password does not match", 0)
	}

	return user, nil
}

func DeleteRefreshTokenOnLogout(userID string) error {
	client, err := database.MongoConnection()
	if err != nil {
		return err
	}

	collection := client.Database(viper.GetString("MONGO_DB")).Collection("refresh_tokens")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"user_id": objectID})
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil
		}
		return err
	}

	return nil
}

func GetMyProfile(userID string, roleCode string) (interface{}, error) {
	user, err := GetSpecUser(userID)
	if err != nil {
		return nil, err
	}

	if user.Picture != "" {
		imageURL, err := generateImageURL(user.Picture)
		if err != nil {
			return nil, err
		}
		user.Picture = imageURL
	}

	result, err := getUserByRoleCode(user, roleCode)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func generateImageURL(imagePath string) (string, error) {
	fileName := filepath.Base(imagePath)
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}

	ext := filepath.Ext(fileName)
	if !contains(allowedExtensions, ext) {
		return "", errors.New("invalid image extension", 0)
	}

	baseURL := "http://" + viper.GetString("BASE_URL") + "/assets/images/"
	return baseURL + fileName, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getUserByRoleCode(user models.User, roleCode string) (models.UserResponse, error) {
	userResponse := models.UserResponse{
		ID:        user.ID,
		Picture:   user.Picture,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Gender:    user.Gender,
		Email:     user.Email,
		Role:      user.Role,
		RoleCode:  user.RoleCode,
		Phone:     user.Phone,
		Address:   user.Address,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}

	switch roleCode {
	case "SA":
		userResponse.Details = user.Details
		return userResponse, nil
	case "AS":
		if details, ok := user.Details.(primitive.D); ok {
			detailsMap := make(map[string]interface{})
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			userResponse.Details = detailsMap
			return userResponse, nil
		} else {
			return models.UserResponse{}, errors.New("school admin details are missing or invalid", 0)
		}
	case "P":
		if details, ok := user.Details.(models.ParentDetails); ok {
			userResponse.Details = details
			return userResponse, nil
		} else {
			return models.UserResponse{}, errors.New("parent details are missing or invalid", 0)
		}
	case "D":
		if details, ok := user.Details.(primitive.D); ok {
			detailsMap := map[string]interface{}{}
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			userResponse.Details = detailsMap
			return userResponse, nil
		} else {
			return models.UserResponse{}, errors.New("driver details are missing or invalid", 0)
		}
	default:
		return models.UserResponse{}, errors.New("role code not found", 0)
	}
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
			return "", errors.New("refresh token not found for the user", 0)
		}
		return "", err
	}

	if time.Now().After(storedToken.ExpiredAt) {
		return "", errors.New("refresh token has expired", 0)
	}

	return storedToken.RefreshToken, nil
}
