package controller

import (
	service2 "BankingLedgerSystem/business"
	"BankingLedgerSystem/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountController handles account APIs
type AccountController struct {
	svc *service2.AccountService
}

func NewAccountController(svc *service2.AccountService) *AccountController {
	return &AccountController{svc: svc}
}

// CreateAccount godoc
// @Summary Create a new account
// @Description Create an account with initial balance
// @Tags Accounts
// @Accept json
// @Produce json
// @Param request body models.CreateAccountRequest true "Account creation request"
// @Success 201 {object} models.SuccessResponse{data=models.Account}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/accounts [post]
func (c *AccountController) CreateAccount(ctx *gin.Context) {
	var req models.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request payload",
			Details: err.Error(),
		})
		return
	}

	account, err := c.svc.CreateAccount(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "could not create account",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.SuccessResponse{
		Code: http.StatusCreated,
		Data: account,
	})
}

// GetAccount godoc
// @Summary Get account details
// @Description Fetch an account by ID
// @Tags Accounts
// @Produce json
// @Param accountNumber path string true "Account Number"
// @Success 200 {object} models.SuccessResponse{data=models.Account}
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /v1/accounts/{accountNumber} [get]
func (c *AccountController) GetAccount(ctx *gin.Context) {
	accountNumber := ctx.Param("accountNumber")

	account, err := c.svc.GetAccount(ctx, accountNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "could not fetch account",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.SuccessResponse{
		Code: http.StatusOK,
		Data: account,
	})
}

// DepositOrWithdraw godoc
// @Summary Deposit or withdraw funds
// @Description Send a deposit/withdraw request, processed asynchronously via Kafka
// @Tags Transactions
// @Accept json
// @Produce json
// @Param request body models.TransactionRequest true "Transaction request"
// @Success 202 {object} models.SuccessResponse{data=string}
// @Failure 400 {object} models.ErrorResponse
// @Router /v1/transactions [post]
func (c *AccountController) DepositOrWithdraw(ctx *gin.Context) {
	var req models.TransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid request payload",
			Details: err.Error(),
		})
		return
	}

	// Serialize request
	payload, _ := json.Marshal(req)

	// Publish to Kafka with accountNumber as key (ordering preserved per account)
	if err := c.svc.PublishTransaction(ctx, []byte(req.AccountNumber), payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "could not publish transaction",
			Details: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusAccepted, models.SuccessResponse{
		Code: http.StatusAccepted,
		Data: "transaction submitted",
	})
}
