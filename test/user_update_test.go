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

func TestIntegration_UpdateUser_Success(t *testing.T) {
	cleanupUser(t, "updatetest")
	t.Cleanup(func() { cleanupUser(t, "updatetest") })

	server := setupTestServer(t)
	registerTestUser(t, server.URL, "updatetest", "rahasia123", "Nama Lama")
	token := loginAndGetToken(t, server.URL, "updatetest", "rahasia123")

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
		if u.Username == "updatetest" {
			targetID = u.ID
		}
	}
	require.NotZero(t, targetID, "targetID tidak boleh 0, artinya user tidak ditemukan di list")

	updateBody := model.UpdateUserRequest{Name: "Nama Baru"}
	payload, _ := json.Marshal(updateBody)

	updateReq, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/users/%d", server.URL, targetID), bytes.NewBuffer(payload))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := client.Do(updateReq)
	assert.Nil(t, err)
	defer updateResp.Body.Close()

	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	var response model.UserResponse
	json.NewDecoder(updateResp.Body).Decode(&response)
	assert.Equal(t, "Nama Baru", response.Name)
	assert.Equal(t, "updatetest", response.Username)
}

func TestIntegration_UpdateUser_NotFound(t *testing.T) {
	server := setupTestServer(t)

	registerTestUser(t, server.URL, "authforupdate", "rahasia123", "Auth For Update")
	t.Cleanup(func() { cleanupUser(t, "authforupdate") })
	token := loginAndGetToken(t, server.URL, "authforupdate", "rahasia123")

	updateBody := model.UpdateUserRequest{Name: "Nama Baru"}
	payload, _ := json.Marshal(updateBody)

	req, _ := http.NewRequest(http.MethodPatch, server.URL+"/api/users/999999999", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Contetn-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestIntegration_UpdateUser_WithoutToken(t *testing.T) {
	server := setupTestServer(t)

	updateBody := model.UpdateUserRequest{Name: "Nama Baru"}
	payload, _ := json.Marshal(updateBody)

	req, _ := http.NewRequest(http.MethodPatch, server.URL+"/api/users/1", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_UpdateUser_ValidationFailed(t *testing.T) {
	server := setupTestServer(t)

	registerTestUser(t, server.URL, "authforvalidation", "rahasia123", "Auth For Validation")
	t.Cleanup(func() { cleanupUser(t, "authforvalidation") })
	token := loginAndGetToken(t, server.URL, "authforvalidation", "rahasia123")

	payload, _ := json.Marshal(map[string]string{})

	req, _ := http.NewRequest(http.MethodPatch, server.URL+"/api/users/1", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

