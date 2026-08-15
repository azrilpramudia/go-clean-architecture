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
	"golang.org/x/crypto/bcrypt"
)

// Register Success Test
func TestRegister_Success(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	userUsecase := usecase.NewUserUsecase(repoMock, gatewayMock, validate, "test-secret", 24)

	request := &model.RegisterUserRequest{
		Username: "burhan",
		Password: "rahasia123",
		Name: "Burhan",
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

// Test Register Username Already Exists
func TestRegister_UsernameAlreadyExists(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	UserUsecase := usecase.NewUserUsecase(repoMock, gatewayMock ,validate, "test-secret", 24)

	request := &model.RegisterUserRequest{
		Username: "burhan",
		Password: "rahasia123",
		Name: "Burhan",
	}

	existingUser := &entity.User{ID: 1, Username: "burhan"}
	repoMock.On("FindByUsername", context.Background(), "burhan").
		Return(existingUser, nil)

	response, err := UserUsecase.Register(context.Background(), request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "username already registered", err.Error())
}

// Test Login Success
func TestLogin_Success(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	userUsecase := usecase.NewUserUsecase(repoMock, gatewayMock, validate, "test-secret", 24)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("rahasia123"), bcrypt.DefaultCost)
	existingUser := &entity.User{
		ID: 1,
		Username: "burhan",
		Password: string(hashed),
		Name: "Burhan",
	}

	request := &model.LoginUserRequest{
		Username: "burhan",
		Password: "rahasia123",
	}

	repoMock.On("FindByUsername", context.Background(), "burhan").
		Return(existingUser, nil)

	response, err := userUsecase.Login(context.Background(), request)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
}

// Test Login Username Not Found
func TestLogin_UsernameNotFound(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	userUsecase := usecase.NewUserUsecase(repoMock, gatewayMock, validate, "test-secret", 24)

	request := &model.LoginUserRequest{
		Username: "tidakada",
		Password: "rahasia123",
	}

	repoMock.On("FindByUsername", context.Background(), "tidakada").
		Return(nil, nil)

	response, err := userUsecase.Login(context.Background(), request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	repoMock.AssertNotCalled(t, "username or password is wrong", err.Error())
}

// Test Login Wrong Password
func TestLogin_WrongPassword(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	userUsecase := usecase.NewUserUsecase(repoMock, gatewayMock, validate, "test-secret", 24)

	hashed, _ := bcrypt.GenerateFromPassword([]byte("passwordbenar"), bcrypt.DefaultCost)
	existingUser := &entity.User{
		ID: 1,
		Username: "burhan",
		Password: string(hashed),
		Name: "Burhan",
	}
	request := &model.LoginUserRequest{
		Username: "burhan",
		Password: "passwordsalah",
	}

	repoMock.On("FindByUsername", context.Background(), "burhan").
		Return(existingUser, nil)

	response, err := userUsecase.Login(context.Background(), request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	assert.Equal(t, "username or password is wrong",err.Error())
}

// Test List Success
func TestList_Success(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	userUsecase := usecase.NewUserUsecase(repoMock, gatewayMock, validate, "test-secret", 24)

	users := []entity.User{
		{ID: 1, Username: "burhan", Name: "Burhan"},
		{ID: 2, Username: "edward", Name: "Edward"},
	}

	repoMock.On("FindAll", context.Background()).Return(users, nil)

	responses, err := userUsecase.List(context.Background())

	assert.Nil(t, err)
	assert.Len(t, responses, 2)
	assert.Equal(t, "burhan", responses[0].Username)
	assert.Equal(t, "edward", responses[1].Username)
}

// Test List Empty Result
func TestList_EmptyResult(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	userUsecase := usecase.NewUserUsecase(repoMock, gatewayMock, validate, "test-secret", 24)

	repoMock.On("FindAll", context.Background()).Return([]entity.User{}, nil)

	responses, err := userUsecase.List(context.Background())

	assert.Nil(t, err)
	assert.Len(t, responses, 0)
}

// Test Register Validation Failed
func TestRegister_ValidationFailed(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	UserUsecase := usecase.NewUserUsecase(repoMock, gatewayMock ,validate, "test-secret", 24)

	request := &model.RegisterUserRequest{}

	response, err := UserUsecase.Register(context.Background(), request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	repoMock.AssertNotCalled(t, "FindByUsername")
	gatewayMock.AssertNotCalled(t, "SendVerificationEmail")
}

// Test Login Validation Failed
func TestLogin_ValidationFailed(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	UserUsecase := usecase.NewUserUsecase(repoMock, gatewayMock ,validate, "test-secret", 24)

	request := &model.LoginUserRequest{}

	response, err := UserUsecase.Login(context.Background(), request)

	assert.Nil(t, response)
	assert.NotNil(t, err)
	repoMock.AssertNotCalled(t, "FindByUsername")
}

func TestDelete_Success(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	UserUsecase := usecase.NewUserUsecase(repoMock, gatewayMock ,validate, "test-secret", 24)

	existingUser := &entity.User{ID: 1, Username: "burhan", Name:"Burhan"}

	repoMock.On("FindByID", context.Background(), int64(1)).
		Return(existingUser, nil)
	repoMock.On("Delete", context.Background(), int64(1)).
		Return(nil)

	err := UserUsecase.Delete(context.Background(), 1)

	assert.Nil(t, err)
	repoMock.AssertExpectations(t)
}

func TestDelete_UserNotFound(t *testing.T) {
	repoMock := new(mocks.UserRepositoryMock)
	gatewayMock := new(gatewaymocks.NotificationGatewayMock)
	validate := validator.New()
	UserUsecase := usecase.NewUserUsecase(repoMock, gatewayMock ,validate, "test-secret", 24)

	repoMock.On("FindByID", context.Background(), int64(999)).
		Return(nil, nil)

		err := UserUsecase.Delete(context.Background(), 999)

		assert.NotNil(t, err)
		assert.Equal(t, "user not found", err.Error())
		repoMock.AssertNotCalled(t, "Delete")
}