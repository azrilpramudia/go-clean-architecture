package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/azrilpramudia/go-clean-architecture/internal/entity"
	"github.com/azrilpramudia/go-clean-architecture/internal/gateway"
	"github.com/azrilpramudia/go-clean-architecture/internal/model"
	"github.com/azrilpramudia/go-clean-architecture/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repository repository.UserRespository
	Gateway gateway.NotificationGateway
	Validate *validator.Validate
	JWTSecret string
	JWTExpiry int
}

func NewUserUsecase(repo repository.UserRespository, gw gateway.NotificationGateway, validate *validator.Validate, jwtSecret string, jwtExpiry int) *UserUsecase {
	return &UserUsecase{
		Repository: repo, 
		Gateway: gw, 
		Validate: validate,
		JWTSecret: jwtSecret,
		JWTExpiry: jwtExpiry,
	}
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

func (u *UserUsecase) Login(ctx context.Context, request *model.LoginUserRequest) (*model.TokenResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, err
	}

	user, err := u.Repository.FindByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("username or password is wrong")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, errors.New("username or password is wrong")
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"username": user.Username,
		"exp": time.Now().Add(time.Duration(u.JWTExpiry) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(u.JWTSecret))
	if err != nil {
		return nil, err
	}

	return &model.TokenResponse{Token: signedToken}, nil
}

func (u *UserUsecase) List(ctx context.Context) ([]model.UserRespone, error) {
	users, err := u.Repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]model.UserRespone, 0, len(users))
	for _, user := range users {
		responses = append(responses, model.UserRespone{
			ID: user.ID,
			Username: user.Username,
			Name: user.Name,
		})
	}
	return responses, nil
}
	
	

	


	

	