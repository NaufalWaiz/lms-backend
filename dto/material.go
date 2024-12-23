package dto

type MaterialResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ClassID   int    `json:"class_id"`
	CreatedAt string `json:"created_at"`
}

type MaterialsResponse struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    []MaterialResponse `json:"data"`
}

type CreateMaterialRequest struct {
	Title      string `json:"title" validate:"required"`
	Content    string `json:"content" validate:"required"`
	ClassID    int    `json:"-"`
	Attachment string `json:"attachment"`
}

type UpdateMaterialRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
	Attachment string `json:"attachment"`
}
