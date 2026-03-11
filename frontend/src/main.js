// ============================================================
// APGUARD Frontend — HTTP + Wails Hybrid SPA
// ============================================================
import './style.css';

// ── Professional SVG Icon System (Lucide-style) ───────────────
const ICONS = {
  'shield':          `<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>`,
  'shield-alert':    `<path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>`,
  'bar-chart':       `<line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/>`,
  'crosshair':       `<circle cx="12" cy="12" r="10"/><line x1="22" y1="12" x2="18" y2="12"/><line x1="6" y1="12" x2="2" y2="12"/><line x1="12" y1="6" x2="12" y2="2"/><line x1="12" y1="22" x2="12" y2="18"/>`,
  'search':          `<circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>`,
  'file-text':       `<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><line x1="10" y1="9" x2="8" y2="9"/>`,
  'clock':           `<circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>`,
  'clipboard-list':  `<path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1" ry="1"/><line x1="9" y1="12" x2="15" y2="12"/><line x1="9" y1="16" x2="13" y2="16"/>`,
  'lock':            `<rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>`,
  'log-in':          `<path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4"/><polyline points="10 17 15 12 10 7"/><line x1="15" y1="12" x2="3" y2="12"/>`,
  'log-out':         `<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/>`,
  'plus':            `<line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>`,
  'trash':           `<polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a1 1 0 0 1 1-1h4a1 1 0 0 1 1 1v2"/><line x1="10" y1="11" x2="10" y2="17"/><line x1="14" y1="11" x2="14" y2="17"/>`,
  'play':            `<polygon points="5 3 19 12 5 21 5 3"/>`,
  'pause':           `<rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/>`,
  'download':        `<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>`,
  'file-down':       `<path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="18" x2="12" y2="12"/><polyline points="9 15 12 18 15 15"/>`,
  'rotate-cw':       `<polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>`,
  'x':               `<line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>`,
  'chevron-right':   `<polyline points="9 18 15 12 9 6"/>`,
  'zap':             `<polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>`,
  'bug':             `<path d="M8 2l1.5 1.5M16 2l-1.5 1.5"/><path d="M9 7h6l1 2H8l1-2z"/><path d="M12 9v13"/><path d="M3 12h18"/><path d="M3 8l2.5 3.5M21 8l-2.5 3.5"/><path d="M3 16l2.5-3.5M21 16l-2.5-3.5"/>`,
  'alert-triangle':  `<path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/>`,
  'alert-circle':    `<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>`,
  'alert-octagon':   `<polygon points="7.86 2 16.14 2 22 7.86 22 16.14 16.14 22 7.86 22 2 16.14 2 7.86 7.86 2"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>`,
  'check-circle':    `<path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/>`,
  'x-circle':        `<circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/>`,
  'loader':          `<line x1="12" y1="2" x2="12" y2="6"/><line x1="12" y1="18" x2="12" y2="22"/><line x1="4.93" y1="4.93" x2="7.76" y2="7.76"/><line x1="16.24" y1="16.24" x2="19.07" y2="19.07"/><line x1="2" y1="12" x2="6" y2="12"/><line x1="18" y1="12" x2="22" y2="12"/><line x1="4.93" y1="19.07" x2="7.76" y2="16.24"/><line x1="16.24" y1="7.76" x2="19.07" y2="4.93"/>`,
  'info':            `<circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="8"/><line x1="12" y1="12" x2="12" y2="16"/>`,
  'scan-line':       `<path d="M3 7V5a2 2 0 0 1 2-2h2"/><path d="M17 3h2a2 2 0 0 1 2 2v2"/><path d="M21 17v2a2 2 0 0 1-2 2h-2"/><path d="M7 21H5a2 2 0 0 1-2-2v-2"/><line x1="7" y1="12" x2="17" y2="12"/>`,
  'map-pin':         `<path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"/><circle cx="12" cy="10" r="3"/>`,
  'activity':        `<polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>`,
};

/**
 * Returns an inline SVG icon string.
 * @param {string} name  - Icon name key from ICONS
 * @param {number} [sz=16] - Width/height in px
 * @param {string} [cls=''] - Optional CSS class(es)
 */
function icon(name, sz = 16, cls = '') {
  const inner = ICONS[name];
  if (!inner) return '';
  const ca = cls ? ` class="${cls}"` : '';
  return `<svg xmlns="http://www.w3.org/2000/svg" width="${sz}" height="${sz}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"${ca} style="vertical-align:middle;flex-shrink:0;display:inline-block">${inner}</svg>`;
}


// ── Mode detection ────────────────────────────────────────────
// In Wails desktop: window.go.main.App is injected.
// In browser (HTTP server): we call /api/* endpoints via fetch.
const GO = window.go?.main?.App ?? null;
const API_BASE = window.location.origin; // same origin as the Go HTTP server

// ── Token storage ─────────────────────────────────────────────
const Token = {
  get: ()    => localStorage.getItem('apguard_token') ?? '',
  set: (t)   => localStorage.setItem('apguard_token', t),
  clear: ()  => localStorage.removeItem('apguard_token'),
};

// ── HTTP helpers ──────────────────────────────────────────────
async function apiFetch(path, options = {}) {
  const token = Token.get();
  const headers = { 'Content-Type': 'application/json', ...(options.headers ?? {}) };
  if (token) headers['Authorization'] = `Bearer ${token}`;
  try {
    const res = await fetch(API_BASE + '/api/' + path, { ...options, headers });
    const text = await res.text();
    let json;
    try { json = JSON.parse(text); } catch {
      return { success: false, error: `Server error (${res.status}) — is the APGUARD server running?` };
    }
    // Server wraps responses as { success, data } or { success, error }
    if (json.success && json.data !== undefined) {
      return { success: true, _data: json.data };
    }
    return json;
  } catch (err) {
    return { success: false, error: 'Cannot connect to APGUARD server. Run: go run ./cmd/server/' };
  }
}

function post(path, body) {
  return apiFetch(path, { method: 'POST', body: JSON.stringify(body) });
}
function get(path) { return apiFetch(path); }
function del(path) { return apiFetch(path, { method: 'DELETE' }); }
function patch(path, body = {}) {
  return apiFetch(path, { method: 'PATCH', body: JSON.stringify(body) });
}

// ── Wails bridge ──────────────────────────────────────────────
async function wailsCall(method, ...args) {
  try {
    return await GO[method](...args);
  } catch (err) {
    return { success: false, error: String(err) };
  }
}

// ── Unified call() — auto-routes Wails vs HTTP ────────────────
async function call(method, ...args) {
  if (GO) return wailsCall(method, ...args);
  return httpCall(method, args);
}

