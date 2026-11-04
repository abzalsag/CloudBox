package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"CloudBox/internal/delivery/http/handler"
	"CloudBox/internal/delivery/http/middleware"
)

func NewRouter(
	l *zap.Logger,
	authH *handler.AuthHandler,
	folderH *handler.FolderHandler,
	fileH *handler.FileHandler,
	tp middleware.TokenParser,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), middleware.Logger(l))

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.Group("/api/v1")

	authGroup := api.Group("/auth")
	authGroup.POST("/register", authH.Register)
	authGroup.POST("/login", authH.Login)
	authGroup.POST("/refresh", authH.Refresh)
	authGroup.POST("/logout", authH.Logout)

	protected := api.Group("", middleware.Auth(tp))
	protected.GET("/me", authH.Me)

	folders := protected.Group("/folders")
	folders.POST("", folderH.Create)
	folders.GET("", folderH.List)
	folders.GET("/:id", folderH.Get)
	folders.PATCH("/:id", folderH.Rename)
	folders.PATCH("/:id/move", folderH.Move)
	folders.DELETE("/:id", folderH.Delete)

	files := protected.Group("/files")
	files.POST("", fileH.Upload)
	files.GET("", fileH.List)
	files.GET("/:id", fileH.Get)
	files.GET("/:id/download", fileH.Download)
	files.PATCH("/:id", fileH.Rename)
	files.PATCH("/:id/move", fileH.Move)
	files.DELETE("/:id", fileH.Delete)

	return r
}
