package middleware

import (
	"net/http"
	"runners/interfaces"
	"runners/models"

	"github.com/gin-gonic/gin"
)

func AuthorizeRequest(ctx *gin.Context, usersService interfaces.UsersService, roles []string) *models.ResponseError {
	accessToken := ctx.Request.Header.Get("Token")
	auth, responseErr := usersService.AuthorizeUser(accessToken, roles)

	if responseErr != nil {
		return responseErr
	}

	if !auth {
		return &models.ResponseError{
			Message: "Insufficient permissions",
			Status:  http.StatusUnauthorized,
		}
	}

	return nil
}
