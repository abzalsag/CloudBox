package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"CloudBox/internal/delivery/http/dto"
	"CloudBox/internal/delivery/http/middleware"
	"CloudBox/internal/delivery/http/response"
	"CloudBox/internal/usecase/trash"
)

type TrashHandler struct {
	uc *trash.Usecase
}

func NewTrashHandler(uc *trash.Usecase) *TrashHandler {
	return &TrashHandler{uc: uc}
}

func (h *TrashHandler) List(c *gin.Context) {
	l, err := h.uc.List(c.Request.Context(), middleware.UserID(c))
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewTrashResponse(l.Folders, l.Files))
}

func (h *TrashHandler) RestoreFile(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	f, err := h.uc.RestoreFile(c.Request.Context(), middleware.UserID(c), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewFileResponse(f))
}

func (h *TrashHandler) RestoreFolder(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	f, err := h.uc.RestoreFolder(c.Request.Context(), middleware.UserID(c), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewFolderResponse(f))
}

func (h *TrashHandler) PurgeFile(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.uc.PurgeFile(c.Request.Context(), middleware.UserID(c), id); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TrashHandler) PurgeFolder(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.uc.PurgeFolder(c.Request.Context(), middleware.UserID(c), id); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
