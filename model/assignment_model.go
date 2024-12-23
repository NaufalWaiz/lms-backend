package model

import (
	"database/sql"
	"time"
)

type Assignment struct {
	ID          int            `json:"id"`
	ClassID     int            `json:"class_id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	DueDate     string         `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	Attachment  sql.NullString `json:"attachment"`
	CreatedBy   sql.NullInt64  `json:"created_by"`
}
