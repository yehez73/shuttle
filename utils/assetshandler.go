package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const MaxFileSize = 10 * 1024 * 1024 // 10 MB

func HandleUploadedFile(c *fiber.Ctx) (string, error) {
	file, err := c.FormFile("picture")
	if err != nil {
		return "", BadRequestResponse(c, "Picture is required", nil)
	}

	if !IsValidImageExtension(file.Filename) {
		return "", BadRequestResponse(c, "Invalid image file extension", nil)
	}

	src, err := file.Open()
	if err != nil {
		return "", ErrorResponse(c, http.StatusInternalServerError, "Failed to open image file", nil)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", ErrorResponse(c, http.StatusInternalServerError, "Failed to read image file", nil)
	}

	if !IsValidImageType(fileBytes) {
		return "", BadRequestResponse(c, "Invalid image file type", nil)
	}

	pictureFileName, err := SavePicture(fileBytes, file.Filename)
	if err != nil {
		return "", ErrorResponse(c, http.StatusInternalServerError, "Failed to save picture", nil)
	}

	return pictureFileName, nil
}

func HandleAssetsOnUpdate(c *fiber.Ctx, existingPicture string) (string, error) {
    if existingPicture != "" {
        err := DeletePicture(existingPicture)
        if err != nil {
            return "", fmt.Errorf("Failed to delete existing picture: %v", err)
        }
    }

    pictureFileName, err := HandleUploadedFile(c)
    if err != nil {
        return "", fmt.Errorf("Failed to upload new picture: %v", err)
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
		fmt.Println("Error decoding image:", err)
		return false
	}

	if img == nil {
		fmt.Println("Invalid image")
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