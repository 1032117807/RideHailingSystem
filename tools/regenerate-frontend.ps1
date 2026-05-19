$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot
$frontend = Join-Path $root "frontend"
$assets = Join-Path $frontend "assets"

function Write-Utf8File {
  param(
    [Parameter(Mandatory = $true)][string]$Path,
    [Parameter(Mandatory = $true)][string]$Content
  )
  [System.IO.File]::WriteAllText($Path, $Content.TrimStart(), [System.Text.UTF8Encoding]::new($false))
}

function Header {
  param([string]$Title, [string]$Page)
@"
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>$Title - TripVerse</title>
  <link rel="stylesheet" href="assets/styles.css">
</head>
<body data-page="$Page">
  <header class="site-header">
    <div class="container nav-shell">
      <a class="brand" href="index.html" aria-label="TripVerse 首页">
        <span class="brand-mark">TV</span>
        <span>TripVerse</span>
      </a>
      <nav class="nav-links" aria-label="主导航"></nav>
      <button class="menu-button" type="button" data-nav-toggle aria-label="打开导航">菜单</button>
    </div>
    <div class="mobile-drawer">
      <div class="container mobile-links"></div>
    </div>
  </header>
"@
}

function Footer {
@"
  <script src="assets/runtime-config.js?v=20260519-blue-ui"></script>
  <script src="assets/app.js?v=20260519-blue-ui"></script>
</body>
</html>
"@
}

function Page {
  param([string]$Title, [string]$Page, [string]$Content)
  "$(Header $Title $Page)`n$Content`n$(Footer)"
}

$styles = @'
:root {
  --color-primary: #3b82f6;
  --color-primary-light: #dbeafe;
  --color-primary-soft: #eff6ff;
  --color-primary-hover: #2563eb;
  --color-bg: #f8fbff;
  --color-surface: #ffffff;
  --color-surface-soft: #f1f7ff;
  --color-border: #dbeafe;
  --color-border-strong: #bfdbfe;
  --color-text-main: #0f172a;
  --color-text-secondary: #475569;
  --color-text-muted: #94a3b8;
  --color-success: #22c55e;
  --color-warning: #f59e0b;
  --color-error: #ef4444;
  --shadow-soft: 0 12px 30px rgba(37, 99, 235, 0.08);
  --shadow-hover: 0 18px 38px rgba(37, 99, 235, 0.14);
  --radius-card: 18px;
  --radius-button: 10px;
  --radius-control: 10px;
  --space-xs: 4px;
  --space-sm: 8px;
  --space-md: 16px;
  --space-lg: 24px;
  --space-xl: 32px;
  --space-2xl: 48px;
  --container: 1180px;
}

* { box-sizing: border-box; }
html { scroll-behavior: smooth; }
body {
  margin: 0;
  min-height: 100vh;
  color: var(--color-text-main);
  background:
    linear-gradient(180deg, #f8fbff 0%, #eef6ff 42%, #f8fbff 100%);
  font-family: Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Microsoft YaHei", Arial, sans-serif;
}
a { color: inherit; text-decoration: none; }
button, input, select, textarea { font: inherit; }
button { cursor: pointer; }
img { display: block; max-width: 100%; }

.container {
  width: min(var(--container), calc(100% - 32px));
  margin: 0 auto;
}
.site-header {
  position: sticky;
  top: 0;
  z-index: 50;
  background: rgba(248, 251, 255, 0.86);
  border-bottom: 1px solid var(--color-border);
  backdrop-filter: blur(18px);
}
.nav-shell {
  min-height: 72px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-lg);
}
.brand {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  font-size: 1.24rem;
  font-weight: 800;
}
.brand-mark {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border-radius: 14px;
  color: #fff;
  background: linear-gradient(135deg, var(--color-primary), #60a5fa);
  box-shadow: 0 12px 24px rgba(59, 130, 246, 0.24);
}
.nav-links, .mobile-links, .tab-strip, .actions, .button-row, .chip-row, .filter-row, .meta-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  align-items: center;
}
.nav-links a, .mobile-links a, .tab-strip a, .nav-user-chip {
  min-height: 40px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0 14px;
  border-radius: 999px;
  color: var(--color-text-secondary);
  transition: background-color 0.2s ease, color 0.2s ease, box-shadow 0.2s ease;
}
.nav-links a:hover, .mobile-links a:hover, .tab-strip a:hover {
  color: var(--color-primary-hover);
  background: var(--color-primary-soft);
}
.nav-links a.is-active, .mobile-links a.is-active, .tab-strip a.is-active {
  color: #fff;
  background: var(--color-primary);
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.18);
}
.nav-user-chip {
  background: var(--color-primary-soft);
  border: 1px solid var(--color-border);
}
.nav-logout {
  min-height: 40px;
  padding: 0 14px;
  color: var(--color-primary-hover);
  background: #fff;
  border: 1px solid var(--color-border-strong);
  border-radius: 999px;
}
.menu-button {
  display: none;
  min-height: 42px;
  padding: 0 16px;
  color: #fff;
  background: var(--color-primary);
  border: 1px solid var(--color-primary);
  border-radius: var(--radius-button);
}
.mobile-drawer { display: none; border-top: 1px solid var(--color-border); }
.mobile-links { padding: 14px 0 18px; align-items: stretch; }
.mobile-links a, .nav-logout-mobile, .nav-user-chip-mobile { width: 100%; }
body.nav-open .mobile-drawer { display: block; }

.page { padding: 28px 0 72px; }
.page-hero { padding: 22px 0 18px; }
.section-block { margin-top: var(--space-lg); }
.hero-grid, .split-grid, .card-grid, .stat-grid, .dashboard-grid, .form-grid, .search-grid, .auth-layout, .chat-layout {
  display: grid;
  gap: 20px;
}
.hero-grid { grid-template-columns: minmax(0, 1.25fr) minmax(320px, 0.85fr); align-items: stretch; }
.dashboard-grid { grid-template-columns: minmax(0, 1.1fr) minmax(320px, 0.9fr); align-items: start; }
.split-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); align-items: start; }
.card-grid { grid-template-columns: repeat(3, minmax(0, 1fr)); }
.stat-grid { grid-template-columns: repeat(4, minmax(0, 1fr)); }
.form-grid, .search-grid { grid-template-columns: repeat(12, minmax(0, 1fr)); }
.auth-layout, .chat-layout { grid-template-columns: minmax(300px, 0.9fr) minmax(0, 1.1fr); }