// ── HTTP method router ────────────────────────────────────────
// Translates Wails method names → HTTP calls and normalises responses.
async function httpCall(method, args) {
  let raw;
  switch (method) {
    // ── Auth
    case 'Login':
      raw = await post('auth/login', { username: args[0], password: args[1] });
      if (raw.success) {
        Token.set(raw.data?.token ?? raw._data?.token ?? raw.token);
        const u = raw.data?.user ?? raw._data?.user ?? raw.user;
        return { success: true, token: Token.get(), user: u };
      }
      return raw;

    case 'Logout':
      await post('auth/logout', {});
      Token.clear();
      return { success: true };

    case 'GetCurrentUser': {
      raw = await get('auth/me');
      if (raw.success) return { success: true, user: raw._data ?? raw };
      return { success: false };
    }

    case 'Register':
      raw = await post('auth/register', { username: args[0], email: args[1], password: args[2], role: args[3] });
      return raw.success ? { success: true, user: raw._data ?? raw.data } : raw;

    // ── Targets
    case 'ListTargets': {
      raw = await get('targets');
      return { success: raw.success, targets: raw._data ?? raw.data ?? [] };
    }
    case 'CreateTarget':
      raw = await post('targets', { url: args[0], name: args[1], description: args[2] });
      return raw.success ? { success: true, target: raw._data ?? raw.data } : raw;

    case 'ToggleTarget':
      raw = await patch(`targets/${args[0]}/toggle`);
      return raw;

    case 'DeleteTarget':
      raw = await del(`targets/${args[0]}`);
      return raw;

    // ── Scans
    case 'StartScan': {
      raw = await post(`targets/${args[0]}/scan`, { profile: args[1] });
      if (raw.success) {
        const sc = raw._data ?? raw.data ?? {};
        return { success: true, scan_id: sc.id, status: 'running' };
      }
      return raw;
    }
    case 'ListScans': {
      raw = await get('scans');
      return { success: raw.success, scans: raw._data ?? raw.data ?? [] };
    }
    case 'GetScan': {
      raw = await get(`scans/${args[0]}`);
      return { success: raw.success, scan: raw._data ?? raw.data };
    }
    case 'GetVulnerabilities': {
      raw = await get(`scans/${args[0]}/vulnerabilities`);
      return { success: raw.success, vulnerabilities: raw._data ?? raw.data ?? [] };
    }
    case 'GetDashboardStats': {
      raw = await get('dashboard/stats');
      return { success: raw.success, stats: raw._data ?? raw.data ?? {} };
    }

    // ── Reports
    case 'GenerateReport': {
      raw = await get(`reports/${args[0]}`);
      return { success: raw.success, report: raw._data ?? raw.data };
    }
    case 'ExportCSV':
      // For CSV we trigger a browser download directly
      window.open(`${API_BASE}/api/reports/${args[0]}/csv?token=${Token.get()}`, '_blank');
      return { success: true, csv: '' };

    // ── Audit
    case 'GetAuditLogs': {
      raw = await get(`audit?limit=${args[0] ?? 200}`);
      return { success: raw.success, logs: raw._data ?? raw.data ?? [] };
    }

    // ── Scheduler
    case 'ListScheduledScans': {
      raw = await get('scheduler');
      return { success: raw.success, tasks: raw._data ?? raw.data ?? [] };
    }
    case 'CreateScheduledScan':
      raw = await post('scheduler', { target_id: args[0], cron_expr: args[1], scan_profile: args[2] });
      return raw.success ? { success: true, task: raw._data ?? raw.data } : raw;

    case 'ToggleScheduledScan':
      raw = await patch(`scheduler/${args[0]}/toggle`);
      return raw;

    case 'DeleteScheduledScan':
      raw = await del(`scheduler/${args[0]}`);
      return raw;

    default:
      console.warn('[APGUARD] Unknown method:', method, args);
      return { success: false, error: `Unknown method: ${method}` };
  }
}

// ── State ─────────────────────────────────────────────────────
const state = {
  user: null,
  page: 'dashboard',
  alert: null,
  alertTimer: null,
  activeScans: new Set(),
  pollTimer: null,
};

// ── Toast Notification ────────────────────────────────────────
function showAlert(msg, type = 'success') {
  state.alert = { msg, type };
  render();
  if (state.alertTimer) clearTimeout(state.alertTimer);
  state.alertTimer = setTimeout(() => {
    state.alert = null;
    render();
  }, 4000);
}

// ── Navigation ────────────────────────────────────────────────
function navigate(page) {
  state.page = page;
  render();
  if (page === 'scans' && state.activeScans.size > 0) startScanPoller();
}

// ── Render ────────────────────────────────────────────────────
function render() {
  const app = document.getElementById('app');
  if (!app) return;

  if (!state.user) {
    app.innerHTML = renderLoginPage();
    bindLoginPage();
    return;
  }

  app.innerHTML = renderAppShell();
  bindAppShell();
  renderPage();
}

// ── Login Page ────────────────────────────────────────────────
function renderLoginPage() {
  return `
    <div class="login-page">
      <div class="login-bg-grid"></div>
      <div class="login-bg-glow"></div>
      <div class="login-card">
        <div class="login-logo">
          <div class="login-logo-icon">${icon('shield', 40)}</div>
          <h1>APGUARD</h1>
          <p>Web Application Security Scanner</p>
        </div>
        ${state.alert ? `<div class="alert alert-${state.alert.type}">${icon('alert-triangle', 14)} ${escHtml(state.alert.msg)}</div>` : ''}
        <form id="login-form">
          <div class="form-group">
            <label class="form-label">Username</label>
            <input id="login-user" class="form-control" type="text" placeholder="admin" autocomplete="username" required/>
          </div>
          <div class="form-group">
            <label class="form-label">Password</label>
            <input id="login-pass" class="form-control" type="password" placeholder="••••••••" autocomplete="current-password" required/>
          </div>
          <button id="login-btn" type="submit" class="btn btn-primary" style="width:100%;justify-content:center;padding:12px;">
            ${icon('log-in', 16)} Sign In to APGUARD
          </button>
        </form>
        <p style="text-align:center;margin-top:20px;font-size:12px;color:var(--text-muted);">
          Default: <span style="color:var(--accent-cyan);font-family:var(--font-mono)">admin / admin123</span>
        </p>
      </div>
    </div>
  `;
}

function bindLoginPage() {
  document.getElementById('login-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = document.getElementById('login-btn');
    const user = document.getElementById('login-user').value.trim();
    const pass = document.getElementById('login-pass').value;

    btn.disabled = true;
    btn.innerHTML = '<span class="spinner"></span> Signing in…';

    const res = await call('Login', user, pass);
    if (res.success) {
      state.user = res.user;
      state.alert = null;
      navigate('dashboard');
    } else {
      state.alert = { msg: res.error, type: 'error' };
      btn.disabled = false;
      btn.innerHTML = `${icon('log-in', 16)} Sign In to APGUARD`;
      render();
    }
  });
}

// ── App Shell ─────────────────────────────────────────────────
const NAV_ITEMS = [
  { id: 'dashboard', icon: 'bar-chart', label: 'Dashboard', section: null },
  { id: 'targets',   icon: 'crosshair', label: 'Targets', section: 'SCANNING' },
  { id: 'scans',     icon: 'scan-line', label: 'Scans', section: null },
  { id: 'reports',   icon: 'file-text', label: 'Reports', section: null },
  { id: 'scheduler', icon: 'clock', label: 'Scheduler', section: 'AUTOMATION' },
  { id: 'audit',     icon: 'clipboard-list', label: 'Audit Logs', section: 'SYSTEM' },
];

const PAGE_TITLES = {
  dashboard: ['Dashboard', 'Security Overview'],
  targets: ['Targets', 'Manage web application targets'],
  scans: ['Scans', 'Run and monitor vulnerability scans'],
  reports: ['Reports', 'Vulnerability reports and exports'],
  scheduler: ['Scheduler', 'Automated scan scheduling'],
  audit: ['Audit Logs', 'System activity trail'],
};

function renderAppShell() {
  const [title, sub] = PAGE_TITLES[state.page] ?? ['APGUARD', ''];
  const activeCount = state.activeScans.size;

  let sections = [];
  let lastSection = undefined;
  let navLinks = '';
  for (const item of NAV_ITEMS) {
    if (item.section !== lastSection) {
      if (item.section) navLinks += `<div class="nav-section">${item.section}</div>`;
      lastSection = item.section;
    }
    const badge = item.id === 'scans' && activeCount > 0 ? `<span class="nav-badge">${activeCount}</span>` : '';
    navLinks += `
      <div class="nav-item ${state.page === item.id ? 'active' : ''}" data-page="${item.id}" id="nav-${item.id}">
        <span class="nav-icon">${icon(item.icon, 16)}</span>
        <span>${item.label}</span>
        ${badge}
      </div>`;
  }

  const initials = state.user?.username?.slice(0,2)?.toUpperCase() ?? 'AP';

  return `
    ${state.alert ? `
      <div style="position:fixed;top:20px;right:24px;z-index:2000;min-width:260px;max-width:400px;animation:slideUp 0.2s ease;">
        <div class="alert alert-${state.alert.type}">${state.alert.type === 'error' ? icon('alert-triangle',14) : icon('check-circle',14)} ${escHtml(state.alert.msg)}</div>
      </div>` : ''}

    <div class="app-layout">
      <aside class="sidebar">
        <div class="sidebar-logo">
          <div class="logo-icon">${icon('shield', 22)}</div>
          <span class="logo-text">APGUARD</span>
          <span class="logo-version">v1.0</span>
        </div>
        <nav class="sidebar-nav">${navLinks}</nav>
        <div class="sidebar-user">
          <div class="user-avatar">${initials}</div>
          <div class="user-info">
            <div class="user-name">${escHtml(state.user?.username ?? '')}</div>
            <div class="user-role">${escHtml(state.user?.role ?? '')}</div>
          </div>
          <button class="btn-logout" id="btn-logout">Exit</button>
        </div>
      </aside>

      <div class="main-content">
        <header class="topbar">
          <span class="page-title">${title}</span>
          <span class="page-subtitle">— ${sub}</span>
          <div class="topbar-actions">
            ${activeCount > 0 ? `<span class="badge badge-running">${icon('zap',12)} ${activeCount} Active Scan${activeCount > 1 ? 's' : ''}</span>` : ''}
            <div class="status-dot" title="System Online"></div>
          </div>
        </header>
        <main class="page-content" id="page-content"></main>
      </div>
    </div>
  `;
}

