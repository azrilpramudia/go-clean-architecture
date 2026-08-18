# Go Clean Architecture

A simple REST API for user management (Register, Login, List, Update, Delete) built following **Clean Architecture** principles, as a backend learning project with Go.

The project structure is based on [golang-clean-architecture](https://github.com/khannedy/golang-clean-architecture) by Eko Kurniawan Khannedy (Programmer Zaman Now), with several adjustments.

## Features

- ✅ Register a new user with input validation
- ✅ Login with JWT authentication
- ✅ Get all users (protected endpoint)
- ✅ Update user data (protected endpoint)
- ✅ Delete user (protected endpoint)
- ✅ Password hashing with bcrypt
- ✅ JWT algorithm validation (prevents *algorithm confusion attack*)
- ✅ Database migrations with `golang-migrate`
- ✅ API documentation with OpenAPI 3.0
- ✅ Unit tests, HTTP-level tests, and integration tests (38+ tests)

## Tech Stack

| Category | Technology |
|---|---|
| Language | Go 1.26 |
| Database | MySQL |
| Auth | JWT (`golang-jwt/jwt/v5`) |
| Validation | `go-playground/validator/v10` |
| Password Hashing | `golang.org/x/crypto/bcrypt` |
| Migration | `golang-migrate/migrate` |
| Testing | `stretchr/testify` |

## Architecture

This project follows the Clean Architecture pattern with clear layer separation:

```
Delivery (HTTP) → Usecase (business logic) → Repository (data access) → Database
                        ↓
                    Gateway (external communication)
```

- **Entity** — pure data representation matching the database structure
- **Model** — data shape for API requests/responses (DTO)
- **Repository** — data access layer, defined as an interface for dependency inversion
- **Usecase** — business logic, orchestrates Repository and Gateway
- **Delivery** — HTTP handlers, receive requests and return responses
- **Gateway** — communication with external systems (e.g. email notification service)
- **Config** — application settings (database, JWT, server) loaded from `config.json`

## Folder Structure

```
go-clean-architecture/
├── api/
│   └── openapi.yaml          # API documentation (OpenAPI 3.0)
├── cmd/
│   └── web/
│       └── main.go           # Application entry point
├── db/
│   └── migrations/           # Database migration files
├── internal/
│   ├── config/                # App configuration & database connection
│   ├── delivery/
│   │   └── http/              # HTTP handlers & middleware
│   ├── entity/                 # Domain models (database structs)
│   ├── gateway/                 # External system communication
│   ├── model/                   # Request/response DTOs
│   ├── repository/              # Data access (interface + implementation)
│   └── usecase/                 # Business logic
├── test/                      # Integration tests
├── config.json.example       # Example configuration
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.26 or later
- MySQL
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)

### 1. Clone the repository

```bash
git clone https://github.com/azrilpramudia/go-clean-architecture.git
cd go-clean-architecture
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Set up configuration

Copy `config.json.example` to `config.json`, then adjust it to your local environment:

```bash
cp config.json.example config.json
```

```json
{
  "database": {
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "your_password_here",
    "name": "myapp"
  },
  "server": {
    "port": 3000
  },
  "jwt": {
    "secret": "your_jwt_secret_here",
    "expiry_hours": 24
  },
  "notification": {
    "base_url": "https://api.emailservice.com"
  }
}
```

> `config.json` is already in `.gitignore` and must not be committed, since it contains credentials.

### 4. Create the database

```bash
mysql -u root -p -e "CREATE DATABASE myapp;"
```

### 5. Run migrations

```bash
migrate -database "mysql://root:your_password@tcp(localhost:3306)/myapp?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up
```

### 6. Run the server

```bash
go run ./cmd/web/main.go
```

The server runs on `http://localhost:3000` (or whichever port is set in `config.json`).

## Running Tests

```bash
# All tests
go test ./... -v

# Unit tests only (no database required)
go test ./internal/... -v

# Integration tests only (requires MySQL running)
go test ./test/... -v

# With coverage report
go test ./... -cover
```

## API Documentation

The full specification is available at [`api/openapi.yaml`](./api/openapi.yaml). You can open it directly in [Swagger Editor](https://editor.swagger.io) to view interactive documentation.

### Endpoint Summary

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| `POST` | `/api/users/register` | ❌ | Register a new user |
| `POST` | `/api/users/login` | ❌ | Login, returns a JWT token |
| `GET` | `/api/users` | ✅ | List all users |
| `PATCH` | `/api/users/{id}` | ✅ | Update user data |
| `DELETE` | `/api/users/{id}` | ✅ | Delete a user |

Endpoints marked with ✅ require the following header:
```
Authorization: Bearer <token>
```

### Example Requests

**Register**
```bash
curl -X POST http://localhost:3000/api/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"azril","password":"secret123","name":"Azril"}'
```

**Login**
```bash
curl -X POST http://localhost:3000/api/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"azril","password":"secret123"}'
```

**Get All Users**
```bash
curl -X GET http://localhost:3000/api/users \
  -H "Authorization: Bearer <token>"
```

**Update User**
```bash
curl -X PATCH http://localhost:3000/api/users/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"New Name"}'
```

**Delete User**
```bash
curl -X DELETE http://localhost:3000/api/users/1 \
  -H "Authorization: Bearer <token>"
```

## Design Notes

- **Passwords** are never returned in API responses (hashed with bcrypt, and `entity.User.Password` is tagged with `json:"-"`)
- **Login error messages** for "user not found" and "wrong password" are intentionally identical (`"username or password is wrong"`) to prevent *user enumeration*
- **Error handling** distinguishes business errors (e.g. `user not found`, `username already registered`), which are returned directly to the client, from technical errors (e.g. database connection failures), which are logged server-side while the client only receives a generic message
- **JWT middleware** validates the signing algorithm to prevent *algorithm confusion attacks*
- `DELETE` responses use `204 No Content` with no body, in line with HTTP standards

## License

This project was built for learning purposes.
