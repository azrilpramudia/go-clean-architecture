package usecase

import (
	"context"
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
	Repository repository.UserRepository
	Gateway gateway.NotificationGateway
	Validate *validator.Validate
	JWTSecret string
	JWTExpiry int
}

func NewUserUsecase(repo repository.UserRepository, gw gateway.NotificationGateway, validate *validator.Validate, jwtSecret string, jwtExpiry int) *UserUsecase {
	return &UserUsecase{
		Repository: repo, 
		Gateway: gw, 
		Validate: validate,
		JWTSecret: jwtSecret,
		JWTExpiry: jwtExpiry,
	}
}

func (u *UserUsecase) Register(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		return nil, err
	}

	existing, err := u.Repository.FindByUsername(ctx, request.Username)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUsernameAlreadyExists
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

	return &model.UserResponse{
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
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, ErrInvalidCredentials
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

func (u *UserUsecase) List(ctx context.Context) ([]model.UserResponse, error) {
	users, err := u.Repository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]model.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, model.UserResponse{
			ID: user.ID,
			Username: user.Username,
			Name: user.Name,
		})
	}
	return responses, nil
}
	
func (u *UserUsecase) Delete(ctx context.Context, id int64) error {
	existing, err := u.Repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}
	
	return u.Repository.Delete(ctx, id)
}
	

	


	

	