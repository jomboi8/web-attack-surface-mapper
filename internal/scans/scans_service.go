package scans

import (
	"database/sql"
	"errors"
	"fmt"
)

type Scan struct {
	ID                   int    `json:"id"`
	TargetID             int    `json:"target_id"`
	TargetURL            string `json:"target_url"`
	Status               string `json:"status"`
	StartTime            string `json:"start_time"`
	EndTime              string `json:"end_time"`
	ScanProfile          string `json:"scan_profile"`
	TotalVulnerabilities int    `json:"total_vulnerabilities"`
	CreatedAt            string `json:"created_at"`
}

type CreateScanRequest struct {
	TargetID    int    `json:"target_id"`
	ScanProfile string `json:"scan_profile"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(req CreateScanRequest) (*Scan, error) {
	if req.TargetID == 0 {
		return nil, errors.New("target_id is required")
	}
	if req.ScanProfile == "" {
		req.ScanProfile = "full"
	}

	res, err := s.db.Exec(
		`INSERT INTO scans (target_id, scan_profile, status) VALUES (?,?,'pending')`,
		req.TargetID, req.ScanProfile,
	)
	if err != nil {
		return nil, fmt.Errorf("creating scan: %w", err)
	}
	id, _ := res.LastInsertId()
	return &Scan{ID: int(id), TargetID: req.TargetID, ScanProfile: req.ScanProfile, Status: "pending"}, nil
}

func (s *Service) List() ([]Scan, error) {
	rows, err := s.db.Query(`
		SELECT s.id, s.target_id, COALESCE(t.url,''), s.status,
		       COALESCE(s.start_time,''), COALESCE(s.end_time,''),
		       s.scan_profile, s.total_vulnerabilities, s.created_at
		FROM scans s LEFT JOIN targets t ON t.id=s.target_id
		ORDER BY s.id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scans []Scan
	for rows.Next() {
		var sc Scan
		if err := rows.Scan(&sc.ID, &sc.TargetID, &sc.TargetURL, &sc.Status,
			&sc.StartTime, &sc.EndTime, &sc.ScanProfile, &sc.TotalVulnerabilities, &sc.CreatedAt); err != nil {
			continue
		}
		scans = append(scans, sc)
	}
	if scans == nil {
		scans = []Scan{}
	}
	return scans, nil
}

func (s *Service) GetByID(scanID int) (*Scan, error) {
	var sc Scan
	err := s.db.QueryRow(`
		SELECT s.id, s.target_id, COALESCE(t.url,''), s.status,
		       COALESCE(s.start_time,''), COALESCE(s.end_time,''),
		       s.scan_profile, s.total_vulnerabilities, s.created_at
		FROM scans s LEFT JOIN targets t ON t.id=s.target_id
		WHERE s.id=?`, scanID,
	).Scan(&sc.ID, &sc.TargetID, &sc.TargetURL, &sc.Status,
		&sc.StartTime, &sc.EndTime, &sc.ScanProfile, &sc.TotalVulnerabilities, &sc.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("scan not found")
	}
	return &sc, err
}

func (s *Service) GetVulnerabilities(scanID int) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(`
		SELECT id, type, severity, endpoint, parameter, payload, description, status, created_at
		FROM vulnerabilities WHERE scan_id=? ORDER BY
		CASE severity WHEN 'CRITICAL' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'MEDIUM' THEN 3 ELSE 4 END
	`, scanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var vulns []map[string]interface{}
	for rows.Next() {
		var id int
		var vtype, severity, endpoint, status, createdAt string
		var parameter, payload, desc sql.NullString
		if err := rows.Scan(&id, &vtype, &severity, &endpoint, &parameter, &payload, &desc, &status, &createdAt); err != nil {
			continue
		}
		vulns = append(vulns, map[string]interface{}{
			"id": id, "type": vtype, "severity": severity,
			"endpoint": endpoint, "parameter": parameter.String,
			"payload": payload.String, "description": desc.String,
			"status": status, "created_at": createdAt,
		})
	}
	if vulns == nil {
		vulns = []map[string]interface{}{}
	}
	return vulns, nil
}

func (s *Service) GetStats() (map[string]interface{}, error) {
	stats := map[string]interface{}{}

	var totalScans, totalTargets, totalVulns int
	s.db.QueryRow(`SELECT COUNT(*) FROM scans`).Scan(&totalScans)
	s.db.QueryRow(`SELECT COUNT(*) FROM targets`).Scan(&totalTargets)
	s.db.QueryRow(`SELECT COUNT(*) FROM vulnerabilities`).Scan(&totalVulns)

	var critical, high, medium, low int
	s.db.QueryRow(`SELECT COUNT(*) FROM vulnerabilities WHERE severity='CRITICAL'`).Scan(&critical)
	s.db.QueryRow(`SELECT COUNT(*) FROM vulnerabilities WHERE severity='HIGH'`).Scan(&high)
	s.db.QueryRow(`SELECT COUNT(*) FROM vulnerabilities WHERE severity='MEDIUM'`).Scan(&medium)
	s.db.QueryRow(`SELECT COUNT(*) FROM vulnerabilities WHERE severity='LOW'`).Scan(&low)

	var activeScans int
	s.db.QueryRow(`SELECT COUNT(*) FROM scans WHERE status='running'`).Scan(&activeScans)

	stats["total_scans"] = totalScans
	stats["total_targets"] = totalTargets
	stats["total_vulnerabilities"] = totalVulns
	stats["active_scans"] = activeScans
	stats["severity_breakdown"] = map[string]int{
		"critical": critical, "high": high, "medium": medium, "low": low,
	}
	return stats, nil
}
