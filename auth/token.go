package auth

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	ownerToken      string
	tokenExpiry     time.Time
	mu              sync.RWMutex
)

func GenerateSecureToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("Failed to generate token:", err)
	}
	return fmt.Sprintf("%x", b)
}

// Call this on server start or reset
func InitOwnerToken(duration time.Duration) {
	mu.Lock()
	defer mu.Unlock()
	ownerToken = GenerateSecureToken()
	tokenExpiry = time.Now().Add(duration)
}

func GetOwnerToken() (string, time.Time) {
	mu.RLock()
	defer mu.RUnlock()
	return ownerToken, tokenExpiry
}

func ValidateToken(token string) bool {
	mu.RLock()
	defer mu.RUnlock()
	if token == "" || ownerToken == "" {
		return false
	}
	if token != ownerToken {
		return false
	}
	if time.Now().After(tokenExpiry) {
		return false
	}
	return true
}

func IsSessionExpired() bool {
	mu.RLock()
	defer mu.RUnlock()
	return time.Now().After(tokenExpiry)
}

// Optional: Reset/rotate the token
func ResetOwnerToken(duration time.Duration) {
	InitOwnerToken(duration)
} 