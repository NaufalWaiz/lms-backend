package model

import "time"

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
    Role      string    `json:"role"`
}
