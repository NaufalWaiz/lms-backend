package service

import (
	"database/sql"
	"fmt"
	"project/dto"
	"project/model"
)

type AssignmentService struct {
	DB *sql.DB
}

func (s *AssignmentService) CreateAssignment(req dto.CreateAssignmentRequest, createdBy int) (*model.Assignment, error) {
    query := `INSERT INTO assignments (class_id, title, description, due_date, attachment, created_by) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, class_id, title, description, due_date, attachment, created_at, created_by`
    var assignment model.Assignment
    err := s.DB.QueryRow(query, req.ClassID, req.Title, req.Description, req.DueDate, req.Attachment, createdBy).Scan(&assignment.ID, &assignment.ClassID, &assignment.Title, &assignment.Description, &assignment.DueDate, &assignment.Attachment, &assignment.CreatedAt, &assignment.CreatedBy)
    if err != nil {
        return nil, err
    }
    return &assignment, nil
}

func (s *AssignmentService) DeleteAssignment(id string) error {
	query := `DELETE FROM assignments WHERE id = $1`
	result, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete assignment: %w", err)
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

// GetAssignmentsByClass - Mengambil tugas berdasarkan class_id
func (s *AssignmentService) GetAssignmentsByClass(classID string) ([]model.Assignment, error) {
    query := `SELECT id, title, description, due_date, class_id, created_at, attachment FROM assignments WHERE class_id = $1`
    rows, err := s.DB.Query(query, classID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch assignments: %w", err)
    }
    defer rows.Close()

    var assignments []model.Assignment
    for rows.Next() {
        var assignment model.Assignment
        if err := rows.Scan(&assignment.ID, &assignment.Title, &assignment.Description, &assignment.DueDate, &assignment.ClassID, &assignment.CreatedAt, &assignment.Attachment); err != nil {
            return nil, fmt.Errorf("failed to scan assignment: %w", err)
        }
        assignments = append(assignments, assignment)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating assignments: %w", err)
    }

    return assignments, nil
}

func (s *AssignmentService) GetAssignments(classID string) ([]model.Assignment, error) {
    query := `SELECT id, title, description, due_date, class_id, created_at, attachment, created_by FROM assignments WHERE class_id = $1`
    rows, err := s.DB.Query(query, classID)
    if err != nil {
        return nil, fmt.Errorf("failed to query assignments: %v", err)
    }
    defer rows.Close()

    var assignments []model.Assignment
    for rows.Next() {
        var assignment model.Assignment
        if err := rows.Scan(&assignment.ID, &assignment.Title, &assignment.Description, &assignment.DueDate, &assignment.ClassID, &assignment.CreatedAt, &assignment.Attachment, &assignment.CreatedBy); err != nil {
            return nil, fmt.Errorf("failed to scan assignment row: %v", err)
        }
        assignments = append(assignments, assignment)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating assignments: %v", err)
    }

    return assignments, nil
}

func (s *AssignmentService) UpdateAssignment(classID, assignmentID int, req dto.UpdateAssignmentRequest) (*model.Assignment, error) {
    query := `
        UPDATE assignments
        SET title = $1, description = $2, due_date = $3, attachment = $4, created_at = NOW()
        WHERE id = $5 AND class_id = $6
        RETURNING id, title, description, due_date, class_id, attachment, created_at
    `

    var assignment model.Assignment
    err := s.DB.QueryRow(query, req.Title, req.Description, req.DueDate, req.Attachment, assignmentID, classID).Scan(
        &assignment.ID,
        &assignment.Title,
        &assignment.Description,
        &assignment.DueDate,
        &assignment.ClassID,
        &assignment.Attachment,
        &assignment.CreatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &assignment, nil
}

func (s *AssignmentService) GetAssignmentsByUserID(userID, classID int) ([]model.Assignment, error) {
    query := `
        SELECT id, class_id, title, description, due_date, created_at, attachment, created_by
        FROM assignments
        WHERE created_by = $1 AND class_id = $2
    `
    rows, err := s.DB.Query(query, userID, classID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch assignments: %w", err)
    }
    defer rows.Close()

    var assignments []model.Assignment
    for rows.Next() {
        var assignment model.Assignment
        if err := rows.Scan(&assignment.ID, &assignment.ClassID, &assignment.Title, &assignment.Description, &assignment.DueDate, &assignment.CreatedAt, &assignment.Attachment, &assignment.CreatedBy); err != nil {
            return nil, fmt.Errorf("failed to scan assignment: %w", err)
        }
        assignments = append(assignments, assignment)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating assignments: %w", err)
    }

    return assignments, nil
}

func (s *AssignmentService) CountAssignmentsCreatedByUser(userID int) (int, error) {
    query := `
        SELECT COUNT(*)
        FROM assignments
        WHERE created_by = $1
    `
    var count int
    err := s.DB.QueryRow(query, userID).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count assignments: %w", err)
    }

    return count, nil
}