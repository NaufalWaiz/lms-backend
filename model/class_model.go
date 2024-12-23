package model

import "time"

type Class struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	JadwalKelas string    `json:"jadwal_kelas"`
	CreatedAt   time.Time `json:"created_at"`
	Teacher     string    `json:"teacher"`
	ClassCode   string    `json:"class_code"`
}
