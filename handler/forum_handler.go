package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"project/dto"
	"project/middleware"
	"project/service"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ForumHandler struct {
	Service *service.ForumService
}

func NewForumHandler(service *service.ForumService) *ForumHandler {
	return &ForumHandler{Service: service}
}

var validate = validator.New()

func (h *ForumHandler) CreateForum(w http.ResponseWriter, r *http.Request) {
    // Ambil username dan role dari context (disimpan oleh middleware)
    username, ok := r.Context().Value("username").(string)
    if !ok || username == "" {
        http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
        return
    }

    authorRole, ok := r.Context().Value("role").(string)
    if !ok || authorRole == "" {
        http.Error(w, "Unauthorized: Missing or invalid role", http.StatusUnauthorized)
        return
    }

    // Decode request payload
    var req dto.CreateForumRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validasi payload
    if err := validate.Struct(req); err != nil {
        http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Panggil service untuk membuat forum
    forum, err := h.Service.CreateForum(req, username, authorRole)
    if err != nil {
        http.Error(w, "Failed to create forum: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Format respons
    forumResponse := struct {
        ID        int       `json:"id"`
        Title     string    `json:"title"`
        Content   string    `json:"content"`
        Author    string    `json:"author"`
        AuthorRole string   `json:"author_role"`
        CreatedAt time.Time `json:"created_at"`
    }{
        ID:        forum.ID,
        Title:     forum.Title,
        Content:   forum.Content,
        Author:    forum.Author,
        AuthorRole: forum.AuthorRole,
		CreatedAt: parseTime(forum.CreatedAt),
    }

    // Kirim respons
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(forumResponse)
}

func (h *ForumHandler) GetForums(w http.ResponseWriter, r *http.Request) {
	forums, err := h.Service.GetForums()
	if err != nil {
		http.Error(w, "Failed to retrieve forums", http.StatusInternalServerError)
		return
	}

	// Membuat response list forum
	var forumResponses []dto.ForumResponse
	for _, forum := range forums {
		forumResponses = append(forumResponses, dto.ForumResponse{
			ID:        forum.ID,
			Title:     forum.Title,
			Content:   forum.Content,
			CreatedAt: forum.CreatedAt,
			Author:    forum.Author,
			AuthorRole: forum.AuthorRole,
		})
	}

	// Menyiapkan response yang lebih terstruktur
	response := dto.ForumListResponse{
		Status:  "success",
		Message: "Forums retrieved successfully",
		Data:    forumResponses,
	}

	// Mengirim response dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}
	}
	return t
}

func (h *ForumHandler) DeleteForum(w http.ResponseWriter, r *http.Request) {
	// Ambil ID forum dari URL parameter
	vars := mux.Vars(r)
	forumID := vars["id"]

	// Lakukan validasi ID (misalnya pastikan ID adalah angka)
	_, err := strconv.Atoi(forumID)
	if err != nil {
		http.Error(w, "Invalid forum ID", http.StatusBadRequest)
		return
	}

	// Panggil service untuk menghapus forum berdasarkan ID
	err = h.Service.DeleteForum(forumID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Forum tidak ditemukan
			http.Error(w, "Forum not found", http.StatusNotFound)
			return
		}
		// Error lain yang terjadi
		http.Error(w, "Failed to delete forum", http.StatusInternalServerError)
		return
	}

	// Kirim respons sukses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Forum deleted successfully",
	})
}

func (h *ForumHandler) GetJWTClaims(w http.ResponseWriter, r *http.Request) {
	// Ambil claims dari context
	claims, ok := r.Context().Value("userClaims").(*middleware.Claims)
	if !ok {
		http.Error(w, "Unauthorized: Unable to extract claims", http.StatusUnauthorized)
		return
	}

	// Kirim response dengan claims dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(claims)
}
