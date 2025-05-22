package http

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ticketly/ticketly_api/internal/core"
	v1 "github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1"
)

func InitRouter(userSrv core.UserService) *gin.Engine {
	router := gin.Default()

	usersHandler := v1.NewUserHandler(userSrv)
	ticketsHandler := v1.NewTicketHandler()

	apiGroup := router.Group("/api/v1/")
	{
		usersHandler.Bind(apiGroup)
		ticketsHandler.Bind(apiGroup)
	}

	return router
}
