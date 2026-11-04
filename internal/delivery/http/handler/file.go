package handler

import (
	"fmt"
	"mime"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"CloudBox/internal/delivery/http/dto"
	"CloudBox/internal/delivery/http/middleware"
	"CloudBox/internal/delivery/http/response"
	"CloudBox/internal/entity"
	"CloudBox/internal/usecase/file"
)

const maxUploadSize = 500 << 20

type FileHandler struct {
	uc *file.Usecase
}

func NewFileHandler(uc *file.Usecase) *FileHandler {
	return &FileHandler{uc: uc}
}

func (h *FileHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	header, err := c.FormFile("file")
	if err != nil {
		response.Error(c, fmt.Errorf("%w: file is required or too large", entity.ErrInvalidInput))
		return
	}

	folderID, err := optionalUUIDForm(c, "folder_id")
	if err != nil {
		response.Error(c, err)
		return
	}

	content, err := header.Open()
	if err != nil {
		response.Error(c, err)
		return
	}
	defer content.Close()

	f, err := h.uc.Upload(c.Request.Context(), file.UploadInput{
		OwnerID:     middleware.UserID(c),
		FolderID:    folderID,
		Name:        header.Filename,
		ContentType: header.Header.Get("Content-Type"),
		Content:     content,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusCreated, dto.NewFileResponse(f))
}

func (h *FileHandler) List(c *gin.Context) {
	folderID, err := optionalUUIDQuery(c, "folder_id")
	if err != nil {
		response.Error(c, err)
		return
	}
	limit, offset := pagination(c)

	items, err := h.uc.List(c.Request.Context(), middleware.UserID(c), folderID, limit, offset)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": dto.NewFileListResponse(items)})
}

func (h *FileHandler) Get(c *gin.Context) {
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
	c.JSON(http.StatusOK, dto.NewFileResponse(f))
}

func (h *FileHandler) Download(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	f, body, err := h.uc.Download(c.Request.Context(), middleware.UserID(c), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	defer body.Close()

	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": f.Name})
	c.DataFromReader(http.StatusOK, f.Size, f.MimeType, body, map[string]string{
		"Content-Disposition": disposition,
	})
}

func (h *FileHandler) Rename(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.RenameFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	f, err := h.uc.Rename(c.Request.Context(), middleware.UserID(c), id, req.Name)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewFileResponse(f))
}

func (h *FileHandler) Move(c *gin.Context) {
	id, err := uuidParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.MoveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badBody(c)
		return
	}
	f, err := h.uc.Move(c.Request.Context(), middleware.UserID(c), id, req.FolderID)
	if err != nil {
		response.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.NewFileResponse(f))
}

func (h *FileHandler) Delete(c *gin.Context) {
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

func optionalUUIDForm(c *gin.Context, name string) (*uuid.UUID, error) {
	raw := c.PostForm(name)
	if raw == "" {
		return nil, nil
	}
	id, err := uuid.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid %s", entity.ErrInvalidInput, name)
	}
	return &id, nil
}
