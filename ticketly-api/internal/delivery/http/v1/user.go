package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1/request"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1/response"
	"github.com/tclutin/ticketly/ticketly_api/internal/service"
	"github.com/tclutin/ticketly/ticketly_api/internal/service/user"
	"net/http"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Bind в будущем прикрутить casdoor middleware
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

	userId, err := h.service.Create(c.Request.Context(), user.RegisterUserDTO{
		ExternalID: req.ExternalID,
		Username:   req.Username,
		Source:     req.Source,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id": userId,
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
