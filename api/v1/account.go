package v1

import (
	"BankingLedgerSystem/models"
	"BankingLedgerSystem/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateAccount godoc
// @Summary Create a new account
// @Description API to create an account with initial balance
// @ID CreateAccount
// @Tags Accounts
// @Accept json
// @Produce json
// @Param request body models.CreateAccountRequest true "Account creation request"
// @Success 201 {object} models.Response{data=models.Account} "Account created successfully"
// @Failure 400 {object} models.Response{error=string} "Invalid input"
// @Failure 500 {object} models.Response{error=string} "Server error"
// @Router /v1/accounts [POST]
func CreateAccount(ctx *gin.Context) {
	var req models.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request payload",
			Details: err.Error(),
		})
		return
	}

	account, err := service.CreateAccount(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request payload",
			Details: err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.SuccessResponse(account))
}

// GetAccount godoc
// @Summary Get account details
// @Description API to fetch account details by ID
// @ID GetAccount
// @Tags Accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} models.Response{data=models.Account} "Account retrieved successfully"
// @Failure 404 {object} models.Response{error=string} "Account not found"
// @Router /v1/accounts/{id} [GET]
func GetAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	account, err := service.GetAccount(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse(account))
}
