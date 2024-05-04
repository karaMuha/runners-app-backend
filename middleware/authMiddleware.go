package middleware

import (
	"net/http"
	"runners/interfaces"
	"runners/models"
)

func AuthorizeRequest(req *http.Request, usersService interfaces.UsersService, roles []string) *models.ResponseError {
	accessToken := req.Header.Get("Token")
	auth, responseErr := usersService.AuthorizeUser(accessToken, roles)

	if responseErr != nil {
		return responseErr
	}

	if !auth {
		return &models.ResponseError{
			Message: "Not authorized",
			Status:  http.StatusUnauthorized,
		}
	}

	return nil
}
