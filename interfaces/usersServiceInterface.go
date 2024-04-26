package interfaces

import "runners/models"

type UsersService interface {
	Login(username string, password string) (string, *models.ResponseError)

	Logout(accessToken string) *models.ResponseError

	AuthorizeUser(accessToken string, expectedRoles []string) (bool, *models.ResponseError)
}
