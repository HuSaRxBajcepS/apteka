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

type AddStockRequest struct {
	MedicineID int `json:"medicine_id"`
	Quantity   int `json:"quantity"`
}

type AddMedicineRequest struct {
	Name                 string  `json:"name"`
	Price                float64 `json:"price"`
	Quantity             int     `json:"quantity"`
	RequiresPrescription bool    `json:"requires_prescription"`
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
	rows, err := h.DB.Query(`
		SELECT m.id, m.name, m.price, m.requires_prescription, COALESCE(s.quantity, 0)
		FROM medicines m
		LEFT JOIN stock s ON s.medicine_id = m.id
		ORDER BY m.id
	`)
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
		var quantity int
		err = rows.Scan(
			&id,
			&name,
			&price,
			&rx,
			&quantity,
		)
		if err != nil {
			http.Error(w, "db error", 500)
			return
		}

		result = append(result, map[string]any{
			"id":       id,
			"name":     name,
			"price":    price,
			"requires": rx,
			"quantity": quantity,
		},
		)
	}

	json.NewEncoder(w).Encode(result)
}

func (h *EmployeeHandler) AddStock(w http.ResponseWriter, r *http.Request) {
	var req AddStockRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.MedicineID <= 0 || req.Quantity <= 0 {
		http.Error(w, "invalid medicine_id or quantity", http.StatusBadRequest)
		return
	}

	_, err = h.DB.Exec(`
		INSERT INTO stock (medicine_id, quantity)
		VALUES ($1, $2)
		ON CONFLICT (medicine_id)
		DO UPDATE SET quantity = stock.quantity + EXCLUDED.quantity
	`, req.MedicineID, req.Quantity)

	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "stock updated",
	})
}

func (h *EmployeeHandler) AddMedicine(w http.ResponseWriter, r *http.Request) {
	var req AddMedicineRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Price <= 0 || req.Quantity < 0 {
		http.Error(w, "invalid medicine data", http.StatusBadRequest)
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var medicineID int
	err = tx.QueryRow(`
		INSERT INTO medicines(name, price, requires_prescription)
		VALUES($1, $2, $3)
		RETURNING id
	`, req.Name, req.Price, req.RequiresPrescription).Scan(&medicineID)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO stock(medicine_id, quantity)
		VALUES($1, $2)
	`, medicineID, req.Quantity)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"message":     "medicine created",
		"medicine_id": medicineID,
	})
}
