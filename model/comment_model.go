package model

import "time"

// Comment represents a comment in the forum system.
type Comment struct {
    ID        int       `json:"id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
    ForumID   int       `json:"forum_id"`
    Author    string    `json:"author"`
    AuthorRole string    `json:"author_role"`
}
