package services

import (
	"apteka/internal/auth"
	"apteka/internal/models"
	"database/sql"
	"errors"
)

type AuthService struct {
	DB *sql.DB
}

func (s *AuthService) Login(email string, password string) (string, string, error) {
	var user models.User
	err := s.DB.QueryRow(`SELECT id, email, password_hash, role FROM users WHERE email=$1`, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", errors.New("invalid credentials")
		}
		return "", "", err
	}

	ok := auth.VerifyPassword(password, user.PasswordHash)
	if !ok {
		return "", "", errors.New("invalid credentials")
	}

	return user.Role, user.Email, nil
}

func (s *AuthService) CreateUser(name string, email string, password string, role string) error {
	hash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`INSERT INTO users(full_name, email, password_hash, role) VALUES ($1,$2,$3,$4)`, name, email, hash, role)
	return err
}
