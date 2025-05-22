package v1

import "github.com/gin-gonic/gin"

type TicketHandler struct {
}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{}
}

func (t *TicketHandler) Bind(router *gin.RouterGroup) {

}

func (t *TicketHandler) Create(c *gin.Context) {}

func (t *TicketHandler) GetAll(c *gin.Context) {}
