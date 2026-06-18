package services

import (
	"database/sql"
	"errors"
)

type StockService struct {
	DB *sql.DB
}

func (s *StockService) GetStock(medicineID int) (int, error) {
	var quantity int
	err := s.DB.QueryRow(`SELECT quantity FROM stock WHERE medicine_id=$1`, medicineID).Scan(&quantity)
	return quantity, err

}

func (s *StockService) Decrease(tx *sql.Tx, medicineID int, quantity int) error {
	var current int
	err := tx.QueryRow(`SELECT quantity FROM stock WHERE medicine_id=$1 FOR UPDATE`, medicineID).Scan(&current)
	if err != nil {
		return err
	}
	if current < quantity {
		return errors.New("Leku nie ma na stanie")
	}
	_, err = tx.Exec(`UPDATE stock SET quantity=quantity-$1 WHERE medicine_id=$2`, quantity, medicineID)
	return err
}

func (s *StockService) Increase(tx *sql.Tx, medicineID int, quantity int) error {
	_, err := tx.Exec(`UPDATE stock SET quantity=quantity+$1 WHERE medicine_id=$2`, quantity, medicineID)
	return err
}
