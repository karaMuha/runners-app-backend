package interfaces

import "runners/models"

type UsersService interface {
	GetUser(username string) (string, *models.ResponseError)

	Logout(accessToken string) *models.ResponseError

	GenerateAccessToken(username string) (string, *models.ResponseError)

	AuthorizeUser(accessToken string, expectedRoles []string) (bool, *models.ResponseError)
}
