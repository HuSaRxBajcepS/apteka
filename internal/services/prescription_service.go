package services

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type PrescriptionService struct {
	DB *sql.DB
}

func GenerateCode() string {
	return fmt.Sprintf("%04d", rand.Intn(10000))
}

func (s *PrescriptionService) Create(doctorID int, patientID int, medicines map[int]int, days int) (string, int, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return "", 0, err
	}
	defer tx.Rollback()
	code := GenerateCode()
	var prescriptionID int
	err = tx.QueryRow(`INSERT INTO prescriptions(doctor_id, patient_id, code) VALUES($1,$2,$3) RETURNING id`, doctorID, patientID, code).Scan(&prescriptionID)
	if err != nil {
		return "", 0, err
	}
	for medicineID, quantity := range medicines {
		_, err = tx.Exec(`INSERT INTO prescription_items(prescription_id, medicine_id, quantity) VALUES($1,$2,$3)`, prescriptionID, medicineID, quantity)
		if err != nil {
			return "", 0, err
		}
	}
	_, err = tx.Exec(`INSERT INTO prescription_status(prescription_id,expires_at) VALUES($1,$2)`, prescriptionID, time.Now().AddDate(0, 0, days))
	if err != nil {
		return "", 0, err
	}
	err = tx.Commit()
	if err != nil {
		return "", 0, err
	}

	return code, prescriptionID, nil
}

func (s *PrescriptionService) Validate(code string) (int, error) {
	var id int
	var completed bool
	var expires time.Time
	err := s.DB.QueryRow(`SELECT p.id, ps.completed, ps.expires_at FROM prescriptions p JOIN prescription_status ps ON ps.prescription_id=p.id WHERE p.code=$1`, code).Scan(
		&id,
		&completed,
		&expires)
	if err != nil {
		return 0, err
	}
	if completed {
		return 0, fmt.Errorf("prescription completed")
	}
	if time.Now().After(expires) {
		return 0, fmt.Errorf("expired")
	}
	return id, nil
}

func (s *PrescriptionService) Complete(tx *sql.Tx, prescriptionID int) error {
	_, err := tx.Exec(`UPDATE prescription_status SET completed=true WHERE prescription_id=$1`, prescriptionID)
	return err
}
