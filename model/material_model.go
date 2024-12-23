package model

import (
	"database/sql"
	"time"
)

type Material struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	ClassID   int       `json:"class_id"`
	CreatedAt time.Time `json:"created_at"`
	Attachment sql.NullString  `json:"attachment"`
}
