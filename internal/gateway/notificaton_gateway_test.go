package gateway_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azrilpramudia/go-clean-architecture/internal/gateway"
	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNotificationEmail_Success(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request ){
		assert.Equal(t, "/send-email", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var received model.SendVerificationEmailRequest
		err := json.NewDecoder(r.Body).Decode(&received)
		assert.Nil(t, err)
		assert.Equal(t, "burhan@example.com", received.Email)
		assert.Equal(t, "123456", received.Code)

		w.WriteHeader(http.StatusOK)
	}))
	defer fakeServer.Close()

	notificationGateway := gateway.NewNotificationGateway(fakeServer.URL)

	err := notificationGateway.SendVerificationEmail(context.Background(), &model.SendVerificationEmailRequest{
		Email: "burhan@example.com",
		Code: "123456",
	})

	assert.Nil(t, err)
}

func TestSendVerificationEmail_ServerError(t *testing.T) {
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer fakeServer.Close()

	notificationGateway := gateway.NewNotificationGateway(fakeServer.URL)

	err := notificationGateway.SendVerificationEmail(context.Background(), &model.SendVerificationEmailRequest{
		Email: "burhan@example.com",
		Code: "123456",
	})

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "status: 500")
}