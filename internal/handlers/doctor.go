package handlers

import (
	"encoding/json"
	"net/http"

	"apteka/internal/services"
)

type DoctorHandler struct {
	Prescription   *services.PrescriptionService
	ExpirationDays int
}

type CreatePrescriptionRequest struct {
	PatientID int `json:"patient_id"`
	Medicines []struct {
		MedicineID int `json:"medicine_id"`
		Quantity   int `json:"quantity"`
	} `json:"medicines"`
}

func (h *DoctorHandler) CreatePrescription(w http.ResponseWriter, r *http.Request) {
	doctorID := r.Context().Value("userID").(int)
	var req CreatePrescriptionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	medicines := make(map[int]int)
	for _, item := range req.Medicines {
		medicines[item.MedicineID] = item.Quantity
	}
	code, id, err := h.Prescription.Create(doctorID, req.PatientID, medicines, h.ExpirationDays)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":   id,
		"code": code,
	},
	)
}
