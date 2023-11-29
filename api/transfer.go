package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Evans-Prah/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)


type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context)  {
	var payload transferRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validAccount(ctx, payload.FromAccountID, payload.Currency){
		return
	}

	if !server.validAccount(ctx, payload.ToAccountID, payload.Currency){
		return
	}

	if !server.sufficientBalance(ctx, payload.FromAccountID, payload.Amount){
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: payload.FromAccountID,
		ToAccountID:   payload.ToAccountID,
		Amount:        payload.Amount,
	}

	response, err := server.store.TransferTransaction(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse {
		Code: http.StatusOK,
		Message: "Money transfered successfully",
		Data: response,
	})
}


func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ApiResponseFunc(http.StatusNotFound, "Unable to fetch details of account, check and try again", nil))
			return false
		}
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, ApiResponseFunc(http.StatusNotFound, err.Error(), nil))
			return false
	}
	return true
}

func (server *Server) sufficientBalance(ctx *gin.Context, accountID int64, transferAmount int64) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ApiResponseFunc(http.StatusNotFound, "Unable to fetch details of account, check and try again", nil))
			return false
		}
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return false
	}

	if transferAmount > account.Balance {
		ctx.JSON(http.StatusBadRequest, ApiResponseFunc(http.StatusBadRequest, "Insufficient account balance for this transfer", nil))
			return false
	}
	return true
}