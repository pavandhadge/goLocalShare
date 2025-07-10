package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pavandhadge/goFileShare/auth"
	"github.com/pavandhadge/goFileShare/handlers"
	"github.com/pavandhadge/goFileShare/utils"
)

const (
	timeout = 15 * time.Second
)

func Start(port, basePath string, isDirMode, isCloudMode bool, shareDuration time.Duration, cloudName, cloudKey, cloudSecret string) {
	if _, err := utils.SecureStat(basePath); err != nil {
		log.Fatalf("Path error: %v", err)
	}

	if isCloudMode {
		if cloudName == "" || cloudKey == "" || cloudSecret == "" {
			log.Fatal("Cloudinary credentials are required in cloud mode. Use --cloud-name, --cloud-key, --cloud-secret.")
		}
		url, err := utils.UploadToCloudinary(basePath, cloudName, cloudKey, cloudSecret, shareDuration)
		if err != nil {
			log.Fatalf("Failed to upload to Cloudinary: %v", err)
		}
		fmt.Printf("File uploaded to Cloudinary. URL: %s\n", url)
		fmt.Printf("This URL will be automatically deleted in %s.\n", shareDuration)
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

	// Check for cert.pem and key.pem in the current directory
	if _, err := os.Stat("cert.pem"); os.IsNotExist(err) {
		log.Fatal("cert.pem not found. Please generate a self-signed certificate using: openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=goFileServer'")
	}
	if _, err := os.Stat("key.pem"); os.IsNotExist(err) {
		log.Fatal("key.pem not found. Please generate a self-signed certificate using: openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj '/CN=goFileServer'")
	}

	// Print startup info with token, duration, and links using the actual IP address
	if !isCloudMode {
		// Get the owner token and expiry
		token, expiry := auth.GetOwnerToken()

		// Try to get the local IP address
		ip := "localhost"
		ifaces, err := net.Interfaces()
		if err == nil {
			for _, iface := range ifaces {
				if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
					continue
				}
				addrs, err := iface.Addrs()
				if err != nil {
					continue
				}
				for _, addr := range addrs {
					var ipStr string
					switch v := addr.(type) {
					case *net.IPNet:
						ipStr = v.IP.String()
					case *net.IPAddr:
						ipStr = v.IP.String()
					}
					if ipStr != "" && !strings.HasPrefix(ipStr, "127.") && !strings.HasPrefix(ipStr, "::1") && !strings.HasPrefix(ipStr, "169.254") && !strings.HasPrefix(ipStr, "fe80:") {
						ip = ipStr
						break
					}
				}
				if ip != "localhost" {
					break
				}
			}
		}

		// Remove colon from port if present
		portNum := port
		if strings.HasPrefix(port, ":") {
			portNum = port[1:]
		}

		// Determine the file name if not in directory mode
		fileName := ""
		if !isDirMode {
			// Try to get the file name from basePath
			fileName = filepath.Base(basePath)
		}

		protocol := "https"
		log.Printf("\n\x1b[1;32mgoFileServer started!\x1b[0m")
		log.Printf("\x1b[1;34mAccess token:\x1b[0m %s", token)
		log.Printf("\x1b[1;34mSession duration:\x1b[0m %s (expires at %s)", shareDuration, expiry.Format(time.RFC1123))
		log.Printf("\x1b[1;34mUser access link:\x1b[0m %s://%s:%s/token", protocol, ip, portNum)
		if isDirMode {
			log.Printf("\x1b[1;34mBrowse link:\x1b[0m %s://%s:%s/browse?token=%s", protocol, ip, portNum, token)
		} else {
			log.Printf("\x1b[1;34mFile link:\x1b[0m %s://%s:%s/file?file=%s&token=%s", protocol, ip, portNum, fileName, token)
		}
		log.Printf("\x1b[1;34mAPI endpoint:\x1b[0m %s://%s:%s/api/files?token=%s", protocol, ip, portNum, token)

		// Open the user access link in the default browser (localhost for browser)
		go func() {
			browserURL := fmt.Sprintf("%s://localhost:%s/", protocol, portNum)
			_ = openBrowser(browserURL)
		}()
	}

	log.Printf("Secure file server started on https://localhost%s", port)
	log.Printf("Session will expire at: %s", time.Now().Add(shareDuration).Format(time.RFC1123))
	if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// openBrowser tries to open the URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string
	// Linux
	cmd = "xdg-open"
	args = []string{url}
	return exec.Command(cmd, args...).Start()
} 