package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"CloudBox/internal/entity"
)

func Error(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, entity.ErrInvalidInput):
		status, msg = http.StatusBadRequest, err.Error()
	case errors.Is(err, entity.ErrUnauthorized):
		status, msg = http.StatusUnauthorized, "unauthorized"
	case errors.Is(err, entity.ErrForbidden):
		status, msg = http.StatusForbidden, "forbidden"
	case errors.Is(err, entity.ErrNotFound):
		status, msg = http.StatusNotFound, "not found"
	case errors.Is(err, entity.ErrAlreadyExists):
		status, msg = http.StatusConflict, "already exists"
	case errors.Is(err, entity.ErrQuotaExceeded):
		status, msg = http.StatusRequestEntityTooLarge, "storage quota exceeded"
	default:
		_ = c.Error(err)
	}

	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}
