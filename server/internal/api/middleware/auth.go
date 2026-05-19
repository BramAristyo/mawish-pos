package middleware

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"strings"

	"github.com/BramAristyo/mawish-pos/server/internal/constant"
	"github.com/BramAristyo/mawish-pos/server/internal/infrastructure/config"
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader(constant.AuthHeader)
		token := strings.Split(auth, " ")

		if auth == "" || len(token) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}

		at, err := jwt.Parse(token[1], func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, usecase_errors.TokenInvalid
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, ok := at.Claims.(jwt.MapClaims)
		if !ok || !at.Valid {
			fmt.Println("Failed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claimMap := make(map[string]any)
		maps.Copy(claimMap, claims)

		userID := claimMap[constant.ClaimUserID]
		userRole := claimMap[constant.ClaimRole]

		c.Set(constant.ClaimUserID, userID)
		c.Set(constant.ClaimRole, userRole)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, constant.CtxUserID, userID)
		ctx = context.WithValue(ctx, constant.CtxRole, userRole)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
