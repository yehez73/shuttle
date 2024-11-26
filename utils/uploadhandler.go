package utils

import (
    "os"
    "path/filepath"
    "fmt"
    "github.com/google/uuid"
)

const MaxFileSize = 10 * 1024 * 1024 // 10 MB

// Fungsi untuk memvalidasi ekstensi file gambar
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

// Fungsi untuk memvalidasi ukuran file
func IsValidFileSize(fileSize int64) bool {
    return fileSize <= MaxFileSize
}

// Fungsi untuk sanitasi nama file agar tidak ada path traversal
func SanitizeFileName(fileName string) string {
    return filepath.Base(fileName) // Hanya mengambil nama file
}

// Fungsi untuk menyimpan file gambar ke server
func SavePicture(fileBytes []byte, fileName string) (string, error) {
    folderPath := "./uploads/profile_pictures"
    err := os.MkdirAll(folderPath, os.ModePerm)
    if err != nil {
        return "", err
    }

    // Menghasilkan nama file unik dengan UUID
    sanitizedFileName := fmt.Sprintf("%s_%s", uuid.New().String(), SanitizeFileName(fileName))
    filePath := filepath.Join(folderPath, sanitizedFileName)

    err = os.WriteFile(filePath, fileBytes, 0644)
    if err != nil {
        return "", err
    }

    return filePath, nil
}