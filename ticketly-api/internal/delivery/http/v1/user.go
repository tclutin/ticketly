package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ticketly/ticketly_api/internal/core"
	"github.com/tclutin/ticketly/ticketly_api/internal/core/user"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1/request"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1/response"
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

	if err := h.service.Create(c.Request.Context(), user.RegisterUserDTO{
		ExternalID: req.ExternalID,
		Username:   req.Username,
		Source:     req.Source,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "created",
	})
}

func (h *UserHandler) GetByExternalId(c *gin.Context) {
	externalID := c.Param("external_id")

	usr, err := h.service.GetByExternalId(c.Request.Context(), externalID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.User{
		UserID:     usr.UserID,
		ExternalID: usr.ExternalID,
		Username:   usr.Username,
		Source:     usr.Source,
		IsBanned:   usr.IsBanned,
		CreatedAt:  usr.CreatedAt,
	})
}
