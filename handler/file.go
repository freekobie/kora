package handler

import (
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/freekobie/kora/model"
	"github.com/freekobie/kora/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileUploader interface {
	UploadFile(ctx *gin.Context, user *model.User, folderID *uuid.UUID, file multipart.File, header *multipart.FileHeader) (*model.File, error)
}

func (h *Handler) CreateFolder(c *gin.Context) {
	var input service.CreateFolderRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	folder, err := h.file.CreateFolder(c.Request.Context(), input, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"folder": folder})
}

// FileUploadHandler handles file uploads.
func (h *Handler) FileUpload(c *gin.Context) {
	idString, ok := c.Get("user_id")
	if !ok {
		slog.Error("failed to fetch user id from context")
		c.JSON(http.StatusInternalServerError, gin.H{"message": ErrServerError.Error()})
		return
	}

	userId := uuid.MustParse(idString.(string))

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	folderIDStr := c.PostForm("folder_id")
	var folderId uuid.UUID
	if folderIDStr != "" {
		id, err := uuid.Parse(folderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder_id"})
			return
		}
		folderId = id
	}

	dbFile, err := h.file.UploadFile(c, userId, folderId, file, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	c.JSON(http.StatusOK, dbFile)
}
