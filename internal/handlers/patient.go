package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type PatientHandler struct {
	DB *sql.DB
}

func (h *PatientHandler) GetPrescription(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	rows, err := h.DB.Query(`SELECT m.name,pi.quantity FROM prescriptions p JOIN prescription_items pi
	ON pi.prescription_id=p.id JOIN medicines m ON m.id=pi.medicine_id WHERE p.code=$1`, code)
	if err != nil {
		http.Error(w, "db error", 500)
		return
	}
	defer rows.Close()
	var result []map[string]any
	for rows.Next() {
		var name string
		var quantity int
		rows.Scan(&name, &quantity)
		result = append(result, map[string]any{
			"name":     name,
			"quantity": quantity,
		},
		)
	}
	json.NewEncoder(w).Encode(result)
}

func (h *PatientHandler) GetOTC(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT id, name, price FROM medicines WHERE requires_prescription=false`)
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
		rows.Scan(&id, &name, &price)
		result = append(result, map[string]any{
			"id":    id,
			"name":  name,
			"price": price,
		},
		)
	}
	json.NewEncoder(w).Encode(result)
}