function bindAppShell() {
  document.querySelectorAll('.nav-item[data-page]').forEach(el => {
    el.addEventListener('click', () => navigate(el.dataset.page));
  });
  document.getElementById('btn-logout')?.addEventListener('click', async () => {
    await call('Logout');
    state.user = null;
    state.page = 'dashboard';
    render();
  });
}

// ── Page Router ───────────────────────────────────────────────
async function renderPage() {
  const container = document.getElementById('page-content');
  if (!container) return;
  container.innerHTML = `<div class="loading-center"><span class="spinner"></span> Loading…</div>`;
  switch (state.page) {
    case 'dashboard': await renderDashboard(container); break;
    case 'targets':   await renderTargets(container);   break;
    case 'scans':     await renderScans(container);     break;
    case 'reports':   await renderReports(container);   break;
    case 'scheduler': await renderScheduler(container); break;
    case 'audit':     await renderAudit(container);     break;
  }
}

// ─────────────────────────────────────────────────────────────
// PAGE: Dashboard
// ─────────────────────────────────────────────────────────────
async function renderDashboard(container) {
  const res = await call('GetDashboardStats');
  if (!res.success) { container.innerHTML = errHTML(res.error); return; }

  const s = res.stats;
  const sev = s.severity_breakdown ?? {};
  const maxSev = Math.max(sev.critical ?? 0, sev.high ?? 0, sev.medium ?? 0, sev.low ?? 0, 1);

  const scansRes = await call('ListScans');
  const recentScans = (scansRes.scans ?? []).slice(0, 5);

  container.innerHTML = `
    <div class="metric-grid">
      <div class="metric-card blue">
        <div class="metric-icon">${icon('search', 26)}</div>
        <div class="metric-value">${s.total_scans ?? 0}</div>
        <div class="metric-label">Total Scans</div>
      </div>
      <div class="metric-card green">
        <div class="metric-icon">${icon('crosshair', 26)}</div>
        <div class="metric-value">${s.total_targets ?? 0}</div>
        <div class="metric-label">Targets</div>
      </div>
      <div class="metric-card red">
        <div class="metric-icon">${icon('bug', 26)}</div>
        <div class="metric-value">${s.total_vulnerabilities ?? 0}</div>
        <div class="metric-label">Vulnerabilities</div>
      </div>
      <div class="metric-card orange">
        <div class="metric-icon">${icon('zap', 26)}</div>
        <div class="metric-value">${s.active_scans ?? 0}</div>
        <div class="metric-label">Active Scans</div>
      </div>
    </div>

    <div class="grid-2">
      <div class="card">
        <div class="card-header">
          <div>
            <div class="card-title">${icon('activity',15)} Severity Breakdown</div>
            <div class="card-subtitle">Vulnerabilities by severity level</div>
          </div>
        </div>
        <div class="sev-bars">
          <div class="sev-bar-row">
            <span class="sev-bar-label" style="color:var(--critical)">Critical</span>
            <div class="sev-bar-track"><div class="sev-bar-fill critical" style="width:${pct(sev.critical,maxSev)}%"></div></div>
            <span class="sev-bar-count">${sev.critical ?? 0}</span>
          </div>
          <div class="sev-bar-row">
            <span class="sev-bar-label" style="color:var(--high)">High</span>
            <div class="sev-bar-track"><div class="sev-bar-fill high" style="width:${pct(sev.high,maxSev)}%"></div></div>
            <span class="sev-bar-count">${sev.high ?? 0}</span>
          </div>
          <div class="sev-bar-row">
            <span class="sev-bar-label" style="color:var(--medium)">Medium</span>
            <div class="sev-bar-track"><div class="sev-bar-fill medium" style="width:${pct(sev.medium,maxSev)}%"></div></div>
            <span class="sev-bar-count">${sev.medium ?? 0}</span>
          </div>
          <div class="sev-bar-row">
            <span class="sev-bar-label" style="color:var(--low)">Low</span>
            <div class="sev-bar-track"><div class="sev-bar-fill low" style="width:${pct(sev.low,maxSev)}%"></div></div>
            <span class="sev-bar-count">${sev.low ?? 0}</span>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <div>
            <div class="card-title">${icon('clock',15)} Recent Scans</div>
            <div class="card-subtitle">Latest scanning activity</div>
          </div>
          <button class="btn btn-ghost btn-sm" onclick="navigate('scans')">View All</button>
        </div>
        ${recentScans.length === 0 ? `<div class="empty-state" style="padding:30px"><div class="empty-icon">${icon('search',36)}</div><p>No scans yet</p></div>` :
          recentScans.map(sc => `
            <div style="display:flex;align-items:center;gap:12px;padding:10px 0;border-bottom:1px solid var(--border)">
              <div style="flex:1;min-width:0;">
                <div style="font-size:13px;font-weight:600;color:var(--text-primary);white-space:nowrap;overflow:hidden;text-overflow:ellipsis">${escHtml(sc.target_url || 'Unknown Target')}</div>
                <div style="font-size:11px;color:var(--text-muted);margin-top:2px">${escHtml(sc.created_at?.slice(0,16) ?? '')}</div>
              </div>
              <div>
                ${scanStatusBadge(sc.status)}
                ${sc.total_vulnerabilities > 0 ? `<span class="badge badge-critical" style="margin-left:4px">${sc.total_vulnerabilities} vulns</span>` : ''}
              </div>
            </div>
          `).join('')}
      </div>
    </div>

    <div class="card mt-16">
      <div class="card-header">
        <div class="card-title">${icon('zap',15)} Quick Scan</div>
        <div class="card-subtitle">Start a scan directly from the dashboard</div>
      </div>
      <div style="display:flex;gap:12px;flex-wrap:wrap;align-items:flex-end;">
        <div style="flex:1;min-width:200px">
          <select id="dash-target-select" class="form-control">
            <option value="">Select target…</option>
          </select>
        </div>
        <div>
          <select id="dash-profile-select" class="form-control">
            <option value="full">Full Scan</option>
            <option value="quick">Quick Scan</option>
          </select>
        </div>
        <button id="dash-scan-btn" class="btn btn-success">${icon('play',14)} Launch Scan</button>
      </div>
    </div>
  `;

  // Populate target select
  const targetsRes = await call('ListTargets');
  const sel = document.getElementById('dash-target-select');
  if (sel && targetsRes.success) {
    (targetsRes.targets ?? []).filter(t => t.enabled).forEach(t => {
      const opt = document.createElement('option');
      opt.value = t.id;
      opt.textContent = t.name + ' — ' + t.url;
      sel.appendChild(opt);
    });
  }

  document.getElementById('dash-scan-btn')?.addEventListener('click', async () => {
    const targetID = parseInt(document.getElementById('dash-target-select')?.value);
    const profile = document.getElementById('dash-profile-select')?.value || 'full';
    if (!targetID) { showAlert('Select a target first', 'error'); return; }
    const res = await call('StartScan', targetID, profile);
    if (res.success) {
      state.activeScans.add(res.scan_id);
      showAlert(`Scan #${res.scan_id} started successfully!`, 'success');
      navigate('scans');
    } else {
      showAlert(res.error, 'error');
    }
  });
}

