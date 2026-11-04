package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"CloudBox/internal/delivery/http/dto"
	"CloudBox/internal/delivery/http/middleware"
	"CloudBox/internal/delivery/http/response"
	"CloudBox/internal/usecase/auth"
)

type AuthHandler struct {
	uc *auth.Usecase
}

func NewAuthHandler(uc *auth.Usecase) *AuthHandler { return &AuthHandler{uc: uc} }

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	user, err := h.uc.Register(c.Request.Context(), auth.RegisterInput{
		Email: req.Email, Password: req.Password, FullName: req.FullName,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewUserResponse(user))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	t, err := h.uc.Login(c.Request.Context(), auth.LoginInput{
		Email: req.Email, Password: req.Password,
	}, requestMeta(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, tokensResponse(t))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	t, err := h.uc.Refresh(c.Request.Context(), req.RefreshToken, requestMeta(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, tokensResponse(t))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	if err := h.uc.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.uc.GetProfile(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

func requestMeta(c *gin.Context) auth.Meta {
	return auth.Meta{UserAgent: c.Request.UserAgent(), IP: c.ClientIP()}
}

func tokensResponse(t *auth.Tokens) dto.TokensResponse {
	return dto.TokensResponse{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    t.ExpiresIn,
	}
}
