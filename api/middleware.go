package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Evans-Prah/simplebank/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("Authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ApiResponseFunc(http.StatusUnauthorized, err.Error(), nil, err.Error()))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) > 2 {
			err := errors.New("Invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ApiResponseFunc(http.StatusUnauthorized, err.Error(), nil, err.Error()))
			return
		}

		authorizationType := fields[0]
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("Unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ApiResponseFunc(http.StatusUnauthorized, err.Error(), nil, err.Error()))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, ApiResponseFunc(http.StatusUnauthorized, err.Error(), nil, err.Error()))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}