package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func GetRouter(ctx context.Context) (*gin.Engine, error) {
	router := gin.New()

	// Middlewares
	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	// Healthcheck
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Accounts
		api.POST("/accounts", createAccountHandler)
		api.GET("/accounts/:id", getAccountHandler)

		// Transactions
		api.POST("/transactions/deposit", depositHandler)
		api.POST("/transactions/withdraw", withdrawHandler)
		api.GET("/accounts/:id/transactions", getTransactionsHandler)
	}

	return router, nil
}
