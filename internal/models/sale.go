package models

import "time"

type Sale struct {
	ID             int       `json:"id"`
	EmployeeID     int       `json:"employee_id"`
	PatientID      *int      `json:"patient_id"`
	PrescriptionID *int      `json:"prescription_id"`
	MedicineID     int       `json:"medicine_id"`
	Quantity       int       `json:"quantity"`
	UnitPrice      float64   `json:"unit_price"`
	TotalPrice     float64   `json:"total_price"`
	SoldAt         time.Time `json:"sold_at"`
}
