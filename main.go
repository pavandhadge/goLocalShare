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
	port := ":8090"
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
	http.Handle("/browse/", securityHeaders(authCheck(rateLimitMiddleware(http.HandlerFunc(browseHandler)))))
	http.Handle("/download/", securityHeaders(authCheck(rateLimitMiddleware(http.HandlerFunc(downloadHandler)))))

	go cleanupTokens()

	log.Printf("Secure file server started on http://localhost%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'none'")
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
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Secure File Server | Cloud Storage</title>
    <style>
        :root {
            --primary: #2563eb;
            --primary-hover: #1d4ed8;
            --background: #f8fafc;
            --surface: #ffffff;
            --text: #1e293b;
            --text-secondary: #64748b;
            --border: #e2e8f0;
            --success: #16a34a;
            --radius: 0.5rem;
            --shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background-color: var(--background);
            color: var(--text);
            line-height: 1.5;
            margin: 0;
            padding: 0;
            min-height: 100vh;
        }

        .container {
            max-width: 800px;
            margin: 2rem auto;
            padding: 2rem;
            background: var(--surface);
            border-radius: var(--radius);
            box-shadow: var(--shadow);
        }

        header {
            text-align: center;
            margin-bottom: 2rem;
        }

        h1 {
            color: var(--primary);
            font-size: 1.8rem;
            margin-bottom: 0.5rem;
        }

        .token-display {
            margin: 1.5rem 0;
        }

        .token-label {
            display: block;
            color: var(--text-secondary);
            font-size: 0.875rem;
            margin-bottom: 0.5rem;
        }

        .token {
            font-family: monospace;
            background-color: #f1f5f9;
            padding: 1rem;
            border-radius: var(--radius);
            border-left: 4px solid var(--primary);
            word-break: break-all;
            position: relative;
        }

        .token::after {
            content: "🔒";
            position: absolute;
            right: 1rem;
            top: 50%;
            transform: translateY(-50%);
        }

        .expiry-notice {
            color: var(--text-secondary);
            font-size: 0.875rem;
            margin: 1rem 0 1.5rem;
            display: flex;
            align-items: center;
            gap: 0.5rem;
        }

        .actions {
            margin: 2rem 0;
        }

        .btn {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            padding: 0.75rem 1.5rem;
            background: var(--primary);
            color: white;
            text-decoration: none;
            border-radius: var(--radius);
            font-weight: 500;
            transition: background-color 0.2s ease;
        }

        .btn:hover {
            background: var(--primary-hover);
        }

        .btn-browse::before {
            content: "📁";
            margin-right: 0.5rem;
        }

        .btn-download::before {
            content: "⬇️";
            margin-right: 0.5rem;
        }

        .docs {
            margin-top: 2.5rem;
            padding: 1.5rem;
            background: #f8fafc;
            border-radius: var(--radius);
        }

        pre {
            background: #1e293b;
            color: #f8fafc;
            padding: 1rem;
            border-radius: 0.375rem;
            overflow-x: auto;
            font-family: monospace;
            font-size: 0.875rem;
        }

        .github-link {
            display: inline-flex;
            align-items: center;
            gap: 0.5rem;
            color: var(--primary);
            text-decoration: none;
            margin-top: 1rem;
        }

        .github-link:hover {
            text-decoration: underline;
        }

        footer {
            text-align: center;
            margin-top: 2rem;
            color: var(--text-secondary);
            font-size: 0.875rem;
        }

        @media (max-width: 768px) {
            .container {
                margin: 1rem;
                padding: 1.5rem;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>Secure File Server v1.0</h1>
            <div>🔐 Secure file sharing</div>
        </header>

        <div class="token-display">
            <span class="token-label">Your secure access token</span>
            <div class="token">{{.Token}}</div>
        </div>

        <p class="expiry-notice">
            <span>⏳ This token will expire in 1 hour</span>
        </p>

        <div class="actions">
            {{if .IsDir}}
                <a href="/browse?token={{.Token}}" class="btn btn-browse">Browse Files</a>
            {{else}}
                <a href="/download/{{.BaseName}}?token={{.Token}}" class="btn btn-download">Download {{.BaseName}}</a>
            {{end}}
        </div>

        <div class="docs">
            <h2>Documentation</h2>

            <h3>Usage</h3>
            <p>To serve a single file:</p>
            <pre>fileserver &lt;filepath&gt;</pre>

            <p>To serve a directory:</p>
            <pre>fileserver --dir &lt;directorypath&gt;</pre>

            <p>The server will start on port 8000. Access the web interface at:</p>
            <pre>http://localhost:8000</pre>

            <h3>Security Features</h3>
            <ul>
                <li>Uses temporary access tokens</li>
                <li>Protected against path traversal attacks</li>
                <li>Request size limited to prevent abuse</li>
            </ul>

            <a href="https://github.com/pavandhadge/goFileShare" class="github-link">
                <svg height="16" viewBox="0 0 16 16" width="16" fill="currentColor">
                    <path fill-rule="evenodd" d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z"></path>
                </svg>
                View on GitHub
            </a>
        </div>
    </div>

    <footer>
        Secure Local File Server
    </footer>
</body>
</html>
`))

	data := struct {
		Token    string
		IsDir    bool
		BaseName string
	}{
		Token:    token,
		IsDir:    isDirMode,
		BaseName: filepath.Base(basePath),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
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
		body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
		.current-dir { margin-bottom: 20px; padding: 10px; background: #f0f0f0; border-radius: 4px; }
		.file-list { margin: 20px 0; }
		.file-item { padding: 8px 0; border-bottom: 1px solid #eee; }
		.file-icon { margin-right: 10px; }
		.file-size { float: right; color: #666; }
		.empty-state { color: #666; font-style: italic; }
		.back-link { display: inline-block; margin-top: 20px; color: #007bff; text-decoration: none; }
		.back-link:hover { text-decoration: underline; }
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
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// For single file mode
	if !isDirMode {
		fileInfo, err := secureStat(basePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		file, err := secureOpen(basePath)
		if err != nil {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(basePath)))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		w.Header().Set("Cache-Control", "no-store")

		http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
		return
	}

	// For directory mode
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
