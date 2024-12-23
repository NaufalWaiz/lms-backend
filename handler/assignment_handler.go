package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"project/dto"
	"project/service"
	"project/utils"
	"strconv"

	"github.com/gorilla/mux"
)

type AssignmentHandler struct {
    Service *service.AssignmentService
}

// CreateAssignment - Membuat tugas baru
func (h *AssignmentHandler) CreateAssignment(w http.ResponseWriter, r *http.Request) {
    // Ambil class_id dari URL
    vars := mux.Vars(r)
    classIDStr := vars["class_id"]

    // Konversi class_id dari string ke integer
    classID, err := strconv.Atoi(classIDStr)
    if err != nil {
        http.Error(w, "Invalid class ID", http.StatusBadRequest)
        return
    }

    // Parse multipart form
    err = r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Get the file from the form
    file, fileHeader, err := r.FormFile("attachment")
    if err != nil && err != http.ErrMissingFile {
        http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
        return
    }

    var fileURL string
    if file != nil {
        defer file.Close()

        // Upload file to Supabase
        fileURL, err = utils.UploadImage(file, fileHeader)
        if err != nil {
            http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
            return
        }
    }

    // Decode payload JSON (hanya untuk title, description, dan due_date)
    var req dto.CreateAssignmentRequest
    req.Title = r.FormValue("title")
    req.Description = r.FormValue("description")
    req.DueDate = r.FormValue("due_date")
    req.ClassID = classID
    req.Attachment = fileURL

    // Ambil user_id dari context
    userID, ok := r.Context().Value("id").(int)
    if !ok || userID == 0 {
        http.Error(w, "Unauthorized: Missing or invalid user ID", http.StatusUnauthorized)
        return
    }

    // Panggil service untuk membuat assignment
    assignment, err := h.Service.CreateAssignment(req, userID)
    if err != nil {
        http.Error(w, "Failed to create assignment: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Berikan respons sukses
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(assignment)
}

// GetAssignmentsByClass - Mendapatkan semua tugas berdasarkan class_id
func (h *AssignmentHandler) GetAssignmentsByClass(w http.ResponseWriter, r *http.Request) {
    classID := r.URL.Query().Get("class_id")
    if classID == "" {
        http.Error(w, "class_id is required", http.StatusBadRequest)
        return
    }

    assignments, err := h.Service.GetAssignmentsByClass(classID)
    if err != nil {
        http.Error(w, "Failed to fetch assignments", http.StatusInternalServerError)
        return
    }

    var response []map[string]interface{}
    for _, assignment := range assignments {
        assignmentMap := map[string]interface{}{
            "id":          assignment.ID,
            "title":       assignment.Title,
            "description": assignment.Description,
            "due_date":    assignment.DueDate,
            "class_id":    assignment.ClassID,
            "created_at":  assignment.CreatedAt,
            "attachment":  assignment.Attachment.String,
        }
        if !assignment.Attachment.Valid {
            assignmentMap["attachment"] = ""
        }
        response = append(response, assignmentMap)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *AssignmentHandler) DeleteAssignment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    assignmentID := vars["id"]

    if err := h.Service.DeleteAssignment(assignmentID); err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Assignment not found", http.StatusNotFound)
        } else {
            http.Error(w, "Failed to delete assignment", http.StatusInternalServerError)
        }
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Assignment deleted successfully",
    })
}

func (h *AssignmentHandler) GetAssignments(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    classID := vars["class_id"] 

    assignments, err := h.Service.GetAssignments(classID) 
    if err != nil {
        fmt.Println(err)
        http.Error(w, "Failed to get assignments", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(assignments)
}

func (h *AssignmentHandler) UpdateAssignment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    classIDStr := vars["class_id"]
    assignmentIDStr := vars["assignment_id"]

    // Konversi ID dari string ke integer
    classID, err := strconv.Atoi(classIDStr)
    if err != nil {
        http.Error(w, "Invalid class ID", http.StatusBadRequest)
        return
    }

    assignmentID, err := strconv.Atoi(assignmentIDStr)
    if err != nil {
        http.Error(w, "Invalid assignment ID", http.StatusBadRequest)
        return
    }

    // Parse multipart form
    err = r.ParseMultipartForm(10 << 20) // 10 MB
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }

    // Get the file from the form
    file, fileHeader, err := r.FormFile("attachment")
    if err != nil && err != http.ErrMissingFile {
        http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
        return
    }

    var fileURL string
    if file != nil {
        defer file.Close()

        // Upload file to Supabase
        fileURL, err = utils.UploadImage(file, fileHeader)
        if err != nil {
            http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
            return
        }
    }

    // Decode payload JSON (hanya untuk title, description, dan due_date)
    var req dto.UpdateAssignmentRequest
    req.Title = r.FormValue("title")
    req.Description = r.FormValue("description")
    req.DueDate = r.FormValue("due_date")
    req.Attachment = fileURL

    // Panggil service untuk memperbarui assignment
    assignment, err := h.Service.UpdateAssignment(classID, assignmentID, req)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Assignment not found", http.StatusNotFound)
        } else {
            fmt.Println(err)
            http.Error(w, "Failed to update assignment", http.StatusInternalServerError)
        }
        return
    }

    // Berikan respons sukses
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(assignment)
}

func (h *AssignmentHandler) GetAssignmentsByUserID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["user_id"]
    classIDStr := vars["class_id"]

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    classID, err := strconv.Atoi(classIDStr)
    if err != nil {
        http.Error(w, "Invalid class ID", http.StatusBadRequest)
        return
    }

    assignments, err := h.Service.GetAssignmentsByUserID(userID, classID)
    if err != nil {
        http.Error(w, "Failed to fetch assignments: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(assignments)
}

func (h *AssignmentHandler) CountAssignmentsCreatedByUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["user_id"]

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    count, err := h.Service.CountAssignmentsCreatedByUser(userID)
    if err != nil {
        http.Error(w, "Failed to count assignments: "+err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]string{"assignment": strconv.Itoa(count)}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}