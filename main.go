package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pavandhadge/goFileShare/server"
	"github.com/pavandhadge/goFileShare/handlers"
)

// CloudinaryConfig holds credentials
type CloudinaryConfig struct {
	CloudName string `json:"cloud_name"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

func getCloudinaryConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "cloudinary.json" // fallback to local dir
	}
	return filepath.Join(home, ".gofileserver_cloudinary.json")
}

func saveCloudinaryConfig(cfg CloudinaryConfig) error {
	path := getCloudinaryConfigPath()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(cfg)
}

func loadCloudinaryConfig() (CloudinaryConfig, error) {
	path := getCloudinaryConfigPath()
	f, err := os.Open(path)
	if err != nil {
		return CloudinaryConfig{}, err
	}
	defer f.Close()
	var cfg CloudinaryConfig
	err = json.NewDecoder(f).Decode(&cfg)
	return cfg, err
}

func main() {
	port := ":8090"

	var durationStr string
	var dirMode bool
	var cloudMode bool
	var cloudName string
	var cloudKey string
	var cloudSecret string
	flag.StringVar(&durationStr, "duration", "1h", "Sharing session duration (e.g. 3h, 30m, 1h30m)")
	flag.BoolVar(&dirMode, "dir", false, "Share a directory instead of a file")
	flag.BoolVar(&cloudMode, "cloud", false, "Upload file to Cloudinary instead of serving locally")
	flag.StringVar(&cloudName, "cloud-name", "", "Cloudinary cloud name (required for --cloud)")
	flag.StringVar(&cloudKey, "cloud-key", "", "Cloudinary API key (required for --cloud)")
	flag.StringVar(&cloudSecret, "cloud-secret", "", "Cloudinary API secret (required for --cloud)")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("Usage: ./server [--duration 3h] [--dir] <path> or ./server [--duration 3h] [--cloud --cloud-name <name> --cloud-key <key> --cloud-secret <secret>] <file>")
	}

	// Use cloudMode directly for the credential logic
	if cloudMode {
		if cloudName == "" || cloudKey == "" || cloudSecret == "" {
			// Try to load from config file
			cfg, err := loadCloudinaryConfig()
			if err == nil && cfg.CloudName != "" && cfg.APIKey != "" && cfg.APISecret != "" {
				cloudName, cloudKey, cloudSecret = cfg.CloudName, cfg.APIKey, cfg.APISecret
			} else {
				log.Fatal("Cloudinary credentials required for --cloud. Provide via flags or set up ~/.gofileserver_cloudinary.json")
			}
		} else {
			// Save to config for future use
			saveCloudinaryConfig(CloudinaryConfig{cloudName, cloudKey, cloudSecret})
		}
	}

	shareDuration, err := time.ParseDuration(durationStr)
	if err != nil || shareDuration < time.Minute {
		log.Fatalf("Invalid duration: %v", durationStr)
	}

	basePath := flag.Arg(0)
	isDirMode := dirMode
	isCloudMode := cloudMode

	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		log.Fatalf("Path resolution error: %v", err)
	}
	basePath = absBasePath

	http.HandleFunc("/file", handlers.FilePageHandler)
	http.HandleFunc("/browse", handlers.BrowsePageHandler)

	server.Start(port, basePath, isDirMode, isCloudMode, shareDuration, cloudName, cloudKey, cloudSecret)
}
