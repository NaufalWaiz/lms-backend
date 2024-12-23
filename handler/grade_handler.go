package handler

import (
    "encoding/json"
    "net/http"
    "project/dto"
    "project/service"
)

type GradeHandler struct {
    Service *service.GradeService
}

func NewGradeHandler(service *service.GradeService) *GradeHandler {
    return &GradeHandler{Service: service}
}

func (h *GradeHandler) CreateGrade(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateGradeRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    grade, err := h.Service.CreateGrade(req)
    if err != nil {
        http.Error(w, "Failed to create grade: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(grade)
}