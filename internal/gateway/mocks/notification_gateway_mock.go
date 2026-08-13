package mocks

import (
	"context"

	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/stretchr/testify/mock"
)

type NotificationGatewayMock struct {
	mock.Mock
}

func (m *NotificationGatewayMock) SendVerificationEmail(ctx context.Context, request *model.SendVerificationEmailRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}