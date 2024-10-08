package api

import (
	"database/sql"
	"net/http"
	"strconv"

	db "github.com/Evans-Prah/simplebank/db/sqlc"
	"github.com/Evans-Prah/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
)


type createAccountPayload struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context)  {
	var payload createAccountPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		validationErrors := formatValidationErrors(err.(validator.ValidationErrors))
		ctx.JSON(http.StatusBadRequest, ApiResponseFunc(http.StatusBadRequest, "Invalid input", nil, validationErrors))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)

	arg := db.CreateAccountParams {
		Owner: authPayload.Username,
		Currency: payload.Currency,
		Balance: 0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		if dbErr, ok := err.(*pq.Error); ok {
			switch (dbErr.Code.Name()) {
			case "foreign_key_violation", "unique_violation":
					ctx.JSON(http.StatusForbidden, errorResponse(err))
					return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
			ctx.JSON(http.StatusNotFound, ApiResponseFunc(http.StatusNotFound, "Unable to fetch details of account, check and try again", nil))
			return
		}
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusForbidden, ApiResponseFunc(http.StatusForbidden, "Account does not belong to authenticated user", nil))
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse {
		Code: http.StatusOK,
		Message: "Account details fetched successfully",
		Data: account,
	})
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)

	arg := db.ListAccountsParams{
		Owner: authPayload.Username,
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

type UpdateAccountRequest struct {
	Balance int64 `json:"balance" binding:"required"`
}


func (server *Server) updateAccount(ctx *gin.Context) {
	// Extract account ID from the route parameters
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ApiResponseFunc(http.StatusBadRequest, "Invalid account ID", nil))
		return
	}

	var req UpdateAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationErrors := formatValidationErrors(err.(validator.ValidationErrors))
		ctx.JSON(http.StatusBadRequest, ApiResponseFunc(http.StatusBadRequest, "Invalid payload", nil, validationErrors))
		return
	}

	account, err := server.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ApiResponseFunc(http.StatusNotFound, "Account not found, check and try again", nil))
			return
		}
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusForbidden, ApiResponseFunc(http.StatusForbidden, "Account does not belong to authenticated user", nil))
		return
	}

	arg := db.UpdateAccountParams {
		ID: id,
		Balance: req.Balance,
	}

	updatedAccount, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    http.StatusOK,
		Message: "Account details updated successfully",
		Data:    updatedAccount,
	})
}

type deleteAccountRequest struct {
	ID 	int64	`uri:"id" binding:"required"`
}


func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	err := ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ApiResponseFunc(http.StatusNotFound, "Account not found, check and try again", nil))
			return
		}
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.TokenPayload)
	if account.Owner != authPayload.Username {
		ctx.JSON(http.StatusForbidden, ApiResponseFunc(http.StatusForbidden, "Account does not belong to authenticated user", nil))
		return
	}

	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, ApiResponse{
		Code:    http.StatusOK,
		Message: "Account deleted successfully",
	})
}
