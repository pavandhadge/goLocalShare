package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/pavandhadge/goFileShare/server"
	"github.com/pavandhadge/goFileShare/handlers"
)

func main() {
	port := ":8090"

	var durationStr string
	var dirMode bool
	var cloudMode bool
	flag.StringVar(&durationStr, "duration", "1h", "Sharing session duration (e.g. 3h, 30m, 1h30m)")
	flag.BoolVar(&dirMode, "dir", false, "Share a directory instead of a file")
	flag.BoolVar(&cloudMode, "cloud", false, "Upload file to Cloudinary instead of serving locally")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("Usage: ./server [--duration 3h] [--dir] <path> or ./server [--duration 3h] [--cloud] <file>")
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

	server.Start(port, basePath, isDirMode, isCloudMode, shareDuration)
}
