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
	cfg := config.Load("config.json")
	db := config.NewDatabase(cfg)
	defer db.Close()

	validate := validator.New()

	userRepository := repository.NewUserRepository(db)
	NotificationGateway := gateway.NewNotificationGateway("https://api.emailprovider.com")
	userUsecase := usecase.NewUserUsecase(userRepository, NotificationGateway, validate)
	userHandler := deliveryhttp.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users/register", userHandler.Register)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}