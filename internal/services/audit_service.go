package services

import (
	"database/sql"
)

type AuditService struct {
	DB *sql.DB
}

func (s *AuditService) Log(userID int, action string) {
	_, _ = s.DB.Exec(`INSERT INTO audit_logs(user_id, action) VALUES($1,$2)`, userID, action)
}
