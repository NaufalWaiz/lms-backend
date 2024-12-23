package service

import (
	"database/sql"
	"fmt"
	"project/model"
)

// CommentService handles operations related to comments.
type CommentService struct {
	DB *sql.DB
}

func NewCommentService(db *sql.DB) *CommentService {
	return &CommentService{DB: db}
}


// CreateComment adds a new comment to the database.
func (s *CommentService) CreateComment(forumID int, content string, author string, authorRole string) (*model.Comment, error) {
    query := `INSERT INTO comments (content, forum_id, author, author_role) VALUES ($1, $2, $3, $4) RETURNING id, content, created_at, forum_id, author, author_role`
    comment := &model.Comment{}
    err := s.DB.QueryRow(query, content, forumID, author, authorRole).Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.ForumID, &comment.Author, &comment.AuthorRole)
    if err != nil {
        return nil, fmt.Errorf("failed to create comment: %v", err)
    }
    return comment, nil
}

// GetCommentsByForum retrieves all comments for a specific forum.
func (s *CommentService) GetComments(forumID int) ([]*model.Comment, error) {
	query := `SELECT id, content, created_at, forum_id, author, author_role FROM comments WHERE forum_id = $1 ORDER BY created_at ASC`
	rows, err := s.DB.Query(query, forumID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %v", err)
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		comment := &model.Comment{}
		if err := rows.Scan(&comment.ID, &comment.Content, &comment.CreatedAt, &comment.ForumID, &comment.Author, &comment.AuthorRole); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %v", err)
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating comments: %v", err)
	}

	return comments, nil
}

func (s *CommentService) DeleteComment(commentID int) error {
	query := `DELETE FROM comments WHERE id = $1`
	result, err := s.DB.Exec(query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("comment with ID %d not found", commentID)
	}

	return nil
}