// ─────────────────────────────────────────────────────────────
// PAGE: Targets
// ─────────────────────────────────────────────────────────────
async function renderTargets(container) {
  const res = await call('ListTargets');
  const targets = res.targets ?? [];

  container.innerHTML = `
    <div class="flex-between mb-24">
      <div>
        <div style="font-size:13px;color:var(--text-muted)">${targets.length} target${targets.length !== 1 ? 's' : ''} registered</div>
      </div>
      <button id="btn-add-target" class="btn btn-primary">${icon('plus',14)} Add Target</button>
    </div>

    ${targets.length === 0 ? `
      <div class="card">
        <div class="empty-state">
          <div class="empty-icon">${icon('crosshair',40)}</div>
          <h3>No targets yet</h3>
          <p>Add your first web application target to begin scanning</p>
          <button class="btn btn-primary mt-16" id="btn-add-target-empty">${icon('plus',14)} Add First Target</button>
        </div>
      </div>
    ` : `
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>Name</th>
              <th>URL</th>
              <th>Description</th>
              <th>Status</th>
              <th>Created</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${targets.map(t => `
              <tr>
                <td style="color:var(--text-muted);font-family:var(--font-mono)">${t.id}</td>
                <td><strong>${escHtml(t.name)}</strong></td>
                <td class="td-mono"><a href="#" style="color:var(--accent-blue);text-decoration:none">${escHtml(t.url)}</a></td>
                <td style="color:var(--text-secondary);font-size:12px">${escHtml(t.description || '—')}</td>
                <td>
                  ${t.enabled
                    ? '<span class="badge badge-success">Active</span>'
                    : '<span class="badge badge-muted">Disabled</span>'}
                </td>
                <td style="color:var(--text-muted);font-size:12px">${escHtml((t.created_at || '').slice(0,10))}</td>
                <td class="td-actions">
                  <button class="btn btn-ghost btn-sm" data-scan="${t.id}" data-url="${escHtml(t.url)}">${icon('search',13)} Scan</button>
                  <button class="btn btn-ghost btn-sm btn-toggle-target" data-id="${t.id}">${t.enabled ? `${icon('pause',13)} Pause` : `${icon('play',13)} Enable`}</button>
                  <button class="btn btn-danger btn-sm btn-del-target" data-id="${t.id}">${icon('trash',13)}</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `}

    <!-- Add Target Modal -->
    <div id="target-modal" class="modal-overlay" style="display:none">
      <div class="modal">
        <div class="modal-header">
          <span class="modal-title">${icon('plus',16)} Add New Target</span>
          <button class="modal-close" id="close-target-modal">${icon('x',14)}</button>
        </div>
        <form id="target-form">
          <div class="form-group">
            <label class="form-label">Target URL *</label>
            <input id="t-url" class="form-control" type="url" placeholder="https://example.com" required/>
          </div>
          <div class="form-group">
            <label class="form-label">Target Name *</label>
            <input id="t-name" class="form-control" type="text" placeholder="My Web App" required/>
          </div>
          <div class="form-group">
            <label class="form-label">Description</label>
            <input id="t-desc" class="form-control" type="text" placeholder="Optional description"/>
          </div>
          <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:8px">
            <button type="button" class="btn btn-ghost" id="cancel-target">Cancel</button>
            <button type="submit" class="btn btn-primary" id="submit-target">Add Target</button>
          </div>
        </form>
      </div>
    </div>
  `;

  // Modal open
  const showModal = () => { document.getElementById('target-modal').style.display = 'flex'; };
  const hideModal = () => { document.getElementById('target-modal').style.display = 'none'; };
  document.getElementById('btn-add-target')?.addEventListener('click', showModal);
  document.getElementById('btn-add-target-empty')?.addEventListener('click', showModal);
  document.getElementById('close-target-modal')?.addEventListener('click', hideModal);
  document.getElementById('cancel-target')?.addEventListener('click', hideModal);

  // Submit target form
  document.getElementById('target-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = document.getElementById('submit-target');
    btn.disabled = true; btn.textContent = 'Adding…';
    const url = document.getElementById('t-url').value.trim();
    const name = document.getElementById('t-name').value.trim();
    const desc = document.getElementById('t-desc').value.trim();
    const res = await call('CreateTarget', url, name, desc);
    if (res.success) {
      showAlert(`Target "${name}" added successfully`, 'success');
      hideModal();
      navigate('targets');
    } else {
      showAlert(res.error, 'error');
      btn.disabled = false; btn.textContent = 'Add Target';
    }
  });

  // Scan buttons
  container.querySelectorAll('[data-scan]').forEach(btn => {
    btn.addEventListener('click', async () => {
      const targetID = parseInt(btn.dataset.scan);
      const res = await call('StartScan', targetID, 'full');
      if (res.success) {
        state.activeScans.add(res.scan_id);
        showAlert(`Scan started on ${btn.dataset.url}`, 'success');
        navigate('scans');
      } else {
        showAlert(res.error, 'error');
      }
    });
  });

  // Toggle buttons
  container.querySelectorAll('.btn-toggle-target').forEach(btn => {
    btn.addEventListener('click', async () => {
      const id = parseInt(btn.dataset.id);
      await call('ToggleTarget', id);
      navigate('targets');
    });
  });

  // Delete buttons
  container.querySelectorAll('.btn-del-target').forEach(btn => {
    btn.addEventListener('click', async () => {
      if (!confirm('Delete this target? This cannot be undone.')) return;
      const id = parseInt(btn.dataset.id);
      const res = await call('DeleteTarget', id);
      if (res.success) { showAlert('Target deleted', 'success'); navigate('targets'); }
      else showAlert(res.error, 'error');
    });
  });
}

// ─────────────────────────────────────────────────────────────
// PAGE: Scans
// ─────────────────────────────────────────────────────────────
async function renderScans(container) {
  const [scansRes, targetsRes] = await Promise.all([call('ListScans'), call('ListTargets')]);
  const scans = scansRes.scans ?? [];
  const targets = targetsRes.targets ?? [];

  const activeCount = scans.filter(s => s.status === 'running').length;
  if (activeCount > 0) startScanPoller();

  container.innerHTML = `
    <div class="flex-between mb-24">
      <div style="font-size:13px;color:var(--text-muted)">${scans.length} scan${scans.length !== 1 ? 's' : ''} total</div>
      <button id="btn-new-scan" class="btn btn-primary">${icon('play',14)} New Scan</button>
    </div>

    ${scans.length === 0 ? `
      <div class="card">
        <div class="empty-state">
          <div class="empty-icon">${icon('scan-line',40)}</div>
          <h3>No scans yet</h3>
          <p>Start your first vulnerability scan against a registered target</p>
        </div>
      </div>
    ` : `
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>#</th>
              <th>Target URL</th>
              <th>Profile</th>
              <th>Status</th>
              <th>Vulnerabilities</th>
              <th>Started</th>
              <th>Duration</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            ${scans.map(sc => `
              <tr>
                <td style="font-family:var(--font-mono);color:var(--text-muted)">${sc.id}</td>
                <td class="td-mono" style="max-width:200px">${escHtml(sc.target_url || 'N/A')}</td>
                <td><span class="tag">${escHtml(sc.scan_profile)}</span></td>
                <td>${scanStatusBadge(sc.status)}</td>
                <td>
                  ${sc.total_vulnerabilities > 0
                    ? `<span class="badge badge-critical">${sc.total_vulnerabilities} found</span>`
                    : `<span style="color:var(--text-muted);font-size:12px">—</span>`}
                </td>
                <td style="color:var(--text-muted);font-size:12px">${formatDate(sc.start_time || sc.created_at)}</td>
                <td style="color:var(--text-muted);font-size:12px">${calcDuration(sc.start_time, sc.end_time)}</td>
                <td class="td-actions">
                  <button class="btn btn-ghost btn-sm btn-view-vulns" data-scan="${sc.id}"
                    ${sc.total_vulnerabilities === 0 ? 'disabled' : ''}>${icon('bug',13)} Vulnerabilities</button>
                  <button class="btn btn-ghost btn-sm btn-view-report" data-scan="${sc.id}">${icon('file-text',13)} Report</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `}

    <!-- New Scan Modal -->
    <div id="scan-modal" class="modal-overlay" style="display:none">
      <div class="modal">
        <div class="modal-header">
          <span class="modal-title">${icon('scan-line',16)} New Vulnerability Scan</span>
          <button class="modal-close" id="close-scan-modal">${icon('x',14)}</button>
        </div>
        <form id="scan-form">
          <div class="form-group">
            <label class="form-label">Target *</label>
            <select id="s-target" class="form-control" required>
              <option value="">Select target…</option>
              ${targets.filter(t => t.enabled).map(t => `<option value="${t.id}">${escHtml(t.name)} — ${escHtml(t.url)}</option>`).join('')}
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Scan Profile</label>
            <select id="s-profile" class="form-control">
              <option value="full">Full Scan (SQLi + XSS + LFI + Misconfiguration)</option>
              <option value="quick">Quick Scan (SQLi + XSS only)</option>
            </select>
          </div>
          <div class="alert alert-info" style="margin-bottom:0">
            Scans run concurrently in background. Results appear in real-time.
          </div>
          <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:16px">
            <button type="button" class="btn btn-ghost" id="cancel-scan">Cancel</button>
            <button type="submit" class="btn btn-success" id="submit-scan">${icon('play',14)} Launch Scan</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Vulnerabilities Modal -->
    <div id="vuln-modal" class="modal-overlay" style="display:none">
      <div class="modal" style="max-width:720px">
        <div class="modal-header">
          <span class="modal-title">${icon('bug',16)} Vulnerabilities</span>
          <button class="modal-close" id="close-vuln-modal">${icon('x',14)}</button>
        </div>
        <div id="vuln-modal-body"><div class="loading-center"><span class="spinner"></span></div></div>
      </div>
    </div>
  `;

  const showScanModal = () => document.getElementById('scan-modal').style.display = 'flex';
  const hideScanModal = () => document.getElementById('scan-modal').style.display = 'none';
  document.getElementById('btn-new-scan')?.addEventListener('click', showScanModal);
  document.getElementById('close-scan-modal')?.addEventListener('click', hideScanModal);
  document.getElementById('cancel-scan')?.addEventListener('click', hideScanModal);

  document.getElementById('scan-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const btn = document.getElementById('submit-scan');
    btn.disabled = true; btn.innerHTML = `${icon('loader',14)} Launching…`;
    const targetID = parseInt(document.getElementById('s-target').value);
    const profile = document.getElementById('s-profile').value;
    if (!targetID) { showAlert('Select a target', 'error'); btn.disabled = false; btn.innerHTML = `${icon('play',14)} Launch Scan`; return; }
    const res = await call('StartScan', targetID, profile);
    if (res.success) {
      state.activeScans.add(res.scan_id);
      showAlert(`Scan #${res.scan_id} launched!`, 'success');
      hideScanModal();
      navigate('scans');
    } else {
      showAlert(res.error, 'error');
      btn.disabled = false; btn.innerHTML = `${icon('play',14)} Launch Scan`;
    }
  });

  // Vuln modal
  const vulnModal = document.getElementById('vuln-modal');
  document.getElementById('close-vuln-modal')?.addEventListener('click', () => vulnModal.style.display = 'none');

  container.querySelectorAll('.btn-view-vulns').forEach(btn => {
    btn.addEventListener('click', async () => {
      const scanID = parseInt(btn.dataset.scan);
      vulnModal.style.display = 'flex';
      document.getElementById('vuln-modal-body').innerHTML = '<div class="loading-center"><span class="spinner"></span> Loading…</div>';
      const res = await call('GetVulnerabilities', scanID);
      const vulns = res.vulnerabilities ?? [];
      document.getElementById('vuln-modal-body').innerHTML = renderVulnTable(vulns);
      bindVulnExpand();
    });
  });

  // Report view
  container.querySelectorAll('.btn-view-report').forEach(btn => {
    btn.addEventListener('click', () => {
      const scanID = parseInt(btn.dataset.scan);
      state.reportScanID = scanID;
      navigate('reports');
    });
  });
}

function renderVulnTable(vulns) {
  if (!vulns.length) return `<div class="empty-state"><div class="empty-icon">${icon('check-circle',40)}</div><h3>No vulnerabilities found</h3><p>The scan completed clean</p></div>`;
  return `
    <div style="overflow-x:auto">
      <table>
        <thead><tr><th>Type</th><th>Severity</th><th>Endpoint</th><th>Parameter</th><th></th></tr></thead>
        <tbody>
          ${vulns.map((v, i) => `
            <tr class="vuln-row" data-idx="${i}" style="cursor:pointer">
              <td><strong>${escHtml(v.type)}</strong></td>
              <td>${severityBadge(v.severity)}</td>
              <td class="td-mono">${escHtml(v.endpoint)}</td>
              <td><span class="tag">${escHtml(v.parameter || '—')}</span></td>
              <td style="color:var(--text-muted)">${icon('chevron-right',14)}</td>
            </tr>
            <tr id="vuln-detail-${i}" style="display:none">
              <td colspan="5" style="padding:0 16px 12px">
                <div class="vuln-detail">Payload:     ${escHtml(v.payload || 'N/A')}
Description: ${escHtml(v.description || 'N/A')}
Status:      ${escHtml(v.status)}
Detected:    ${escHtml(v.created_at || '')}</div>
              </td>
            </tr>
          `).join('')}
        </tbody>
      </table>
    </div>
  `;
}

function bindVulnExpand() {
  document.querySelectorAll('.vuln-row').forEach(row => {
    row.addEventListener('click', () => {
      const det = document.getElementById(`vuln-detail-${row.dataset.idx}`);
      if (det) det.style.display = det.style.display === 'none' ? 'table-row' : 'none';
    });
  });
}

// ─────────────────────────────────────────────────────────────
// PAGE: Reports
// ─────────────────────────────────────────────────────────────
async function renderReports(container) {
  const scansRes = await call('ListScans');
  const scans = (scansRes.scans ?? []).filter(s => s.status === 'completed');

  // Auto-select scan if coming from scan page
  let selectedScan = state.reportScanID ?? (scans[0]?.id);
  state.reportScanID = null;

  container.innerHTML = `
    <div class="card mb-24">
      <div class="card-header">
        <div class="card-title">${icon('file-text',15)} Scan Report Viewer</div>
      </div>
      <div style="display:flex;gap:12px;align-items:flex-end;flex-wrap:wrap">
        <div style="flex:1;min-width:200px">
          <label class="form-label">Select Completed Scan</label>
          <select id="report-scan-select" class="form-control">
            <option value="">Choose a scan…</option>
            ${scans.map(s => `<option value="${s.id}" ${s.id === selectedScan ? 'selected' : ''}>Scan #${s.id} — ${escHtml(s.target_url)} (${escHtml(s.created_at?.slice(0,10) ?? '')})</option>`).join('')}
          </select>
        </div>
        <button id="btn-load-report" class="btn btn-primary">${icon('bar-chart',14)} Generate Report</button>
        <button id="btn-export-csv" class="btn btn-ghost">${icon('download',14)} Export CSV</button>
        <button id="btn-export-pdf" class="btn btn-ghost" style="color:var(--critical)">${icon('file-down',14)} Export PDF</button>
      </div>
    </div>
    <div id="report-body"></div>
  `;

  const loadReport = async () => {
    const scanID = parseInt(document.getElementById('report-scan-select')?.value);
    if (!scanID) { showAlert('Select a scan first', 'error'); return; }
    const rb = document.getElementById('report-body');
    rb.innerHTML = '<div class="loading-center"><span class="spinner"></span> Generating report…</div>';
    const res = await call('GenerateReport', scanID);
    if (!res.success) { rb.innerHTML = errHTML(res.error); return; }
    rb.innerHTML = renderReportView(res.report);
    bindVulnExpand();
  };

  document.getElementById('btn-load-report')?.addEventListener('click', loadReport);

  document.getElementById('btn-export-csv')?.addEventListener('click', async () => {
    const scanID = parseInt(document.getElementById('report-scan-select')?.value);
    if (!scanID) { showAlert('Select a scan first', 'error'); return; }
    const res = await call('ExportCSV', scanID);
    if (!res.success) { showAlert(res.error, 'error'); return; }
    downloadText(`apguard-scan-${scanID}.csv`, res.csv);
    showAlert('CSV exported successfully', 'success');
  });

  document.getElementById('btn-export-pdf')?.addEventListener('click', async () => {
    const scanID = parseInt(document.getElementById('report-scan-select')?.value);
    if (!scanID) { showAlert('Select a scan first', 'error'); return; }
    const btn = document.getElementById('btn-export-pdf');
    btn.disabled = true;
    btn.innerHTML = `${icon('loader',14)} Building PDF…`;
    try {
      const res = await call('GenerateReport', scanID);
      if (!res.success) { showAlert(res.error ?? 'Failed to generate report', 'error'); return; }
      exportPDF(res.report, scanID);
      showAlert('PDF print dialog opened — choose "Save as PDF"', 'success');
    } finally {
      btn.disabled = false;
      btn.innerHTML = `${icon('file-down',14)} Export PDF`;
    }
  });

  // Auto-load if scan preselected
  if (selectedScan) await loadReport();
}

function renderReportView(report) {
  const sev = report.severity_counts ?? {};
  const maxSev = Math.max(...Object.values(sev), 1);
  const riskColor = { CRITICAL: 'var(--critical)', HIGH: 'var(--high)', MEDIUM: 'var(--medium)', LOW: 'var(--low)', INFO: 'var(--accent-blue)' };

  return `
    <div class="metric-grid mb-24">
      <div class="metric-card ${report.risk_level === 'CRITICAL' ? 'red' : report.risk_level === 'HIGH' ? 'orange' : 'green'}">
        <div class="metric-icon">${icon('crosshair', 26)}</div>
        <div class="metric-value" style="color:${riskColor[report.risk_level] ?? 'var(--accent-green)'};">${report.risk_level}</div>
        <div class="metric-label">Risk Level</div>
      </div>
      <div class="metric-card blue">
        <div class="metric-icon">${icon('bug', 26)}</div>
        <div class="metric-value">${report.total_vulns ?? 0}</div>
        <div class="metric-label">Vulnerabilities</div>
      </div>
      <div class="metric-card red">
        <div class="metric-icon">${icon('alert-octagon',26)}</div>
        <div class="metric-value">${sev.CRITICAL ?? 0}</div>
        <div class="metric-label">Critical</div>
      </div>
      <div class="metric-card orange">
        <div class="metric-icon">${icon('alert-triangle',26)}</div>
        <div class="metric-value">${sev.HIGH ?? 0}</div>
        <div class="metric-label">High</div>
      </div>
    </div>

    <div class="grid-2 mb-24">
      <div class="card">
        <div class="card-header"><div class="card-title">${icon('map-pin',15)} Scan Details</div></div>
        <table style="width:100%;font-size:13px">
          <tr><td style="color:var(--text-secondary);padding:6px 0;width:40%">Scan ID</td><td style="font-family:var(--font-mono)">#${report.scan_id}</td></tr>
          <tr><td style="color:var(--text-secondary);padding:6px 0">Target</td><td class="td-mono">${escHtml(report.target_url)}</td></tr>
          <tr><td style="color:var(--text-secondary);padding:6px 0">Status</td><td>${scanStatusBadge(report.status)}</td></tr>
          <tr><td style="color:var(--text-secondary);padding:6px 0">Profile</td><td><span class="tag">${escHtml(report.scan_profile ?? 'full')}</span></td></tr>
          <tr><td style="color:var(--text-secondary);padding:6px 0">Started</td><td style="font-size:12px">${formatDate(report.start_time)}</td></tr>
          <tr><td style="color:var(--text-secondary);padding:6px 0">Ended</td><td style="font-size:12px">${formatDate(report.end_time)}</td></tr>
          <tr><td style="color:var(--text-secondary);padding:6px 0">Generated</td><td style="font-size:12px">${formatDate(report.generated_at)}</td></tr>
        </table>
      </div>

      <div class="card">
        <div class="card-header"><div class="card-title">${icon('activity',15)} Severity Distribution</div></div>
        <div class="sev-bars">
          <div class="sev-bar-row">
            <span class="sev-bar-label" style="color:var(--critical)">Critical</span>
            <div class="sev-bar-track"><div class="sev-bar-fill critical" style="width:${pct(sev.CRITICAL,maxSev)}%"></div></div>
            <span class="sev-bar-count">${sev.CRITICAL ?? 0}</span>
          </div>
          <div class="sev-bar-row">
            <span class="sev-bar-label" style="color:var(--high)">High</span>
            <div class="sev-bar-track"><div class="sev-bar-fill high" style="width:${pct(sev.HIGH,maxSev)}%"></div></div>
            <span class="sev-bar-count">${sev.HIGH ?? 0}</span>
          </div>
          <div class="sev-bar-row">
            <span class="sev-bar-label" style="color:var(--medium)">Medium</span>
            <div class="sev-bar-track"><div class="sev-bar-fill medium" style="width:${pct(sev.MEDIUM,maxSev)}%"></div></div>
            <span class="sev-bar-count">${sev.MEDIUM ?? 0}</span>
          </div>
          <div class="sev-bar-row">
            <span class="sev-bar-label" style="color:var(--low)">Low</span>
            <div class="sev-bar-track"><div class="sev-bar-fill low" style="width:${pct(sev.LOW,maxSev)}%"></div></div>
            <span class="sev-bar-count">${sev.LOW ?? 0}</span>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="card-header">
        <div class="card-title">${icon('bug',15)} Vulnerability Details</div>
        <span style="font-size:12px;color:var(--text-muted)">Click row to expand</span>
      </div>
      ${renderVulnTable(report.vulnerabilities ?? [])}
    </div>
  `;
}

// ─────────────────────────────────────────────────────────────
// PAGE: Scheduler
// ─────────────────────────────────────────────────────────────
async function renderScheduler(container) {
  const [tasksRes, targetsRes] = await Promise.all([call('ListScheduledScans'), call('ListTargets')]);
  const tasks = tasksRes.tasks ?? [];
  const targets = targetsRes.targets ?? [];

  container.innerHTML = `
    <div class="flex-between mb-24">
      <div style="font-size:13px;color:var(--text-muted)">${tasks.length} scheduled task${tasks.length !== 1 ? 's' : ''}</div>
      <button id="btn-add-schedule" class="btn btn-primary">${icon('plus',14)} Schedule Scan</button>
    </div>

    ${tasks.length === 0 ? `
      <div class="card">
        <div class="empty-state">
          <div class="empty-icon">${icon('clock',40)}</div>
          <h3>No scheduled scans</h3>
          <p>Schedule automatic scans to run at recurring intervals</p>
        </div>
      </div>
    ` : `
      <div class="table-wrapper">
        <table>
          <thead>
            <tr><th>#</th><th>Target</th><th>Schedule</th><th>Profile</th><th>Status</th><th>Last Run</th><th>Actions</th></tr>
          </thead>
          <tbody>
            ${tasks.map(t => `
              <tr>
                <td style="font-family:var(--font-mono);color:var(--text-muted)">${t.id}</td>
                <td class="td-mono">${escHtml(t.target_url || 'N/A')}</td>
                <td><span class="tag">${escHtml(t.cron_expr)}</span></td>
                <td><span class="tag">${escHtml(t.scan_profile)}</span></td>
                <td>${t.enabled ? '<span class="badge badge-success">Active</span>' : '<span class="badge badge-muted">● Paused</span>'}</td>
                <td style="color:var(--text-muted);font-size:12px">${t.last_run ? formatDate(t.last_run) : 'Never'}</td>
                <td class="td-actions">
                  <button class="btn btn-ghost btn-sm btn-toggle-schedule" data-id="${t.id}">${t.enabled ? `${icon('pause',13)} Pause` : `${icon('play',13)} Resume`}</button>
                  <button class="btn btn-danger btn-sm btn-del-schedule" data-id="${t.id}">${icon('trash',13)}</button>
                </td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `}

    <!-- Schedule Modal -->
    <div id="schedule-modal" class="modal-overlay" style="display:none">
      <div class="modal">
        <div class="modal-header">
          <span class="modal-title">${icon('clock',16)} Schedule Automated Scan</span>
          <button class="modal-close" id="close-schedule-modal">${icon('x',14)}</button>
        </div>
        <form id="schedule-form">
          <div class="form-group">
            <label class="form-label">Target *</label>
            <select id="sch-target" class="form-control" required>
              <option value="">Select target…</option>
              ${targets.filter(t => t.enabled).map(t => `<option value="${t.id}">${escHtml(t.name)} — ${escHtml(t.url)}</option>`).join('')}
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Schedule</label>
            <select id="sch-cron" class="form-control">
              <option value="@hourly">Every hour</option>
              <option value="@daily">Every day at midnight</option>
              <option value="@weekly">Every week (Sunday)</option>
              <option value="@monthly">Monthly (1st day)</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">Scan Profile</label>
            <select id="sch-profile" class="form-control">
              <option value="full">Full Scan</option>
              <option value="quick">Quick Scan</option>
            </select>
          </div>
          <div style="display:flex;gap:8px;justify-content:flex-end;margin-top:8px">
            <button type="button" class="btn btn-ghost" id="cancel-schedule">Cancel</button>
            <button type="submit" class="btn btn-primary">Save Schedule</button>
          </div>
        </form>
      </div>
    </div>
  `;

  const showModal = () => document.getElementById('schedule-modal').style.display = 'flex';
  const hideModal = () => document.getElementById('schedule-modal').style.display = 'none';
  document.getElementById('btn-add-schedule')?.addEventListener('click', showModal);
  document.getElementById('close-schedule-modal')?.addEventListener('click', hideModal);
  document.getElementById('cancel-schedule')?.addEventListener('click', hideModal);

  document.getElementById('schedule-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const targetID = parseInt(document.getElementById('sch-target').value);
    const cron = document.getElementById('sch-cron').value;
    const profile = document.getElementById('sch-profile').value;
    if (!targetID) { showAlert('Select a target', 'error'); return; }
    const res = await call('CreateScheduledScan', targetID, cron, profile);
    if (res.success) { showAlert('Scheduled scan created', 'success'); hideModal(); navigate('scheduler'); }
    else showAlert(res.error, 'error');
  });

  container.querySelectorAll('.btn-toggle-schedule').forEach(btn => {
    btn.addEventListener('click', async () => {
      await call('ToggleScheduledScan', parseInt(btn.dataset.id));
      navigate('scheduler');
    });
  });

  container.querySelectorAll('.btn-del-schedule').forEach(btn => {
    btn.addEventListener('click', async () => {
      if (!confirm('Delete scheduled task?')) return;
      await call('DeleteScheduledScan', parseInt(btn.dataset.id));
      showAlert('Schedule deleted', 'success');
      navigate('scheduler');
    });
  });
}

// ─────────────────────────────────────────────────────────────
// PAGE: Audit Logs
// ─────────────────────────────────────────────────────────────
async function renderAudit(container) {
  const res = await call('GetAuditLogs', 200);
  const logs = res.logs ?? [];

  container.innerHTML = `
    <div class="flex-between mb-24">
      <div style="font-size:13px;color:var(--text-muted)">${logs.length} log entries</div>
      <button class="btn btn-ghost btn-sm" onclick="navigate('audit')">${icon('rotate-cw',13)} Refresh</button>
    </div>
    ${logs.length === 0 ? `
      <div class="card">
        <div class="empty-state">
          <div class="empty-icon">${icon('clipboard-list',40)}</div>
          <h3>No audit logs yet</h3>
          <p>All user actions will appear here</p>
        </div>
      </div>
    ` : `
      <div class="table-wrapper">
        <table>
          <thead><tr><th>#</th><th>User ID</th><th>Action</th><th>Details</th><th>Timestamp</th></tr></thead>
          <tbody>
            ${logs.map(l => `
              <tr>
                <td style="font-family:var(--font-mono);color:var(--text-muted)">${l.id}</td>
                <td style="font-family:var(--font-mono);font-size:12px">${l.user_id}</td>
                <td><span class="tag">${escHtml(l.action)}</span></td>
                <td style="font-size:12px;color:var(--text-secondary)">${escHtml(l.details || '—')}</td>
                <td style="font-size:12px;color:var(--text-muted)">${formatDate(l.timestamp)}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>
      </div>
    `}
  `;
}

// ─────────────────────────────────────────────────────────────
// Scan Poller
// ─────────────────────────────────────────────────────────────
function startScanPoller() {
  if (state.pollTimer) return;
  state.pollTimer = setInterval(async () => {
    if (state.page !== 'scans') { stopScanPoller(); return; }
    const res = await call('ListScans');
    const activeCount = (res.scans ?? []).filter(s => s.status === 'running').length;
    if (activeCount === 0) {
      state.activeScans.clear();
      stopScanPoller();
    }
    renderPage();
  }, 3000);
}

function stopScanPoller() {
  clearInterval(state.pollTimer);
  state.pollTimer = null;
}

// ─────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────
function escHtml(s) {
  return String(s ?? '').replace(/&/g,'&amp;').replace(/</g,'&lt;').replace(/>/g,'&gt;').replace(/"/g,'&quot;').replace(/'/g,'&#39;');
}

function pct(val, max) { return Math.round(((val ?? 0) / max) * 100); }

function severityBadge(s) {
  const map = { CRITICAL: 'badge-critical', HIGH: 'badge-high', MEDIUM: 'badge-medium', LOW: 'badge-low' };
  return `<span class="badge ${map[s] ?? 'badge-muted'}">${escHtml(s)}</span>`;
}

function scanStatusBadge(status) {
  const map = { running: 'badge-running', completed: 'badge-completed', pending: 'badge-pending', failed: 'badge-failed' };
  const statusIcons = { running: icon('zap',12), completed: icon('check-circle',12), pending: icon('clock',12), failed: icon('x-circle',12) };
  return `<span class="badge ${map[status] ?? 'badge-muted'}">${statusIcons[status] ?? ''} ${escHtml(status)}</span>`;
}

function formatDate(s) {
  if (!s) return '—';
  const d = new Date(s);
  if (isNaN(d)) return s;
  return d.toLocaleString(undefined, { year:'numeric', month:'short', day:'numeric', hour:'2-digit', minute:'2-digit' });
}

function calcDuration(start, end) {
  if (!start || !end) return '—';
  const ms = new Date(end) - new Date(start);
  if (ms < 0 || isNaN(ms)) return '—';
  const s = Math.floor(ms / 1000);
  if (s < 60) return `${s}s`;
  return `${Math.floor(s/60)}m ${s%60}s`;
}

function errHTML(msg) {
  return `<div class="alert alert-error">${icon('alert-triangle',14)} ${escHtml(msg || 'An error occurred')}</div>`;
}

function downloadText(filename, text) {
  const blob = new Blob([text], { type: 'text/csv' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url; a.download = filename; a.click();
  URL.revokeObjectURL(url);
}

// ─────────────────────────────────────────────────────────────
// PDF Export — generates a styled print document via iframe
// ─────────────────────────────────────────────────────────────
function exportPDF(report, scanID) {
  const html = printReportHTML(report, scanID);

  // Remove any previously created print frame
  let frame = document.getElementById('apg-print-frame');
  if (frame) frame.remove();

  frame = document.createElement('iframe');
  frame.id = 'apg-print-frame';
  frame.style.cssText = 'position:fixed;top:-9999px;left:-9999px;width:210mm;height:297mm;border:none;visibility:hidden;';
  document.body.appendChild(frame);

  frame.onload = () => {
    frame.contentWindow.focus();
    frame.contentWindow.print();
    // Clean up after a delay to allow the print dialog to open
    setTimeout(() => frame.remove(), 5000);
  };

  // Write the HTML into the iframe
  const doc = frame.contentDocument ?? frame.contentWindow.document;
  doc.open();
  doc.write(html);
  doc.close();
}

function printReportHTML(report, scanID) {
  const sev = report.severity_counts ?? {};
  const vulns = report.vulnerabilities ?? [];
  const riskColors = {
    CRITICAL: '#ff3b5c', HIGH: '#ff8c42', MEDIUM: '#ffcc00', LOW: '#00c897', INFO: '#3d7aff'
  };
  const riskColor = riskColors[report.risk_level] ?? '#3d7aff';

  const sevBadgeStyle = (s) => {
    const bg = { CRITICAL: '#ff3b5c', HIGH: '#ff8c42', MEDIUM: '#ffcc00', LOW: '#00c897' };
    return `background:${bg[s] ?? '#888'};color:#000;padding:2px 8px;border-radius:4px;font-size:11px;font-weight:700;`;
  };

  const vulnRows = vulns.map((v, i) => `
    <tr style="background:${i % 2 === 0 ? '#f9f9f9' : '#fff'}">
      <td style="padding:8px;border:1px solid #ddd;font-size:12px;font-weight:600">${escHtml(v.type)}</td>
      <td style="padding:8px;border:1px solid #ddd;text-align:center">
        <span style="${sevBadgeStyle(v.severity)}">${escHtml(v.severity)}</span>
      </td>
      <td style="padding:8px;border:1px solid #ddd;font-family:monospace;font-size:11px;word-break:break-all">${escHtml(v.endpoint)}</td>
      <td style="padding:8px;border:1px solid #ddd;font-family:monospace;font-size:11px">${escHtml(v.parameter || '—')}</td>
      <td style="padding:8px;border:1px solid #ddd;font-family:monospace;font-size:10px;word-break:break-all;color:#555">${escHtml(v.payload || '—')}</td>
      <td style="padding:8px;border:1px solid #ddd;font-size:11px;color:#333">${escHtml(v.description || '—')}</td>
    </tr>`).join('');

  const generatedAt = new Date().toLocaleString();

  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8"/>
  <title>APGUARD Report — Scan #${report.scan_id}</title>
  <style>
    @page { size: A4; margin: 18mm 15mm; }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body { font-family: 'Segoe UI', Arial, sans-serif; color: #1a1a2e; background: #fff; font-size: 13px; }

    /* Header */
    .header { display: flex; align-items: center; justify-content: space-between;
              border-bottom: 3px solid #3d7aff; padding-bottom: 14px; margin-bottom: 22px; }
    .header-logo { display: flex; align-items: center; gap: 10px; }
    .logo-shield { width: 44px; height: 44px; background: linear-gradient(135deg,#3d7aff,#00d4ff);
                   border-radius: 10px; display: flex; align-items: center; justify-content: center;
                   font-size: 22px; }
    .logo-text { font-size: 24px; font-weight: 800; letter-spacing: -0.5px; color: #3d7aff; }
    .logo-sub  { font-size: 11px; color: #666; margin-top: 2px; }
    .header-meta { text-align: right; font-size: 11px; color: #666; line-height: 1.6; }

    /* Section titles */
    h2 { font-size: 14px; font-weight: 700; color: #1a1a2e; margin: 20px 0 10px;
         border-left: 4px solid #3d7aff; padding-left: 10px; }

    /* KPI grid */
    .kpi-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 20px; }
    .kpi { border: 1px solid #e0e0e0; border-radius: 8px; padding: 12px; text-align: center; }
    .kpi-value { font-size: 28px; font-weight: 800; line-height: 1; }
    .kpi-label { font-size: 11px; color: #666; margin-top: 4px; text-transform: uppercase; letter-spacing: .5px; }

    /* Details table */
    .detail-table { width: 100%; border-collapse: collapse; margin-bottom: 20px; }
    .detail-table td { padding: 7px 12px; border: 1px solid #e0e0e0; font-size: 12px; }
    .detail-table td:first-child { font-weight: 600; width: 30%; background: #f4f6fb; color: #444; }

    /* Severity bars */
    .sev-row { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
    .sev-name { width: 70px; font-size: 12px; font-weight: 600; }
    .sev-track { flex: 1; height: 10px; background: #eee; border-radius: 5px; overflow: hidden; }
    .sev-fill  { height: 100%; border-radius: 5px; }
    .sev-count { width: 30px; text-align: right; font-size: 12px; font-weight: 700; }

    /* Vuln table */
    .vuln-table { width: 100%; border-collapse: collapse; font-size: 12px; }
    .vuln-table th { background: #1a1a2e; color: #fff; padding: 9px 8px; text-align: left;
                     font-size: 11px; font-weight: 600; text-transform: uppercase; letter-spacing: .4px; }
    .vuln-table td { vertical-align: top; }

    /* Footer */
    .footer { margin-top: 30px; border-top: 1px solid #e0e0e0; padding-top: 10px;
              font-size: 10px; color: #999; display: flex; justify-content: space-between; }

    @media print {
      body { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
    }
  </style>
</head>
<body>

  <!-- HEADER -->
  <div class="header">
    <div class="header-logo">
      <div class="logo-shield">&#x1F6E1;</div>
      <div>
        <div class="logo-text">APGUARD</div>
        <div class="logo-sub">Web Application Security Scanner</div>
      </div>
    </div>
    <div class="header-meta">
      <strong>Vulnerability Scan Report</strong><br/>
      Scan ID: #${report.scan_id}<br/>
      Generated: ${generatedAt}
    </div>
  </div>

  <!-- KPI CARDS -->
  <div class="kpi-grid">
    <div class="kpi" style="border-color:${riskColor}">
      <div class="kpi-value" style="color:${riskColor}">${report.risk_level ?? 'N/A'}</div>
      <div class="kpi-label">Risk Level</div>
    </div>
    <div class="kpi">
      <div class="kpi-value" style="color:#ff3b5c">${sev.CRITICAL ?? 0}</div>
      <div class="kpi-label">Critical</div>
    </div>
    <div class="kpi">
      <div class="kpi-value" style="color:#ff8c42">${sev.HIGH ?? 0}</div>
      <div class="kpi-label">High</div>
    </div>
    <div class="kpi">
      <div class="kpi-value" style="color:#3d7aff">${report.total_vulns ?? 0}</div>
      <div class="kpi-label">Total Vulns</div>
    </div>
  </div>

  <!-- SCAN DETAILS -->
  <h2>Scan Details</h2>
  <table class="detail-table">
    <tr><td>Target URL</td><td style="font-family:monospace">${escHtml(report.target_url ?? '')}</td></tr>
    <tr><td>Status</td><td>${escHtml(report.status ?? '')}</td></tr>
    <tr><td>Scan Profile</td><td>${escHtml(report.scan_profile ?? 'full')}</td></tr>
    <tr><td>Started</td><td>${escHtml(report.start_time ?? '—')}</td></tr>
    <tr><td>Ended</td><td>${escHtml(report.end_time ?? '—')}</td></tr>
    <tr><td>Risk Score</td><td>${report.risk_score ?? 0}</td></tr>
  </table>

  <!-- SEVERITY BREAKDOWN -->
  <h2>Severity Distribution</h2>
  ${(() => {
    const maxSev = Math.max(sev.CRITICAL??0, sev.HIGH??0, sev.MEDIUM??0, sev.LOW??0, 1);
    const bars = [
      { label:'Critical', key:'CRITICAL', color:'#ff3b5c' },
      { label:'High',     key:'HIGH',     color:'#ff8c42' },
      { label:'Medium',   key:'MEDIUM',   color:'#ffcc00' },
      { label:'Low',      key:'LOW',      color:'#00c897' },
    ];
    return bars.map(b => `
      <div class="sev-row">
        <span class="sev-name" style="color:${b.color}">${b.label}</span>
        <div class="sev-track">
          <div class="sev-fill" style="width:${Math.round(((sev[b.key]??0)/maxSev)*100)}%;background:${b.color}"></div>
        </div>
        <span class="sev-count" style="color:${b.color}">${sev[b.key]??0}</span>
      </div>`).join('');
  })()}

  <!-- VULNERABILITIES -->
  <h2 style="margin-top:22px">Vulnerability Details (${vulns.length})</h2>
  ${vulns.length === 0
    ? '<p style="color:#888;font-size:13px;padding:16px 0">✓ No vulnerabilities detected in this scan.</p>'
    : `<table class="vuln-table">
        <thead><tr>
          <th style="width:18%">Type</th>
          <th style="width:10%;text-align:center">Severity</th>
          <th style="width:22%">Endpoint</th>
          <th style="width:12%">Parameter</th>
          <th style="width:18%">Payload</th>
          <th>Description</th>
        </tr></thead>
        <tbody>${vulnRows}</tbody>
      </table>`}

  <!-- FOOTER -->
  <div class="footer">
    <span>APGUARD v1.0 — Confidential Security Report</span>
    <span>Scan #${report.scan_id} — ${escHtml(report.target_url ?? '')}</span>
  </div>

</body>
</html>`;
}


// ─────────────────────────────────────────────────────────────
// Bootstrap
// ─────────────────────────────────────────────────────────────
async function init() {
  // Try to restore a session on startup — works in both modes:
  //   Wails desktop: backend restored the session from ~/.apguard/session.json
  //   HTTP browser:  we have a JWT in localStorage
  if (GO || Token.get()) {
    const res = await call('GetCurrentUser');
    if (res.success && res.user) {
      state.user = res.user;
    } else if (!GO) {
      Token.clear(); // token expired/invalid (HTTP mode only)
    }
  }
  render();
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init);
} else {
  init();
}
