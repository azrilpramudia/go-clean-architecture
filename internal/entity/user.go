package entity

import (
	"time"
)

type User struct {
	ID int64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Name string `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}