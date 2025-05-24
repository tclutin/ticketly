package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1/request"
	"github.com/tclutin/ticketly/ticketly_api/internal/service"
	"github.com/tclutin/ticketly/ticketly_api/internal/service/ticket"
	"net/http"
	"strconv"
)

type TicketHandler struct {
	service service.TicketService
}

func NewTicketHandler(service service.TicketService) *TicketHandler {
	return &TicketHandler{
		service: service,
	}
}

func (t *TicketHandler) Bind(router *gin.RouterGroup) {
	ticketsGroup := router.Group("/tickets")
	{
		ticketsGroup.GET("", t.GetAll)
		ticketsGroup.POST("", t.Create)
		ticketsGroup.POST("/:ticket_id/close", t.Close)
	}
}

func (t *TicketHandler) Create(c *gin.Context) {
	var req request.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	ticketId, err := t.service.Create(c.Request.Context(), ticket.CreateTicketDTO{
		UserID:  req.UserID,
		Type:    req.Type,
		Content: req.Content,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ticket_id": ticketId,
	})

}

func (t *TicketHandler) GetAll(c *gin.Context) {
	tickets, err := t.service.GetAll(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
	}

	c.JSON(http.StatusOK, tickets)
}

func (t *TicketHandler) Close(c *gin.Context) {
	ticketStr := c.Param("ticket_id")
	operatorId := 228

	ticketId, err := strconv.ParseUint(ticketStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket_id"})
		return
	}

	var req request.CloseTicketRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "invalid request body: " + err.Error(),
		})
		return
	}

	if err = t.service.Close(c.Request.Context(), ticket.CloseTicketDTO{
		TicketID:   ticketId,
		OperatorID: uint64(operatorId),
		Message:    req.Content,
	}); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "closed",
	})
}
