package model

type Grade struct {
    ID      int `json:"id"`
    UserID  int `json:"user_id"`
    ClassID int `json:"class_id"`
    Grade   int `json:"grade"`
}