package service

import (
	"database/sql"
	"fmt"
	// "project/config"
	"project/dto"
	"project/model"
)

type MaterialService struct {
	DB *sql.DB
}

func (s *MaterialService) CreateMaterial(req dto.CreateMaterialRequest) (*model.Material, error) {
    query := `INSERT INTO materials (title, content, class_id, attachment) VALUES ($1, $2, $3, $4) RETURNING id, title, content, class_id, attachment, created_at`
    var material model.Material
    err := s.DB.QueryRow(query, req.Title, req.Content, req.ClassID, req.Attachment).Scan(&material.ID, &material.Title, &material.Content, &material.ClassID, &material.Attachment, &material.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &material, nil
}


func (s *MaterialService) DeleteMaterial(id string) error {
	query := `DELETE FROM materials WHERE id = $1`
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete material: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *MaterialService) GetMaterialsByClass(classID string) ([]model.Material, error) {
    query := `SELECT id, title, content, class_id, created_at, attachment FROM materials WHERE class_id = $1`
    rows, err := s.DB.Query(query, classID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch materials: %w", err)
    }
    defer rows.Close()

    var materials []model.Material
    for rows.Next() {
        var material model.Material
        if err := rows.Scan(&material.ID, &material.Title, &material.Content, &material.ClassID, &material.CreatedAt, &material.Attachment); err != nil {
            return nil, fmt.Errorf("failed to scan material: %w", err)
        }
        materials = append(materials, material)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating materials: %w", err)
    }

    return materials, nil
}

func (s *MaterialService) GetMaterials(classID string) ([]model.Material, error) {
    query := `SELECT id, title, content, class_id, created_at, attachment FROM materials WHERE class_id = $1`
    rows, err := s.DB.Query(query, classID)
    if err != nil {
        return nil, fmt.Errorf("failed to query materials: %v", err)
    }
    defer rows.Close()

    var materials []model.Material
    for rows.Next() {
        var material model.Material
        if err := rows.Scan(&material.ID, &material.Title, &material.Content, &material.ClassID, &material.CreatedAt, &material.Attachment); err != nil {
            return nil, fmt.Errorf("failed to scan material row: %v", err)
        }
        materials = append(materials, material)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating materials: %v", err)
    }

    return materials, nil
}

func (s *MaterialService) UpdateMaterial(classID, materialID int, req dto.UpdateMaterialRequest) (*model.Material, error) {
    query := `
        UPDATE materials
        SET title = $1, content = $2, attachment = $3, created_at = NOW()
        WHERE id = $4 AND class_id = $5
        RETURNING id, title, content, class_id, attachment, created_at
    `

    var material model.Material
    err := s.DB.QueryRow(query, req.Title, req.Content, req.Attachment, materialID, classID).Scan(
        &material.ID,
        &material.Title,
        &material.Content,
        &material.ClassID,
        &material.Attachment,
        &material.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &material, nil
}