.hero-copy, .hero-panel, .panel, .stat-card, .table-card, .banner-card, .auth-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-card);
  box-shadow: var(--shadow-soft);
}
.hero-copy { padding: 38px; }
.hero-panel, .panel, .table-card, .banner-card, .auth-card, .stat-card { padding: var(--space-lg); }
.panel, .stat-card, .info-card, .list-item, .ticket-card, .order-card, .seat-card, .message {
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease, background-color 0.2s ease;
}
.panel:hover, .stat-card:hover, .ticket-card:hover, .order-card:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-hover);
  border-color: var(--color-border-strong);
}
.eyebrow, .badge, .tag, .mini-chip {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  border-radius: 999px;
  font-weight: 700;
}
.eyebrow {
  padding: 8px 12px;
  color: #1d4ed8;
  background: var(--color-primary-soft);
  border: 1px solid var(--color-border);
}
.badge { padding: 7px 12px; color: #1e40af; background: #dbeafe; }
.tag, .mini-chip { padding: 7px 12px; color: var(--color-primary-hover); background: var(--color-primary-soft); border: 1px solid var(--color-border); font-size: 0.9rem; }
.headline {
  margin: 18px 0 12px;
  max-width: 820px;
  font-size: clamp(2rem, 5vw, 3.6rem);
  line-height: 1.1;
  letter-spacing: 0;
}
.section-title { margin: 12px 0 8px; font-size: clamp(1.75rem, 3vw, 2.35rem); line-height: 1.18; letter-spacing: 0; }
.subhead { margin: 0; font-size: 1.15rem; line-height: 1.35; font-weight: 700; }
.lede, .muted, .list-meta, .table-note, .detail-copy p {
  color: var(--color-text-secondary);
  line-height: 1.72;
}
.row-between, .ticket-top, .order-top, .pricing-row, .table-row {
  display: flex;
  justify-content: space-between;
  gap: 14px;
  align-items: center;
  flex-wrap: wrap;
}
.list-stack, .timeline, .search-card, .chat-window, .auth-side, .stack-tight, .detail-copy {
  display: grid;
  gap: 14px;
}

.button {
  min-height: 44px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 0 18px;
  border-radius: var(--radius-button);
  border: 1px solid transparent;
  font-weight: 700;
  transition: transform 0.2s ease, background-color 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease, color 0.2s ease;
}
.button:hover { transform: translateY(-1px); }
.button:focus-visible, input:focus-visible, select:focus-visible, textarea:focus-visible, .seat:focus-visible, .menu-button:focus-visible {
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.18);
}
.button-primary {
  color: #fff;
  background: var(--color-primary);
  border-color: var(--color-primary);
  box-shadow: 0 8px 20px rgba(59, 130, 246, 0.22);
}
.button-primary:hover { background: var(--color-primary-hover); border-color: var(--color-primary-hover); }
.button-secondary, .button-dark {
  color: var(--color-primary-hover);
  background: var(--color-primary-soft);
  border-color: var(--color-border-strong);
}
.button-ghost {
  color: var(--color-text-main);
  background: #fff;
  border-color: var(--color-border);
}
.button-danger {
  color: #fff;
  background: var(--color-error);
  border-color: var(--color-error);
}
.button[disabled], button[disabled] { opacity: 0.58; cursor: not-allowed; transform: none; }

.field { display: grid; gap: 8px; }
.field label { color: var(--color-text-secondary); font-size: 0.95rem; font-weight: 600; }
.field input, .field select, .field textarea, .table-inline-select {
  width: 100%;
  min-height: 42px;
  padding: 0 12px;
  color: var(--color-text-main);
  background: #fff;
  border: 1px solid var(--color-border-strong);
  border-radius: var(--radius-control);
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
}
.field textarea { min-height: 120px; padding: 12px; resize: vertical; }
.field input:focus, .field select:focus, .field textarea:focus, .table-inline-select:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}
.checkbox-row { display: flex !important; align-items: center; gap: 10px; min-height: 42px; }
.checkbox-row input { width: 18px; min-height: auto; height: 18px; }
.span-3 { grid-column: span 3; }
.span-4 { grid-column: span 4; }
.span-5 { grid-column: span 5; }
.span-6 { grid-column: span 6; }
.span-7 { grid-column: span 7; }
.span-8 { grid-column: span 8; }
.span-12 { grid-column: 1 / -1; }

.info-card, .list-item, .ticket-card, .order-card, .seat-card, .message {
  padding: 18px;
  border: 1px solid var(--color-border);
  border-radius: 16px;
  background: #fff;
}
.info-card { background: var(--color-primary-soft); }
.empty-state {
  padding: 34px 26px;
  text-align: center;
  color: var(--color-text-secondary);
  background: #fff;
  border: 1px dashed var(--color-border-strong);
  border-radius: var(--radius-card);
}
.skeleton {
  min-height: 18px;
  border-radius: 999px;
  background: linear-gradient(90deg, #eaf3ff, #f8fbff, #eaf3ff);
  background-size: 200% 100%;
  animation: shimmer 1.2s infinite linear;
}
@keyframes shimmer { from { background-position: 200% 0; } to { background-position: -200% 0; } }
.fade-in { animation: fadeIn 0.24s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(6px); } to { opacity: 1; transform: translateY(0); } }

