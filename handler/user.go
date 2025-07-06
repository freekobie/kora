package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/freekobie/kora/model"
	"github.com/freekobie/kora/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateUser godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		object	true	"User registration info"
//	@Success		201		{object}	UserResponse
//	@Failure		400		{object}	Response
//	@Failure		409		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/auth/register [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=8,max=20"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	user, err := h.user.CreateUser(c.Request.Context(), input.Name, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateUser) {
			c.JSON(http.StatusConflict, Response{Status: http.StatusConflict, Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{
			Message: "internal server error",
			Status:  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, UserResponse{Status: http.StatusCreated, User: *user})
}

// VerifyUser godoc
//
//	@Summary		Verify user email
//	@Description	Verify a user's email with a code
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			verification	body		object	true	"Verification info"
//	@Success		200				{object}	UserResponse
//	@Failure		400				{object}	Response
//	@Failure		500				{object}	Response
//	@Router			/auth/verify [post]
func (h *Handler) VerifyUser(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	user, err := h.user.VerifyUser(c.Request.Context(), input.Code, input.Email)
	if err != nil {
		if errors.Is(err, service.ErrInvalidToken) {
			c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{Status: http.StatusOK, User: *user})
}

// RequestVerificationCode godoc
//
//	@Summary		Request verification email
//	@Description	Request a new verification code for a user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			email	body		object	true	"User email"
//	@Success		202		{object}	Response
//	@Failure		400		{object}	Response
//	@Failure		404		{object}	Response
//	@Failure		422		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/auth/verify/request [post]
func (h *Handler) RequestVerificationCode(c *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	err := h.user.ResendVerificationEmail(c.Request.Context(), input.Email)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Message: err.Error()})
			return
		} else if strings.Contains(err.Error(), "user already verified") {
			c.JSON(http.StatusUnprocessableEntity, Response{Status: http.StatusUnprocessableEntity, Message: err.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: ErrServerError.Error()})
		return
	}

	c.JSON(http.StatusAccepted, Response{Status: http.StatusAccepted, Message: fmt.Sprintf("new verification code has ben sent to '%s'", input.Email)})

}

// LoginUser godoc
//
//	@Summary		Login user
//	@Description	Authenticate user and return session tokens
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		object	true	"User credentials"
//	@Success		200			{object}	SessionResponse
//	@Failure		400			{object}	Response
//	@Failure		401			{object}	Response
//	@Failure		500			{object}	Response
//	@Router			/auth/login [post]
func (h *Handler) LoginUser(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	session, err := h.user.NewSession(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrFailedOperation) {
			c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: ErrServerError.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, Response{Status: http.StatusUnauthorized, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SessionResponse{Status: http.StatusOK, Session: *session})
}

// GetUserAccessToken godoc
//
//	@Summary		Refresh access token
//	@Description	Get a new access token using a refresh token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			refreshToken	body		object	true	"Refresh token"
//	@Success		200				{object}	AccessResponse
//	@Failure		400				{object}	Response
//	@Failure		401				{object}	Response
//	@Failure		500				{object}	Response
//	@Router			/auth/access [post]
func (h *Handler) GetUserAccessToken(c *gin.Context) {

	var input struct {
		RefresToken string `json:"refreshToken" binding:"required,jwt"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	access, err := h.user.RefreshSession(c.Request.Context(), input.RefresToken)
	if err != nil {
		if errors.Is(err, service.ErrFailedOperation) {
			c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: ErrServerError.Error()})
			return
		}
		c.JSON(http.StatusUnauthorized, Response{Status: http.StatusUnauthorized, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, AccessResponse{Status: http.StatusOK, Access: *access})
}

// GetUser godoc
//
//	@Summary		Get user by ID
//	@Description	Get user details by user ID
//	@Tags			users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	UserResponse
//	@Failure		400	{object}	Response
//	@Failure		404	{object}	Response
//	@Failure		500	{object}	Response
//	@Router			/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	userId, err := getUUIDparam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "invalid id"})
		return
	}

	user, err := h.user.FetchUser(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{Status: http.StatusOK, User: *user})
}

// UpdateUserData godoc
//
//	@Summary		Update user data
//	@Description	Update user profile information
//	@Tags			users
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		object	true	"User update info"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	Response
//	@Failure		422		{object}	Response
//	@Failure		500		{object}	Response
//	@Router			/users/profile [patch]
func (h *Handler) UpdateUserData(c *gin.Context) {
	var input map[string]any
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}

	idString, ok := c.Get("user_id")
	if !ok {
		slog.Error("failed to fetch user id from context")
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: ErrServerError.Error()})
		return
	}

	id := uuid.MustParse(idString.(string))
	input["id"] = id

	user, err := h.user.UpdateUser(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrFailedOperation) {
			c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: ErrServerError.Error()})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, Response{Status: http.StatusUnprocessableEntity, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, UserResponse{Status: http.StatusOK, User: *user})
}

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	Response
//	@Failure		400	{object}	Response
//	@Failure		404	{object}	Response
//	@Failure		500	{object}	Response
//	@Router			/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	userId, err := getUUIDparam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Message: "invalid id"})
		return
	}

	err = h.user.DeleteUser(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Message: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Message: "internal server error"})
		return
	}

	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Message: "user deleted successfully"})
}
