package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"apteka/internal/services"
)

type EmployeeHandler struct {
	Sale *services.SaleService
	DB   *sql.DB
}

type SellRequest struct {
	MedicineID       int     `json:"medicine_id"`
	Quantity         int     `json:"quantity"`
	PatientID        *int    `json:"patient_id"`
	PrescriptionCode *string `json:"prescription_code"`
}

func (h *EmployeeHandler) SellMedicine(w http.ResponseWriter, r *http.Request) {
	employeeID := r.Context().Value("userID").(int)
	var req SellRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", 400)
		return
	}

	err = h.Sale.Sell(employeeID, req.PatientID, req.MedicineID, req.Quantity, req.PrescriptionCode)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "sale completed"},
	)
}

func (h *EmployeeHandler) GetMedicines(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, price, requires_prescription FROM medicines`)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}

	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var id int
		var name string
		var price float64
		var rx bool
		rows.Scan(
			&id,
			&name,
			&price,
			&rx,
		)

		result = append(result, map[string]any{
			"id":       id,
			"name":     name,
			"price":    price,
			"requires": rx,
		},
		)
	}

	json.NewEncoder(w).Encode(result)
}
