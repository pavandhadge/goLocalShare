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

func UploadToCloudinary(filePath string) (string, error) {
	cloudName := "du1jbnyp0"
	apiKey := "476468415943614"
	apiSecret := "RjmUy0N30VGpxpKM6TnXRZyUCFs"

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return "", fmt.Errorf("Cloudinary credentials not set in environment variables")
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

	go func(publicID string) {
		time.Sleep(1 * time.Hour)
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