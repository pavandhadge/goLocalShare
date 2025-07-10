package server

import (
	"log"
	"net"
	"net/http"
	"strings"
	"time"
	"github.com/pavandhadge/goFileShare/auth"
)

const (
	rateLimit       = 100 // requests per minute
	bruteForceDelay = 2 * time.Second
)

var (
	rateLimiter    = make(map[string]int)
	lastFailedAuth = make(map[string]time.Time)
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Auth-Token")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline';")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(next http.Handler) http.Handler {
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

func AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow /token and / (landing) without auth
		if r.URL.Path == "/token" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("X-Auth-Token")
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		if !auth.ValidateToken(token) {
			logMsg := "[AuthCheck] Invalid or missing token for path: " + r.URL.Path
			if strings.HasPrefix(r.URL.Path, "/api/") {
				log.Println(logMsg + " - returning 401 JSON")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"unauthorized","message":"Invalid or expired token","code":401}`))
				return
			}
			accept := r.Header.Get("Accept")
			if strings.Contains(accept, "text/html") {
				log.Println(logMsg + " - redirecting to /token")
				http.Redirect(w, r, "/token", http.StatusSeeOther)
				return
			}
			log.Println(logMsg + " - returning 401 JSON (default)")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"unauthorized","message":"Invalid or expired token","code":401}`))
			return
		}

		next.ServeHTTP(w, r)
	})
} 