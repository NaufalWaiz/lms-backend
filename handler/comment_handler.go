package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project/dto"
	"project/service"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type CommentHandler struct {
	CommentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{CommentService: commentService}
}

// CreateComment handles POST /forums/{forumID}/comments requests.
func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
    log.Printf("Method: %s, URL: %s", r.Method, r.URL.Path)

    // Ambil username dari context (disimpan oleh middleware)
    username, ok := r.Context().Value("username").(string)
    if !ok || username == "" {
        http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
        return
    }

    // Ambil role dari context (disimpan oleh middleware)
    authorRole, ok := r.Context().Value("role").(string)
    if !ok || authorRole == "" {
        http.Error(w, "Unauthorized: Missing or invalid role", http.StatusUnauthorized)
        return
    }

    // Ambil forumID dari URL
    vars := mux.Vars(r)
    forumIDStr, ok := vars["forumID"]
    if !ok {
        http.Error(w, "Invalid forum ID", http.StatusBadRequest)
        return
    }

    forumID, err := strconv.Atoi(forumIDStr)
    if err != nil {
        http.Error(w, "Invalid forum ID", http.StatusBadRequest)
        return
    }

    // Decode request payload
    var req dto.CreateCommentRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validasi request (sesuai kebutuhan)
    if req.Content == "" {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Panggil service untuk menyimpan komentar
    comment, err := h.CommentService.CreateComment(forumID, req.Content, username, authorRole)
    if err != nil {
        http.Error(w, "Failed to create comment: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Format respons
    commentResponse := struct {
        ID        int       `json:"id"`
        Content   string    `json:"content"`
        Author    string    `json:"author"`
        AuthorRole string   `json:"author_role"`
        ForumID   int       `json:"forum_id"`
        CreatedAt time.Time `json:"created_at"`
    }{
        ID:        comment.ID,
        Content:   comment.Content,
        Author:    comment.Author,
        AuthorRole: comment.AuthorRole,
        ForumID:   comment.ForumID,
        CreatedAt: comment.CreatedAt,
    }

    // Kirim respons
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(commentResponse)
}

// GetComments handles GET /forums/{forum_id}/comments requests.
func (h *CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	forumID, err := strconv.Atoi(vars["forum_id"])
	if err != nil {
		http.Error(w, "Invalid forum ID", http.StatusBadRequest)
		return
	}

	comments, err := h.CommentService.GetComments(forumID)
	if err != nil {
		http.Error(w, "Failed to retrieve comments", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

func (h *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    commentID, err := strconv.Atoi(vars["comment_id"])
    if err != nil {
        log.Printf("Error converting comment_id: %v", vars["comment_id"])
        http.Error(w, "Invalid comment ID", http.StatusBadRequest)
        return
    }

    log.Printf("Deleting comment with ID: %d", commentID)

    err = h.CommentService.DeleteComment(commentID)
    if err != nil {
        if err.Error() == fmt.Sprintf("comment with ID %d not found", commentID) {
            http.Error(w, "Comment not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Comment deleted successfully",
	})
}

