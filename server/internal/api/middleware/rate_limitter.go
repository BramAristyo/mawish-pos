package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/BramAristyo/mawish-pos/server/pkg/logger"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter(log *logger.ZapLogger) gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Background cleanup to prevent memory leaks
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			mu.Lock()
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, exists := clients[ip]; !exists {
			// Allow 10 requests per second with a burst of 20
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(10), 20)}
		}
		cl := clients[ip]
		cl.lastSeen = time.Now()
		mu.Unlock()

		if !cl.limiter.Allow() {
			log.Warn("Rate limit exceeded", "ip", ip, "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
