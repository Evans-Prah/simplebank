package api

import (
	"database/sql"
	"net/http"

	db "github.com/Evans-Prah/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountPayload struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=GHS USD EUR GBP"`
}

func (server *Server) createAccount(ctx *gin.Context)  {
	var payload createAccountPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams {
		Owner: payload.Owner,
		Currency: payload.Currency,
		Balance: 0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID 	int64	`uri:"id" binding:"required,min=1"`
}


func (server *Server) getAccount(ctx *gin.Context)  {
	var req getAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}


type getAccountsRequest struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	PageSize  int32 `form:"page_size" binding:"required,min=2,max=50"`
}

func (server *Server) getAccounts(ctx *gin.Context)  {
	var req getAccountsRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit: req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}