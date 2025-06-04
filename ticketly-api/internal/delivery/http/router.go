package http

import (
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/middleware"
	v1 "github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1"
	"github.com/tclutin/ticketly/ticketly_api/internal/service"
)

func InitRouter(
	userSrv service.UserService,
	ticketSrv service.TicketService,
	operatorSrv service.OperatorService,
	casdoorClient *casdoorsdk.Client,
) *gin.Engine {
	router := gin.Default()

	usersHandler := v1.NewUserHandler(userSrv)
	ticketsHandler := v1.NewTicketHandler(ticketSrv)
	operatorsHandler := v1.NewOperatorHandler(operatorSrv, ticketSrv)
	authHandler := v1.NewAuthHandler(operatorSrv, casdoorClient)

	router.Use(
		gin.Recovery(),
		middleware.CORSMiddleware(),
		middleware.ErrorMiddleware(),
	)

	apiGroup := router.Group("/api")
	{
		v1Group := apiGroup.Group("/v1")
		{
			authHandler.Bind(v1Group)
			usersHandler.Bind(v1Group)
			operatorsHandler.Bind(v1Group, casdoorClient)
			ticketsHandler.Bind(v1Group, operatorSrv, casdoorClient)
		}
	}

	return router
}
