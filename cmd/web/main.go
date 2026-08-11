package main

import (
	"database/sql"
	"log"
	"net/http"

	deliveryhttp "github.com/azrilpramudia/go-clean-architecture/internal/delivery/http"
	"github.com/azrilpramudia/go-clean-architecture/internal/repository"
	"github.com/azrilpramudia/go-clean-architecture/internal/usecase"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:MeowHX@01@tcp(localhost:3306)/myapp?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	validate := validator.New()

	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository, validate)
	userHandler := deliveryhttp.NewUserHandler(userUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/users/register", userHandler.Register)

	log.Println("server running on :3000")
	log.Fatal(http.ListenAndServe(":3000", mux))
}