package usecase

import (
	"context"
	"errors"

	"github.com/azrilpramudia/go-clean-architecture/internal/entity"
	"github.com/azrilpramudia/go-clean-architecture/internal/gateway"
	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/azrilpramudia/go-clean-architecture/internal/repository"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repository repository.UserRespository
	Gateway gateway.NotificationGateway
	Validate *validator.Validate
}

func NewUserUsecase(repo repository.UserRespository, gw gateway.NotificationGateway, validate *validator.Validate) *UserUsecase {
	return &UserUsecase{Repository: repo, Gateway: gw, Validate: validate }
}

func (u *UserUsecase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserRespone, error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, err
	}
	
	existing, _ := u.Repository.FindByUsername(ctx, request.Username)
	if existing != nil {
		return nil, errors.New("username already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: request.Username,
		Password: string(hashed),
		Name: request.Name,
	}

	if err := u.Repository.Save(ctx, user); err != nil {
		return nil, err
	}

	_ = u.Gateway.SendVerificationEmail(ctx, &model.SendVerificationEmailRequest{
		Email: user.Username,
		Code: "123456",
	})

	return &model.UserRespone{
		ID: user.ID,
		Username: user.Username,
		Name: user.Name,
	}, nil
}