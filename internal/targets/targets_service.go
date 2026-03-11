package targets

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Target struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	CreatedBy   int    `json:"created_by"`
	CreatedAt   string `json:"created_at"`
}

type CreateTargetRequest struct {
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   int    `json:"created_by"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(req CreateTargetRequest) (*Target, error) {
	if req.URL == "" || req.Name == "" {
		return nil, errors.New("url and name are required")
	}
	if _, err := url.ParseRequestURI(req.URL); err != nil {
		return nil, errors.New("invalid url format")
	}

	res, err := s.db.Exec(
		`INSERT INTO targets (url, name, description, created_by) VALUES (?,?,?,?)`,
		req.URL, req.Name, req.Description, req.CreatedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("creating target: %w", err)
	}
	id, _ := res.LastInsertId()
	return &Target{
		ID: int(id), URL: req.URL, Name: req.Name,
		Description: req.Description, Enabled: true,
		CreatedBy: req.CreatedBy, CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *Service) List() ([]Target, error) {
	rows, err := s.db.Query(`SELECT id, url, name, description, enabled, created_by, created_at FROM targets ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var targets []Target
	for rows.Next() {
		var t Target
		var enabled int
		if err := rows.Scan(&t.ID, &t.URL, &t.Name, &t.Description, &enabled, &t.CreatedBy, &t.CreatedAt); err != nil {
			continue
		}
		t.Enabled = enabled == 1
		targets = append(targets, t)
	}
	if targets == nil {
		targets = []Target{}
	}
	return targets, nil
}

func (s *Service) GetByID(id int) (*Target, error) {
	var t Target
	var enabled int
	err := s.db.QueryRow(
		`SELECT id, url, name, description, enabled, created_by, created_at FROM targets WHERE id=?`, id,
	).Scan(&t.ID, &t.URL, &t.Name, &t.Description, &enabled, &t.CreatedBy, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("target not found")
	}
	if err != nil {
		return nil, err
	}
	t.Enabled = enabled == 1
	return &t, nil
}

func (s *Service) Toggle(id int) error {
	_, err := s.db.Exec(`UPDATE targets SET enabled = NOT enabled WHERE id=?`, id)
	return err
}

func (s *Service) Delete(id int) error {
	res, err := s.db.Exec(`DELETE FROM targets WHERE id=?`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("target not found")
	}
	return nil
}
