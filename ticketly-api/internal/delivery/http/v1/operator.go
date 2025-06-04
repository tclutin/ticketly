package v1

import (
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/middleware"
	"github.com/tclutin/ticketly/ticketly_api/internal/service"
	"net/http"
)

type OperatorHandler struct {
	service       service.OperatorService
	ticketService service.TicketService
}

func NewOperatorHandler(service service.OperatorService, ticketService service.TicketService) *OperatorHandler {
	return &OperatorHandler{
		service:       service,
		ticketService: ticketService,
	}
}

func (o *OperatorHandler) Bind(router *gin.RouterGroup, client *casdoorsdk.Client) {
	operatorsGroup := router.Group("/operators", middleware.AuthMiddleware(o.service, client))
	{
		ticketGroup := operatorsGroup.Group("/tickets")
		{
			ticketGroup.GET("/connections", o.GetActiveConnections)
		}
	}
}

func (o *OperatorHandler) GetActiveConnections(c *gin.Context) {
	operatorId, ok := c.Get("operator_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "operator_id not found in context"})
		return
	}

	connections, err := o.ticketService.GetActiveConnections(c.Request.Context(), operatorId.(uint64))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, connections)
}
