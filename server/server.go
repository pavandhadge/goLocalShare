package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pavandhadge/goFileShare/auth"
	"github.com/pavandhadge/goFileShare/handlers"
	"github.com/pavandhadge/goFileShare/utils"
)

const (
	timeout = 15 * time.Second
)

func Start(port, basePath string, isDirMode, isCloudMode bool, shareDuration time.Duration) {
	if _, err := utils.SecureStat(basePath); err != nil {
		log.Fatalf("Path error: %v", err)
	}

	if isCloudMode {
		url, err := utils.UploadToCloudinary(basePath)
		if err != nil {
			log.Fatalf("Failed to upload to Cloudinary: %v", err)
		}
		fmt.Printf("File uploaded to Cloudinary. URL: %s\n", url)
		fmt.Println("This URL will be automatically deleted in 1 hour.")
		return
	}

	handlers.Init(basePath, isDirMode, shareDuration)
	auth.InitOwnerToken(shareDuration)

	// Session expiry middleware
	sessionMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if auth.IsSessionExpired() {
				handlers.SessionEndedHandler(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Minimal API and UI routes
	http.Handle("/api/files", sessionMiddleware(CORS(AuthCheck(RateLimitMiddleware(http.HandlerFunc(handlers.FilesAPIHandler))))))
	http.Handle("/api/download/", sessionMiddleware(CORS(AuthCheck(RateLimitMiddleware(http.HandlerFunc(handlers.DownloadAPIHandler))))))
	http.Handle("/", sessionMiddleware(SecurityHeaders(http.HandlerFunc(handlers.HomeHandler))))
	http.Handle("/token", sessionMiddleware(SecurityHeaders(http.HandlerFunc(handlers.TokenEntryPageHandler))))
	http.Handle("/test", sessionMiddleware(SecurityHeaders(http.HandlerFunc(handlers.TestPageHandler))))

	// Add /download/ legacy route: redirect to /api/download/
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		// Strip /download/ and redirect to /api/download/
		newPath := "/api" + r.URL.Path
		if r.URL.RawQuery != "" {
			newPath += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, newPath, http.StatusTemporaryRedirect)
	})

	// No more legacy /browse, /browse/, /download/ routes

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  2 * timeout,
		Handler:      nil,
	}

	log.Printf("Secure file server started on http://localhost%s", port)
	log.Printf("API endpoints available at http://localhost%s/api", port)
	log.Printf("Session will expire at: %s", time.Now().Add(shareDuration).Format(time.RFC1123))
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
} 