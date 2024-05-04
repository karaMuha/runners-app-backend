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

func (us UsersService) GetUser(username string) (string, *models.ResponseError) {
	if strings.TrimSpace(username) == "" {
		return "", &models.ResponseError{
			Message: "Invalid username or password",
			Status:  http.StatusBadRequest,
		}
	}

	queryResult := us.usersRepository.QueryGetUser(username)

	var password string
	err := queryResult.Scan(&password)

	switch err {
	case nil:
		return password, nil
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

	return us.usersRepository.QueryRemoveAccessToken(accessToken)
}

func (us UsersService) AuthorizeUser(accessToken string, expectedRoles []string) (bool, *models.ResponseError) {
	if accessToken == "" {
		return false, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusUnauthorized,
		}
	}

	queryResult := us.usersRepository.QueryGetUserRole(accessToken)

	var role string
	err := queryResult.Scan(&role)

	if err != nil {
		return false, &models.ResponseError{
			Message: "User in not logged in",
			Status:  http.StatusUnauthorized,
		}
	}

	if role == "" {
		return false, &models.ResponseError{
			Message: "User has no role",
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

func (us UsersService) GenerateAccessToken(username string) (string, *models.ResponseError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.ResponseError{
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
	}

	accessToken := base64.StdEncoding.EncodeToString(hash)

	us.usersRepository.QuerySetAccessToken(accessToken, username)

	return accessToken, nil
}
