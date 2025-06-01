package middleware

import (
	"errors"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tclutin/ticketly/ticketly_api/internal/service"
	coreerrors "github.com/tclutin/ticketly/ticketly_api/internal/service/errors"
	"github.com/tclutin/ticketly/ticketly_api/internal/service/operator"
	"net/http"
	"strings"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		switch {
		case errors.Is(err, coreerrors.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, coreerrors.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, coreerrors.ErrTicketNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, coreerrors.ErrTicketAlreadyClosedOrInProgress):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, coreerrors.ErrTicketWrongType):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, coreerrors.ErrTicketWrongStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, coreerrors.ErrOperatorNotAssigned):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, coreerrors.ErrActiveTicketAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, coreerrors.ErrOperatorNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}

		c.Abort()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func AuthMiddleware(srv service.OperatorService, casdoorClient *casdoorsdk.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is empty"})
			return
		}

		items := strings.Split(bearerToken, " ")
		if len(items) != 2 || items[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is invalid"})
			return
		}

		token := items[1]

		claims, err := casdoorClient.ParseJwtToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		operatorId, err := srv.SyncOperator(c.Request.Context(), operator.CreateOperatorDTO{
			CasdooID: uuid.MustParse(claims.User.Id),
			Email:    claims.User.Email,
			Name:     claims.User.Name,
		})

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Set("operator_id", operatorId)

		c.Next()
	}
}
