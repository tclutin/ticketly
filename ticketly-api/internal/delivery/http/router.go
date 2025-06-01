package http

import (
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/middleware"
	v1 "github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1"
	"github.com/tclutin/ticketly/ticketly_api/internal/service"
)

func InitRouter(userSrv service.UserService, ticketSrv service.TicketService, casdoorClient *casdoorsdk.Client) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.ErrorMiddleware(), gin.Recovery(), middleware.CORSMiddleware())

	usersHandler := v1.NewUserHandler(userSrv)
	ticketsHandler := v1.NewTicketHandler(ticketSrv)

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/signin", func(c *gin.Context) {
			code := c.Query("code")
			state := c.Query("state")

			token, err := casdoorClient.GetOAuthToken(code, state)
			if err != nil {
				_ = c.Error(err)
				return
			}

			jwtToken, err := casdoorClient.ParseJwtToken(token.AccessToken)
			if err != nil {
				_ = c.Error(err)
				return
			}

			fmt.Println(jwtToken.User.Email, jwtToken.User.Name, jwtToken.User.Id, jwtToken.User.FirstName, jwtToken.User.LastName)
		})

		v1Group := apiGroup.Group("/v1")
		{
			usersHandler.Bind(v1Group)
			ticketsHandler.Bind(v1Group)
		}
	}

	return router
}
