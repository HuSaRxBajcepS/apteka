package services

import (
	"database/sql"
	"time"
)

type PatientService struct {
	DB *sql.DB
}

type PatientMedicineResponse struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

type PatientPrescriptionResponse struct {
	ID        int                       `json:"id"`
	Code      string                    `json:"code"`
	CreatedAt time.Time                 `json:"created_at"`
	ExpiresAt time.Time                 `json:"expires_at"`
	Completed bool                      `json:"completed"`
	Medicines []PatientMedicineResponse `json:"medicines"`
}

type PatientMeResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func (s *PatientService) MyPrescriptions(patientID int) ([]PatientPrescriptionResponse, error) {
	rows, err := s.DB.Query(`
		SELECT p.id, p.code, p.created_at, ps.expires_at, ps.completed
		FROM prescriptions p
		JOIN prescription_status ps ON ps.prescription_id = p.id
		WHERE p.patient_id = $1
		ORDER BY p.created_at DESC
	`, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescriptions []PatientPrescriptionResponse

	for rows.Next() {
		var prescription PatientPrescriptionResponse

		err = rows.Scan(
			&prescription.ID,
			&prescription.Code,
			&prescription.CreatedAt,
			&prescription.ExpiresAt,
			&prescription.Completed,
		)
		if err != nil {
			return nil, err
		}

		items, err := s.DB.Query(`
			SELECT m.name, pi.quantity
			FROM prescription_items pi
			JOIN medicines m ON m.id = pi.medicine_id
			WHERE pi.prescription_id = $1
			ORDER BY m.name
		`, prescription.ID)
		if err != nil {
			return nil, err
		}

		for items.Next() {
			var medicine PatientMedicineResponse

			err = items.Scan(&medicine.Name, &medicine.Quantity)
			if err != nil {
				items.Close()
				return nil, err
			}

			prescription.Medicines = append(prescription.Medicines, medicine)
		}

		items.Close()

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, rows.Err()
}

func (s *PatientService) Me(patientID int) (*PatientMeResponse, error) {
	var patient PatientMeResponse

	err := s.DB.QueryRow(`
		SELECT id, full_name, email, role
		FROM users
		WHERE id = $1
	`, patientID).Scan(
		&patient.ID,
		&patient.FullName,
		&patient.Email,
		&patient.Role,
	)

	if err != nil {
		return nil, err
	}

	return &patient, nil
}
