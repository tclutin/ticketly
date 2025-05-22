package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ticketly/ticketly_api/internal/core"
	coreerrors "github.com/tclutin/ticketly/ticketly_api/internal/core/errors"
	"github.com/tclutin/ticketly/ticketly_api/internal/core/user"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1/request"
	"log/slog"
	"net/http"
)

type UserHandler struct {
	service core.UserService
}

func NewUserHandler(service core.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Bind(router *gin.RouterGroup) {
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("", h.Register)
		usersGroup.GET("/:external_id", h.GetByExternalId)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	err := h.service.Create(c.Request.Context(), user.RegisterUserDTO{
		ExternalID: req.ExternalID,
		Username:   req.Username,
		Source:     req.Source,
	})

	if err != nil {
		if errors.Is(err, coreerrors.ErrUserAlreadyExists) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		slog.Error("error fuck", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "successfully",
	})
}

func (h *UserHandler) GetByExternalId(c *gin.Context) {
	externalID := c.Param("external_id")

	usr, err := h.service.GetByExternalId(c.Request.Context(), externalID)
	if err != nil {
		if errors.Is(err, coreerrors.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		slog.Error("error fuck", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(200, usr)
}
