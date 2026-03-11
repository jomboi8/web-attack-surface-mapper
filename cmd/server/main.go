package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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

// ──────────────────────────────────────────────────────────────
// Services (global singletons for the server)
// ──────────────────────────────────────────────────────────────
var (
	authSvc      *auth.Service
	targetSvc    *targets.Service
	scanSvc      *scans.Service
	reportSvc    *reports.ReportService
	auditSvc     *audit.Service
	schedulerSvc *scheduler.Scheduler
	scanEngine   *scanner.Engine
)

// in-memory session: map[token]user
var sessions = map[string]*auth.User{}

// ──────────────────────────────────────────────────────────────
// Main
// ──────────────────────────────────────────────────────────────
func main() {
	home, _ := os.UserHomeDir()
	dataDir := filepath.Join(home, ".apguard")
	logDir := filepath.Join(dataDir, "logs")

	if err := logger.Init(logDir); err != nil {
		log.Println("Logger init:", err)
	}

	db, err := database.Init(dataDir)
	if err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	authSvc = auth.NewService(db)
	targetSvc = targets.NewService(db)
	scanSvc = scans.NewService(db)
	reportSvc = reports.NewService(db)
	auditSvc = audit.NewService(db)
	scanEngine = scanner.NewEngine(db)
	schedulerSvc = scheduler.NewScheduler(db, nil)

	// Create default admin if needed
	users, _ := authSvc.ListUsers()
	if len(users) == 0 {
		authSvc.Register(auth.RegisterRequest{Username: "admin", Email: "admin@apguard.local", Password: "admin123", Role: "admin"})
		log.Println("[APGUARD] Default admin user created: admin / admin123")
	}

	mux := http.NewServeMux()
	registerRoutes(mux)

	// Serve frontend static files
	frontendDir := filepath.Join(".", "frontend", "dist")
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		// fallback to src for dev
		frontendDir = filepath.Join(".", "frontend")
	}
	mux.Handle("/", corsMiddleware(http.FileServer(http.Dir(frontendDir))))

	addr := ":9272"
	log.Printf("[APGUARD] Server running at http://localhost%s", addr)
	fmt.Printf("\n\033[1;36m╔══════════════════════════════════════╗\033[0m\n")
	fmt.Printf("\033[1;36m  ║  APGUARD Security Scanner v1.0       ║\033[0m\n")
	fmt.Printf("\033[1;36m  ║  http://localhost%s                  ║\033[0m\n", addr)
	fmt.Printf("\033[1;36m  ╚══════════════════════════════════════╝\033[0m\n\n")

	go func() {
		time.Sleep(800 * time.Millisecond)
		openBrowser("http://localhost" + addr)
	}()

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// ──────────────────────────────────────────────────────────────
// Route Registration
// ──────────────────────────────────────────────────────────────
func registerRoutes(mux *http.ServeMux) {
	api := func(path string, h http.HandlerFunc) {
		mux.Handle("/api/"+path, corsMiddleware(h))
	}

	// Auth
	api("auth/login", handleLogin)
	api("auth/logout", handleLogout)
	api("auth/me", handleMe)
	api("auth/register", handleRegister)

	// Targets
	api("targets", handleTargets)
	api("targets/", handleTargetByID) // /api/targets/:id and /api/targets/:id/toggle

	// Scans
	api("scans", handleScans)
	api("scans/", handleScanByID) // /api/scans/:id, /api/scans/:id/vulns

	// Dashboard
	api("dashboard/stats", handleDashboardStats)

	// Reports
	api("reports/", handleReports) // /api/reports/:id, /api/reports/:id/csv

	// Audit
	api("audit", handleAudit)

	// Scheduler
	api("scheduler", handleScheduler)
	api("scheduler/", handleSchedulerByID)
}

