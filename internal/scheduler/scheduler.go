package scheduler

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

// ScheduledTask stores info about a scheduled scan job.
type ScheduledTask struct {
	ID          int    `json:"id"`
	TargetID    int    `json:"target_id"`
	TargetURL   string `json:"target_url"`
	CronExpr    string `json:"cron_expr"`
	ScanProfile string `json:"scan_profile"`
	Enabled     bool   `json:"enabled"`
	LastRun     string `json:"last_run"`
}

// ScanRunner is an interface for things that can launch a scan.
type ScanRunner interface {
	LaunchScan(targetID int, targetURL, profile string)
}

type Scheduler struct {
	db     *sql.DB
	runner ScanRunner
	mu     sync.Mutex
	stopCh chan struct{}
}

func NewScheduler(db *sql.DB, runner ScanRunner) *Scheduler {
	return &Scheduler{db: db, runner: runner, stopCh: make(chan struct{})}
}

// ─── CRUD ──────────────────────────────────────────────────────────────────

func (s *Scheduler) CreateTask(targetID int, cronExpr, profile string) (*ScheduledTask, error) {
	res, err := s.db.Exec(
		`INSERT INTO scheduled_scans (target_id, cron_expr, scan_profile) VALUES (?,?,?)`,
		targetID, cronExpr, profile,
	)
	if err != nil {
		return nil, fmt.Errorf("creating scheduled task: %w", err)
	}
	id, _ := res.LastInsertId()
	return &ScheduledTask{
		ID: int(id), TargetID: targetID,
		CronExpr: cronExpr, ScanProfile: profile, Enabled: true,
	}, nil
}

