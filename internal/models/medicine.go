package models

type Medicine struct {
	ID                   int     `json:"id"`
	Name                 string  `json:"name"`
	Price                float64 `json:"price"`
	RequiresPrescription bool    `json:"requires_prescription"`
}
