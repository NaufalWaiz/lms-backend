package dto

type CreateForumRequest struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Author  string `json:"author"`
	AuthorRole string `json:"author_role"`
}

type ForumResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Author    string `json:"author"`
	AuthorRole string `json:"author_role"`
}

type ForumListResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    []ForumResponse 	`json:"data"`
}

type CommentRequest struct {
	Content string `json:"content" validate:"required"`
	Author  string `json:"author"`
	AuthorRole string `json:"author_role"`
}

type CreateCommentRequest struct {
    Content string `json:"content"`
    ForumID int    `json:"forum_id"`
}

type UserResponse struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Role     string `json:"role"`
}

type MembersResponse struct {
    Status  string         `json:"status"`
    Message string         `json:"message"`
    Data    []UserResponse `json:"data"`
}