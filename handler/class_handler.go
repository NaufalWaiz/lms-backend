package handler

import (
	"database/sql"
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

type ClassHandler struct {
	Service *service.ClassService
}

func (h *ClassHandler) CreateClass(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	class, err := h.Service.CreateClass(req)
	if err != nil {
		http.Error(w, "Failed to create class", http.StatusInternalServerError)
		return
	}

	// Format JSON dengan indentasi
	response, err := json.MarshalIndent(class, "", "  ")
	if err != nil {
		http.Error(w, "Failed to format response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (h *ClassHandler) DeleteClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classID := vars["id"]

	if err := h.Service.DeleteClass(classID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete class", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Class deleted successfully",
	})
}

func (h *ClassHandler) GetClasses(w http.ResponseWriter, r *http.Request) {
	classes, err := h.Service.GetClasses() // Panggil service untuk mengambil kelas
	if err != nil {
		http.Error(w, "Failed to get classes", http.StatusInternalServerError)
		return
	}

	// Kembalikan data dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

func (h *ClassHandler) GetClassByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classID := vars["id"]

	class, err := h.Service.GetClassByID(classID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get class", http.StatusInternalServerError)
		}
		return
	}

	// Kirimkan response JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}

func (h *ClassHandler) UpdateClass(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	classID := vars["id"]

	// Decode JSON request body
	var req dto.UpdateClassRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call service to update class
	updatedClass, err := h.Service.UpdateClass(classID, req)
	if err != nil {
		if err.Error() == "class not found" {
			http.Error(w, "Class not found", http.StatusNotFound)
		} else {
			fmt.Println(err)
			http.Error(w, "Failed to update class", http.StatusInternalServerError)
		}
		return
	}

	// Return updated class as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedClass)
}

func (h *ClassHandler) JoinClass(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    classIDStr := vars["class_id"]

    classID, err := strconv.Atoi(classIDStr)
    if err != nil {
        http.Error(w, "Invalid class ID", http.StatusBadRequest)
        return
    }

    userID, ok := r.Context().Value("id").(int)
    if !ok || userID == 0 {
        log.Printf("Failed to retrieve userID from context: %v", r.Context().Value("id"))
        http.Error(w, "Unauthorized: Missing or invalid user ID", http.StatusUnauthorized)
        return
    }


    username, ok := r.Context().Value("username").(string)
    if !ok || username == "" {
        http.Error(w, "Unauthorized: Missing or invalid token", http.StatusUnauthorized)
        return
    }

    role, ok := r.Context().Value("role").(string)
    if !ok || role == "" {
        http.Error(w, "Unauthorized: Missing or invalid role", http.StatusUnauthorized)
        return
    }

    log.Printf("UserID: %d, Username: %s, Role: %s is joining class ID: %d", userID, username, role, classID)

    var req struct {
        ClassCode string `json:"class_code"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    err = h.Service.JoinClass(userID, req.ClassCode)
    if err != nil {
        if err.Error() == "class not found" {
            http.Error(w, "Class not found", http.StatusNotFound)
        } else if err.Error() == "user not found" {
            http.Error(w, "User not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    response := map[string]interface{}{
        "status":     "success",
        "message":    "Successfully joined the class",
        "class_id":   classID,
        "class_code": req.ClassCode,
        "user_id":    userID,
        "username":   username,
        "role":       role,
        "timestamp":  time.Now().Format(time.RFC3339),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func (h *ClassHandler) ManageClassMembers(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    classID, err := strconv.Atoi(vars["class_id"])
    if err != nil {
        http.Error(w, "Invalid class ID", http.StatusBadRequest)
        return
    }

    var req struct {
        UserID int `json:"user_id"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    var action string
    if r.Method == http.MethodPost {
        err = h.Service.AddMember(classID, req.UserID)
        action = "added"
    } else if r.Method == http.MethodDelete {
        err = h.Service.RemoveMember(classID, req.UserID)
        action = "removed"
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Membuat response yang lebih informatif
    response := map[string]interface{}{
        "status":   "success",
        "message":  fmt.Sprintf("User %d successfully %s to class %d", req.UserID, action, classID),
        "class_id": classID,
        "user_id":  req.UserID,
        "action":   action,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func (h *ClassHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    classIDStr := vars["class_id"]

    classID, err := strconv.Atoi(classIDStr)
    if err != nil {
        http.Error(w, "Invalid class ID", http.StatusBadRequest)
        return
    }

    members, err := h.Service.GetMembers(classID)
    if err != nil {
        http.Error(w, "Failed to get members: "+err.Error(), http.StatusInternalServerError)
        return
    }

    var userResponses []dto.UserResponse
    for _, member := range members {
        userResponses = append(userResponses, dto.UserResponse{
            ID:       member.ID,
            Username: member.Username,
            Role:     member.Role,
        })
    }

    response := dto.MembersResponse{
        Status:  "success",
        Message: "Members retrieved successfully",
        Data:    userResponses,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *ClassHandler) GetClassesByStudentID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    studentIDStr := vars["student_id"]

    studentID, err := strconv.Atoi(studentIDStr)
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }

    classes, err := h.Service.GetClassesByStudentID(studentID)
    if err != nil {
        http.Error(w, "Failed to get classes: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(classes)
}

func (h *ClassHandler) CountClassesByUserID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["user_id"]

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    count, err := h.Service.CountClassesByUserID(userID)
    if err != nil {
        http.Error(w, "Failed to count classes: "+err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]string{"class_count": strconv.Itoa(count)}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}