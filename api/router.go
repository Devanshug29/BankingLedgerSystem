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

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

func RegisterRoutes(router *gin.Engine, config *config.Config, c *controller.AccountController) {
	v1router := router.Group("/v1")
	v1router.POST("/accounts", c.CreateAccount)
	v1router.GET("/accounts/:accountNumber", c.GetAccount)
	v1router.POST("/transactions", c.DepositOrWithdraw)

}
