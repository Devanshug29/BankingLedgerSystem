package api

import (
	controller "BankingLedgerSystem/api/v1"
	"BankingLedgerSystem/internal/config"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func GetRouter(ctx context.Context) *gin.Engine {
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//// API routes
	//api := router.Group("/api/v1")
	//{
	//	// Accounts
	//	api.POST("/accounts", createAccountHandler)
	//	api.GET("/accounts/:id", getAccountHandler)
	//
	//	// Transactions
	//	api.POST("/transactions/deposit", depositHandler)
	//	api.POST("/transactions/withdraw", withdrawHandler)
	//	api.GET("/accounts/:id/transactions", getTransactionsHandler)
	//}

	return router
}

func RegisterRoutes(router *gin.Engine, config *config.Config, c *controller.AccountController) {
	v1router := router.Group("/v1/accounts")
	v1router.POST("", c.CreateAccount)
	v1router.GET("/:accountNumber", c.GetAccount)

}
