package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"shuttle/logger"
	"shuttle/errors"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

const MaxFileSize = 10 * 1024 * 1024 // 10 MB

func HandleUploadedFile(c *fiber.Ctx) (string, error) {
	file, err := c.FormFile("picture")
	if err != nil {
		logger.LogError(err, "Failed to get file", nil)
		return "", errors.New("picture is required", 400)
	}

	if !IsValidImageExtension(file.Filename) {
		return "", errors.New("invalid image file extension", 400)
	}

	src, err := file.Open()
	if err != nil {
		logger.LogError(err, "Failed to open file", nil)
		return "", errors.New("something went wrong, please try again later", 500)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		logger.LogError(err, "Failed to read file", nil)
		return "", errors.New("something went wrong, please try again later", 500)
	}

	if !IsValidImageType(fileBytes) {
		return "", errors.New("invalid image file type", 400)
	}

	pictureFileName, err := SavePicture(fileBytes, file.Filename)
	if err != nil {
		logger.LogError(err, "Failed to save picture", nil)
		return "", errors.New("something went wrong, please try again later", 500)
	}
	
	return pictureFileName, nil
}

func HandleAssetsOnUpdate(c *fiber.Ctx, existingPicture string) (string, error) {
	println("existingPicture: ", existingPicture)
    if existingPicture != "" {
        err := DeletePicture(existingPicture)
        if err != nil {
            return "", err
        }
    }

    pictureFileName, err := HandleUploadedFile(c)
    if err != nil {
        return "", err
    }

    return pictureFileName, nil
}

func IsValidImageExtension(fileName string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png"}
	ext := filepath.Ext(fileName)
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

func IsValidImageType(fileBytes []byte) bool {
	img, err := imaging.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return false
	}

	if img == nil {
		return false
	}

	return true
}

func IsValidFileSize(fileSize int64) bool {
	return fileSize <= MaxFileSize
}

func SanitizeFileName(fileName string) string {
	return filepath.Base(fileName)
}

func SavePicture(fileBytes []byte, fileName string) (string, error) {
	folderPath := "./assets/images"

	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	sanitizedFileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fileName))

	filePath := filepath.Join(folderPath, sanitizedFileName)

	err = os.WriteFile(filePath, fileBytes, 0644)
	if err != nil {
		return "", err
	}

	return sanitizedFileName, nil
}

func DeletePicture(fileName string) error {
    if fileName == "" {
        return nil
    }

    filePath := filepath.Join("./assets/images", fileName)
    err := os.Remove(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }

    return nil
}

func GenerateImageAssetsURL(imagePath string) (string, error) {
	fileName := filepath.Base(imagePath)
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}

	ext := filepath.Ext(fileName)
	if !contains(allowedExtensions, ext) {
		return "", fmt.Errorf("invalid image extension")
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