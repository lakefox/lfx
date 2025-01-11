package ipban

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// BanManager manages banned IPs
type BanManager struct {
	bannedIPs   map[string]time.Time
	mu          sync.RWMutex
	banDuration time.Duration
}

var bm *BanManager

// NewBanManager creates a new BanManager with a specified ban duration
func NewBanManager(banDuration time.Duration) {
	bm = &BanManager{
		bannedIPs:   make(map[string]time.Time),
		banDuration: banDuration,
	}
}

// BanIP bans an IP address for the configured duration
func Ban(ip string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.bannedIPs[ip] = time.Now().Add(bm.banDuration)
}

// IsBanned checks if an IP is currently banned
func IsBanned(ip string) bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	expiry, exists := bm.bannedIPs[ip]
	if !exists {
		return false
	}
	if time.Now().After(expiry) {
		removeIP(ip)
		return false
	}
	return true
}

// Middleware to block banned IPs
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetIP(r)
		if IsBanned(ip) {
			http.Error(w, "Your IP is banned.", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// getIP extracts the IP address from the HTTP request
func GetIP(r *http.Request) string {
	// Check for IP in X-Forwarded-For header
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	// Check for IP in the X-Real-IP header
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// Default to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// removeIP removes an IP from the banned list
func removeIP(ip string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	delete(bm.bannedIPs, ip)
}
