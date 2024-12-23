package service

import (
    "database/sql"
    "fmt"
    "project/dto"
)

type RapotService struct {
    DB *sql.DB
}

func NewRapotService(db *sql.DB) *RapotService {
    return &RapotService{DB: db}
}

func (s *RapotService) GetRapotByUserID(userID int) ([]dto.RapotResponse, error) {
    query := `
        SELECT c.name, g.grade
        FROM classes c
        JOIN grades g ON c.id = g.class_id
        WHERE g.user_id = $1
    `
    rows, err := s.DB.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch rapot: %v", err)
    }
    defer rows.Close()

    var rapots []dto.RapotResponse
    for rows.Next() {
        var rapot dto.RapotResponse
        if err := rows.Scan(&rapot.ClassName, &rapot.Grade); err != nil {
            return nil, fmt.Errorf("failed to scan rapot: %v", err)
        }
        rapots = append(rapots, rapot)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rapots: %v", err)
    }

    return rapots, nil
}