package usecase_test

import (
	"context"
	"testing"

	"github.com/azrilpramudia/go-clean-architecture/internal/entity"
	gatewaymocks "github.com/azrilpramudia/go-clean-architecture/internal/gateway/mocks"
	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/azrilpramudia/go-clean-architecture/internal/repository/mocks"
	"github.com/azrilpramudia/go-clean-architecture/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister_Success(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	userUsecase := usecase.NewUserUsecase(repoMock, gatewayMock, validate)

	request := &model.RegisterUserRequest{
		Username: "burhan",
		Password: "burhan@123",
		Name: "Burhan Snake",
	}

	repoMock.On("FindByUsername", context.Background(), "burhan").
		Return(nil, nil)

	repoMock.On("Save", context.Background(), mock.AnythingOfType("*entity.User")).
		Return(nil)

	gatewayMock.On("SendVerificationEmail", context.Background(), mock.AnythingOfType("*model.SendVerificationEmailRequest")).
		Return(nil)

	response, err := userUsecase.Register(context.Background(), request)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "burhan", response.Username)
	repoMock.AssertExpectations(t)
}

func TestRegister_UsernameAlreadyExists(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	UserUsecase := usecase.NewUserUsecase(repoMock, gatewayMock ,validate)

	request := &model.RegisterUserRequest{
		Username: "burhan",
		Password: "burhan@123",
		Name: "Burhan Snake",
	}

	existingUser := &entity.User{ID: 1, Username: "burhan"}
	repoMock.On("FindByUsername", context.Background(), "burhan").
		Return(existingUser, nil)

	response, err := UserUsecase.Register(context.Background(), request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "username already registered", err.Error())
}

func TestRegister_ValidationFailed(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	userUsecase := usecase.NewUserUsecase(repoMock, gatewayMock, validate)

	request := &model.RegisterUserRequest{}

	response, err := userUsecase.Register(context.Background(), request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	repoMock.AssertNotCalled(t, "FindByUsername")
}