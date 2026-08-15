package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azrilpramudia/go-clean-architecture/internal/config"
	deliveryhttp "github.com/azrilpramudia/go-clean-architecture/internal/delivery/http"
	gatewaymocks "github.com/azrilpramudia/go-clean-architecture/internal/gateway/mocks"
	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/azrilpramudia/go-clean-architecture/internal/repository"
	"github.com/azrilpramudia/go-clean-architecture/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestServer(t *testing.T) *httptest.Server {
	cfg := config.Load("../config.json")
	db := config.NewDatabase(cfg)
	t.Cleanup(func() {db.Close() })

	validate := validator.New()

	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	gatewayMock.On("SendVerificationEmail", mock.Anything, mock.Anything).Return(nil)

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, gatewayMock, validate, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	userHandler := deliveryhttp.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users/register", userHandler.Register)
	mux.HandleFunc("POST /api/users/login", userHandler.Login)
	mux.HandleFunc("GET /api/users", deliveryhttp.AuthMiddleware(cfg.JWT.Secret, userHandler.List))

	server := httptest.NewServer(mux)
	t.Cleanup(func() { server.Close() })

	return server
}

func cleanupUser(t *testing.T, username string) {
	cfg := config.Load("../config.json")
	db := config.NewDatabase(cfg)
	defer db.Close()

	_, err := db.ExecContext(context.Background(), "DELETE FROM users WHERE username = ?", username)
	if err != nil {
		t.Fatalf("failed to cleanup test user: %v", err)
	}
}

func TestIntegration_RegisterUser_Success(t *testing.T) {
	cleanupUser(t, "integrationtest")
	t.Cleanup(func() {cleanupUser(t, "integrationtest") })

	server := setupTestServer(t)

	requestBody := model.RegisterUserRequest{
		Username: "integrationtest",
		Password: "rahasia123",
		Name: "Integration Test",
	}
	payload, _ := json.Marshal(requestBody)

	resp, err := http.Post(server.URL+"/api/users/register", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response model.UserRespone
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.Nil(t, err)
	assert.Equal(t, "integrationtest", response.Username)
	assert.Equal(t, "Integration Test", response.Name)
	assert.NotZero(t, response.ID)
}

func TestIntegration_RegisterUser_DuplicateUsername(t *testing.T) {
	cleanupUser(t, "duplicatetest")
	t.Cleanup(func() { cleanupUser(t, "duplicatetest") })

	server := setupTestServer(t)

	requestBody := model.RegisterUserRequest{
		Username: "duplicatetest",
		Password: "rahasia123",
		Name: "Duplicate Test",
	}
	payload, _ := json.Marshal(requestBody)

	resp1, err := http.Post(server.URL+"/api/users/register", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	resp1.Body.Close()
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)

	resp2, err := http.Post(server.URL+"/api/users/register", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	resp2.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp2.StatusCode)
}

func TestIntegration_Register_ValidationFailed(t *testing.T) {
	server := setupTestServer(t)

	payload, _ := json.Marshal(map[string]string{})

	resp, err := http.Post(server.URL+"/api/users/register", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}