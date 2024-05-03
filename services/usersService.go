package services

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"runners/models"
	"runners/repositories"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	usersRepository *repositories.UsersRepository
}

func NewUsersService(usersRepository *repositories.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (us UsersService) GetUser(username string, password string) (string, *models.ResponseError) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return "", &models.ResponseError{
			Message: "Invalid username or password",
			Status:  http.StatusBadRequest,
		}
	}

	queryResult := us.usersRepository.QueryGetUserId(username, password)

	var userId string
	err := queryResult.Scan(&userId)

	switch err {
	case nil:
		return userId, nil
	case sql.ErrNoRows:
		return "", nil
	default:
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
}

func (us UsersService) Logout(accessToken string) *models.ResponseError {
	if accessToken == "" {
		return &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}

	return us.usersRepository.RemoveAccessToken(accessToken)
}

func (us UsersService) AuthorizeUser(accessToken string, expectedRoles []string) (bool, *models.ResponseError) {
	if accessToken == "" {
		return false, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusUnauthorized,
		}
	}

	role, responseErr := us.usersRepository.GetUserRole(accessToken)
	if responseErr != nil {
		return false, responseErr
	}

	if role == "" {
		return false, &models.ResponseError{
			Message: "Failed to authorize user",
			Status:  http.StatusUnauthorized,
		}
	}

	for _, expectedRole := range expectedRoles {
		if expectedRole == role {
			return true, nil
		}
	}

	return false, nil
}

func (us UsersService) GenerateAccessToken(username string, userId string) (string, *models.ResponseError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.ResponseError{
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
	}

	accessToken := base64.StdEncoding.EncodeToString(hash)

	us.usersRepository.SetAccessToken(accessToken, userId)

	return accessToken, nil
}