// ──────────────────────────────────────────────────────────────
// Auth Handlers
// ──────────────────────────────────────────────────────────────
func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w); return
	}
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonErr(w, "invalid request body", http.StatusBadRequest); return
	}
	resp, err := authSvc.Login(req)
	if err != nil {
		jsonErr(w, err.Error(), http.StatusUnauthorized); return
	}
	sessions[resp.Token] = &resp.User
	auditSvc.Record(resp.User.ID, "LOGIN", "User logged in via HTTP")
	jsonOK(w, map[string]any{"token": resp.Token, "user": resp.User})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	token := bearerToken(r)
	if u, ok := sessions[token]; ok {
		auditSvc.Record(u.ID, "LOGOUT", "User logged out")
		delete(sessions, token)
	}
	jsonOK(w, map[string]any{"message": "logged out"})
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	if u == nil { jsonErr(w, "unauthorized", http.StatusUnauthorized); return }
	jsonOK(w, u)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { methodNotAllowed(w); return }
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonErr(w, "invalid request", http.StatusBadRequest); return
	}
	u, err := authSvc.Register(req)
	if err != nil { jsonErr(w, err.Error(), http.StatusBadRequest); return }
	jsonOK(w, u)
}

// ──────────────────────────────────────────────────────────────
// Target Handlers
// ──────────────────────────────────────────────────────────────
func handleTargets(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ts, err := targetSvc.List()
		if err != nil { jsonErr(w, err.Error(), 500); return }
		jsonOK(w, ts)
	case http.MethodPost:
		u := currentUser(r)
		var req targets.CreateTargetRequest
		json.NewDecoder(r.Body).Decode(&req)
		if u != nil { req.CreatedBy = u.ID }
		t, err := targetSvc.Create(req)
		if err != nil { jsonErr(w, err.Error(), http.StatusBadRequest); return }
		if u != nil { auditSvc.Record(u.ID, "CREATE_TARGET", "Created: "+req.URL) }
		jsonOK(w, t)
	default:
		methodNotAllowed(w)
	}
}

func handleTargetByID(w http.ResponseWriter, r *http.Request) {
	// Path: /api/targets/:id  or  /api/targets/:id/toggle  or  /api/targets/:id/scan
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/targets/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil { jsonErr(w, "invalid id", http.StatusBadRequest); return }
	action := ""
	if len(parts) > 1 { action = parts[1] }

	switch action {
	case "toggle":
		if err := targetSvc.Toggle(id); err != nil { jsonErr(w, err.Error(), 500); return }
		jsonOK(w, map[string]any{"toggled": true})
	case "scan":
		u := currentUser(r)
		t, err := targetSvc.GetByID(id)
		if err != nil { jsonErr(w, err.Error(), 404); return }
		body := struct{ Profile string `json:"profile"` }{}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Profile == "" { body.Profile = "full" }
		sc, err := scanSvc.Create(scans.CreateScanRequest{TargetID: id, ScanProfile: body.Profile})
		if err != nil { jsonErr(w, err.Error(), 500); return }
		if u != nil { auditSvc.Record(u.ID, "START_SCAN", fmt.Sprintf("Scan #%d on %s", sc.ID, t.URL)) }
		go scanEngine.RunScan(sc.ID, t.URL, body.Profile)
		jsonOK(w, sc)
	case "":
		if r.Method == http.MethodDelete {
			u := currentUser(r)
			if err := targetSvc.Delete(id); err != nil { jsonErr(w, err.Error(), 500); return }
			if u != nil { auditSvc.Record(u.ID, "DELETE_TARGET", fmt.Sprintf("Deleted target #%d", id)) }
			jsonOK(w, map[string]any{"deleted": true})
		} else {
			methodNotAllowed(w)
		}
	default:
		http.NotFound(w, r)
	}
}

// ──────────────────────────────────────────────────────────────
// Scan Handlers
// ──────────────────────────────────────────────────────────────
func handleScans(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ss, err := scanSvc.List()
		if err != nil { jsonErr(w, err.Error(), 500); return }
		jsonOK(w, ss)
	case http.MethodPost:
		var req scans.CreateScanRequest
		json.NewDecoder(r.Body).Decode(&req)
		t, err := targetSvc.GetByID(req.TargetID)
		if err != nil { jsonErr(w, "target not found", 404); return }
		sc, err := scanSvc.Create(req)
		if err != nil { jsonErr(w, err.Error(), 500); return }
		u := currentUser(r)
		if u != nil { auditSvc.Record(u.ID, "START_SCAN", fmt.Sprintf("Scan #%d on %s", sc.ID, t.URL)) }
		go scanEngine.RunScan(sc.ID, t.URL, req.ScanProfile)
		jsonOK(w, sc)
	default:
		methodNotAllowed(w)
	}
}

