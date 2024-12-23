package dto

type CreateClassRequest struct {
	Name        string `json:"name" validate:"required"`
	JadwalKelas string `json:"jadwal_kelas" validate:"required"`
	Teacher     string `json:"teacher"`
	ClassCode   string `json:"class_code"`
}

type AssignmentsResponse struct {
	Status  string               `json:"status"`
	Message string               `json:"message"`
	Data    []AssignmentResponse `json:"data"`
}

type ClassResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	JadwalKelas string `json:"jadwal_kelas"`
	CreatedAt   string `json:"created_at"`
	Teacher     string `json:"teacher"`
	ClassCode   string `json:"class_code"`
}

type UpdateClassRequest struct {
	Name        string `json:"name"`
	JadwalKelas string `json:"jadwal_kelas"`
	Teacher     string `json:"teacher"`
	ClassCode   string `json:"class_code"`
}
