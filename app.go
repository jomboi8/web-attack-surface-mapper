package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"apguard/internal/audit"
	"apguard/internal/auth"
	"apguard/internal/reports"
	"apguard/internal/scanner"
	"apguard/internal/scans"
	"apguard/internal/scheduler"
	"apguard/internal/targets"
	"apguard/pkg/database"
	"apguard/pkg/logger"
)

// App struct — all public methods are exposed to the Wails frontend.
type App struct {
	ctx context.Context

	authService      *auth.Service
	targetService    *targets.Service
	scanService      *scans.Service
	reportService    *reports.ReportService
	auditService     *audit.Service
	schedulerService *scheduler.Scheduler
	scanEngine       *scanner.Engine

	// Track the currently authenticated user
	currentUser *auth.User
	runningMu    sync.Mutex
	runningScans map[int]bool

	// Path where the session token is persisted on disk
	sessionFile string
}

// sessionData is the JSON structure saved to disk for desktop session persistence.
type sessionData struct {
	Token string `json:"token"`
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{runningScans: make(map[int]bool)}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Determine data directory
	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".apguard")
	logDir := filepath.Join(dataDir, "logs")
	a.sessionFile = filepath.Join(dataDir, "session.json")

	// Init logger
	if err := logger.Init(logDir); err != nil {
		fmt.Println("Logger init failed:", err)
	}

	// Init database
	db, err := database.Init(dataDir)
	if err != nil {
		logger.Error("Database init failed: %v", err)
		return
	}

	// Init services
	a.authService = auth.NewService(db)
	a.targetService = targets.NewService(db)
	a.scanService = scans.NewService(db)
	a.reportService = reports.NewService(db)
	a.auditService = audit.NewService(db)
	a.scanEngine = scanner.NewEngine(db)

	// Scheduler (with scan launcher) — start the cron daemon
	a.schedulerService = scheduler.NewScheduler(db, a)
	go a.schedulerService.Start()

	// Create default admin user if no users exist
	a.ensureDefaultAdmin()

	// Restore session from disk so the user stays logged in
	a.restoreSession()

	logger.Info("APGUARD started successfully")
}

// restoreSession tries to load and validate a saved JWT from disk.
func (a *App) restoreSession() {
	data, err := os.ReadFile(a.sessionFile)
	if err != nil {
		return // no saved session
	}
	var sd sessionData
	if err := json.Unmarshal(data, &sd); err != nil || sd.Token == "" {
		return
	}
	claims, err := auth.ValidateJWT(sd.Token)
	if err != nil {
		logger.Info("Saved session token expired or invalid — requiring fresh login")
		os.Remove(a.sessionFile)
		return
	}
	// Extract user id from claims and reload user from DB
	subRaw, _ := claims["sub"]
	userID := 0
	switch v := subRaw.(type) {
	case float64:
		userID = int(v)
	case int:
		userID = v
	}
	if userID == 0 {
		return
	}
	u, err := a.authService.GetUserByID(userID)
	if err != nil {
		return
	}
	a.currentUser = u
	logger.Info("Session restored for user: %s", u.Username)
}

// saveSession writes the JWT token to disk.
func (a *App) saveSession(token string) {
	data, _ := json.Marshal(sessionData{Token: token})
	os.WriteFile(a.sessionFile, data, 0600)
}

// clearSession removes the persisted session file.
func (a *App) clearSession() {
	if a.sessionFile != "" {
		os.Remove(a.sessionFile)
	}
}

func (a *App) ensureDefaultAdmin() {
	users, err := a.authService.ListUsers()
	if err != nil || len(users) > 0 {
		return
	}
	_, err = a.authService.Register(auth.RegisterRequest{
		Username: "admin",
		Email:    "admin@apguard.local",
		Password: "admin123",
		Role:     "admin",
	})
	if err != nil {
		logger.Warn("Could not create default admin: %v", err)
	} else {
		logger.Info("Default admin user created (admin / admin123)")
	}
}

// LaunchScan implements scheduler.ScanRunner interface.
func (a *App) LaunchScan(targetID int, targetURL, profile string) {
	scan, err := a.scanService.Create(scans.CreateScanRequest{TargetID: targetID, ScanProfile: profile})
	if err != nil {
		logger.Error("LaunchScan: creating scan record: %v", err)
		return
	}
	go a.scanEngine.RunScan(scan.ID, targetURL, profile)
}

// ─── AUTH ───────────────────────────────────────────────────────────────────

func (a *App) Login(username, password string) map[string]interface{} {
	resp, err := a.authService.Login(auth.LoginRequest{Username: username, Password: password})
	if err != nil {
		return errResp(err.Error())
	}
	a.currentUser = &resp.User
	// Persist the session so the desktop app remembers the login after restart.
	a.saveSession(resp.Token)
	if a.auditService != nil {
		a.auditService.Record(resp.User.ID, "LOGIN", "User logged in")
	}
	return map[string]interface{}{
		"success": true,
		"token":   resp.Token,
		"user": map[string]interface{}{
			"id":       resp.User.ID,
			"username": resp.User.Username,
			"email":    resp.User.Email,
			"role":     resp.User.Role,
		},
	}
}

func (a *App) Logout() map[string]interface{} {
	if a.currentUser != nil && a.auditService != nil {
		a.auditService.Record(a.currentUser.ID, "LOGOUT", "User logged out")
	}
	a.currentUser = nil
	// Remove the persisted session so the user isn't auto-logged-in next time.
	a.clearSession()
	return map[string]interface{}{"success": true}
}

