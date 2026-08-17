package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/azrilpramudia/go-clean-architecture/internal/config"
	deliveryhttp "github.com/azrilpramudia/go-clean-architecture/internal/delivery/http"
	"github.com/azrilpramudia/go-clean-architecture/internal/gateway"
	"github.com/azrilpramudia/go-clean-architecture/internal/repository"
	"github.com/azrilpramudia/go-clean-architecture/internal/usecase"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("%v", err)
	}
	db := config.NewDatabase(cfg)
	defer db.Close()

	validate := validator.New()
	notificationGateway := gateway.NewNotificationGateway(cfg.Notification.BaseURL)

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, notificationGateway, validate, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	userHandler := deliveryhttp.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users/register", userHandler.Register)
	mux.HandleFunc("POST /api/users/login", userHandler.Login)
	mux.HandleFunc("GET /api/users", deliveryhttp.AuthMiddleware(cfg.JWT.Secret, userHandler.List))
	mux.HandleFunc("DELETE /api/users/{id}", deliveryhttp.AuthMiddleware(cfg.JWT.Secret, userHandler.Delete))

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}