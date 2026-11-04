package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"CloudBox/internal/delivery/http/response"
	"CloudBox/internal/entity"
)

func badBody(c *gin.Context) {
	response.Error(c, fmt.Errorf("%w: invalid request body", entity.ErrInvalidInput))
}
