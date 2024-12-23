package dto

type CreateAssignmentRequest struct {
    ClassID     int    `json:"class_id" validate:"required"`
    Title       string `json:"title" validate:"required"`
    Description string `json:"description"`
    DueDate     string `json:"due_date"`
    Attachment  string `json:"attachment"`
}

type AssignmentResponse struct {
    ID          int    `json:"id"`
    ClassID     int    `json:"class_id"`
    Title       string `json:"title"`
    Description string `json:"description"`
    DueDate     string `json:"due_date"`
    CreatedAt   string `json:"created_at"`
}

type UpdateAssignmentRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    DueDate     string `json:"due_date"`
    Attachment  string `json:"attachment"`
}