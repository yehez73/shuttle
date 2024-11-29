package services

import (
	"context"
	"errors"
	"log"
	"path/filepath"
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

func DeleteRefreshTokenOnLogout(userID string) error {
	client, err := database.MongoConnection()
	if err != nil {
		log.Print(err)
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
        log.Print(err)
        return nil, err
    }

    imageURL := GenerateImageURL(user.Picture)
    user.Picture = imageURL

    result, err := getUserByRoleCode(user, roleCode)
    if err != nil {
        log.Print(err)
        return nil, err
    }

    return result, nil
}

func GenerateImageURL(imagePath string) string {
    fileName := filepath.Base(imagePath)
    allowedExtensions := []string{".jpg", ".jpeg", ".png"}
    
    ext := filepath.Ext(fileName)
    if !contains(allowedExtensions, ext) {
        return ""
    }

    baseURL := "http://" + viper.GetString("BASE_URL") + "/assets/images/"
    return baseURL + fileName
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

func getUserByRoleCode(user models.User, roleCode string) (interface{}, error) {
    switch roleCode {
	case "SA":
		return user, nil
	case "AS":
		if details, ok := user.Details.(primitive.D); ok {
			detailsMap := make(map[string]interface{})
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			user.Details = detailsMap
			
			return user, nil
		} else {
			return nil, errors.New("school admin details are missing or invalid")
		}
	case "P":
		if details, ok := user.Details.(models.ParentDetails); ok {
			return models.UserResponse{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Role:      user.Role,
				RoleCode:  user.RoleCode,
				Status:    user.Status,
				Details:   details,
				CreatedAt: user.CreatedAt,
				CreatedBy: user.CreatedBy,
				UpdatedAt: user.UpdatedAt,
				UpdatedBy: user.UpdatedBy,
			}, nil
		} else {
			return nil, errors.New("parent details are missing or invalid")
		}
	case "D":
		if details, ok := user.Details.(primitive.D); ok {
			detailsMap := map[string]interface{}{}
			for _, elem := range details {
				detailsMap[elem.Key] = elem.Value
			}
			user.Details = detailsMap

			return user, nil
		} else {
			return nil, errors.New("driver details are missing or invalid")
		}
	default:
		return nil, errors.New("role code not found")
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
