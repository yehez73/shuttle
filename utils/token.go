package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"shuttle/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var jwtSecret []byte
var encryptionKey []byte
var client *mongo.Client

func init() {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	jwtSecret = []byte(viper.GetString("JWT_SECRET"))
	encryptionKey = []byte(viper.GetString("ENCRYPTION_KEY"))

	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("MONGO_URI")))
	if err != nil {
		panic(err)
	}
}

// Signed Access Token
func GenerateToken(userId, name, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"name":   name,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 2).Unix(), // 2 hours expiration
	})

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	encryptedToken, err := encryptToken(signedToken)
	if err != nil {
		return "", err
	}

	return encryptedToken, nil
}

// Same, but with 15 days expiration time and for reissuing access token
func GenerateRefreshToken(userId, name, role string) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"name":   name,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24 * 15).Unix(), // 15 days expiration
	})

	signedRefreshToken, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	encryptedRefreshToken, err := encryptToken(signedRefreshToken)
	if err != nil {
		return "", err
	}

	return encryptedRefreshToken, nil
}

// AES encryption for tokens
func encryptToken(token string) (string, error) {
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encryptedToken := gcm.Seal(nonce, nonce, []byte(token), nil)
	return base64.URLEncoding.EncodeToString(encryptedToken), nil
}

func decryptToken(encryptedToken string) (string, error) {
	encryptedBytes, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedBytes) < nonceSize {
		return "", errors.New("malformed encrypted token")
	}

	nonce, ciphertext := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]
	decryptedToken, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decryptedToken), nil
}

func ValidateToken(encryptedToken string) (jwt.MapClaims, error) {
	decryptedToken, err := decryptToken(encryptedToken)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(decryptedToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func SaveRefreshToken(userId, refreshToken string) error {
	collection := client.Database(viper.GetString("MONGO_DB")).Collection("refresh_tokens")
	expiration := time.Now().Add(time.Hour * 24 * 15)

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(context.Background(), models.RefreshToken{
		UserID:       objectId,
		RefreshToken: refreshToken,
		ExpiredAt:    expiration,
	})
	return err
}

func GetUserIDFromToken(encryptedToken string) (string, error) {
	decryptedToken, err := decryptToken(encryptedToken)
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(decryptedToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["userId"].(string), nil
	}
	return "", err
}

var InvalidTokens = make(map[string]struct{})

func InvalidateToken(token string) {
	InvalidTokens[token] = struct{}{}
}