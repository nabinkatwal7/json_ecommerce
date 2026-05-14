package api

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	mu sync.Mutex
	m  map[string]*rate.Limiter
	lim rate.Limit
	burst int
}

func newIPLimiter(rps float64, burst int) *ipLimiter {
	if rps <= 0 {
		rps = 20
	}
	if burst <= 0 {
		burst = 40
	}
	return &ipLimiter{
		m:     make(map[string]*rate.Limiter),
		lim:   rate.Limit(rps),
		burst: burst,
	}
}

func (l *ipLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	lim, ok := l.m[ip]
	if !ok {
		lim = rate.NewLimiter(l.lim, l.burst)
		l.m[ip] = lim
	}
	return lim.Allow()
}

func clientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return c.ClientIP()
}

func (l *ipLimiter) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !l.allow(clientIP(c)) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
