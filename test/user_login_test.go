package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/stretchr/testify/assert"
)

func registerTestUser(t *testing.T, serverURL, username, password, name string) {
	requestBody := model.RegisterUserRequest{
		Username: username,
		Password: password,
		Name: name,
	}
	payload, _ := json.Marshal(requestBody)

	resp, err := http.Post(serverURL+"/api/users/register", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestIntegration_LoginSuccess(t *testing.T) {
	cleanupUser(t, "logintest")
	t.Cleanup(func() { cleanupUser(t, "logintest") })

	server := setupTestServer(t)
	registerTestUser(t, server.URL, "logintest", "rahasia123", "Login Test")

	requestBody := model.LoginUserRequest{
		Username: "logintest",
		Password: "rahasia123",
	}
	payload, _ := json.Marshal(requestBody)

	resp, err := http.Post(server.URL+"/api/users/login", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response model.TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.Nil(t, err)
	assert.NotEmpty(t, response.Token)
}

func TestIntegration_Login_WrongPassword(t *testing.T) {
	cleanupUser(t, "wrongpasstest")
	t.Cleanup(func() { cleanupUser(t, "wrongpasstest") })

	server := setupTestServer(t)
	registerTestUser(t, server.URL, "wrongpasstest", "passwordbenar", "Wrong Pass Test")

	requestBody := model.LoginUserRequest{
		Username: "wrongpasstest",
		Password: "passwordsalah",
	}
	payload, _ := json.Marshal(requestBody)

	resp, err := http.Post(server.URL+"/api/users/login", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_Login_UsernameNotFound(t *testing.T) {
	server := setupTestServer(t)

	requestBody := model.LoginUserRequest{
		Username: "usernameyangtidakada",
		Password: "apapun123",
	}

	payload, _ := json.Marshal(requestBody)

	resp, err := http.Post(server.URL+"/api/users/login", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_GetAllUsers_WithValidToken(t *testing.T) {
	cleanupUser(t, "listtest")
	t.Cleanup(func() { cleanupUser(t, "listtest") })

	server := setupTestServer(t)
	registerTestUser(t, server.URL, "listtest", "rahasia123", "List Test")

	loginBody := model.LoginUserRequest{Username: "listtest", Password: "rahasia123"}
	loginPayload, _ := json.Marshal(loginBody)
	loginResp, err := http.Post(server.URL+"/api/users/login", "application/json", bytes.NewBuffer(loginPayload))
	assert.Nil(t, err)
	defer loginResp.Body.Close()

	var TokenResponse model.TokenResponse
	json.NewDecoder(loginResp.Body).Decode(&TokenResponse)

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+TokenResponse.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var users []model.UserResponse
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.Nil(t, err)
	assert.NotEmpty(t, users)
}

func TestIntegration_GetAllUsers_WithoutToken(t *testing.T) {
	server := setupTestServer(t)

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/users", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_GetAllUsers_WithInvalidToken(t *testing.T) {
	server := setupTestServer(t)

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/users", nil)
	req.Header.Set("Authorization", "Bearer token-ngasal-tidak-valid")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}