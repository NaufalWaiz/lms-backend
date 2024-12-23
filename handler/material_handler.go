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

	// "strings"

	"github.com/gorilla/mux"
)

type MaterialHandler struct {
	Service *service.MaterialService
}

// CreateMaterial - Membuat materi baru
func (h *MaterialHandler) CreateMaterial(w http.ResponseWriter, r *http.Request) {
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
    if err != nil {
        http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Upload file to Supabase
    fileURL, err := utils.UploadImage(file, fileHeader)
    if err != nil {
        http.Error(w, "Failed to upload file: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Decode payload JSON (hanya untuk title dan content)
    var req dto.CreateMaterialRequest
    req.Title = r.FormValue("title")
    req.Content = r.FormValue("content")
    req.ClassID = classID
    req.Attachment = fileURL

    // Panggil service untuk membuat material
    material, err := h.Service.CreateMaterial(req)
    if err != nil {
        http.Error(w, "Failed to create material: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Berikan respons sukses
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(material)
}


// GetMaterialsByClass - Mendapatkan semua materi berdasarkan class_id
func (h *MaterialHandler) GetMaterialsByClass(w http.ResponseWriter, r *http.Request) {
    classID := r.URL.Query().Get("class_id")
    if classID == "" {
        http.Error(w, "class_id is required", http.StatusBadRequest)
        return
    }

    materials, err := h.Service.GetMaterialsByClass(classID)
    if err != nil {
        http.Error(w, "Failed to fetch materials", http.StatusInternalServerError)
        return
    }

    var response []map[string]interface{}
    for _, material := range materials {
        materialMap := map[string]interface{}{
            "id":        material.ID,
            "title":     material.Title,
            "content":   material.Content,
            "class_id":  material.ClassID,
            "created_at": material.CreatedAt,
            "attachment": material.Attachment.String,
        }
        if !material.Attachment.Valid {
            materialMap["attachment"] = ""
        }
        response = append(response, materialMap)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// DeleteMaterial - Menghapus materi berdasarkan ID
func (h *MaterialHandler) DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	materialID := vars["id"]

	if err := h.Service.DeleteMaterial(materialID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Material not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete material", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Material deleted successfully",
	})
}

func (h *MaterialHandler) GetMaterials(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    classID := vars["class_id"] // Ambil class_id dari URL parameter

    materials, err := h.Service.GetMaterials(classID) // Panggil service untuk mengambil materi berdasarkan class_id
    if err != nil {
        http.Error(w, "Failed to get materials", http.StatusInternalServerError)
        return
    }

    // Kembalikan data dalam format JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(materials)
}

func (h *MaterialHandler) UpdateMaterial(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    classIDStr := vars["class_id"]
    materialIDStr := vars["material_id"]

    // Konversi ID dari string ke integer
    classID, err := strconv.Atoi(classIDStr)
    if err != nil {
        http.Error(w, "Invalid class ID", http.StatusBadRequest)
        return
    }

    materialID, err := strconv.Atoi(materialIDStr)
    if err != nil {
        http.Error(w, "Invalid material ID", http.StatusBadRequest)
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

    // Decode payload JSON (hanya untuk title dan content)
    var req dto.UpdateMaterialRequest
    req.Title = r.FormValue("title")
    req.Content = r.FormValue("content")
    req.Attachment = fileURL

    // Panggil service untuk memperbarui material
    material, err := h.Service.UpdateMaterial(classID, materialID, req)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Material not found", http.StatusNotFound)
        } else {
            fmt.Println(err)
            http.Error(w, "Failed to update material", http.StatusInternalServerError)
        }
        return
    }

    // Berikan respons sukses
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(material)
}