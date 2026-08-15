package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	deliveryhttp "github.com/azrilpramudia/go-clean-architecture/internal/delivery/http"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret-key"

func generateTestToken(t *testing.T, secret string, expired bool) string {
	exp := time.Now().Add(1 * time.Hour)
	if expired {
		exp = time.Now().Add(-1 * time.Hour)
	}

	claims := jwt.MapClaims{
		"sub": 1,
		"exp": exp.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	assert.Nil(t, err)
	return signed
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	nextCalled := false
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}

	handler := deliveryhttp.AuthMiddleware(testSecret, dummyHandler)

	token := generateTestToken(t, testSecret, false)
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	assert.True(t, nextCalled, "next handler harusnya terpanggil kalau token valid")
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	nextCalled := false
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := deliveryhttp.AuthMiddleware(testSecret, dummyHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	assert.False(t, nextCalled, "next handler tidak boleh terpanggil kalau header kosong")
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthMiddleware_MissingBearerPrefix(t *testing.T) {
	nextCalled := false
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := deliveryhttp.AuthMiddleware(testSecret, dummyHandler)

	token := generateTestToken(t, testSecret, false)
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", token)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthMiddleware_NoSpaceAfterBearer(t *testing.T) {
	nextCalled := false
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}

	handler := deliveryhttp.AuthMiddleware(testSecret, dummyHandler)

	token := generateTestToken(t, testSecret, false)
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAuthMiddleware_InvalidTokenFormat(t *testing.T) {
	nextCalled := false
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := deliveryhttp.AuthMiddleware(testSecret, dummyHandler)

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer toke-ngasal-bukan-jwt")
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	nextCalled := false
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := deliveryhttp.AuthMiddleware(testSecret, dummyHandler)

	token := generateTestToken(t, "secret-yang-salah", false)
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	nextCalled := false
	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}

	handler := deliveryhttp.AuthMiddleware(testSecret, dummyHandler)

	token := generateTestToken(t, testSecret, true)
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	recorder := httptest.NewRecorder()

	handler(recorder, req)

	assert.False(t, nextCalled)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}