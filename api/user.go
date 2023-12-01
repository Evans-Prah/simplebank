package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	db "github.com/Evans-Prah/simplebank/db/sqlc"
	"github.com/Evans-Prah/simplebank/db/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type createUserPayload struct {
	Username string `json:"username" binding:"required,customUsername"`
	Fullname string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type createUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var payload createUserPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		validationErrors := formatValidationErrors(err.(validator.ValidationErrors))
		ctx.JSON(http.StatusBadRequest, ApiResponseFunc(http.StatusBadRequest, "Validation Errors", nil, validationErrors))
		return
	}

	existingUserArgs := db.GetUserByUsernameOrEmailParams{
		Username: payload.Username,
		Email:    payload.Email,
	}

	existingUser, err := server.store.GetUserByUsernameOrEmail(ctx, existingUserArgs)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusFailedDependency, errorResponse(err))
		return
	}

	if existingUser.Username == payload.Username {
		ctx.JSON(http.StatusConflict, ApiResponseFunc(http.StatusConflict, "Username already exists", nil, nil))
		return
	}

	if existingUser.Email == payload.Email {
		ctx.JSON(http.StatusConflict, ApiResponseFunc(http.StatusConflict, "Email already exists", nil, nil))
		return
	}

	hashedPassword, hashErr := util.HashPassword(payload.Password)
	if hashErr != nil {
		ctx.JSON(http.StatusInternalServerError, ApiResponseFunc(http.StatusFailedDependency, "Could not create user, try again", nil, hashErr.Error()))
		return
	}

	arg := db.CreateUserParams{
		Username:       payload.Username,
		FullName:       payload.Fullname,
		Email:          payload.Email,
		HashedPassword: hashedPassword,
	}

	user, createErr := server.store.CreateUser(ctx, arg)
	if createErr != nil {
		ctx.JSON(http.StatusInternalServerError, ApiResponseFunc(http.StatusBadRequest, "Something bad happened when creating user, try again in a few minutes", nil, createErr.Error()))
		return
	}

	createUserResponseDto := createUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, ApiResponseFunc(http.StatusCreated, "User created successfully", createUserResponseDto))

}


type loginUserPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	Username string `json:"username"`
	AccessToken string `json:"access_token"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var payload loginUserPayload
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		validationErrors := formatValidationErrors(err.(validator.ValidationErrors))
		ctx.JSON(http.StatusBadRequest, ApiResponseFunc(http.StatusBadRequest, "Validation Errors", nil, validationErrors))
		return
	}

	existingUserArgs := db.GetUserByUsernameOrEmailParams{
		Username: payload.Username,
		Email:    payload.Username,
	}

	existingUser, err := server.store.GetUserByUsernameOrEmail(ctx, existingUserArgs)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusInternalServerError, ApiResponseFunc(http.StatusInternalServerError, "Something bad happened, try again later", nil, nil))
		return
	}

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, ApiResponseFunc(http.StatusNotFound, "Invalid credentials, check and try again", nil, nil))
		return
	}

	err = util.CheckPassword(payload.Password, existingUser.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ApiResponseFunc(http.StatusUnauthorized, "Invalid credentials, check and try again", nil, nil))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(existingUser.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusFailedDependency, ApiResponseFunc(http.StatusFailedDependency, "Something bad happened, try again later", nil, nil))
		return
	}

	response := loginUserResponse {
		Username: existingUser.Username,
		AccessToken: accessToken,
	}

	ctx.JSON(http.StatusOK, ApiResponseFunc(http.StatusOK, "User login successful", response))

}