func (a *App) Register(username, email, password, role string) map[string]interface{} {
	user, err := a.authService.Register(auth.RegisterRequest{
		Username: username, Email: email, Password: password, Role: role,
	})
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "user": user}
}

func (a *App) GetCurrentUser() map[string]interface{} {
	if a.currentUser == nil {
		return map[string]interface{}{"success": false, "error": "not authenticated"}
	}
	return map[string]interface{}{
		"success": true,
		"token":   "", // token already stored client-side or in session file
		"user": map[string]interface{}{
			"id":       a.currentUser.ID,
			"username": a.currentUser.Username,
			"email":    a.currentUser.Email,
			"role":     a.currentUser.Role,
		},
	}
}

// ─── TARGETS ────────────────────────────────────────────────────────────────

func (a *App) CreateTarget(url, name, description string) map[string]interface{} {
	createdBy := 0
	if a.currentUser != nil {
		createdBy = a.currentUser.ID
	}
	t, err := a.targetService.Create(targets.CreateTargetRequest{
		URL: url, Name: name, Description: description, CreatedBy: createdBy,
	})
	if err != nil {
		return errResp(err.Error())
	}
	if a.auditService != nil && a.currentUser != nil {
		a.auditService.Record(a.currentUser.ID, "CREATE_TARGET", "Created target: "+url)
	}
	return map[string]interface{}{"success": true, "target": t}
}

func (a *App) ListTargets() map[string]interface{} {
	ts, err := a.targetService.List()
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "targets": ts}
}

func (a *App) ToggleTarget(id int) map[string]interface{} {
	if err := a.targetService.Toggle(id); err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true}
}

func (a *App) DeleteTarget(id int) map[string]interface{} {
	if err := a.targetService.Delete(id); err != nil {
		return errResp(err.Error())
	}
	if a.auditService != nil && a.currentUser != nil {
		a.auditService.Record(a.currentUser.ID, "DELETE_TARGET", fmt.Sprintf("Deleted target ID %d", id))
	}
	return map[string]interface{}{"success": true}
}

// ─── SCANS ───────────────────────────────────────────────────────────────────

func (a *App) StartScan(targetID int, profile string) map[string]interface{} {
	target, err := a.targetService.GetByID(targetID)
	if err != nil {
		return errResp(err.Error())
	}
	if !target.Enabled {
		return errResp("target is disabled")
	}

	scan, err := a.scanService.Create(scans.CreateScanRequest{TargetID: targetID, ScanProfile: profile})
	if err != nil {
		return errResp(err.Error())
	}

	if a.auditService != nil && a.currentUser != nil {
		a.auditService.Record(a.currentUser.ID, "START_SCAN", fmt.Sprintf("Scan #%d started on %s", scan.ID, target.URL))
	}

	// Run scan in background
	go func() {
		a.runningMu.Lock()
		a.runningScans[scan.ID] = true
		a.runningMu.Unlock()

		a.scanEngine.RunScan(scan.ID, target.URL, profile)

		a.runningMu.Lock()
		delete(a.runningScans, scan.ID)
		a.runningMu.Unlock()
	}()

	return map[string]interface{}{"success": true, "scan_id": scan.ID, "status": "running"}
}

func (a *App) ListScans() map[string]interface{} {
	ss, err := a.scanService.List()
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "scans": ss}
}

func (a *App) GetScan(scanID int) map[string]interface{} {
	sc, err := a.scanService.GetByID(scanID)
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "scan": sc}
}

func (a *App) GetVulnerabilities(scanID int) map[string]interface{} {
	vulns, err := a.scanService.GetVulnerabilities(scanID)
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "vulnerabilities": vulns}
}

func (a *App) GetDashboardStats() map[string]interface{} {
	stats, err := a.scanService.GetStats()
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "stats": stats}
}

// ─── REPORTS ─────────────────────────────────────────────────────────────────

func (a *App) GenerateReport(scanID int) map[string]interface{} {
	rpt, err := a.reportService.GenerateScanReport(scanID)
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "report": rpt}
}

func (a *App) ExportCSV(scanID int) map[string]interface{} {
	csv, err := a.reportService.ExportCSV(scanID)
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "csv": csv}
}

// ─── AUDIT ───────────────────────────────────────────────────────────────────

func (a *App) GetAuditLogs(limit int) map[string]interface{} {
	logs, err := a.auditService.List(limit)
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "logs": logs}
}

// ─── SCHEDULER ────────────────────────────────────────────────────────────────

func (a *App) CreateScheduledScan(targetID int, cronExpr, profile string) map[string]interface{} {
	task, err := a.schedulerService.CreateTask(targetID, cronExpr, profile)
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "task": task}
}

func (a *App) ListScheduledScans() map[string]interface{} {
	tasks, err := a.schedulerService.ListTasks()
	if err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true, "tasks": tasks}
}

func (a *App) DeleteScheduledScan(id int) map[string]interface{} {
	if err := a.schedulerService.DeleteTask(id); err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true}
}

func (a *App) ToggleScheduledScan(id int) map[string]interface{} {
	if err := a.schedulerService.ToggleTask(id); err != nil {
		return errResp(err.Error())
	}
	return map[string]interface{}{"success": true}
}

// ─── HELPERS ─────────────────────────────────────────────────────────────────

func errResp(msg string) map[string]interface{} {
	return map[string]interface{}{"success": false, "error": msg}
}
