package services

import (
	"net/http"
	"runners/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRunnerInvalidFirstName(t *testing.T) {
	runner := &models.Runner{
		LastName: "Smith",
		Age:      30,
		Country:  "United States",
	}

	responseErr := validateRunner(runner)

	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid first name", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerInvalidLastName(t *testing.T) {
	runner := &models.Runner{
		FirstName: "Adam",
		Age:       30,
		Country:   "United States",
	}

	responseErr := validateRunner(runner)

	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid last name", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerInvalidAge(t *testing.T) {
	runner := &models.Runner{
		FirstName: "Adam",
		LastName:  "Smith",
		Country:   "United States",
	}

	responseErr := validateRunner(runner)

	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid age", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerInvalidCountry(t *testing.T) {
	runner := &models.Runner{
		FirstName: "Adam",
		LastName:  "Smith",
		Age:       30,
	}

	responseErr := validateRunner(runner)

	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid country", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerWhiteSpaceFirstName(t *testing.T) {
	runner := &models.Runner{
		FirstName: " ",
		LastName:  "Smith",
		Age:       30,
		Country:   "United States",
	}

	responseErr := validateRunner(runner)

	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid first name", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerValidRunner(t *testing.T) {
	runner := &models.Runner{
		FirstName: "Adam",
		LastName:  "Smith",
		Age:       30,
		Country:   "United States",
	}

	responseErr := validateRunner(runner)

	assert.Nil(t, responseErr)
}

func TestValidateRunnerIdEmptyUuid(t *testing.T) {
	runnerId := ""

	responseErr := validateRunnerId(runnerId)

	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid runner ID", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerIdInvalidUuid(t *testing.T) {
	runnerId := "e5280c8b-093d-457a-a535-"

	responseErr := validateRunnerId(runnerId)

	assert.NotEmpty(t, responseErr)
	assert.Equal(t, "Invalid runner ID", responseErr.Message)
	assert.Equal(t, http.StatusBadRequest, responseErr.Status)
}

func TestValidateRunnerIdValidUuid(t *testing.T) {
	runnerId := "e5280c8b-093d-457a-a535-2127326cd1b2"

	responseErr := validateRunnerId(runnerId)

	assert.Nil(t, responseErr)
}
