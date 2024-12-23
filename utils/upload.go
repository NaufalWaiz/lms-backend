package utils

import (
    "bytes"
    "fmt"
    "mime/multipart"
    "time"
	
    "project/config"
)

func UploadImage(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
    // Baca konten file
    buf := new(bytes.Buffer)
    _, err := buf.ReadFrom(file)
    if err != nil {
        return "", fmt.Errorf("failed to read file: %v", err)
    }

    // Buat nama unik untuk file
    fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), fileHeader.Filename)
    filePath := fmt.Sprintf("uploads/%s", fileName)

    // Unggah file ke bucket Supabase
    response := config.SupabaseClient.Storage.From("assets").Upload(filePath, buf)
    if response.Message != "" {
        return "", fmt.Errorf("failed to upload file: %v", response.Message)
    }

    // Buat URL file
    imageURL := fmt.Sprintf("%s/storage/v1/object/public/assets/%s", config.SupabaseURL, filePath)
    return imageURL, nil
}