func (s *Scheduler) ListTasks() ([]ScheduledTask, error) {
	rows, err := s.db.Query(`
		SELECT ss.id, ss.target_id, COALESCE(t.url,''), ss.cron_expr, ss.scan_profile, ss.enabled, COALESCE(ss.last_run,'')
		FROM scheduled_scans ss LEFT JOIN targets t ON t.id=ss.target_id ORDER BY ss.id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []ScheduledTask
	for rows.Next() {
		var t ScheduledTask
		var enabled int
		if err := rows.Scan(&t.ID, &t.TargetID, &t.TargetURL, &t.CronExpr, &t.ScanProfile, &enabled, &t.LastRun); err != nil {
			continue
		}
		t.Enabled = enabled == 1
		tasks = append(tasks, t)
	}
	if tasks == nil {
		tasks = []ScheduledTask{}
	}
	return tasks, nil
}

func (s *Scheduler) DeleteTask(id int) error {
	_, err := s.db.Exec(`DELETE FROM scheduled_scans WHERE id=?`, id)
	return err
}

func (s *Scheduler) ToggleTask(id int) error {
	_, err := s.db.Exec(`UPDATE scheduled_scans SET enabled = NOT enabled WHERE id=?`, id)
	return err
}

// UpdateLastRun marks the last run time for a scheduled task.
func (s *Scheduler) UpdateLastRun(taskID int) {
	s.db.Exec(`UPDATE scheduled_scans SET last_run=CURRENT_TIMESTAMP WHERE id=?`, taskID)
	log.Printf("[Scheduler] Task #%d last run updated", taskID)
}

// Stop signals the scheduler loop to stop.
func (s *Scheduler) Stop() {
	select {
	case s.stopCh <- struct{}{}:
	default:
	}
}

// ─── Cron Daemon ───────────────────────────────────────────────────────────
//
// Start runs in a goroutine and wakes up every minute to check which scheduled
// tasks are due. It supports the simple presets used in the UI:
//   @hourly   – fire at the top of each hour
//   @daily    – fire at 00:00 each day
//   @weekly   – fire at 00:00 every Sunday
//   @monthly  – fire at 00:00 on the 1st of each month
//   @midnight – alias for @daily
//
// For custom 5-field cron expressions we do a best-effort minute-level match.

func (s *Scheduler) Start() {
	if s.runner == nil {
		log.Println("[Scheduler] No runner provided — scheduler daemon inactive")
		return
	}
	log.Println("[Scheduler] Cron daemon started")

	// Align to the next whole minute
	now := time.Now()
	waitUntilNextMinute := time.Until(now.Truncate(time.Minute).Add(time.Minute))
	select {
	case <-time.After(waitUntilNextMinute):
	case <-s.stopCh:
		return
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		s.checkAndRun()
		select {
		case <-ticker.C:
		case <-s.stopCh:
			log.Println("[Scheduler] Cron daemon stopped")
			return
		}
	}
}

func (s *Scheduler) checkAndRun() {
	tasks, err := s.ListTasks()
	if err != nil {
		log.Printf("[Scheduler] ListTasks error: %v", err)
		return
	}
	now := time.Now()
	for _, task := range tasks {
		if !task.Enabled {
			continue
		}
		if isDue(task.CronExpr, task.LastRun, now) {
			log.Printf("[Scheduler] Firing task #%d (%s) on target %s", task.ID, task.CronExpr, task.TargetURL)
			s.UpdateLastRun(task.ID)
			if s.runner != nil {
				go s.runner.LaunchScan(task.TargetID, task.TargetURL, task.ScanProfile)
			}
		}
	}
}

// isDue returns true if the given cron expression should fire *now* (at time t),
// taking into account the last time it ran (lastRunStr).
func isDue(expr, lastRunStr string, t time.Time) bool {
	// Parse last run time; if never run, treat as zero time.
	var lastRun time.Time
	if lastRunStr != "" {
		// SQLite CURRENT_TIMESTAMP format: "2006-01-02 15:04:05"
		for _, layout := range []string{
			"2006-01-02 15:04:05",
			time.RFC3339,
		} {
			if parsed, err := time.ParseInLocation(layout, lastRunStr, time.Local); err == nil {
				lastRun = parsed
				break
			}
		}
	}

	switch expr {
	case "@hourly":
		// Fire once per hour — at the top of each hour.
		scheduled := t.Truncate(time.Hour)
		return !scheduled.Before(lastRun.Truncate(time.Hour).Add(time.Hour)) && t.Minute() == 0

	case "@daily", "@midnight":
		// Fire once per day at 00:00.
		if t.Hour() != 0 || t.Minute() != 0 {
			return false
		}
		today := t.Truncate(24 * time.Hour)
		return lastRun.Before(today)

	case "@weekly":
		// Fire at 00:00 every Sunday.
		if t.Weekday() != time.Sunday || t.Hour() != 0 || t.Minute() != 0 {
			return false
		}
		weekStart := t.Truncate(24 * time.Hour)
		return lastRun.Before(weekStart)

	case "@monthly":
		// Fire at 00:00 on the 1st of each month.
		if t.Day() != 1 || t.Hour() != 0 || t.Minute() != 0 {
			return false
		}
		monthStart := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
		return lastRun.Before(monthStart)

	default:
		// Basic 5-field cron: minute hour dom month dow
		return matchCron5(expr, t, lastRun)
	}
}

// matchCron5 does a minimal 5-field cron match (minute, hour, dom, month, dow).
// Supports * and exact values. Ranges, lists and step values are not supported.
func matchCron5(expr string, t time.Time, lastRun time.Time) bool {
	// Prevent re-firing within the same minute.
	if !lastRun.IsZero() && t.Sub(lastRun) < time.Minute {
		return false
	}
	fields := splitFields(expr)
	if len(fields) != 5 {
		return false
	}
	return matchField(fields[0], t.Minute()) &&
		matchField(fields[1], t.Hour()) &&
		matchField(fields[2], t.Day()) &&
		matchField(fields[3], int(t.Month())) &&
		matchField(fields[4], int(t.Weekday()))
}

func matchField(field string, value int) bool {
	if field == "*" {
		return true
	}
	var n int
	if _, err := fmt.Sscanf(field, "%d", &n); err != nil {
		return false
	}
	return n == value
}

func splitFields(s string) []string {
	var fields []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ' ' || s[i] == '\t' {
			if i > start {
				fields = append(fields, s[start:i])
			}
			start = i + 1
		}
	}
	return fields
}

// CronDescription returns a human-readable description for common cron presets.
func CronDescription(expr string) string {
	descriptions := map[string]string{
		"@hourly":   "Every hour",
		"@daily":    "Every day at midnight",
		"@weekly":   "Every week on Sunday",
		"@monthly":  "First day of every month",
		"@midnight": "Every day at midnight",
	}
	if d, ok := descriptions[expr]; ok {
		return d
	}
	return expr
}
