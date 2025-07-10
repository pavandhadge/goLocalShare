package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadToCloudinary uploads a file to Cloudinary using user credentials and deletes it after the given duration.
func UploadToCloudinary(filePath, cloudName, apiKey, apiSecret string, duration time.Duration) (string, error) {
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return "", fmt.Errorf("Cloudinary credentials not set")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", fmt.Errorf("failed to initialize Cloudinary: %v", err)
	}

	ctx := context.Background()

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	publicID := fmt.Sprintf("temp_%d_%s", time.Now().Unix(), filepath.Base(filePath))

	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       publicID,
		UniqueFilename: Bool(false),
		Overwrite:      Bool(true),
		ResourceType:   "auto",
	})

	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	// Delete the file after the specified duration
	go func(publicID string) {
		time.Sleep(duration)
		_, err := cld.Upload.Destroy(ctx, uploader.DestroyParams{
			PublicID: publicID,
			Type:     "upload",
		})
		if err != nil {
			log.Printf("Failed to delete Cloudinary file %s: %v", publicID, err)
		} else {
			log.Printf("Successfully deleted Cloudinary file: %s", publicID)
		}
	}(uploadResult.PublicID)

	return uploadResult.SecureURL, nil
}

func Bool(b bool) *bool {
	return &b
} 