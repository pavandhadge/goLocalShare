package main

import (
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	requestSizeLimit  = 10 << 20 // 10MB
	timeout           = 15 * time.Second
	authTokenLength   = 32
	authTokenValidity = 1 * time.Hour
	rateLimit         = 100 // requests per minute
	bruteForceDelay   = 2 * time.Second
)

var (
	basePath       string
	isDirMode      bool
	authTokens     = make(map[string]time.Time)
	rateLimiter    = make(map[string]int)
	lastFailedAuth = make(map[string]time.Time)
)

func main() {
	port := ":8080"
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("Usage: ./server <filepath> or ./server --dir <directory>")
	}
	if args[0] == "--dir" {
		if len(args) < 2 {
			log.Fatal("Missing directory path after --dir")
		}
		basePath = args[1]
		isDirMode = true
	} else {
		basePath = args[0]
		isDirMode = false
	}

	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		log.Fatalf("Path resolution error: %v", err)
	}
	basePath = absBasePath

	if _, err := secureStat(basePath); err != nil {
		log.Fatalf("Path error: %v", err)
	}

	srv := &http.Server{
		Addr:         port,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  2 * timeout,
		Handler:      nil,
	}

	http.Handle("/", securityHeaders(authCheck(rateLimitMiddleware(http.HandlerFunc(homeHandler)))))
	http.Handle("/browse", securityHeaders(authCheck(rateLimitMiddleware(http.HandlerFunc(browseHandler)))))
	http.Handle("/download/", securityHeaders(authCheck(rateLimitMiddleware(http.HandlerFunc(downloadHandler)))))

	go cleanupTokens()

	log.Printf("Secure file server started on http://localhost%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'none'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		if lastTry, exists := lastFailedAuth[ip]; exists {
			if time.Since(lastTry) < bruteForceDelay {
				http.Error(w, "Too many attempts", http.StatusTooManyRequests)
				return
			}
		}

		if rateLimiter[ip] >= rateLimit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		rateLimiter[ip]++

		time.AfterFunc(time.Minute, func() {
			rateLimiter[ip]--
		})

		next.ServeHTTP(w, r)
	})
}

func authCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		token := r.URL.Query().Get("token")
		if token == "" {
			token = r.Header.Get("X-Auth-Token")
		}

		if !validateToken(token) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			lastFailedAuth[ip] = time.Now()
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	token := generateSecureToken()
	authTokens[token] = time.Now().Add(authTokenValidity)

	tmpl := template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>Secure File Server</title>
	<style>
		body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
		.container { background-color: #f9f9f9; padding: 30px; border-radius: 8px; margin-top: 20px; }
		.token { word-break: break-all; background-color: #e6f7ff; padding: 10px; margin: 15px 0; }
		.btn { display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white;
			text-decoration: none; border-radius: 4px; margin-top: 10px; }
	</style>
</head>
<body>
	<h1>Secure File Access</h1>
	<div class="container">
		<p>Your secure access token:</p>
		<div class="token">{{.Token}}</div>
		<p>This token will expire in 1 hour.</p>

		{{if .IsDir}}
			<a href="/browse?token={{.Token}}" class="btn">Browse Files</a>
		{{else}}
			<a href="/download/{{.FileName}}?token={{.Token}}" class="btn">Download {{.FileName}}</a>
		{{end}}
	</div>
</body>
</html>
`))

	data := struct {
		Token    string
		IsDir    bool
		FileName string
	}{
		Token:    token,
		IsDir:    isDirMode,
		FileName: filepath.Base(basePath),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func browseHandler(w http.ResponseWriter, r *http.Request) {
	if !isDirMode {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.New("browse").Parse(`
		<!DOCTYPE html>
<html>
<head>
			<title>Browse Files</title>
			<style>
				:root {
					--primary-color: #4361ee;
					--secondary-color: #3f37c9;
					--background-color: #f8f9fa;
					--card-bg: #ffffff;
					--text-color: #333333;
					--text-light: #6c757d;
					--border-radius: 8px;
					--box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
					--transition: all 0.3s ease;
				}

				body {
					font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
					max-width: 900px;
					margin: 0 auto;
					padding: 20px;
					background-color: var(--background-color);
					color: var(--text-color);
					line-height: 1.6;
				}

				h1 {
					color: var(--primary-color);
					margin-bottom: 1rem;
					font-weight: 600;
				}

				.current-dir {
					background-color: var(--card-bg);
					padding: 12px 16px;
					border-radius: var(--border-radius);
					box-shadow: var(--box-shadow);
					margin-bottom: 20px;
					font-size: 0.95rem;
					border-left: 4px solid var(--primary-color);
				}

				.current-dir strong {
					color: var(--secondary-color);
					font-weight: 500;
				}

				.file-list {
					margin-top: 25px;
					background-color: var(--card-bg);
					border-radius: var(--border-radius);
					box-shadow: var(--box-shadow);
					overflow: hidden;
				}

				.file-item {
					padding: 14px 20px;
					display: flex;
					align-items: center;
					transition: var(--transition);
					border-bottom: 1px solid rgba(0, 0, 0, 0.05);
				}

				.file-item:last-child {
					border-bottom: none;
				}

				.file-item:hover {
					background-color: rgba(67, 97, 238, 0.05);
					transform: translateX(5px);
				}

				.file-icon {
					margin-right: 12px;
					font-size: 1.2rem;
					width: 24px;
					text-align: center;
				}

				.file-item a {
					color: var(--text-color);
					text-decoration: none;
					flex-grow: 1;
					transition: var(--transition);
				}

				.file-item a:hover {
					color: var(--primary-color);
				}

				.file-size {
					color: var(--text-light);
					font-size: 0.85rem;
					margin-left: auto;
					font-family: 'Courier New', monospace;
				}

				.back-link {
					display: inline-block;
					margin-top: 25px;
					padding: 10px 16px;
					background-color: var(--primary-color);
					color: white;
					text-decoration: none;
					border-radius: var(--border-radius);
					transition: var(--transition);
					font-weight: 500;
				}

				.back-link:hover {
					background-color: var(--secondary-color);
					transform: translateY(-2px);
					box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
				}

				.back-link::before {
					content: "←";
					margin-right: 5px;
				}

				.empty-state {
					text-align: center;
					padding: 40px 20px;
					color: var(--text-light);
				}

				/* Animation for empty state */
				@keyframes fadeIn {
					from { opacity: 0; transform: translateY(10px); }
					to { opacity: 1; transform: translateY(0); }
				}

				.file-list {
					animation: fadeIn 0.4s ease-out;
				}

				/* Responsive design */
				@media (max-width: 768px) {
					body {
						padding: 15px;
					}

					.file-item {
						padding: 12px 15px;
					}
				}
			</style>
</head>
<body>
			<h1>Available Files</h1>
			<div class="current-dir">
				Current directory: <strong>{{.CurrentDir}}</strong>
			</div>

			<div class="file-list">
				{{if .Files}}
					{{range .Files}}
					<div class="file-item">
						<span class="file-icon">{{if .IsDir}}📁{{else}}📄{{end}}</span>
						<a href="{{if .IsDir}}/browse/{{.Name}}?token={{$.Token}}{{else}}/download/{{.Name}}?token={{$.Token}}{{end}}">
							{{.Name}}{{if .IsDir}}/{{end}}
						</a>
						{{if not .IsDir}}<span class="file-size">{{.Size}}</span>{{end}}
					</div>
					{{end}}
				{{else}}
					<div class="empty-state">
						This directory is empty
					</div>
				{{end}}
			</div>

			<a href="/?token={{.Token}}" class="back-link">Back to Home</a>
</body>
</html>
`))

	currentDir := basePath
	if strings.HasPrefix(r.URL.Path, "/browse/") {
		relPath := strings.TrimPrefix(r.URL.Path, "/browse/")
		currentDir = filepath.Join(basePath, relPath)
	}

	files, err := os.ReadDir(currentDir)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	type fileInfo struct {
		Name  string
		IsDir bool
		Size  string
	}

	var fileList []fileInfo
	for _, file := range files {
		info, _ := file.Info()
		size := ""
		if !file.IsDir() {
			size = fmt.Sprintf("%.2f MB", float64(info.Size())/1024/1024)
		}
		fileList = append(fileList, fileInfo{
			Name:  file.Name(),
			IsDir: file.IsDir(),
			Size:  size,
		})
	}

	token := r.URL.Query().Get("token")
	data := struct {
		CurrentDir string
		Files      []fileInfo
		Token      string
	}{
		CurrentDir: currentDir,
		Files:      fileList,
		Token:      token,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	requestedPath := strings.TrimPrefix(r.URL.Path, "/download/")
	cleanPath, err := securePath(basePath, requestedPath)
	if err != nil {
		log.Printf("Path traversal attempt from %s: %v", r.RemoteAddr, err)
		http.Error(w, "Invalid path", http.StatusForbidden)
		return
	}

	fileInfo, err := secureStat(cleanPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if fileInfo.IsDir() {
		token := r.URL.Query().Get("token")
		http.Redirect(w, r, "/browse/"+requestedPath+"?token="+token, http.StatusSeeOther)
		return
	}

	file, err := secureOpen(cleanPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(cleanPath)))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("Cache-Control", "no-store")

	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func securePath(base, requested string) (string, error) {
	if strings.Contains(requested, "../") || strings.Contains(requested, "~/") ||
		strings.Contains(requested, "..\\") || strings.Contains(requested, "\\") {
		return "", errors.New("path traversal attempt")
	}

	joined := filepath.Join(base, requested)
	absPath, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(absPath, base) {
		return "", errors.New("path outside base directory")
	}

	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(realPath, base) {
		return "", errors.New("symlink points outside base directory")
	}

	return realPath, nil
}

func secureStat(path string) (os.FileInfo, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}

	if !fi.Mode().IsRegular() && !fi.Mode().IsDir() {
		return nil, errors.New("special files not allowed")
	}

	return fi, nil
}

func secureOpen(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|0x20000, 0)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		file.Close()
		return nil, errors.New("symlinks not allowed")
	}

	return file, nil
}

func generateSecureToken() string {
	b := make([]byte, authTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("Failed to generate token:", err)
	}
	return fmt.Sprintf("%x", b)
}

func validateToken(token string) bool {
	for validToken, expiry := range authTokens {
		if subtle.ConstantTimeCompare([]byte(token), []byte(validToken)) == 1 {
			if time.Now().Before(expiry) {
				return true
			}
			delete(authTokens, validToken)
			return false
		}
	}
	return false
}

func cleanupTokens() {
	for range time.Tick(5 * time.Minute) {
		for token, expiry := range authTokens {
			if time.Now().After(expiry) {
				delete(authTokens, token)
			}
		}
	}
}
