package services

import (
	"database/sql"
	"errors"
)

type SaleService struct {
	DB           *sql.DB
	Stock        *StockService
	Prescription *PrescriptionService
	Audit        *AuditService
}

func (s *SaleService) Sell(employeeID int, patientID *int, medicineID int, quantity int, prescriptionCode *string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var price float64
	var requires bool
	err = tx.QueryRow(`SELECT price, requires_prescription FROM medicines WHERE id=$1`, medicineID).Scan(
		&price,
		&requires)
	if err != nil {
		return err
	}
	var prescriptionID *int
	if requires {
		if prescriptionCode == nil {
			return errors.New("Wymaga recepty")
		}
		id, err := s.Prescription.Validate(*prescriptionCode)
		if err != nil {
			return err
		}
		prescriptionID = &id
	}

	err = s.Stock.Decrease(tx, medicineID, quantity)
	if err != nil {
		return err
	}
	total := price * float64(quantity)
	_, err = tx.Exec(`INSERT INTO sales(employee_id,patient_id,prescription_id,medicine_id,quantity,unit_price,total_price)
	VALUES($1,$2,$3,$4,$5,$6,$7)`, employeeID, patientID, prescriptionID, medicineID, quantity, price, total)
	if err != nil {
		return err
	}
	if prescriptionID != nil {
		err = s.Prescription.Complete(tx, *prescriptionID)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	go s.Audit.Log(employeeID, "Sprzedane")

	return nil
}
