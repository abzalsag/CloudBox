package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"CloudBox/internal/delivery/http/dto"
	"CloudBox/internal/delivery/http/middleware"
	"CloudBox/internal/delivery/http/response"
	"CloudBox/internal/entity"
	"CloudBox/internal/usecase/folder"
)

type FolderHandler struct {
	uc *folder.Usecase
}

func NewFolderHandler(uc *folder.Usecase) *FolderHandler {
	return &FolderHandler{uc: uc}
}

func (h *FolderHandler) Create(c *gin.Context) {
	var req dto.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	f, err := h.uc.Create(c.Request.Context(), middleware.UserID(c), req.ParentID, req.Name)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewFolderResponse(f))
}

func (h *FolderHandler) List(c *gin.Context) {
	parentID, err := optionalUUIDQuery(c, "parent_id")
	if err != nil {
		response.Error(c, err)
		return
	}
	limit, offset := pagination(c)

	items, err := h.uc.List(c.Request.Context(), middleware.UserID(c), parentID, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": dto.NewFolderListResponse(items)})
}

func (h *FolderHandler) Get(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	f, err := h.uc.Get(c.Request.Context(), middleware.UserID(c), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewFolderResponse(f))
}

func (h *FolderHandler) Rename(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.RenameFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	f, err := h.uc.Rename(c.Request.Context(), middleware.UserID(c), id, req.Name)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewFolderResponse(f))
}

func (h *FolderHandler) Move(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.MoveFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	f, err := h.uc.Move(c.Request.Context(), middleware.UserID(c), id, req.ParentID)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewFolderResponse(f))
}

func (h *FolderHandler) Delete(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.uc.Delete(c.Request.Context(), middleware.UserID(c), id); err != nil {
		response.Error(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func uuidParam(c *gin.Context, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: invalid %s", entity.ErrInvalidInput, name)
	}
	return id, nil
}

func optionalUUIDQuery(c *gin.Context, name string) (*uuid.UUID, error) {
	raw := c.Query(name)
	if raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid %s", entity.ErrInvalidInput, name)
	}
	return &id, nil
}

func pagination(c *gin.Context) (limit, offset int) {
	limit, offset = 50, 0
	if v, err := strconv.Atoi(c.Query("limit")); err == nil && v > 0 {
		limit = min(v, 100)
	}
	if v, err := strconv.Atoi(c.Query("offset")); err == nil && v > 0 {
		offset = v
	}
	return limit, offset
}
