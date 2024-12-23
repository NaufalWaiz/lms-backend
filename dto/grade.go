package dto

type CreateGradeRequest struct {
    UserID  int `json:"user_id" validate:"required"`
    ClassID int `json:"class_id" validate:"required"`
    Grade   int `json:"grade" validate:"required"`
}

type GradeResponse struct {
    ID      int `json:"id"`
    UserID  int `json:"user_id"`
    ClassID int `json:"class_id"`
    Grade   int `json:"grade"`
}