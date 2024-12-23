package service

import (
	"database/sql"
	"fmt"
	"project/dto"
	"project/model"
	"strconv"
)

type ForumService struct {
	DB *sql.DB
}

func NewForumService(db *sql.DB) *ForumService {
	return &ForumService{DB: db}
}

func (s *ForumService) CreateForum(req dto.CreateForumRequest, author, authorRole string) (*model.Forum, error) {
    query := `INSERT INTO forums (title, content, author, author_role) VALUES ($1, $2, $3, $4) RETURNING id, title, content, author, created_at, author_role`
    var forum model.Forum
    err := s.DB.QueryRow(query, req.Title, req.Content, author, authorRole).Scan(&forum.ID, &forum.Title, &forum.Content, &forum.Author, &forum.CreatedAt, &forum.AuthorRole)
    if err != nil {
        return nil, fmt.Errorf("failed to create forum: %w", err)
    }
    return &forum, nil
}

func (s *ForumService) GetForums() ([]model.Forum, error) {
	query := `SELECT id, title, content, author, created_at, author_role FROM forums`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forums []model.Forum
	for rows.Next() {
		var forum model.Forum
		if err := rows.Scan(&forum.ID, &forum.Title, &forum.Content, &forum.Author, &forum.CreatedAt, &forum.AuthorRole); err != nil {
			return nil, err
		}
		forums = append(forums, forum)
	}
	return forums, nil
}

func (s *ForumService) DeleteForum(id string) error {
	// Konversi ID ke integer jika diperlukan
	forumID, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("invalid forum ID format")
	}

	query := `DELETE FROM forums WHERE id = $1`
	result, err := s.DB.Exec(query, forumID)
	if err != nil {
		return fmt.Errorf("error deleting forum: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // Forum tidak ditemukan
	}

	return nil
}
