package models

import "time"

type Prescription struct {
	ID        int       `json:"id"`
	DoctorID  int       `json:"doctor_id"`
	PatientID int       `json:"patient_id"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

type PrescriptionItem struct {
	ID             int
	PrescriptionID int
	MedicineID     int
	Quantity       int
}
