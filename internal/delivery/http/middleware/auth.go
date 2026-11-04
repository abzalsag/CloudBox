package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"CloudBox/internal/delivery/http/response"
	"CloudBox/internal/entity"
)

type TokenParser interface {
	Parse(raw string) (uuid.UUID, error)
}

const userIDKey = "user_id"

func Auth(p TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		scheme, tok, ok := strings.Cut(c.GetHeader("Authorization"), " ")
		if !ok || !strings.EqualFold(scheme, "Bearer") || tok == "" {
			response.Error(c, entity.ErrUnauthorized)
			return
		}
		id, err := p.Parse(tok)
		if err != nil {
			response.Error(c, entity.ErrUnauthorized)
			return
		}
		c.Set(userIDKey, id)
		c.Next()
	}
}

func UserID(c *gin.Context) uuid.UUID {
	return c.MustGet(userIDKey).(uuid.UUID)
}
