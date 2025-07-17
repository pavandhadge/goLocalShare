package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
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
	certFile := "cert.pem"
	keyFile := "key.pem"
	if _, err := os.Stat(certFile); os.IsNotExist(err) || os.IsNotExist(func() error { _, err := os.Stat(keyFile); return err }()) {
		log.Printf("cert.pem or key.pem not found, generating self-signed certificate...")
		err := generateSelfSignedCert(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed to generate self-signed certificate: %v", err)
		}
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

// generateSelfSignedCert creates a self-signed certificate and key for localhost usage.
func generateSelfSignedCert(certFile, keyFile string) error {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	now := time.Now()
	tmpl := x509.Certificate{
		SerialNumber: big.NewInt(now.UnixNano()),
		Subject: pkix.Name{
			Organization: []string{"goFileServer"},
			CommonName:   "localhost",
		},
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(365 * 24 * time.Hour), // valid for 1 year
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &priv.PublicKey, priv)
	if err != nil {
		return err
	}
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return err
	}
	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	privBytes := x509.MarshalPKCS1PrivateKey(priv)
	if err := pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}); err != nil {
		return err
	}
	return nil
} 