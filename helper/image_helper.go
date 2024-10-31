package helper

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
)


func ValidateImageFile(file *multipart.FileHeader) error {
    allowedFormats := map[string]bool{
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
        ".heic":  true, // Add any additional formats you want to support
    }

    // Get the file extension
    ext := strings.ToLower(filepath.Ext(file.Filename))
    if !allowedFormats[ext] {
        return errors.New("unsupported file format. Please upload a JPG, JPEG, PNG, or GIF file")
    }

    return nil
}