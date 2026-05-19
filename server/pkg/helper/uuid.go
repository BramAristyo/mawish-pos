package helper

import (
	"github.com/BramAristyo/mawish-pos/server/pkg/usecase_errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ParseUUID(c *gin.Context) (uuid.UUID, error) {
	id := c.Param("id")
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, usecase_errors.InvalidID
	}
	return parsed, nil
}
