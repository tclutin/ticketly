package v1

import (
	"fmt"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/middleware"
	"github.com/tclutin/ticketly/ticketly_api/internal/delivery/http/v1/response"
	"github.com/tclutin/ticketly/ticketly_api/internal/service"
	"github.com/tclutin/ticketly/ticketly_api/internal/service/operator"
	"net/http"
)

type AuthHandler struct {
	service       service.OperatorService
	casdoorClient *casdoorsdk.Client
}

func NewAuthHandler(service service.OperatorService, casdoorClient *casdoorsdk.Client) *AuthHandler {
	return &AuthHandler{
		service:       service,
		casdoorClient: casdoorClient,
	}
}

func (a *AuthHandler) Bind(router *gin.RouterGroup) {
	authGroup := router.Group("/auth")
	{
		authGroup.GET("/callback", a.OAuthCallback)
		authGroup.GET("/who", middleware.AuthMiddleware(a.service, a.casdoorClient), a.Who)
	}
}

func (a *AuthHandler) OAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	token, err := a.casdoorClient.GetOAuthToken(code, state)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OAuth failed: " + err.Error()})
		return
	}

	claims, err := a.casdoorClient.ParseJwtToken(token.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
		return
	}

	_, err = a.service.SyncOperator(c.Request.Context(), operator.CreateOperatorDTO{
		CasdooID: uuid.MustParse(claims.User.Id),
		Email:    claims.User.Email,
		Name:     claims.User.Name,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	fmt.Println(token.AccessToken, token.Expiry)
	//желательно сгенерировать новый токен, а не использовать токен casdoor
	//потом же в authmiddleware не нужно будет делать +1 запрос GetByCasdoorId

	c.JSON(http.StatusOK, gin.H{
		"access_token": token.AccessToken,
		"expire_at":    token.Expiry,
	})
}

func (a *AuthHandler) Who(c *gin.Context) {
	operatorId, ok := c.Get("operator_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "operator_id not found in context"})
		return
	}

	opr, err := a.service.GetById(c.Request.Context(), operatorId.(uint64))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response.Operator{
		OperatorID: opr.OperatorID,
		CasdoorID:  opr.CasdoorID,
		Email:      opr.Email,
		Name:       opr.Name,
		CreatedAt:  opr.CreatedAt,
		UpdatedAt:  opr.UpdatedAt,
	})
}