func handleScanByID(w http.ResponseWriter, r *http.Request) {
	// /api/scans/:id  or  /api/scans/:id/vulnerabilities
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/scans/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil { jsonErr(w, "invalid id", http.StatusBadRequest); return }
	action := ""
	if len(parts) > 1 { action = parts[1] }

	switch action {
	case "vulnerabilities", "vulns":
		vulns, err := scanSvc.GetVulnerabilities(id)
		if err != nil { jsonErr(w, err.Error(), 500); return }
		jsonOK(w, vulns)
	default:
		sc, err := scanSvc.GetByID(id)
		if err != nil { jsonErr(w, err.Error(), 404); return }
		jsonOK(w, sc)
	}
}

// ──────────────────────────────────────────────────────────────
// Dashboard
// ──────────────────────────────────────────────────────────────
func handleDashboardStats(w http.ResponseWriter, r *http.Request) {
	stats, err := scanSvc.GetStats()
	if err != nil { jsonErr(w, err.Error(), 500); return }
	jsonOK(w, stats)
}

// ──────────────────────────────────────────────────────────────
// Reports
// ──────────────────────────────────────────────────────────────
func handleReports(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/reports/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil { jsonErr(w, "invalid id", http.StatusBadRequest); return }
	sub := ""
	if len(parts) > 1 { sub = parts[1] }

	if sub == "csv" {
		csv, err := reportSvc.ExportCSV(id)
		if err != nil { jsonErr(w, err.Error(), 500); return }
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=apguard-scan-%d.csv", id))
		w.Write([]byte(csv))
		return
	}
	rpt, err := reportSvc.GenerateScanReport(id)
	if err != nil { jsonErr(w, err.Error(), 500); return }
	jsonOK(w, rpt)
}

// ──────────────────────────────────────────────────────────────
// Audit
// ──────────────────────────────────────────────────────────────
func handleAudit(w http.ResponseWriter, r *http.Request) {
	limit := 200
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	logs, err := auditSvc.List(limit)
	if err != nil { jsonErr(w, err.Error(), 500); return }
	jsonOK(w, logs)
}

// ──────────────────────────────────────────────────────────────
// Scheduler
// ──────────────────────────────────────────────────────────────
func handleScheduler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := schedulerSvc.ListTasks()
		if err != nil { jsonErr(w, err.Error(), 500); return }
		jsonOK(w, tasks)
	case http.MethodPost:
		var body struct {
			TargetID int    `json:"target_id"`
			CronExpr string `json:"cron_expr"`
			Profile  string `json:"scan_profile"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		task, err := schedulerSvc.CreateTask(body.TargetID, body.CronExpr, body.Profile)
		if err != nil { jsonErr(w, err.Error(), 500); return }
		jsonOK(w, task)
	default:
		methodNotAllowed(w)
	}
}

func handleSchedulerByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/scheduler/"), "/")
	id, err := strconv.Atoi(parts[0])
	if err != nil { jsonErr(w, "invalid id", http.StatusBadRequest); return }
	action := ""
	if len(parts) > 1 { action = parts[1] }

	switch action {
	case "toggle":
		if err := schedulerSvc.ToggleTask(id); err != nil { jsonErr(w, err.Error(), 500); return }
		jsonOK(w, map[string]any{"toggled": true})
	default:
		if r.Method == http.MethodDelete {
			if err := schedulerSvc.DeleteTask(id); err != nil { jsonErr(w, err.Error(), 500); return }
			jsonOK(w, map[string]any{"deleted": true})
		} else {
			methodNotAllowed(w)
		}
	}
}

// ──────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────
func jsonOK(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"success": true, "data": data})
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{"success": false, "error": msg})
}

func methodNotAllowed(w http.ResponseWriter) {
	jsonErr(w, "method not allowed", http.StatusMethodNotAllowed)
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	return r.URL.Query().Get("token")
}

func currentUser(r *http.Request) *auth.User {
	return sessions[bearerToken(r)]
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if r.Method == http.MethodOptions { w.WriteHeader(http.StatusOK); return }
		next.ServeHTTP(w, r)
	})
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "cmd"
		args = []string{"/c", "start", url}
	}
	exec.Command(cmd, args...).Start()
}
