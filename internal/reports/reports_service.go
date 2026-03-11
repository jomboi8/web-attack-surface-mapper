package reports

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"strings"
	"time"
)

type ReportService struct {
	db *sql.DB
}

func NewService(db *sql.DB) *ReportService {
	return &ReportService{db: db}
}

// GenerateScanReport produces a structured report for a given scan.
func (r *ReportService) GenerateScanReport(scanID int) (map[string]interface{}, error) {
	var scanURL, status, startTime, endTime, scanProfile string
	var totalVulns int
	err := r.db.QueryRow(`
		SELECT COALESCE(t.url,''), s.status, COALESCE(s.start_time,''),
			   COALESCE(s.end_time,''), s.scan_profile, s.total_vulnerabilities
		FROM scans s LEFT JOIN targets t ON t.id=s.target_id WHERE s.id=?`, scanID,
	).Scan(&scanURL, &status, &startTime, &endTime, &scanProfile, &totalVulns)
	if err != nil {
		return nil, fmt.Errorf("scan not found: %w", err)
	}

	// Vulnerability breakdown
	rows, err := r.db.Query(`
		SELECT id, type, severity, endpoint, parameter, payload, description, status, created_at
		FROM vulnerabilities WHERE scan_id=? ORDER BY
		CASE severity WHEN 'CRITICAL' THEN 1 WHEN 'HIGH' THEN 2 WHEN 'MEDIUM' THEN 3 ELSE 4 END`, scanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vulns []map[string]interface{}
	severityCounts := map[string]int{"CRITICAL": 0, "HIGH": 0, "MEDIUM": 0, "LOW": 0}
	typeCounts := map[string]int{}

	for rows.Next() {
		var id int
		var vtype, severity, endpoint, vstatus, createdAt string
		var parameter, payload, desc sql.NullString
		if err := rows.Scan(&id, &vtype, &severity, &endpoint, &parameter, &payload, &desc, &vstatus, &createdAt); err != nil {
			continue
		}
		vulns = append(vulns, map[string]interface{}{
			"id": id, "type": vtype, "severity": severity, "endpoint": endpoint,
			"parameter": parameter.String, "payload": payload.String,
			"description": desc.String, "status": vstatus, "created_at": createdAt,
		})
		severityCounts[severity]++
		typeCounts[vtype]++
	}

	if vulns == nil {
		vulns = []map[string]interface{}{}
	}

	riskScore := severityCounts["CRITICAL"]*10 + severityCounts["HIGH"]*5 + severityCounts["MEDIUM"]*2 + severityCounts["LOW"]*1
	riskLevel := "LOW"
	switch {
	case riskScore >= 20:
		riskLevel = "CRITICAL"
	case riskScore >= 10:
		riskLevel = "HIGH"
	case riskScore >= 4:
		riskLevel = "MEDIUM"
	}

	return map[string]interface{}{
		"scan_id":          scanID,
		"target_url":       scanURL,
		"status":           status,
		"start_time":       startTime,
		"end_time":         endTime,
		"scan_profile":     scanProfile,
		"total_vulns":      totalVulns,
		"severity_counts":  severityCounts,
		"type_counts":      typeCounts,
		"risk_score":       riskScore,
		"risk_level":       riskLevel,
		"vulnerabilities":  vulns,
		"generated_at":     time.Now().Format(time.RFC3339),
	}, nil
}

// ExportCSV exports vulnerabilities of a scan as CSV string.
func (r *ReportService) ExportCSV(scanID int) (string, error) {
	rows, err := r.db.Query(`
		SELECT id, type, severity, endpoint, COALESCE(parameter,''), COALESCE(payload,''), COALESCE(description,''), status, created_at
		FROM vulnerabilities WHERE scan_id=?`, scanID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var sb strings.Builder
	w := csv.NewWriter(&sb)
	w.Write([]string{"ID", "Type", "Severity", "Endpoint", "Parameter", "Payload", "Description", "Status", "Created At"})

	for rows.Next() {
		var id int
		var vtype, severity, endpoint, parameter, payload, desc, status, createdAt string
		if err := rows.Scan(&id, &vtype, &severity, &endpoint, &parameter, &payload, &desc, &status, &createdAt); err != nil {
			continue
		}
		w.Write([]string{fmt.Sprintf("%d", id), vtype, severity, endpoint, parameter, payload, desc, status, createdAt})
	}
	w.Flush()
	return sb.String(), nil
}
