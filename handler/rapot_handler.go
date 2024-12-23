package handler

import (
    "encoding/json"
    "net/http"
    "project/service"
    "strconv"

    "github.com/gorilla/mux"
)

type RapotHandler struct {
    Service *service.RapotService
}

func NewRapotHandler(service *service.RapotService) *RapotHandler {
    return &RapotHandler{Service: service}
}

func (h *RapotHandler) GetRapotByUserID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userIDStr := vars["user_id"]

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    rapots, err := h.Service.GetRapotByUserID(userID)
    if err != nil {
        http.Error(w, "Failed to fetch rapot: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(rapots)
}