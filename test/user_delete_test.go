package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loginAndGetToken(t *testing.T, serverURL, username, password string) string {
	loginBody := model.LoginUserRequest{Username: username, Password: password}
	payload, _ := json.Marshal(loginBody)

	resp, err := http.Post(serverURL+"/api/users/login", "application/json", bytes.NewBuffer(payload))
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var tokenResponse model.TokenResponse
	json.NewDecoder(resp.Body).Decode(&tokenResponse)
	return tokenResponse.Token
}

func TestIntegration_DeleteUser_Success(t *testing.T) {
	cleanupUser(t, "deletetest")
	t.Cleanup(func() { cleanupUser(t, "deletetest")})

	server := setupTestServer(t)
	registerTestUser(t, server.URL, "deletetest", "rahasia123", "Delete Test")
	token := loginAndGetToken(t, server.URL, "deletetest", "rahasia123")

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	listResp, err := client.Do(req)
	assert.Nil(t, err)
	defer listResp.Body.Close()

	var users []model.UserResponse
	json.NewDecoder(listResp.Body).Decode(&users)

	var targetID int64
	for _, u := range users {
		if u.Username == "deletetest" {
			targetID = u.ID
		}
	}
	require.NotZero(t, targetID, "targetID tidak boleh 0, artinya user tidak ditemukan di list")

	deleteReq, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/users/%d", server.URL, targetID), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteResp, err := client.Do(deleteReq)
	assert.Nil(t, err)
	defer deleteResp.Body.Close()

	assert.Equal(t, http.StatusNoContent, deleteResp.StatusCode)
}

func TestIntegration_DeleteUser_NotFound(t *testing.T) {
	server := setupTestServer(t)

	registerTestUser(t, server.URL, "authfordelete", "rahasia123", "Auth For Delete")
	t.Cleanup(func() { cleanupUser(t, "authfordelete") })
	token := loginAndGetToken(t, server.URL, "authfordelete", "rahasia123")

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/users/999999999", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_DeleteUser_WithoutToken(t *testing.T) {
	server := setupTestServer(t)
	
	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/users/1", nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_DeleteUser_InvalidIDFormat(t *testing.T) {
	server := setupTestServer(t)

	registerTestUser(t, server.URL, "authforinvalidid", "rahasia123", "Auth For Invalid ID")
	t.Cleanup(func() { cleanupUser(t, "authforinvalidid") })
	token := loginAndGetToken(t, server.URL, "authforinvalidid", "rahasia123")

	req, _ := http.NewRequest(http.MethodDelete, server.URL+"/api/users/bukan-angka", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

