package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/Evans-Prah/simplebank/db/sqlc"
	"github.com/Evans-Prah/simplebank/token"
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

	fromAccount, valid := server.validAccount(ctx, payload.FromAccountID, payload.Currency)
	if !valid{
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)
	if fromAccount.Owner != authPayload.Username {
		ctx.JSON(http.StatusForbidden, ApiResponseFunc(http.StatusForbidden, "Account does not belong to authenticated user", nil))
		return
	}

	_, valid = server.validAccount(ctx, payload.FromAccountID, payload.Currency)
	if !valid{
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


func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ApiResponseFunc(http.StatusNotFound, "Unable to fetch details of account, check and try again", nil))
			return account, false
		}
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, ApiResponseFunc(http.StatusNotFound, err.Error(), nil))
			return account, false
	}
	return account, true
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