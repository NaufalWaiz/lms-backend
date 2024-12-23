package service

import (
    "database/sql"
    "fmt"
    "project/dto"
    "project/model"
)

type GradeService struct {
    DB *sql.DB
}

func NewGradeService(db *sql.DB) *GradeService {
    return &GradeService{DB: db}
}

func (s *GradeService) CreateGrade(req dto.CreateGradeRequest) (*model.Grade, error) {
    query := `INSERT INTO grades (user_id, class_id, grade) VALUES ($1, $2, $3) RETURNING id, user_id, class_id, grade`
    var grade model.Grade
    err := s.DB.QueryRow(query, req.UserID, req.ClassID, req.Grade).Scan(&grade.ID, &grade.UserID, &grade.ClassID, &grade.Grade)
    if err != nil {
        return nil, fmt.Errorf("failed to create grade: %v", err)
    }
    return &grade, nil
}