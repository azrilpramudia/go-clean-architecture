package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/azrilpramudia/go-clean-architecture/internal/model"
)

type NotificationGateway interface {
	SendVerificationEmail(ctx context.Context, request *model.SendVerificationEmailRequest) error
}

type NotificationGatewayImpl struct {
	BaseURL string
	Client *http.Client
}

func NewNotificationGateway(baseURL string) NotificationGateway {
	return &NotificationGatewayImpl{
		BaseURL: baseURL,
		Client: &http.Client{},
	}
}

func (g *NotificationGatewayImpl) SendVerificationEmail(ctx context.Context, request *model.SendVerificationEmailRequest) error {
	payload, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/send-email", g.BaseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send varification email, status: %d", resp.StatusCode)
	}

	return nil
}