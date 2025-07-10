package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pavandhadge/goFileShare/utils"
	"mime"
)

type FileInfo struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	IsDirectory  bool      `json:"isDirectory"`
	Size         int64     `json:"size,omitempty"`
	ModTime      time.Time `json:"modTime"`
	SizeFormatted string   `json:"sizeFormatted,omitempty"`
}

type FilesResponse struct {
	Files      []FileInfo `json:"files,omitempty"`
	CurrentDir string     `json:"currentDir"`
	BasePath   string     `json:"basePath"`
	BaseInfo   FileInfo   `json:"baseInfo"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func FilesAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("X-Auth-Token")
	log.Printf("[FilesAPIHandler] Token: '%s'", token)

	// Extract path from URL
	path := strings.TrimPrefix(r.URL.Path, "/api/files")
	if path == "" {
		path = "/"
	}
	log.Printf("[FilesAPIHandler] Requested path: '%s'", path)

	// Validate path is within base directory
	cleanPath, err := utils.SecurePath(BasePath, path)
	log.Printf("[FilesAPIHandler] Resolved path: '%s'", cleanPath)
	if err != nil {
		response := ErrorResponse{
			Error:   "Invalid path",
			Message: "The requested path is not accessible",
			Code:    http.StatusForbidden,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	stat, err := os.Stat(cleanPath)
	if err != nil {
		response := ErrorResponse{
			Error:   "Not found",
			Message: "The requested file or directory does not exist",
			Code:    http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	baseInfo := FileInfo{
		Name:        stat.Name(),
		Path:        path,
		IsDirectory: stat.IsDir(),
		ModTime:     stat.ModTime(),
	}
	if !stat.IsDir() {
		baseInfo.Size = stat.Size()
		baseInfo.SizeFormatted = formatFileSize(stat.Size())
	}

	var fileList []FileInfo
	if stat.IsDir() {
		files, err := os.ReadDir(cleanPath)
		if err != nil {
			response := ErrorResponse{
				Error:   "Directory read error",
				Message: "Could not read the requested directory",
				Code:    http.StatusInternalServerError,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response)
			return
		}

		for _, file := range files {
			info, err := file.Info()
			if err != nil {
				continue // Skip files we can't get info for
			}

			fileInfo := FileInfo{
				Name:        file.Name(),
				Path:        filepath.Join(path, file.Name()),
				IsDirectory: file.IsDir(),
				ModTime:     info.ModTime(),
			}

			if !file.IsDir() {
				fileInfo.Size = info.Size()
				fileInfo.SizeFormatted = formatFileSize(info.Size())
			}

			fileList = append(fileList, fileInfo)
		}
	}

	response := FilesResponse{
		Files:      fileList,
		CurrentDir: path,
		BasePath:   BasePath,
		BaseInfo:   baseInfo,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DownloadAPIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("X-Auth-Token")
	log.Printf("[DownloadAPIHandler] Token: '%s'", token)

	requestedPath := strings.TrimPrefix(r.URL.Path, "/api/download/")
	log.Printf("[DownloadAPIHandler] requestedPath: '%s'", requestedPath)

	// Check if base is a file or directory
	baseInfo, err := utils.SecureStat(BasePath)
	if err != nil {
		http.Error(w, "Base path error", http.StatusInternalServerError)
		return
	}

	var cleanPath string
	if baseInfo.IsDir() {
		// Directory mode: resolve relative to base
		cleanPath, err = utils.SecurePath(BasePath, requestedPath)
		if err != nil {
			log.Printf("[DownloadAPIHandler] SecurePath error: %v", err)
			http.Error(w, "Invalid path", http.StatusForbidden)
			return
		}
	} else {
		// Single-file mode: only allow empty or filename
		baseFile := filepath.Base(BasePath)
		if requestedPath == "" || requestedPath == baseFile {
			cleanPath = BasePath
		} else {
			log.Printf("[DownloadAPIHandler] Single-file mode: requested '%s', allowed '%s'", requestedPath, baseFile)
			http.NotFound(w, r)
			return
		}
	}
	log.Printf("[DownloadAPIHandler] cleanPath: '%s'", cleanPath)

	fileInfo, err := utils.SecureStat(cleanPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if fileInfo.IsDir() {
		http.Error(w, "Cannot download a directory", http.StatusBadRequest)
		return
	}

	file, err := utils.SecureOpen(cleanPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	defer file.Close()

	mimeType := mime.TypeByExtension(filepath.Ext(cleanPath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)

	if strings.HasPrefix(mimeType, "image/") {
		w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(cleanPath)))
	} else {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(cleanPath)))
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Cache-Control", "no-store")

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
} 