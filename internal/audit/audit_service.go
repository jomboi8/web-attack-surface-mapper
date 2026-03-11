package audit

import (
	"database/sql"
)

type Log struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	Timestamp string `json:"timestamp"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Record(userID int, action, details string) error {
	_, err := s.db.Exec(
		`INSERT INTO audit_logs (user_id, action, details) VALUES (?,?,?)`,
		userID, action, details,
	)
	return err
}

func (s *Service) List(limit int) ([]Log, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query(
		`SELECT id, user_id, action, details, timestamp FROM audit_logs ORDER BY id DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var logs []Log
	for rows.Next() {
		var l Log
		var details sql.NullString
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &details, &l.Timestamp); err != nil {
			continue
		}
		l.Details = details.String
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []Log{}
	}
	return logs, nil
}