.stat-label { color: var(--color-text-secondary); font-size: 0.95rem; }
.stat-value { margin-top: 10px; color: #1d4ed8; font-size: clamp(1.8rem, 3vw, 2.5rem); font-weight: 800; }
.stat-delta { margin-top: 12px; display: inline-flex; padding: 7px 10px; border-radius: 999px; color: #166534; background: #dcfce7; font-size: 0.88rem; }
.route-line { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.route-city { font-size: 1.2rem; font-weight: 800; }
.route-divider { flex: 1; min-width: 60px; height: 2px; border-radius: 999px; background: linear-gradient(90deg, var(--color-primary), #93c5fd); }
.price-pill { display: inline-flex; align-items: baseline; gap: 4px; color: #1d4ed8; font-size: 2rem; font-weight: 800; }
.price-pill small { color: var(--color-text-secondary); font-size: 0.9rem; }
.timeline-item { position: relative; padding-left: 44px; }
.timeline-item::before { content: ""; position: absolute; left: 17px; top: 24px; bottom: -14px; width: 2px; background: var(--color-border-strong); }
.timeline-item:last-child::before { display: none; }
.timeline-dot { position: absolute; left: 8px; top: 16px; width: 20px; height: 20px; border-radius: 999px; background: var(--color-primary); box-shadow: 0 0 0 5px var(--color-primary-soft); }

.table-card { overflow: hidden; }
.table-scroll { overflow-x: auto; }
table, .table { width: 100%; border-collapse: separate; border-spacing: 0; }
th, td {
  padding: 14px 16px;
  text-align: left;
  white-space: nowrap;
  border-bottom: 1px solid #e0f2fe;
}
th {
  color: #1e3a8a;
  background: var(--color-primary-soft);
  font-weight: 700;
}
tbody tr:hover td { background: #f8fbff; }
.table-actions { display: grid; gap: 10px; min-width: 190px; }
.table-actions .button { min-height: 38px; }

.seat-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.seat {
  padding: 16px 10px;
  border: 1px solid var(--color-border-strong);
  border-radius: 14px;
  background: #fff;
  text-align: center;
}
.seat:hover, .seat.is-selected {
  color: var(--color-primary-hover);
  background: var(--color-primary-soft);
  border-color: var(--color-primary);
}
.seat.is-disabled { opacity: 0.5; cursor: not-allowed; }
.pricing-box {
  display: grid;
  gap: 12px;
  padding: 20px;
  color: var(--color-text-main);
  background: var(--color-primary-soft);
  border: 1px solid var(--color-border);
  border-radius: 16px;
}
.pricing-total { padding-top: 12px; border-top: 1px solid var(--color-border-strong); font-size: 1.12rem; }
.ai-session-shell { display: grid; gap: 12px; margin-top: 18px; padding: 14px; border: 1px solid var(--color-border); border-radius: 16px; background: var(--color-primary-soft); }
.ai-session-tabs { display: flex; gap: 10px; overflow-x: auto; }
.ai-session-tab { flex: 0 0 auto; min-height: 38px; padding: 0 14px; border: 1px solid var(--color-border-strong); border-radius: 999px; background: #fff; color: var(--color-text-secondary); }
.ai-session-tab.is-active { color: #fff; background: var(--color-primary); }
.message.user { background: #fff; border-color: var(--color-border-strong); }
.message.ai { background: var(--color-primary-soft); }
.toast {
  position: fixed;
  right: 20px;
  bottom: 20px;
  z-index: 80;
  max-width: min(420px, calc(100% - 40px));
  padding: 14px 18px;
  color: #fff;
  background: var(--color-primary-hover);
  border-radius: 14px;
  box-shadow: var(--shadow-hover);
  opacity: 0;
  pointer-events: none;
  transform: translateY(10px);
  transition: 0.24s ease;
}
.toast.is-visible { opacity: 1; transform: translateY(0); }
@media (max-width: 1100px) {
  .hero-grid, .dashboard-grid, .split-grid, .auth-layout, .chat-layout { grid-template-columns: 1fr; }
  .card-grid, .stat-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 900px) {
  .nav-links { display: none; }
  .menu-button { display: inline-flex; }
  .span-3, .span-4, .span-5, .span-6, .span-7, .span-8 { grid-column: span 6; }
}
@media (max-width: 680px) {
  .container { width: min(var(--container), calc(100% - 20px)); }
  .page { padding-bottom: 48px; }
  .hero-copy, .hero-panel, .panel, .table-card, .banner-card, .auth-card, .stat-card { padding: 20px; }
  .headline { font-size: clamp(2rem, 10vw, 2.8rem); }
  .card-grid, .stat-grid, .seat-grid { grid-template-columns: 1fr; }
  .form-grid, .search-grid { grid-template-columns: repeat(6, minmax(0, 1fr)); }
  .span-3, .span-4, .span-5, .span-6, .span-7, .span-8, .span-12 { grid-column: 1 / -1; }
  .actions .button, .button-row .button, .filter-row .button { width: 100%; }
}
'@

$pages = @{}

$pages["index.html"] = @{ Title = "首页"; Page = "home"; Content = @'
  <main class="page">
    <section class="page-hero">
      <div class="container hero-grid">
        <div class="hero-copy fade-in">
          <span class="eyebrow">一站式出行票务平台</span>
          <h1 class="headline">从城市通勤到跨城出行，一次搜索就能安排清楚。</h1>
          <p class="lede">TripVerse 统一乘客购票、司机发车、管理审核和 AI 出行助手，采用浅蓝色企业级界面，让复杂票务流程更清晰、更可靠。</p>
          <div class="actions">
            <a class="button button-primary" href="search.html">立即购票</a>
            <a class="button button-secondary" href="ai-assistant.html">体验 AI 助手</a>
          </div>
        </div>
        <aside class="hero-panel">
          <div class="search-card">
            <span class="badge">快速搜索</span>
            <div class="form-grid">
              <div class="field span-6"><label for="home-start">出发地</label><input id="home-start" value="杭州"></div>
              <div class="field span-6"><label for="home-end">目的地</label><input id="home-end" value="苏州"></div>
              <div class="field span-6"><label for="home-date">日期</label><input id="home-date" type="date" value="2026-05-19"></div>
              <div class="field span-6"><label for="home-pref">偏好</label><select id="home-pref"><option>高铁优先</option><option>最快到达</option><option>价格最低</option></select></div>
            </div>
            <a class="button button-primary" href="search.html">查看班次</a>
            <div class="chip-row"><span class="tag">低价提醒</span><span class="tag">电子票</span><span class="tag">AI 推荐</span></div>
          </div>
        </aside>
      </div>
    </section>
    <section class="section-block">
      <div class="container card-grid">
        <article class="panel"><span class="badge">乘客端</span><h2 class="subhead">搜索、下单、支付、售后闭环</h2><p class="muted">按真实班次数据展示结果，支持订单状态、退款改签和个人通知。</p><a class="button button-ghost" href="orders.html">查看订单</a></article>
        <article class="panel"><span class="badge">司机端</span><h2 class="subhead">班次发布与收益管理</h2><p class="muted">司机可发布车次、核验订单、查看收入，并用 AI 生成运营方案。</p><a class="button button-ghost" href="driver-dashboard.html">进入司机端</a></article>
        <article class="panel"><span class="badge">管理端</span><h2 class="subhead">用户、订单、风控和模型治理</h2><p class="muted">后台集中处理退款、风控日志、知识库、MCP 工具和 Token 配额。</p><a class="button button-ghost" href="admin-dashboard.html">进入管理端</a></article>
      </div>
    </section>
    <section class="section-block">
      <div class="container stat-grid">
        <div class="stat-card"><div class="stat-label">今日搜索</div><div class="stat-value">82,640</div><div class="stat-delta">较昨日 +12%</div></div>
        <div class="stat-card"><div class="stat-label">支付成功率</div><div class="stat-value">98.4%</div><div class="stat-delta">链路稳定</div></div>
        <div class="stat-card"><div class="stat-label">AI 响应</div><div class="stat-value">1.6s</div><div class="stat-delta">命中缓存 41%</div></div>
        <div class="stat-card"><div class="stat-label">移动端访问</div><div class="stat-value">67%</div><div class="stat-delta">响应式适配</div></div>
      </div>
    </section>
  </main>
'@ }

$pages["search.html"] = @{ Title = "票务搜索"; Page = "search"; Content = @'
  <main class="page">
    <div class="container">
      <div class="tab-strip"><a class="is-active" href="search.html">票务搜索</a><a href="trip-detail.html">班次详情</a><a href="checkout.html">下单确认</a><a href="payment.html">支付</a></div>
      <section class="page-hero">
        <div class="hero-copy">
          <span class="eyebrow">乘客端 / 搜索</span>
          <h1 class="section-title">查找合适的出行班次</h1>
          <p class="lede">输入起点、终点和日期后，页面会调用后端票务搜索接口；搜索、加载、空状态和错误状态都在同一区域清晰反馈。</p>
          <div class="form-grid section-block" data-ticket-search-form>
            <div class="field span-3"><label for="startCity">出发地</label><input id="startCity" name="startCity" value="杭州"></div>
            <div class="field span-3"><label for="endCity">目的地</label><input id="endCity" name="endCity" value="苏州"></div>
            <div class="field span-3"><label for="date">日期</label><input id="date" name="date" type="date" value="2026-05-19"></div>
            <div class="field span-3"><label for="preference">偏好</label><select id="preference" name="preference"><option>优先推荐</option><option>最快到达</option><option>最低价格</option></select></div>
            <div class="field span-3"><label class="checkbox-row"><input name="allowTransfer" type="checkbox" checked><span>允许一次中转</span></label></div>
          </div>
        </div>
      </section>
      <section class="dashboard-grid">
        <article class="panel">
          <div class="row-between"><div><h2 class="subhead">搜索结果</h2><p class="muted">结果由接口返回，卡片会保留真实 ticketId 进入详情页。</p></div><span class="mini-chip" data-ticket-search-count>加载中</span></div>
          <div class="list-stack section-block" data-ticket-search-results><div class="info-card"><strong>正在加载班次</strong><p class="muted">请稍候，正在请求后端搜索接口。</p><div class="skeleton"></div></div></div>
        </article>
        <aside class="panel"><h2 class="subhead">搜索提示</h2><div class="list-stack section-block"><div class="info-card"><strong>完整链路</strong><p class="muted">搜索结果 -> 查看详情 -> 下单确认 -> 创建订单并支付。</p></div><div class="info-card"><strong>空状态处理</strong><p class="muted">若暂无班次，可更换日期或保留中转选项重新搜索。</p></div></div></aside>
      </section>
    </div>
  </main>
'@ }

$pages["trip-detail.html"] = @{ Title = "班次详情"; Page = "trip-detail"; Content = @'
  <main class="page"><div class="container">
    <div class="tab-strip"><a href="search.html">票务搜索</a><a class="is-active" href="trip-detail.html">班次详情</a><a data-ticket-detail-checkout-tab href="checkout.html">下单确认</a></div>
    <section class="page-hero"><div class="hero-copy"><span class="eyebrow">班次详情</span><h1 class="section-title" data-ticket-detail-title>正在加载班次信息</h1><p class="lede">查看路线、停靠点、余票、票价和订单入口。</p><div class="chip-row" data-ticket-detail-tags><span class="tag">加载中</span></div></div></section>
    <section class="dashboard-grid">
      <article class="panel"><div class="route-line" data-ticket-detail-route><span class="route-city">出发站</span><span class="route-divider"></span><span class="route-city">到达站</span></div><div class="list-stack section-block" data-ticket-detail-stops><div class="info-card">正在加载停靠站</div></div></article>
      <aside class="panel"><h2 class="subhead">购票信息</h2><div class="pricing-box section-block"><div class="pricing-row"><span>参考票价</span><strong data-ticket-detail-price>--</strong></div><div class="pricing-row"><span>余票</span><strong data-ticket-detail-seat>--</strong></div><div class="pricing-row"><span>状态</span><strong data-ticket-detail-status>加载中</strong></div></div><a class="button button-primary section-block" data-ticket-detail-checkout-link href="checkout.html">选择此班次</a></aside>
    </section>
  </div></main>
'@ }

$pages["checkout.html"] = @{ Title = "下单确认"; Page = "checkout"; Content = @'
  <main class="page"><div class="container">
    <div class="tab-strip"><a data-checkout-back-link href="trip-detail.html">返回班次</a><a class="is-active" href="checkout.html">下单确认</a><a href="payment.html">支付</a></div>
    <section class="dashboard-grid">
      <article class="panel"><span class="eyebrow">乘客信息</span><h1 class="section-title" data-checkout-title>确认订单</h1><p class="lede" data-checkout-lede>填写乘车人信息并确认座席数量。</p><div class="form-grid section-block">
        <div class="field span-4"><label for="passengerName">乘车人</label><input id="passengerName" name="passengerName" value="张三"></div>
        <div class="field span-4"><label for="idCard">证件号</label><input id="idCard" name="idCard" value="330100199001011234"></div>
        <div class="field span-4"><label for="phone">手机号</label><input id="phone" name="phone" value="13800000000"></div>
        <div class="field span-4"><label for="seatType">座席</label><select id="seatType" name="seatType" data-seat-type><option data-multiplier="1" value="standard">标准座</option><option data-multiplier="1.4" value="business">商务座</option></select></div>
        <div class="field span-4"><label for="ticketCount">数量</label><input id="ticketCount" name="ticketCount" data-ticket-count type="number" min="1" value="1"></div>
      </div><div class="button-row section-block"><button class="button button-primary" type="button" data-submit-order>提交订单</button><button class="button button-secondary" type="button" data-save-draft>保存草稿</button></div></article>
      <aside class="panel"><h2 class="subhead">费用预览</h2><div class="pricing-box section-block"><div data-checkout-route>等待班次数据</div><div class="pricing-row"><span>基础票价</span><strong data-base-price data-checkout-base-price data-base-price="0">--</strong></div><div class="pricing-row"><span>出发时间</span><strong data-checkout-departure>--</strong></div><div class="pricing-row"><span>剩余座位</span><strong data-checkout-seat-available>--</strong></div><div class="pricing-row pricing-total"><span>合计</span><strong data-total-output>--</strong></div></div><a class="button button-ghost section-block" data-checkout-trip-link href="trip-detail.html">查看班次</a><div data-checkout-actions></div></aside>
    </section>
  </div></main>
'@ }

$pages["orders.html"] = @{ Title = "我的订单"; Page = "orders"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">乘客端 / 订单</span><h1 class="section-title">我的订单</h1><p class="lede">订单列表支持加载、空状态、取消、退款和详情跳转。</p></div></section><section class="dashboard-grid"><article class="panel"><div class="chip-row" data-order-summary><span class="mini-chip">加载中</span></div><div class="list-stack section-block" data-order-list><div class="info-card">正在加载订单</div></div></article><aside class="panel"><h2 class="subhead">售后说明</h2><p class="muted">订单状态会随支付、退款审核和司机核验同步更新。</p></aside></section></div></main>
'@ }

$pages["order-detail.html"] = @{ Title = "订单详情"; Page = "order-detail"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">订单详情</span><h1 class="section-title" data-order-detail-title>正在加载订单</h1><div class="chip-row"><span class="mini-chip" data-order-detail-status>加载中</span></div></div></section><section class="dashboard-grid"><article class="panel"><h2 class="subhead">订单进度</h2><div class="timeline section-block" data-order-detail-timeline><div class="info-card">正在加载时间线</div></div><div class="button-row section-block" data-order-detail-actions></div></article><aside class="panel"><h2 class="subhead">订单信息</h2><div class="list-stack section-block" data-order-detail-meta></div><div class="pricing-box section-block" data-order-detail-pricing></div></aside></section></div></main>
'@ }

$pages["payment.html"] = @{ Title = "订单支付"; Page = "payment"; Content = @'
  <main class="page"><div class="container"><section class="dashboard-grid"><article class="panel"><span class="eyebrow">支付中心</span><h1 class="section-title" data-payment-title>确认支付</h1><p class="lede" data-payment-lede>创建支付单后可模拟支付成功，便于联调订单状态。</p><div class="button-row section-block" data-payment-actions><button class="button button-primary" type="button" data-payment-create>创建支付</button></div></article><aside class="panel"><h2 class="subhead">支付摘要</h2><div class="list-stack section-block"><div data-payment-summary class="info-card">等待订单信息</div><div data-payment-order class="info-card">订单号 --</div><div data-payment-amount class="pricing-box">金额 --</div></div></aside></section></div></main>
'@ }

$pages["profile.html"] = @{ Title = "个人中心"; Page = "profile"; Content = @'
  <main class="page"><div class="container"><section class="dashboard-grid"><article class="panel"><span class="eyebrow">账号资料</span><h1 class="section-title">个人中心</h1><div class="form-grid section-block"><div class="field span-6"><label for="nickname">昵称</label><input id="nickname" name="nickname" value="TripVerse 用户"></div><div class="field span-6"><label for="phone2">手机号</label><input id="phone2" name="phone" value="13800000000"></div><div class="field span-6"><label for="email">邮箱</label><input id="email" name="email" value="user@example.com"></div><div class="field span-6"><label for="verified">实名状态</label><select id="verified" name="realNameVerified"><option value="true">已实名</option><option value="false">未实名</option></select></div><div class="field span-6"><label for="notificationText">通知文案</label><input id="notificationText" name="notificationText" value="系统消息"></div></div><button class="button button-primary section-block" type="button" data-save-profile>保存资料</button></article><aside class="panel"><div class="row-between"><h2 class="subhead">通知</h2><span class="mini-chip" data-notification-unread>0 未读</span></div><button class="button button-secondary section-block" type="button" data-notification-read-all>全部已读</button><div class="list-stack section-block" data-notification-list><div class="info-card">正在加载通知</div></div></aside></section></div></main>
'@ }

$pages["login.html"] = @{ Title = "登录"; Page = "login"; Content = @'
  <main class="page"><div class="container auth-layout auth-layout-single"><section class="auth-card"><span class="eyebrow">欢迎回来</span><h1 class="section-title">登录 TripVerse</h1><form class="form-grid section-block" data-auth-form="login"><div class="field span-12"><label for="role">登录身份</label><select id="role" name="role"><option value="passenger">乘客</option><option value="driver">司机</option><option value="admin">管理员</option></select></div><div class="field span-12"><label for="loginMode">登录方式</label><select id="loginMode" name="loginMode" data-login-mode><option value="password">手机号 + 密码</option><option value="code">邮箱验证码</option></select></div><div class="field span-12" data-phone-group><label for="phone">手机号</label><input id="phone" name="phone" value="13800000000"></div><div class="field span-12" data-password-group><label for="password">密码</label><input id="password" name="password" type="password" value="password123"></div><div class="field span-12" data-email-group><label for="email">邮箱</label><input id="email" name="email" value="user@example.com"></div><div class="field span-12" data-code-group><label for="emailCode">验证码</label><input id="emailCode" name="emailCode" value="123456"></div><button class="button button-secondary span-12" type="button" data-send-code>发送验证码</button><button class="button button-primary span-12" type="submit">登录</button></form><div class="button-row section-block"><a class="button button-ghost" href="register.html">创建新账号</a></div></section></div></main>
'@ }

$pages["register.html"] = @{ Title = "注册"; Page = "register"; Content = @'
  <main class="page"><div class="container auth-layout auth-layout-single"><section class="auth-card"><span class="eyebrow">创建账号</span><h1 class="section-title">注册 TripVerse</h1><form class="form-grid section-block" data-auth-form="register"><div class="field span-12"><label for="role">注册身份</label><select id="role" name="role"><option value="passenger">乘客</option><option value="driver">司机</option></select></div><div class="field span-6" data-phone-group><label for="phone">手机号</label><input id="phone" name="phone" value="13900000000"></div><div class="field span-6" data-email-group><label for="email">邮箱</label><input id="email" name="email" value="new@example.com"></div><div class="field span-6" data-code-group><label for="emailCode">验证码</label><input id="emailCode" name="emailCode" value="123456"></div><div class="field span-6" data-password-group><label for="password">密码</label><input id="password" name="password" type="password" value="password123"></div><button class="button button-secondary span-6" type="button" data-send-code>发送验证码</button><button class="button button-primary span-6" type="submit">注册</button></form><div class="button-row section-block"><a class="button button-ghost" href="login.html">已有账号登录</a></div></section></div></main>
'@ }

$pages["ai-assistant.html"] = @{ Title = "AI 助手"; Page = "ai"; Content = @'
  <main class="page"><div class="container chat-layout"><aside class="panel"><span class="eyebrow">AI 出行助手</span><h1 class="section-title">自然语言规划出行</h1><div class="list-stack section-block"><button class="button button-secondary" type="button" data-ai-suggestion="明天早上从杭州到苏州，帮我找最快方案">最快方案</button><button class="button button-secondary" type="button" data-ai-suggestion="今晚去上海虹桥，帮我控制预算">低价方案</button></div></aside><section class="panel"><div class="chat-window" data-ai-chat><div class="message ai"><strong>AI 助手</strong><div>你好，我可以帮你搜索车票、解释退改规则并整理出行方案。</div></div></div><form class="form-grid section-block" data-ai-form><div class="field span-12"><label for="aiText">输入问题</label><textarea id="aiText" name="message">明天早上杭州到苏州有哪些推荐？</textarea></div><button class="button button-primary span-12" type="submit">发送</button></form></section></div></main>
'@ }

$pages["driver-dashboard.html"] = @{ Title = "司机工作台"; Page = "driver"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">司机端</span><h1 class="section-title">运营工作台</h1><p class="lede">今日车次、核验、退款提醒和待发班次集中展示。</p></div></section><section class="stat-grid"><div class="stat-card"><div class="stat-label">今日班次</div><div class="stat-value" data-driver-dashboard-today-trips>--</div></div><div class="stat-card"><div class="stat-label">已完成</div><div class="stat-value" data-driver-dashboard-completed-trips>--</div></div><div class="stat-card"><div class="stat-label">已售票</div><div class="stat-value" data-driver-dashboard-sold-tickets>--</div></div><div class="stat-card"><div class="stat-label">上座率</div><div class="stat-value" data-driver-dashboard-seat-rate>--</div></div></section><section class="dashboard-grid section-block"><article class="panel"><h2 class="subhead">待发班次</h2><div class="list-stack section-block" data-driver-dashboard-upcoming-list></div></article><aside class="panel"><div class="pricing-box"><div class="pricing-row"><span>今日收入</span><strong data-driver-dashboard-today-income>--</strong></div><div class="pricing-row"><span>待核验</span><strong data-driver-dashboard-pending-verify>--</strong></div><div class="pricing-row"><span>退款数</span><strong data-driver-dashboard-refund-count>--</strong></div></div><div class="list-stack section-block" data-driver-dashboard-alert-list></div></aside></section></div></main>
'@ }

$pages["driver-publish.html"] = @{ Title = "发布班次"; Page = "driver"; Content = @'
  <main class="page"><div class="container dashboard-grid"><section class="panel"><span class="eyebrow">司机端 / 发布班次</span><h1 class="section-title">创建新的出行班次</h1><form class="form-grid section-block" data-driver-publish-form><div class="field span-6"><label for="start">出发地</label><input id="start" name="start" value="杭州"></div><div class="field span-6"><label for="end">目的地</label><input id="end" name="end" value="苏州"></div><div class="field span-6"><label for="depart">出发时间</label><input id="depart" name="depart" type="datetime-local" value="2026-05-20T09:00"></div><div class="field span-6"><label for="arrival">到达时间</label><input id="arrival" name="arrival" type="datetime-local" value="2026-05-20T10:30"></div><div class="field span-4"><label for="seats">座位数</label><input id="seats" name="seats" type="number" value="36"></div><div class="field span-4"><label for="price">票价（分）</label><input id="price" name="price" type="number" value="8900"></div><div class="field span-4"><label for="vehicleType">车型</label><select id="vehicleType" name="vehicleType"><option value="car">汽车</option><option value="rail">高铁</option></select></div><div class="field span-12"><label for="stops">停靠点</label><textarea id="stops" name="stops">嘉兴,上海虹桥</textarea></div><button class="button button-primary span-12" type="submit" data-submit-trip>发布班次</button></form></section><aside class="panel"><h2 class="subhead">AI 填充</h2><p class="muted">可从司机 AI 页面生成草稿后自动回填。</p><button class="button button-secondary" type="button" data-fill-trip>填充示例</button><div class="toast" data-toast>已处理</div></aside></div></main>
'@ }

$pages["driver-trips.html"] = @{ Title = "我的班次"; Page = "driver"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">司机端 / 班次</span><h1 class="section-title">我的班次</h1></div></section><section class="table-card"><div class="table-scroll"><table><thead><tr><th>路线</th><th>出发</th><th>余票</th><th>状态</th><th>操作</th></tr></thead><tbody data-driver-trip-list><tr><td colspan="5">正在加载班次</td></tr></tbody></table></div></section></div></main>
'@ }

$pages["driver-trip-detail.html"] = @{ Title = "司机班次详情"; Page = "driver"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">班次详情</span><h1 class="section-title" data-driver-trip-title>正在加载班次</h1><div class="chip-row" data-driver-trip-chips></div></div></section><section class="dashboard-grid"><article class="panel"><h2 class="subhead">班次信息</h2><div class="list-stack section-block"><div class="info-card" data-driver-trip-time></div><div class="info-card" data-driver-trip-stops></div><div class="info-card" data-driver-trip-seat-info></div></div><button class="button button-secondary section-block" type="button" data-driver-open-verification>刷新核验</button></article><aside class="panel"><h2 class="subhead">订单核验</h2><div class="chip-row" data-driver-order-summary></div><div class="list-stack section-block" data-driver-trip-orders></div></aside></section></div></main>
'@ }

$pages["driver-income.html"] = @{ Title = "司机收益"; Page = "driver"; Content = @'
  <main class="page"><div class="container"><section class="stat-grid"><div class="stat-card"><div class="stat-label">今日收入</div><div class="stat-value" data-driver-income-today>--</div></div><div class="stat-card"><div class="stat-label">待结算</div><div class="stat-value" data-driver-income-pending-settle>--</div></div><div class="stat-card"><div class="stat-label">本周收入</div><div class="stat-value" data-driver-income-week>--</div></div><div class="stat-card"><div class="stat-label">客单价</div><div class="stat-value" data-driver-income-avg-order>--</div></div></section><section class="dashboard-grid section-block"><article class="panel"><h2 class="subhead">热门路线</h2><div class="list-stack section-block" data-driver-income-route-list></div></article><aside class="panel"><h2 class="subhead">经营建议</h2><div class="mini-chip">退款率 <span data-driver-income-refund-rate>--</span></div><div class="list-stack section-block" data-driver-income-suggestion-list></div></aside></section></div></main>
'@ }

$pages["driver-ai.html"] = @{ Title = "司机 AI"; Page = "driver"; Content = @'
  <main class="page"><div class="container dashboard-grid"><section class="panel"><span class="eyebrow">司机 AI</span><h1 class="section-title">生成发车方案</h1><form class="form-grid section-block" data-driver-ai-form><div class="field span-12"><label for="driverAiPrompt">运营需求</label><textarea id="driverAiPrompt" name="prompt">明天上午杭州到苏州，帮我生成一班高上座率车次。</textarea></div><button class="button button-primary span-12" type="submit">生成方案</button></form></section><aside class="panel"><h2 class="subhead">生成结果</h2><div class="list-stack section-block" data-driver-ai-result><div class="info-card">等待生成</div></div><div class="list-stack section-block" data-driver-ai-validation></div><div class="button-row"><button class="button button-secondary" type="button" data-driver-ai-apply>应用到表单</button><button class="button button-primary" type="button" data-driver-ai-publish>发布班次</button></div></aside></div></main>
'@ }

$adminDashboard = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">管理端</span><h1 class="section-title">运营管理工作台</h1><div class="chip-row" data-admin-dashboard-summary><span class="mini-chip">加载中</span></div></div></section><section class="stat-grid"><div class="stat-card"><div class="stat-label">用户总数</div><div class="stat-value" data-admin-dashboard-total-users>--</div></div><div class="stat-card"><div class="stat-label">活跃用户</div><div class="stat-value" data-admin-dashboard-active-users>--</div></div><div class="stat-card"><div class="stat-label">待退款</div><div class="stat-value" data-admin-dashboard-pending-refund>--</div></div><div class="stat-card"><div class="stat-label">已退款</div><div class="stat-value" data-admin-dashboard-refunded>--</div></div></section><section class="dashboard-grid section-block"><article class="panel"><h2 class="subhead">角色分布</h2><div class="list-stack section-block"><div class="info-card">乘客 <strong data-admin-dashboard-passenger-count>--</strong></div><div class="info-card">司机 <strong data-admin-dashboard-driver-count>--</strong></div><div class="info-card">管理员 <strong data-admin-dashboard-admin-count>--</strong></div><div class="info-card">拒绝退款 <strong data-admin-dashboard-rejected>--</strong></div></div></article><aside class="panel"><h2 class="subhead">待办提醒</h2><p class="muted" data-admin-dashboard-pending-note>正在加载</p><div class="empty-state section-block" data-admin-dashboard-empty-state>暂无异常</div></aside></section></div></main>
'@
$pages["admin-dashboard.html"] = @{ Title = "管理工作台"; Page = "admin"; Content = $adminDashboard }

$pages["admin-users.html"] = @{ Title = "用户管理"; Page = "admin"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">管理端 / 用户</span><h1 class="section-title">用户管理</h1><form class="filter-row section-block" data-admin-user-filter><div class="field"><label for="keyword">关键词</label><input id="keyword" name="keyword" placeholder="手机号 / 昵称"></div><div class="field"><label for="role">角色</label><select id="role" name="role"><option value="">全部</option><option value="passenger">乘客</option><option value="driver">司机</option><option value="admin">管理员</option></select></div><button class="button button-secondary" type="button" data-admin-user-reset>重置</button></form><div class="chip-row" data-admin-user-summary></div></div></section><section class="table-card"><div class="table-scroll"><table><thead><tr><th>用户</th><th>手机号</th><th>角色</th><th>状态</th><th>实名</th><th>创建时间</th><th>操作</th></tr></thead><tbody data-admin-user-list><tr><td colspan="7">正在加载用户</td></tr></tbody></table></div><div class="empty-state section-block" data-admin-user-empty-state>暂无用户数据</div></section></div></main>
'@ }

$pages["admin-orders.html"] = @{ Title = "订单审核"; Page = "admin"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">管理端 / 订单</span><h1 class="section-title">订单与退款审核</h1><form class="filter-row section-block" data-admin-order-filter><div class="field"><label for="refundStatus">退款状态</label><select id="refundStatus" name="refundStatus"><option value="">全部</option><option value="pending">待审核</option><option value="approved">已通过</option><option value="rejected">已拒绝</option></select></div><div class="field"><label for="reviewNote">审核备注</label><input id="reviewNote" name="reviewNote" value="符合平台规则"></div></form><div class="chip-row" data-admin-order-summary></div></div></section><section class="panel"><div class="list-stack" data-admin-order-list><div class="info-card">正在加载订单</div></div></section></div></main>
'@ }

$pages["admin-tokens.html"] = @{ Title = "Token 配额"; Page = "admin"; Content = @'
  <main class="page"><div class="container dashboard-grid"><section class="panel"><span class="eyebrow">管理端 / Token</span><h1 class="section-title">Token 配额管理</h1><form class="filter-row section-block" data-admin-token-filter-form><div class="field"><label for="tokenKeyword">用户</label><input id="tokenKeyword" name="keyword" placeholder="搜索用户"></div></form><div class="list-stack section-block" data-admin-token-user-list><div class="info-card">正在加载用户配额</div></div></section><aside class="panel"><h2 class="subhead">使用明细</h2><div class="list-stack section-block" data-admin-token-detail-list></div></aside></div></main>
'@ }

$pages["admin-risk.html"] = @{ Title = "风控日志"; Page = "admin"; Content = @'
  <main class="page"><div class="container dashboard-grid"><section class="panel"><span class="eyebrow">管理端 / 风控</span><h1 class="section-title">风险事件监控</h1><form class="filter-row section-block" data-admin-risk-filter-form><div class="field"><label for="riskType">类型</label><input id="riskType" name="type" placeholder="登录 / 支付 / 退款"></div><div class="field"><label for="riskStatus">状态</label><select id="riskStatus" name="status"><option value="">全部</option><option value="open">待处理</option><option value="closed">已关闭</option></select></div></form><div class="chip-row section-block" data-admin-risk-summary></div></section><aside class="panel"><h2 class="subhead">日志列表</h2><div class="list-stack section-block" data-admin-risk-list><div class="info-card">正在加载风控日志</div></div></aside></div></main>
'@ }

$pages["admin-knowledge.html"] = @{ Title = "知识库"; Page = "admin"; Content = @'
  <main class="page"><div class="container dashboard-grid"><section class="panel"><span class="eyebrow">管理端 / 知识库</span><h1 class="section-title">文档管理</h1><form class="form-grid section-block" data-knowledge-upload-form><div class="field span-6"><label for="docTitle">标题</label><input id="docTitle" name="title" value="退改签规则"></div><div class="field span-6"><label for="docType">类型</label><input id="docType" name="type" value="policy"></div><div class="field span-12"><label for="docContent">内容</label><textarea id="docContent" name="content">乘客可在发车前申请退款。</textarea></div><button class="button button-primary span-12" type="submit">上传文档</button></form><form class="filter-row section-block" data-knowledge-filter-form><div class="field"><label for="knowledgeKeyword">筛选</label><input id="knowledgeKeyword" name="keyword"></div><button class="button button-secondary" type="button" data-knowledge-filter-reset>重置</button></form><div class="list-stack section-block" data-knowledge-document-list></div></section><aside class="panel"><h2 class="subhead">语义检索</h2><form class="form-grid section-block" data-knowledge-search-form><div class="field span-12"><label for="query">问题</label><textarea id="query" name="query">退款需要多久到账？</textarea></div><button class="button button-primary span-12" type="submit">检索</button></form><div class="list-stack section-block" data-knowledge-search-result></div><div class="list-stack section-block" data-knowledge-detail></div></aside></div></main>
'@ }

$pages["admin-models.html"] = @{ Title = "模型治理"; Page = "admin"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">管理端 / 模型</span><h1 class="section-title">模型治理</h1><p class="lede">查看对话模型、风控模型和票务推荐模型的健康状态。</p></div></section><section class="card-grid"><article class="panel"><span class="badge">Chat</span><h2 class="subhead">对话助手</h2><p class="muted">可用，延迟稳定。</p><button class="button button-secondary" type="button" data-toast="模型已刷新">刷新</button></article><article class="panel"><span class="badge">Risk</span><h2 class="subhead">风控分类</h2><p class="muted">监控退款、登录与支付异常。</p><button class="button button-secondary" type="button" data-toast="已查看详情">查看</button></article><article class="panel"><span class="badge">Recommend</span><h2 class="subhead">票务推荐</h2><p class="muted">结合价格、时长和余票排序。</p><button class="button button-secondary" type="button" data-toast="已切换版本">切换版本</button></article></section><div class="toast"></div></div></main>
'@ }

$pages["admin-mcp.html"] = @{ Title = "MCP 工具"; Page = "admin"; Content = @'
  <main class="page"><div class="container"><section class="page-hero"><div class="hero-copy"><span class="eyebrow">管理端 / MCP</span><h1 class="section-title">工具连接管理</h1><p class="lede">统一查看外部工具可用性、权限和最近调用记录。</p></div></section><section class="card-grid"><article class="panel"><span class="badge">地图服务</span><h2 class="subhead">路线规划</h2><p class="muted">用于计算预计到达时间。</p><button class="button button-secondary" type="button" data-toast="工具已启用">启用</button></article><article class="panel"><span class="badge">支付服务</span><h2 class="subhead">支付查询</h2><p class="muted">用于同步支付单状态。</p><button class="button button-secondary" type="button" data-toast="已测试连接">测试连接</button></article><article class="panel"><span class="badge">消息服务</span><h2 class="subhead">通知投递</h2><p class="muted">用于订单、退款和核验提醒。</p><button class="button button-secondary" type="button" data-toast="配置已保存">保存配置</button></article></section><div class="toast"></div></div></main>
'@ }

foreach ($entry in $pages.GetEnumerator()) {
  $pageSpec = $entry.Value
  Write-Utf8File (Join-Path $frontend $entry.Key) (Page $pageSpec.Title $pageSpec.Page $pageSpec.Content)
}
Write-Utf8File (Join-Path $assets "styles.css") $styles

Write-Host ('Regenerated {0} frontend pages and shared styles.' -f $pages.Count)
