(function bootstrapTripVerse(window, document) {
  "use strict";

  const STORAGE_KEY = "tripverse_auth";
  const defaultRuntime = {
    appName: "TripVerse",
    apiBaseUrl: "/api",
    timeout: 10000,
    useMock: true,
    debug: true,
  };

  const runtimeConfig = Object.freeze({
    ...defaultRuntime,
    ...(window.__TRIPVERSE_RUNTIME__ || {}),
  });

  const ROUTES = Object.freeze({
    home: "index.html",
    auth: Object.freeze({
      login: "login.html",
      register: "register.html",
    }),
    passenger: Object.freeze({
      search: "search.html",
      tripDetail: "trip-detail.html",
      checkout: "checkout.html",
      payment: "payment.html",
      orders: "orders.html",
      orderDetail: "order-detail.html",
      profile: "profile.html",
      aiAssistant: "ai-assistant.html",
    }),
    driver: Object.freeze({
      dashboard: "driver-dashboard.html",
      publish: "driver-publish.html",
      trips: "driver-trips.html",
      tripDetail: "driver-trip-detail.html",
      income: "driver-income.html",
      ai: "driver-ai.html",
    }),
    admin: Object.freeze({
      dashboard: "admin-dashboard.html",
      users: "admin-users.html",
      orders: "admin-orders.html",
      tokens: "admin-tokens.html",
      risk: "admin-risk.html",
      models: "admin-models.html",
      knowledge: "admin-knowledge.html",
      mcp: "admin-mcp.html",
    }),
  });

  const API_ENDPOINTS = Object.freeze({
    auth: Object.freeze({
      sendEmailCode: "/auth/email/send",
      register: "/auth/register",
      loginByPassword: "/auth/login/password",
      loginByCode: "/auth/login/code",
      logout: "/auth/logout",
      me: "/auth/me",
    }),
    user: Object.freeze({
      profile: "/users/profile",
      verifyRealName: "/users/verify-real-name",
      accountStatus: "/users/account/status",
    }),
    passenger: Object.freeze({
      searchTickets: "/tickets/search",
      ticketDetail: "/tickets/:ticketId",
      createOrder: "/orders",
      myOrders: "/orders/my",
      orderDetail: "/orders/:orderId",
      orderTicket: "/orders/:orderId/ticket",
      cancelOrder: "/orders/:orderId/cancel",
      refundOrder: "/orders/:orderId/refund",
      createPayment: "/payments/create",
      paymentStatus: "/payments/:paymentId/status",
      mockPaymentSuccess: "/payments/:paymentId/mock-success",
      aiChat: "/ai/chat",
    }),
    driver: Object.freeze({
      dashboard: "/driver/dashboard",
      trips: "/driver/trips",
      tripDetail: "/driver/trips/:tripId",
      income: "/driver/income",
      verifyOrder: "/driver/orders/:orderId/verify",
      verifyTicket: "/driver/tickets/verify",
      aiCreateTrip: "/ai/driver/create-trip",
    }),
    admin: Object.freeze({
      dashboard: "/admin/dashboard",
      users: "/admin/users",
      updateUser: "/admin/users/:userId",
      userSummary: "/admin/users/summary",
      orders: "/admin/orders",
      approveRefund: "/admin/orders/:orderId/refund/approve",
      rejectRefund: "/admin/orders/:orderId/refund/reject",
      tokens: "/admin/tokens",
      riskLogs: "/admin/risk/logs",
      riskLogDetail: "/admin/risk/logs/:eventId",
      knowledge: "/admin/knowledge",
      knowledgeUpload: "/admin/knowledge/upload",
      knowledgeSearch: "/admin/knowledge/search",
      knowledgeDetail: "/admin/knowledge/:documentId",
      knowledgeReindex: "/admin/knowledge/:documentId/reindex",
    }),
    notification: Object.freeze({
      my: "/notifications/my",
      unreadCount: "/notifications/unread-count",
      markRead: "/notifications/:notificationId/read",
      markAllRead: "/notifications/read-all",
    }),
  });

  const ROLE_LABELS = Object.freeze({
    passenger: "乘客",
    driver: "司机",
    admin: "管理员",
  });

  const NAV_CONFIG = Object.freeze({
    guest: [
      { key: "home", label: "首页", href: ROUTES.home },
      { key: "search", label: "购票", href: ROUTES.passenger.search },
      { key: "ai", label: "AI 助手", href: ROUTES.passenger.aiAssistant },
      { key: "login", label: "登录", href: ROUTES.auth.login },
      { key: "register", label: "注册", href: ROUTES.auth.register },
    ],
    passenger: [
      { key: "home", label: "首页", href: ROUTES.home },
      { key: "search", label: "购票", href: ROUTES.passenger.search },
      { key: "orders", label: "订单", href: ROUTES.passenger.orders },
      { key: "profile", label: "个人中心", href: ROUTES.passenger.profile },
      { key: "ai", label: "AI 助手", href: ROUTES.passenger.aiAssistant },
    ],
    driver: [
      { key: "home", label: "首页", href: ROUTES.home },
      { key: "search", label: "购票", href: ROUTES.passenger.search },
      { key: "driver", label: "司机端", href: ROUTES.driver.dashboard },
      { key: "profile", label: "个人中心", href: ROUTES.passenger.profile },
      { key: "ai", label: "AI 助手", href: ROUTES.passenger.aiAssistant },
    ],
    admin: [
      { key: "home", label: "首页", href: ROUTES.home },
      { key: "admin-dashboard", label: "工作台", href: ROUTES.admin.dashboard },
      { key: "admin-users", label: "用户管理", href: ROUTES.admin.users },
      { key: "admin-orders", label: "订单审核", href: ROUTES.admin.orders },
      { key: "admin-tokens", label: "Token 配额", href: ROUTES.admin.tokens },
      { key: "admin-risk", label: "风控日志", href: ROUTES.admin.risk },
      { key: "admin-knowledge", label: "知识库", href: ROUTES.admin.knowledge },
      { key: "admin-models", label: "模型治理", href: ROUTES.admin.models },
      { key: "admin-mcp", label: "MCP 工具", href: ROUTES.admin.mcp },
      { key: "profile", label: "个人中心", href: ROUTES.passenger.profile },
      { key: "ai", label: "AI 助手", href: ROUTES.passenger.aiAssistant },
    ],
  });

  const AUTH_REQUIRED_PAGES = new Set([
    ROUTES.passenger.orders,
    ROUTES.passenger.orderDetail,
    ROUTES.passenger.checkout,
    ROUTES.passenger.payment,
    ROUTES.passenger.profile,
    ROUTES.driver.dashboard,
    ROUTES.driver.publish,
    ROUTES.driver.trips,
    ROUTES.driver.tripDetail,
    ROUTES.driver.income,
    ROUTES.driver.ai,
    ROUTES.admin.dashboard,
    ROUTES.admin.users,
    ROUTES.admin.orders,
    ROUTES.admin.tokens,
    ROUTES.admin.risk,
    ROUTES.admin.models,
    ROUTES.admin.knowledge,
    ROUTES.admin.mcp,
  ]);

  const PASSENGER_ONLY_PAGES = new Set([
    ROUTES.passenger.orders,
    ROUTES.passenger.orderDetail,
    ROUTES.passenger.checkout,
    ROUTES.passenger.payment,
  ]);

  const DRIVER_ONLY_PAGES = new Set([
    ROUTES.driver.dashboard,
    ROUTES.driver.publish,
    ROUTES.driver.trips,
    ROUTES.driver.tripDetail,
    ROUTES.driver.income,
    ROUTES.driver.ai,
  ]);

  const ADMIN_ONLY_PAGES = new Set([
    ROUTES.admin.dashboard,
    ROUTES.admin.users,
    ROUTES.admin.orders,
    ROUTES.admin.tokens,
    ROUTES.admin.risk,
    ROUTES.admin.models,
    ROUTES.admin.knowledge,
    ROUTES.admin.mcp,
  ]);

  function joinUrl(baseUrl, path) {
    const safeBase = String(baseUrl || "").replace(/\/+$/, "");
    const safePath = String(path || "").replace(/^\/+/, "");
    return safePath ? `${safeBase}/${safePath}` : safeBase;
  }

  function resolvePath(template, pathParams) {
    if (!template) {
      return "";
    }

    return Object.entries(pathParams || {}).reduce((current, [key, value]) => {
      return current.replace(`:${key}`, encodeURIComponent(String(value)));
    }, template);
  }

  function buildQuery(params) {
    const query = new URLSearchParams();
    Object.entries(params || {}).forEach(([key, value]) => {
      if (value === undefined || value === null || value === "") {
        return;
      }
      query.append(key, String(value));
    });

    const result = query.toString();
    return result ? `?${result}` : "";
  }

  async function request(path, options) {
    const settings = {
      method: "GET",
      data: undefined,
      body: undefined,
      query: undefined,
      headers: {},
      pathParams: {},
      ...options,
    };

    const resolvedPath = resolvePath(path, settings.pathParams);
    const url = joinUrl(runtimeConfig.apiBaseUrl, resolvedPath) + buildQuery(settings.query);

    if (runtimeConfig.useMock) {
      return Promise.resolve({
        success: true,
        mock: true,
        method: settings.method,
        url,
        data: settings.data || null,
      });
    }

    const auth = readAuth();
    const controller = new AbortController();
    const timer = window.setTimeout(() => controller.abort(), runtimeConfig.timeout);
    const isFormDataBody = settings.body instanceof window.FormData;
    const fetchOptions = {
      method: settings.method,
      headers: {
        ...(auth?.token ? { Authorization: `Bearer ${auth.token}` } : {}),
        ...(isFormDataBody ? {} : { "Content-Type": "application/json" }),
        ...settings.headers,
      },
      signal: controller.signal,
    };

    if (settings.body !== undefined && settings.method !== "GET") {
      fetchOptions.body = settings.body;
    } else if (settings.data !== undefined && settings.method !== "GET") {
      fetchOptions.body = JSON.stringify(settings.data);
    }

    try {
      const response = await window.fetch(url, fetchOptions);
      let payload = null;
      const responseType = response.headers.get("content-type") || "";
      if (responseType.includes("application/json")) {
        payload = await response.json();
      } else {
        const text = await response.text();
        payload = text ? { message: text } : null;
      }

      if (!response.ok) {
        const error = new Error(payload?.message || `HTTP ${response.status}`);
        error.status = response.status;
        error.payload = payload;
        throw error;
      }
      return payload;
    } finally {
      window.clearTimeout(timer);
    }
  }

  const api = Object.freeze({
    request,
    get(path, query, options) {
      return request(path, { ...options, method: "GET", query });
    },
    post(path, data, options) {
      return request(path, { ...options, method: "POST", data });
    },
    put(path, data, options) {
      return request(path, { ...options, method: "PUT", data });
    },
    patch(path, data, options) {
      return request(path, { ...options, method: "PATCH", data });
    },
    delete(path, options) {
      return request(path, { ...options, method: "DELETE" });
    },
  });

  const TripVerse = {
    config: runtimeConfig,
    routes: ROUTES,
    endpoints: API_ENDPOINTS,
    api,
    utils: {
      joinUrl,
      resolvePath,
      buildQuery,
    },
  };

  window.TripVerse = TripVerse;

  function readAuth() {
    try {
      const raw = window.localStorage.getItem(STORAGE_KEY);
      if (!raw) {
        return null;
      }
      return JSON.parse(raw);
    } catch (error) {
      return null;
    }
  }

  function saveAuth(authState) {
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(authState));
    TripVerse.auth = authState;
  }

  function clearAuth() {
    window.localStorage.removeItem(STORAGE_KEY);
    TripVerse.auth = null;
  }

  function escapeHtml(value) {
    return String(value ?? "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#39;");
  }

  function getQueryParam(key) {
    return new URLSearchParams(window.location.search).get(key) || "";
  }

  function buildTicketDetailHref(ticketId) {
    return `${ROUTES.passenger.tripDetail}?ticketId=${encodeURIComponent(ticketId || "")}`;
  }

  function buildCheckoutHref(ticketId) {
    return `${ROUTES.passenger.checkout}?ticketId=${encodeURIComponent(ticketId || "")}`;
  }

  function buildOrderDetailHref(orderId) {
    return `${ROUTES.passenger.orderDetail}?orderId=${encodeURIComponent(orderId || "")}`;
  }

  function renderStationLinks(stops, ticketId, hrefOverride, title) {
    const names = (Array.isArray(stops) ? stops : [])
      .map((stop) => {
        if (typeof stop === "string") {
          return stop;
        }
        return stop?.stopName || stop?.name || "";
      })
      .map((name) => String(name || "").trim())
      .filter(Boolean);

    if (!names.length || (!ticketId && !hrefOverride)) {
      return "";
    }

    const href = hrefOverride || buildCheckoutHref(ticketId);
    const linkTitle = title || "进入该班次下单";
    return `
      <div class="station-link-row">
        ${names.map((name) => `<a class="tag station-link" href="${href}" title="${escapeHtml(linkTitle)}">${escapeHtml(name)}</a>`).join("")}
      </div>
    `;
  }

  function saveLastTicketId(ticketId) {
    if (!ticketId) {
      return;
    }
    try {
      window.sessionStorage.setItem("tripverse_last_ticket_id", String(ticketId));
    } catch (error) {
      // Ignore storage failures.
    }
  }

  function readLastTicketId() {
    try {
      return window.sessionStorage.getItem("tripverse_last_ticket_id") || "";
    } catch (error) {
      return "";
    }
  }

  function formatShortDateTime(value) {
    if (!value) {
      return "--";
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return String(value);
    }

    const month = `${date.getMonth() + 1}`.padStart(2, "0");
    const day = `${date.getDate()}`.padStart(2, "0");
    const hours = `${date.getHours()}`.padStart(2, "0");
    const minutes = `${date.getMinutes()}`.padStart(2, "0");
    return `${month}-${day} ${hours}:${minutes}`;
  }

  function formatFullDateTime(value) {
    if (!value) {
      return "--";
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return String(value);
    }

    const year = date.getFullYear();
    const month = `${date.getMonth() + 1}`.padStart(2, "0");
    const day = `${date.getDate()}`.padStart(2, "0");
    const hours = `${date.getHours()}`.padStart(2, "0");
    const minutes = `${date.getMinutes()}`.padStart(2, "0");
    return `${year}-${month}-${day} ${hours}:${minutes}`;
  }

  function formatLocalInputDateTime(value) {
    if (!value) {
      return "";
    }

    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return String(value);
    }

    const year = date.getFullYear();
    const month = `${date.getMonth() + 1}`.padStart(2, "0");
    const day = `${date.getDate()}`.padStart(2, "0");
    const hours = `${date.getHours()}`.padStart(2, "0");
    const minutes = `${date.getMinutes()}`.padStart(2, "0");
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  }

  function formatMoneyFromCent(value) {
    const amount = Number(value || 0) / 100;
    if (Number.isInteger(amount)) {
      return `¥${amount}`;
    }
    return `¥${amount.toFixed(2).replace(/\.00$/, "").replace(/(\.\d)0$/, "$1")}`;
  }

  function formatPriceCent(value) {
    return formatMoneyFromCent(value);
  }

  function mapTripStatus(status) {
    switch (status) {
      case "published":
        return "售票中";
      case "closed":
        return "已关闭";
      case "draft":
        return "草稿";
      default:
        return status || "--";
    }
  }

  function mapOrderStatus(orderStatus, payStatus) {
    if (orderStatus === "pending_payment" || payStatus === "unpaid") {
      return "待支付";
    }
    if (orderStatus === "pending_verification") {
      return "待出发";
    }
    if (orderStatus === "completed") {
      return "已完成";
    }
    if (orderStatus === "cancelled") {
      return "已取消";
    }
    return orderStatus || "--";
  }

  function mapRefundStatus(refundStatus) {
    if (refundStatus === "requested") {
      return "退款申请中";
    }
    if (refundStatus === "refunded") {
      return "已退款";
    }
    if (refundStatus === "rejected") {
      return "退款已驳回";
    }
    return "无退款";
  }

  function mapSeatType(seatType) {
    if (seatType === "standard") {
      return "标准座";
    }
    return seatType || "--";
  }

  function toRfc3339FromLocal(value) {
    if (!value) {
      return "";
    }
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return "";
    }
    return date.toISOString();
  }

  function getCurrentFileName() {
    const path = window.location.pathname.split("/").pop();
    return path || ROUTES.home;
  }

  function getRoleHome(role) {
    if (role === "driver") {
      return ROUTES.driver.dashboard;
    }
    if (role === "admin") {
      return ROUTES.admin.dashboard;
    }
    return ROUTES.home;
  }

  function redirectTo(target) {
    if (!target) {
      return;
    }
    if (getCurrentFileName() === target) {
      return;
    }
    window.location.href = target;
  }

  function requireAccess() {
    const currentFile = getCurrentFileName();
    const auth = readAuth();
    TripVerse.auth = auth;

    if ((currentFile === ROUTES.auth.login || currentFile === ROUTES.auth.register) && auth) {
      redirectTo(getRoleHome(auth.role));
      return false;
    }

    if (!AUTH_REQUIRED_PAGES.has(currentFile)) {
      return true;
    }

    if (!auth) {
      redirectTo(ROUTES.auth.login);
      return false;
    }

    if (PASSENGER_ONLY_PAGES.has(currentFile) && auth.role !== "passenger") {
      redirectTo(getRoleHome(auth.role));
      return false;
    }

    if (DRIVER_ONLY_PAGES.has(currentFile) && auth.role !== "driver") {
      redirectTo(getRoleHome(auth.role));
      return false;
    }

    if (ADMIN_ONLY_PAGES.has(currentFile) && auth.role !== "admin") {
      redirectTo(getRoleHome(auth.role));
      return false;
    }

    return true;
  }

  function createToast() {
    let toast = document.querySelector(".toast");
    if (!toast) {
      toast = document.createElement("div");
      toast.className = "toast";
      document.body.appendChild(toast);
    }
    return toast;
  }

  const toast = createToast();

  function showToast(message) {
    toast.textContent = message;
    toast.classList.add("is-visible");
    window.clearTimeout(showToast.timer);
    showToast.timer = window.setTimeout(() => {
      toast.classList.remove("is-visible");
    }, 2200);
  }

  TripVerse.showToast = showToast;

  function friendlyErrorMessage(error, fallback) {
    if (!window.navigator.onLine) {
      return "网络已断开，请检查网络连接";
    }

    const status = Number(error?.status || 0);
    const rawMessage = String(error?.message || error?.payload?.message || "").trim();
    const lowerMessage = rawMessage.toLowerCase();

    if (lowerMessage.includes("invalid phone or password")) {
      return "用户不存在或密码错误";
    }
    if (lowerMessage.includes("invalid email or code")) {
      return "用户不存在或验证码错误";
    }
    if (lowerMessage.includes("user not found")) {
      return "用户不存在";
    }
    if (lowerMessage.includes("password")) {
      return "密码错误";
    }
    if (lowerMessage.includes("code")) {
      return "验证码错误";
    }
    if (status === 400 || status === 401 || status === 403) {
      return rawMessage || fallback || "账号或密码错误";
    }
    if (status >= 500) {
      return "服务器异常，请稍后再试";
    }
    if (error?.name === "AbortError" || lowerMessage.includes("failed to fetch")) {
      return "服务器连接失败，请稍后再试";
    }

    return rawMessage || fallback || "操作失败，请稍后再试";
  }

  function renderNavigation() {
    const auth = readAuth();
    const role = auth?.role || "guest";
    const items = NAV_CONFIG[role] || NAV_CONFIG.guest;
    const pageId = document.body.dataset.page;

    function renderItems(includeAuthMeta, isMobile) {
      const itemHtml = items
        .map((item) => {
          const activeClass = item.href === getCurrentFileName() || item.key === pageId ? "is-active" : "";
          const closeAttr = isMobile ? " data-nav-close" : "";
          return `<a class="${activeClass}" data-nav="${item.key}"${closeAttr} href="${item.href}">${item.label}</a>`;
        })
        .join("");

      if (!auth || !includeAuthMeta) {
        return itemHtml;
      }

      return `${itemHtml}<span class="nav-user-chip">${auth.nickname || auth.phone || "已登录"} / ${ROLE_LABELS[auth.role] || auth.role}</span><button class="nav-logout" type="button" data-logout>退出</button>`;
    }

    const desktopNav = document.querySelector(".nav-links");
    const mobileNav = document.querySelector(".mobile-links");

    if (desktopNav) {
      desktopNav.innerHTML = renderItems(true, false);
    }

    if (mobileNav) {
      mobileNav.innerHTML = renderItems(false, true);
      if (auth) {
        mobileNav.innerHTML += `<span class="nav-user-chip nav-user-chip-mobile">${auth.nickname || auth.phone || "已登录"} / ${ROLE_LABELS[auth.role] || auth.role}</span><button class="nav-logout nav-logout-mobile" type="button" data-logout>退出</button>`;
      }
    }
  }

  function initNavigation() {
    const body = document.body;

    document.querySelectorAll("[data-nav-toggle]").forEach((button) => {
      button.addEventListener("click", () => {
        body.classList.toggle("nav-open");
      });
    });

    document.querySelectorAll("[data-nav-close]").forEach((link) => {
      link.addEventListener("click", () => {
        body.classList.remove("nav-open");
      });
    });
  }

  function initLogoutActions() {
    document.querySelectorAll("[data-logout]").forEach((button) => {
      button.addEventListener("click", async () => {
        try {
          await api.post(API_ENDPOINTS.auth.logout, {});
        } catch (error) {
          // Ignore mock or network failures during local demo logout.
        }
        clearAuth();
        showToast("已退出登录");
        window.setTimeout(() => redirectTo(ROUTES.auth.login), 300);
      });
    });
  }

  function initToastTriggers() {
    document.querySelectorAll("[data-toast]").forEach((button) => {
      button.addEventListener("click", () => {
        showToast(button.dataset.toast || "已处理");
      });
    });
  }

  function initSeatPicker() {
    const seatButtons = document.querySelectorAll("[data-seat]");
    if (!seatButtons.length) {
      return;
    }

    seatButtons.forEach((button) => {
      if (button.classList.contains("is-disabled")) {
        return;
      }

      button.addEventListener("click", () => {
        seatButtons.forEach((seat) => seat.classList.remove("is-selected"));
        button.classList.add("is-selected");
        const target = document.querySelector("[data-seat-output]");
        if (target) {
          target.textContent = button.dataset.seat || "";
        }
      });
    });
  }

  function initCheckoutSummary() {
    const quantityInput = document.querySelector("[data-ticket-count]");
    const seatTypeInput = document.querySelector("[data-seat-type]");
    const totalOutput = document.querySelector("[data-total-output]");
    const basePriceNode = document.querySelector("[data-base-price]");
    const feeNode = document.querySelector("[data-service-fee]");

    function updateCheckoutTotal() {
      if (!quantityInput || !seatTypeInput || !totalOutput || !basePriceNode) {
        return;
      }

      const quantity = Number(quantityInput.value || 1);
      const currentOption = seatTypeInput.selectedOptions[0];
      const multiplier = Number(currentOption?.dataset.multiplier || 1);
      const basePrice = Number(basePriceNode.dataset.basePrice || 0);
      const serviceFee = Number(feeNode?.dataset.fee || 0);
      const total = quantity * Math.round(basePrice * multiplier) + serviceFee;

      totalOutput.textContent = formatPriceCent(total);
    }

    [quantityInput, seatTypeInput].forEach((node) => {
      if (!node) {
        return;
      }

      node.addEventListener("input", updateCheckoutTotal);
      node.addEventListener("change", updateCheckoutTotal);
    });

    updateCheckoutTotal();
  }

  function deprecatedLegacyPassengerAi() {
    const aiForm = document.querySelector("[data-ai-form]");
    if (!aiForm) {
      return;
    }

    const aiInput = aiForm.querySelector("textarea");
    const aiChat = document.querySelector("[data-ai-chat]");

    aiForm.addEventListener("submit", (event) => {
      event.preventDefault();
      if (!aiInput || !aiChat) {
        return;
      }

      const value = aiInput.value.trim();
      if (!value) {
        return;
      }

      const userMessage = document.createElement("div");
      userMessage.className = "message user";
      userMessage.innerHTML = `<strong>你</strong><div>${escapeHtml(value)}</div>`;
      aiChat.appendChild(userMessage);

      const aiMessage = document.createElement("div");
      aiMessage.className = "message ai";
      aiMessage.innerHTML =
        "<strong>AI 助手</strong><div>我已根据你的条件生成建议：优先推荐早班高铁，其次是低价大巴，并已提醒你注意换乘时间和退改规则。</div>";
      aiChat.appendChild(aiMessage);

      aiInput.value = "";
      aiChat.scrollTop = aiChat.scrollHeight;
    });
  }

  function deprecatedLegacyDriverDraftGenerator() {
    const fillTripButton = document.querySelector("[data-fill-trip]");
    if (!fillTripButton) {
      return;
    }

    fillTripButton.addEventListener("click", () => {
      redirectTo(ROUTES.driver.ai);
      return;
    });
  }

  // Legacy phone-code auth flow removed. Keep a no-op stub to avoid accidental reuse.
  function deprecatedLegacyAuthForms() {
    return;
  }

  function renderProfileAccountModule() {
    return;
  }

  function initAuthForms() {
    const authForm = document.querySelector("[data-auth-form]");
    if (!authForm) {
      return;
    }

    const modeSelect = authForm.querySelector("[name='loginMode']");
    const passwordGroup = authForm.querySelector("[data-password-group]");
    const codeGroups = Array.from(authForm.querySelectorAll("[data-code-group]"));
    const phoneGroups = Array.from(authForm.querySelectorAll("[data-phone-group]"));
    const emailGroups = Array.from(authForm.querySelectorAll("[data-email-group]"));

    function syncLoginMode() {
      const isRegister = authForm.dataset.authForm === "register";
      const mode = String(modeSelect?.value || "password");

      if (passwordGroup) {
        passwordGroup.style.display = mode === "password" || isRegister ? "grid" : "none";
      }

      codeGroups.forEach((node) => {
        node.style.display = mode === "code" || isRegister ? "grid" : "none";
      });

      phoneGroups.forEach((node) => {
        node.style.display = mode === "password" || isRegister ? "grid" : "none";
      });

      emailGroups.forEach((node) => {
        node.style.display = mode === "code" || isRegister ? "grid" : "none";
      });
    }

    if (modeSelect) {
      modeSelect.addEventListener("change", syncLoginMode);
    }
    syncLoginMode();

    authForm.addEventListener("submit", async (event) => {
      event.preventDefault();

      const formData = new window.FormData(authForm);
      const formType = authForm.dataset.authForm;
      const role = String(formData.get("role") || "passenger");
      const phone = String(formData.get("phone") || "").trim();
      const email = String(formData.get("email") || "").trim();
      const nickname = String(formData.get("nickname") || "").trim();
      const password = String(formData.get("password") || "").trim();
      const emailCode = String(formData.get("emailCode") || "").trim();
      const loginMode = String(formData.get("loginMode") || "password");

      if (formType === "register" && !phone) {
        showToast("请输入手机号");
        return;
      }

      if ((formType === "register" || loginMode === "code") && !email) {
        showToast("请输入邮箱");
        return;
      }

      if (loginMode === "password" && !phone) {
        showToast("请输入手机号");
        return;
      }

      if (loginMode === "password" && !password) {
        showToast("请填写密码");
        return;
      }

      if (loginMode === "code" && !emailCode) {
        showToast("请填写邮箱验证码");
        return;
      }

      try {
        if (formType === "register") {
          if (!emailCode || !password) {
            showToast("请补全邮箱验证码和密码");
            return;
          }

          const registerResponse = await api.post("/auth/register", {
            role,
            phone,
            email,
            emailCode,
            password,
            nickname,
          });

          const registerData = registerResponse?.data || {};
          const registerUser = registerData.user || {};

          saveAuth(
            runtimeConfig.useMock
              ? {
                  token: `mock-token-${Date.now()}`,
                  role,
                  phone,
                  email,
                  nickname: nickname || (role === "driver" ? "新司机" : "新乘客"),
                  status: "active",
                }
              : {
                  token: registerData.token || "",
                  role: registerUser.role || role,
                  phone: registerUser.phone || phone,
                  email: registerUser.email || email,
                  nickname: registerUser.nickname || nickname || (role === "driver" ? "新司机" : "新乘客"),
                  status: registerUser.status || "active",
                }
          );

          showToast("注册成功，已自动登录");
          window.setTimeout(() => redirectTo(getRoleHome(role)), 300);
          return;
        }

        const loginResponse = await api.post(
          loginMode === "password" ? "/auth/login/password" : "/auth/login/code",
          loginMode === "password"
            ? { role, phone, password }
            : { role, email, emailCode }
        );

        const loginData = loginResponse?.data || {};
        const loginUser = loginData.user || {};

        saveAuth(
          runtimeConfig.useMock
            ? {
                token: `mock-token-${Date.now()}`,
                role,
                phone,
                email,
                nickname: role === "driver" ? "李师傅" : "张明",
                status: "active",
              }
            : {
                token: loginData.token || "",
                role: loginUser.role || role,
                phone: loginUser.phone || phone,
                email: loginUser.email || email,
                nickname: loginUser.nickname || (role === "driver" ? "李师傅" : "张明"),
                status: loginUser.status || "active",
              }
        );

        showToast("登录成功");
        window.setTimeout(() => redirectTo(getRoleHome(role)), 300);
      } catch (error) {
        showToast(friendlyErrorMessage(error, formType === "register" ? "注册失败" : "登录失败"));
      }
    });

    document.querySelectorAll("[data-send-code]").forEach((button) => {
      button.addEventListener("click", async () => {
        const emailInput = authForm.querySelector("[name='email']");
        const email = String(emailInput?.value || "").trim();
        if (!email) {
          showToast("请先输入邮箱");
          return;
        }

        try {
          await api.post(API_ENDPOINTS.auth.sendEmailCode, {
            email,
            scene: authForm.dataset.authForm === "register" ? "register" : "login",
          });

          showToast("验证码已发送到邮箱");
        } catch (error) {
          showToast(friendlyErrorMessage(error, "验证码发送失败"));
        }
      });
    });
  }

  function initDriverTripsPage() {
    if (getCurrentFileName() !== ROUTES.driver.trips) {
      return;
    }

    const tbody = document.querySelector("[data-driver-trip-list]");
    if (!tbody) {
      return;
    }

    api.get(API_ENDPOINTS.driver.trips)
      .then((result) => {
        const trips = result?.data || [];
        if (!trips.length) {
          tbody.innerHTML = '<tr><td colspan="6">No trips yet. Create your first trip.</td></tr>';
          return;
        }

        tbody.innerHTML = trips.map((trip) => {
          const sold = Math.max((trip.seatTotal || 0) - (trip.seatAvailable || 0), 0);
          return `
            <tr>
              <td>${escapeHtml(trip.startCity)} → ${escapeHtml(trip.endCity)}</td>
              <td>${escapeHtml(formatShortDateTime(trip.departureTime))}</td>
              <td>${sold} / ${trip.seatTotal || 0}</td>
              <td>${escapeHtml(formatPriceCent(trip.priceCent))}</td>
              <td>${escapeHtml(mapTripStatus(trip.status))}</td>
              <td><a href="${ROUTES.driver.tripDetail}?tripId=${trip.id}">查看详情</a></td>
            </tr>
          `;
        }).join("");
      })
      .catch((error) => {
        tbody.innerHTML = `<tr><td colspan="6">${escapeHtml(error.message || "Failed to load trips")}</td></tr>`;
      });
  }

  function initDriverPublishPage() {
    if (getCurrentFileName() !== ROUTES.driver.publish) {
      return;
    }

    const form = document.querySelector("[data-driver-publish-form]");
    const submitButton = document.querySelector("[data-submit-trip]");
    if (!form || !submitButton) {
      return;
    }

    const departInput = form.querySelector("[name='depart']");
    if (departInput && !form.querySelector("[name='arrival']")) {
      const wrapper = document.createElement("div");
      wrapper.className = "field span-6";
      wrapper.innerHTML = `
        <label>到达时间</label>
        <input name="arrival" type="datetime-local">
      `;
      departInput.closest(".field")?.insertAdjacentElement("afterend", wrapper);
    }

    submitButton.addEventListener("click", async () => {
      const departureTime = form.querySelector("[name='depart']")?.value || "";
      const arrivalValue = form.querySelector("[name='arrival']")?.value || "";
      const startCity = form.querySelector("[name='start']")?.value || "";
      const endCity = form.querySelector("[name='end']")?.value || "";
      const seatTotal = Number(form.querySelector("[name='seats']")?.value || 0);
      const priceCent = Number(form.querySelector("[name='price']")?.value || 0);
      const vehicleType = form.querySelector("[name='vehicleType']")?.value || "car";
      const stopsRaw = form.querySelector("[name='stops']")?.value || "";

      const stops = stopsRaw
        .split(/[，。、,]/)
        .map((item) => item.trim())
        .filter(Boolean)
        .map((stopName, index) => ({
          stopOrder: index + 1,
          stopName,
        }));

      const fallbackArrival = departureTime
        ? new Date(new Date(departureTime).getTime() + 2 * 60 * 60 * 1000).toISOString()
        : "";

      try {
        const result = await api.post(API_ENDPOINTS.driver.trips, {
          vehicleType,
          startCity,
          endCity,
          departureTime: toRfc3339FromLocal(departureTime),
          arrivalTime: toRfc3339FromLocal(arrivalValue) || fallbackArrival,
          seatTotal,
          priceCent,
          stops,
        });

        const trip = result?.data;
        showToast("Trip created");
        if (trip?.id) {
          window.setTimeout(() => redirectTo(`${ROUTES.driver.tripDetail}?tripId=${trip.id}`), 300);
        }
      } catch (error) {
        showToast(error.message || "Create trip failed");
      }
    });
  }

  function hardenDriverPublishPage() {
    if (getCurrentFileName() !== ROUTES.driver.publish) {
      return;
    }

    const form = document.querySelector("[data-driver-publish-form]");
    const oldButton = document.querySelector("[data-submit-trip]");
    if (!form || !oldButton) {
      return;
    }

    const departInput = form.querySelector("[name='depart']");
    if (departInput && !form.querySelector("[name='arrival']")) {
      const wrapper = document.createElement("div");
      wrapper.className = "field span-6";
      wrapper.innerHTML = `
        <label>到达时间</label>
        <input name="arrival" type="datetime-local">
      `;
      departInput.closest(".field")?.insertAdjacentElement("afterend", wrapper);
    }

    const submitButton = oldButton.cloneNode(true);
    oldButton.replaceWith(submitButton);

    submitButton.addEventListener("click", async () => {
      const departureTime = form.querySelector("[name='depart']")?.value || "";
      const arrivalValue = form.querySelector("[name='arrival']")?.value || "";
      const startCity = String(form.querySelector("[name='start']")?.value || "").trim();
      const endCity = String(form.querySelector("[name='end']")?.value || "").trim();
      const seatTotal = Number(form.querySelector("[name='seats']")?.value || 0);
      const priceCent = Number(form.querySelector("[name='price']")?.value || 0);
      const vehicleType = form.querySelector("[name='vehicleType']")?.value || "car";
      const stopsRaw = String(form.querySelector("[name='stops']")?.value || "").trim();

      if (!departureTime) {
        showToast("请选择出发时间");
        return;
      }
      if (!startCity || !endCity) {
        showToast("请～写起点和终点");
        return;
      }
      if (startCity === endCity) {
        showToast("起点和终点不能相同");
        return;
      }
      if (seatTotal <= 0) {
        showToast("总座位数必须大于 0");
        return;
      }
      if (priceCent <= 0) {
        showToast("票价必须大于 0");
        return;
      }

      const stops = stopsRaw
        ? stopsRaw
            .split(/[，。、,]/)
            .map((item) => item.trim())
            .filter(Boolean)
            .map((stopName, index) => ({
              stopOrder: index + 1,
              stopName,
            }))
        : [];

      const fallbackArrival = new Date(new Date(departureTime).getTime() + 2 * 60 * 60 * 1000).toISOString();

      try {
        submitButton.disabled = true;
        submitButton.textContent = "提交?..";

        const result = await api.post(API_ENDPOINTS.driver.trips, {
          vehicleType,
          startCity,
          endCity,
          departureTime: toRfc3339FromLocal(departureTime),
          arrivalTime: toRfc3339FromLocal(arrivalValue) || fallbackArrival,
          seatTotal,
          priceCent,
          stops,
        });

        const trip = result?.data;
        showToast("班次创建成功");
        window.setTimeout(() => {
          if (trip?.id) {
            redirectTo(`${ROUTES.driver.tripDetail}?tripId=${trip.id}`);
          } else {
            redirectTo(ROUTES.driver.trips);
          }
        }, 300);
      } catch (error) {
        showToast(error.message || "创建班次失败");
      } finally {
        submitButton.disabled = false;
        submitButton.textContent = "提交发布";
      }
    });
  }

  function hardenDriverPublishPageV2() {
    if (getCurrentFileName() !== ROUTES.driver.publish) {
      return;
    }

    const form = document.querySelector("[data-driver-publish-form]");
    const oldButton = document.querySelector("[data-submit-trip]");
    if (!form || !oldButton) {
      return;
    }

    const departInput = form.querySelector("[name='depart']");
    if (departInput && !form.querySelector("[name='arrival']")) {
      const wrapper = document.createElement("div");
      wrapper.className = "field span-6";
      wrapper.innerHTML = `
        <label>到达时间</label>
        <input name="arrival" type="datetime-local">
      `;
      departInput.closest(".field")?.insertAdjacentElement("afterend", wrapper);
    }

    const submitButton = oldButton.cloneNode(true);
    oldButton.replaceWith(submitButton);

    submitButton.addEventListener("click", async () => {
      const departureTime = form.querySelector("[name='depart']")?.value || "";
      const arrivalValue = form.querySelector("[name='arrival']")?.value || "";
      const startCity = String(form.querySelector("[name='start']")?.value || "").trim();
      const endCity = String(form.querySelector("[name='end']")?.value || "").trim();
      const seatTotal = Number(form.querySelector("[name='seats']")?.value || 0);
      const priceCent = Number(form.querySelector("[name='price']")?.value || 0);
      const vehicleType = form.querySelector("[name='vehicleType']")?.value || "car";
      const stopsRaw = String(form.querySelector("[name='stops']")?.value || "").trim();

      if (!departureTime) {
        showToast("请选择出发时间");
        return;
      }
      if (!startCity || !endCity) {
        showToast("请～写起点和终点");
        return;
      }
      if (startCity === endCity) {
        showToast("起点和终点不能相同");
        return;
      }
      if (seatTotal <= 0) {
        showToast("总座位数必须大于 0");
        return;
      }
      if (priceCent <= 0) {
        showToast("票价必须大于 0");
        return;
      }

      const departureDate = new Date(departureTime);
      const arrivalDate = arrivalValue
        ? new Date(arrivalValue)
        : new Date(new Date(departureTime).getTime() + 2 * 60 * 60 * 1000 + 15 * 60 * 1000);
      if (Number.isNaN(departureDate.getTime())) {
        showToast("出发时间格式不正?");
        return;
      }
      if (Number.isNaN(arrivalDate.getTime())) {
        showToast("到达时间格式不正?");
        return;
      }
      if (arrivalDate.getTime() <= departureDate.getTime()) {
        showToast("到达时间必须晚于出发时间");
        return;
      }

      const stops = stopsRaw
        ? splitDriverStops(stopsRaw).map((stopName, index) => ({
            stopOrder: index + 1,
            stopName,
          }))
        : [];

      const duplicateEndpointStop = stops.find((item) => item.stopName === startCity || item.stopName === endCity);
      if (duplicateEndpointStop) {
        showToast("途经站点不能和起点或终点重复");
        return;
      }

      const stopNameSet = new Set();
      for (const stop of stops) {
        if (stopNameSet.has(stop.stopName)) {
          showToast("途经站点不能重复");
          return;
        }
        stopNameSet.add(stop.stopName);
      }

      try {
        submitButton.disabled = true;
        submitButton.textContent = "提交?..";

        const result = await api.post(API_ENDPOINTS.driver.trips, {
          vehicleType,
          startCity,
          endCity,
          departureTime: toRfc3339FromLocal(departureTime),
          arrivalTime: toRfc3339FromLocal(arrivalValue) || arrivalDate.toISOString(),
          seatTotal,
          priceCent,
          stops,
        });

        const trip = result?.data;
        showToast("班次创建成功");
        window.setTimeout(() => {
          if (trip?.id) {
            redirectTo(`${ROUTES.driver.tripDetail}?tripId=${trip.id}`);
          } else {
            redirectTo(ROUTES.driver.trips);
          }
        }, 300);
      } catch (error) {
        showToast(error.message || "创建班次失败");
      } finally {
        submitButton.disabled = false;
        submitButton.textContent = "提交发布";
      }
    });
  }

  function initDriverTripDetailPage() {
    if (getCurrentFileName() !== ROUTES.driver.tripDetail) {
      return;
    }

    const tripId = getQueryParam("tripId");
    if (!tripId) {
      return;
    }

    const title = document.querySelector(".page-hero .section-title");
    const chips = document.querySelector("[data-driver-trip-chips]") || document.querySelector(".page-hero .chip-row");
    const detailMetaBoxes = Array.from(document.querySelectorAll(".panel .list-item .list-meta"));
    const timeBox = document.querySelector("[data-driver-trip-time]") || detailMetaBoxes[0];
    const stopsBox = document.querySelector("[data-driver-trip-stops]") || detailMetaBoxes[1];
    const seatInfoBox = document.querySelector("[data-driver-trip-seat-info]") || detailMetaBoxes[2];

    api.get(API_ENDPOINTS.driver.tripDetail, undefined, { pathParams: { tripId } })
      .then((result) => {
        const trip = result?.data;
        if (!trip) {
          throw new Error("Trip not found");
        }

        if (title) {
          title.textContent = `${trip.startCity} → ${trip.endCity}`;
        }

        if (chips) {
          const sold = Math.max((trip.seatTotal || 0) - (trip.seatAvailable || 0), 0);
          const occupancy = trip.seatTotal ? Math.round((sold / trip.seatTotal) * 100) : 0;
          chips.innerHTML = `
            <span class="mini-chip">${escapeHtml(mapTripStatus(trip.status))}</span>
            <span class="mini-chip">上座?${occupancy}%</span>
          `;
        }

        if (timeBox) {
          timeBox.innerHTML = `
            <span>${escapeHtml(formatFullDateTime(trip.departureTime))}</span>
            <span>预计到达 ${escapeHtml(formatFullDateTime(trip.arrivalTime))}</span>
          `;
        }

        if (stopsBox) {
          const stops = trip.stops || [];
          stopsBox.innerHTML = stops.length
            ? stops.map((stop) => `<span>${escapeHtml(stop.stopName)}</span>`).join("")
            : "<span>No stops</span>";
        }

        if (seatInfoBox) {
          const sold = Math.max((trip.seatTotal || 0) - (trip.seatAvailable || 0), 0);
          seatInfoBox.innerHTML = `
            <span>${escapeHtml(formatPriceCent(trip.priceCent))}</span>
            <span>${trip.seatTotal || 0} 座 / 已售 ${sold}</span>
          `;
        }
      })
      .catch((error) => {
        if (title) {
          title.textContent = error.message || "Failed to load trip";
        }
      });
  }

    function initTicketSearchPage() {
    if (getCurrentFileName() !== ROUTES.passenger.search) {
      return;
    }

    const form = document.querySelector("[data-ticket-search-form]");
    const resultsBox = document.querySelector("[data-ticket-search-results]");
    const countBox = document.querySelector("[data-ticket-search-count]");
    if (!form || !resultsBox) {
      return;
    }

    const runSearch = async () => {
      const startCity = form.querySelector("[name='startCity']")?.value || "";
      const endCity = form.querySelector("[name='endCity']")?.value || "";
      const date = form.querySelector("[name='date']")?.value || "";
      const allowTransfer = Boolean(form.querySelector("[name='allowTransfer']")?.checked);

      try {
        const result = await api.get(API_ENDPOINTS.passenger.searchTickets, {
          startCity,
          endCity,
          date,
          allowTransfer,
        });

        const trips = Array.isArray(result?.data) ? result.data : [];
        if (countBox) {
          countBox.textContent = `共 ${trips.length} 个结果`;
        }

        if (!trips.length) {
          resultsBox.innerHTML = '<div class="info-card"><strong>暂无班次</strong><p class="muted">可以尝试换日期，或者保留“允许一次中转”后再次搜索。</p></div>';
          return;
        }

        resultsBox.innerHTML = trips.map((trip) => {
          if (trip.kind === "suggestion") {
            const suggestions = Array.isArray(trip.suggestions) ? trip.suggestions : [];
            return `
              <div class="ticket-card">
                <div>
                  <strong>可尝试的中转方向</strong>
                  <p class="muted">当天没有精确可购的直达或中转组合，可以从这些中转城市重新搜索。</p>
                  <div class="list-stack section-block">
                    ${suggestions.length ? suggestions.map((suggestion, index) => `
                      <div class="info-card">
                        <strong>方案 ${index + 1}：${escapeHtml(suggestion.route || "--")}</strong>
                        <p class="muted">${escapeHtml(suggestion.reason || "")}</p>
                        <p class="muted">中转城市：${escapeHtml(suggestion.transferCity || "--")}，前段 ${Number(suggestion.firstLegCount || 0)} 班，后段 ${Number(suggestion.secondLegCount || 0)} 班</p>
                      </div>
                    `).join("") : `<div class="info-card"><p class="muted">暂无候选路线。</p></div>`}
                  </div>
                </div>
              </div>
            `;
          }

          const legs = Array.isArray(trip.legs) ? trip.legs : [];
          const transferBlock = trip.kind === "transfer" && legs.length
            ? `
              <div class="list-stack section-block">
                ${legs.map((leg, index) => `
                  <div class="info-card">
                    <strong>第 ${index + 1} 段：${escapeHtml(leg.startCity || "")} -> ${escapeHtml(leg.endCity || "")}</strong>
                    <p class="muted">${escapeHtml(formatShortDateTime(leg.departureTime))} - ${escapeHtml(formatShortDateTime(leg.arrivalTime))}</p>
                    <p class="muted">${escapeHtml(leg.vehicleType || "--")}，票价 ${escapeHtml(formatPriceCent(leg.priceCent || 0))}，余票 ${Number(leg.seatAvailable || 0)}</p>
                    <div class="button-row section-block">
                      <a class="button button-secondary" href="${buildTicketDetailHref(leg.tripId)}">查看第 ${index + 1} 段</a>
                      <a class="button button-primary" href="${buildCheckoutHref(leg.tripId)}">购买第 ${index + 1} 段</a>
                    </div>
                  </div>
                `).join("")}
                <div class="info-card">
                  <strong>中转信息</strong>
                  <p class="muted">经 ${escapeHtml(trip.transferCity || "--")} 中转，候车约 ${Number(trip.transferWaitMinute || 0)} 分钟。</p>
                </div>
              </div>
            `
            : "";

          const actionButtons = trip.kind === "transfer" && legs.length
            ? legs.map((leg, index) => `
                <a class="button button-secondary" href="${buildTicketDetailHref(leg.tripId)}">查看第 ${index + 1} 段</a>
                <a class="button button-primary" href="${buildCheckoutHref(leg.tripId)}">购买第 ${index + 1} 段</a>
              `).join("")
            : `
                <a class="button button-secondary" href="${buildTicketDetailHref(trip.id)}">查看详情</a>
                <a class="button button-primary" href="${buildCheckoutHref(trip.id)}">去下单</a>
              `;
          const stationLinks = trip.kind !== "transfer" ? renderStationLinks(trip.stops, trip.id) : "";

          return `
            <div class="ticket-card">
              <div class="ticket-top">
                <div>
                  <div class="route-line">
                    <span class="route-city">${escapeHtml(`${trip.startCity} ${formatShortDateTime(trip.departureTime)}`)}</span>
                    <span class="route-divider"></span>
                    <span class="route-city">${escapeHtml(`${trip.endCity} ${formatShortDateTime(trip.arrivalTime)}`)}</span>
                  </div>
                  <div class="list-meta">
                    <span>${trip.kind === "transfer" ? `一次中转 ${escapeHtml(trip.transferCity || "--")}` : "直达"}</span>
                    <span>${escapeHtml(trip.vehicleType || "car")}</span>
                    <span>余票 ${trip.seatAvailable || 0}</span>
                    <span>${escapeHtml(mapTripStatus(trip.status))}</span>
                  </div>
                  ${stationLinks}
                  ${transferBlock}
                </div>
                <div class="price-pill">${escapeHtml(formatPriceCent(trip.priceCent))}</div>
              </div>
              <div class="button-row">
                <span class="tag">${trip.kind === "transfer" ? "支持一次中转" : `${trip.seatAvailable || 0} seats left`}</span>
                ${actionButtons}
              </div>
            </div>
          `;
        }).join("");
      } catch (error) {
        resultsBox.innerHTML = `<div class="info-card"><strong>搜索失败</strong><p class="muted">${escapeHtml(error.message || "Search failed")}</p></div>`;
      }
    };

    form.querySelectorAll("input, select").forEach((node) => {
      node.addEventListener("change", runSearch);
    });

    runSearch();
  }

  function syncDriverTripOrderStatus() {
    if (getCurrentFileName() !== ROUTES.driver.tripDetail) {
      return;
    }

    const tripId = getQueryParam("tripId");
    if (!tripId) {
      return;
    }

    const orderStatusBox = document.querySelectorAll(".split-grid > aside.panel .list-stack")[0];
    if (!orderStatusBox) {
      return;
    }

    api.get(API_ENDPOINTS.driver.tripDetail, undefined, { pathParams: { tripId } })
      .then((result) => {
        const trip = result?.data;
        const summary = trip?.orderSummary;
        if (!summary) {
          return;
        }

        orderStatusBox.innerHTML = `
          <div class="info-card">
            <strong>待核销 ${summary.pendingVerificationCount || 0} 单</strong>
            <p class="muted">${escapeHtml(summary.pendingVerificationNote || "暂无待核销说明。")}</p>
          </div>
          <div class="info-card">
            <strong>退款申请 ${summary.refundRequestCount || 0} 笔</strong>
            <p class="muted">${escapeHtml(summary.refundRequestNote || "暂无退款申请。")}</p>
          </div>
        `;
      })
      .catch(() => {
        // Keep the existing static content if syncing fails.
      });
  }

  function initTicketDetailPage() {
    if (getCurrentFileName() !== ROUTES.passenger.tripDetail) {
      return;
    }

    const ticketId = getQueryParam("ticketId");
    if (!ticketId) {
      return;
    }

    saveLastTicketId(ticketId);

    const title = document.querySelector("[data-ticket-detail-title]");
    const routeBox = document.querySelector("[data-ticket-detail-route]");
    const tagsBox = document.querySelector("[data-ticket-detail-tags]");
    const stopsBox = document.querySelector("[data-ticket-detail-stops]");
    const priceBox = document.querySelector("[data-ticket-detail-price]");
    const seatBox = document.querySelector("[data-ticket-detail-seat]");
    const statusBox = document.querySelector("[data-ticket-detail-status]");
    const checkoutLinks = Array.from(document.querySelectorAll("[data-ticket-detail-checkout-link], [data-ticket-detail-checkout-tab]"));

    checkoutLinks.forEach((link) => {
      link.href = `${ROUTES.passenger.checkout}?ticketId=${encodeURIComponent(ticketId)}`;
    });

    api.get(API_ENDPOINTS.passenger.ticketDetail, undefined, { pathParams: { ticketId } })
      .then((result) => {
        const trip = result?.data;
        if (!trip) {
          throw new Error("Ticket not found");
        }

        if (title) {
          title.textContent = `${trip.startCity} → ${trip.endCity}`;
        }

        if (routeBox) {
          routeBox.innerHTML = `
            <span class="route-city">${escapeHtml(formatFullDateTime(trip.departureTime))}</span>
            <span class="route-divider"></span>
            <span class="route-city">${escapeHtml(formatFullDateTime(trip.arrivalTime))}</span>
          `;
        }

        if (tagsBox) {
          tagsBox.innerHTML = `
            <span class="tag">${escapeHtml(formatPriceCent(trip.priceCent))}</span>
            <span class="tag">余票 ${trip.seatAvailable || 0}</span>
            <span class="tag">${escapeHtml(mapTripStatus(trip.status))}</span>
          `;
        }

        if (stopsBox) {
          const stops = trip.stops || [];
          stopsBox.innerHTML = stops.length
            ? stops.map((stop) => `
                <div class="timeline-item">
                  <span class="timeline-dot"></span>
                  <strong><a class="station-text-link" href="${buildCheckoutHref(ticketId)}">${escapeHtml(stop.stopName)}</a></strong>
                  <p class="muted">${escapeHtml(formatFullDateTime(stop.planArrivalTime || stop.planDepartureTime || trip.departureTime))}</p>
                </div>
              `).join("")
            : `
                <div class="timeline-item">
                  <span class="timeline-dot"></span>
                  <strong><a class="station-text-link" href="${buildCheckoutHref(ticketId)}">${escapeHtml(trip.startCity)}</a></strong>
                  <p class="muted">${escapeHtml(formatFullDateTime(trip.departureTime))}</p>
                </div>
                <div class="timeline-item">
                  <span class="timeline-dot"></span>
                  <strong><a class="station-text-link" href="${buildCheckoutHref(ticketId)}">${escapeHtml(trip.endCity)}</a></strong>
                  <p class="muted">${escapeHtml(formatFullDateTime(trip.arrivalTime))}</p>
                </div>
              `;
        }
      })
      .catch((error) => {
        if (title) {
          title.textContent = error.message || "Failed to load ticket";
        }
      });
  }

  function initCheckoutPage() {
    return;
  }

  function initOrdersPage() {
    return;
  }

  function initOrdersPage() {
    return;
  }

  function initOrderDetailPage() {
    return;
  }

  function initOrderDetailPage() {
    return;
  }

  function initProfilePage() {
    if (getCurrentFileName() !== ROUTES.passenger.profile) {
      return;
    }

    const saveButton = document.querySelector("[data-save-profile]");
    const nicknameInput = document.querySelector("[name='nickname']");
    const phoneInput = document.querySelector("[name='phone']");
    let emailInput = document.querySelector("[name='email']");
    const verifiedInput = document.querySelector("[name='realNameVerified']");

    if (!emailInput) {
      const notificationField = document.querySelector("[name='notificationText']")?.closest(".field");
      if (notificationField) {
        const emailField = document.createElement("div");
        emailField.className = "field span-12";
        emailField.innerHTML = `
          <label>Email</label>
          <input name="email" type="email" value="">
        `;
        notificationField.insertAdjacentElement("beforebegin", emailField);
        emailInput = emailField.querySelector("[name='email']");
      }
    }

    Promise.all([
      api.get(API_ENDPOINTS.user.profile),
      api.get(API_ENDPOINTS.user.accountStatus),
    ]).then(([profileResult, statusResult]) => {
      const user = profileResult?.data || {};
      const status = statusResult?.data || {};

      if (nicknameInput) nicknameInput.value = user.nickname || "";
      if (phoneInput) phoneInput.value = user.phone || "";
      if (emailInput) {
        if (!emailInput.name) {
          emailInput.name = "email";
        }
        emailInput.value = user.email || "";
      }
      if (verifiedInput) verifiedInput.value = status.realNameVerified ? "已实名" : "未实名";

      const auth = readAuth();
      if (auth) {
        saveAuth({
          ...auth,
          role: user.role || auth.role,
          phone: user.phone || auth.phone,
          email: user.email || auth.email,
          nickname: user.nickname || auth.nickname,
          status: user.status || auth.status,
        });
      }
    }).catch((error) => {
      showToast(error.message || "Load profile failed");
    });

    if (saveButton) {
      saveButton.addEventListener("click", async () => {
        try {
          const result = await api.put(API_ENDPOINTS.user.profile, {
            nickname: nicknameInput?.value || "",
            email: emailInput?.value || "",
          });

          const user = result?.data || {};
          const auth = readAuth();
          if (auth) {
            saveAuth({
              ...auth,
              phone: user.phone || auth.phone,
              email: user.email || auth.email,
              nickname: user.nickname || auth.nickname,
            });
          }

          showToast("Profile saved");
          renderNavigation();
          initLogoutActions();
        } catch (error) {
          showToast(error.message || "Save failed");
        }
      });
    }
  }

  function initPaymentPage() {
    return;
  }

  function initDriverTripDetailPage() {
    if (getCurrentFileName() !== ROUTES.driver.tripDetail) {
      return;
    }

    const tripId = getQueryParam("tripId");
    if (!tripId) {
      return;
    }

    const title = document.querySelector("[data-driver-trip-title]");
    const chips = document.querySelector("[data-driver-trip-chips]");
    const timeBox = document.querySelector("[data-driver-trip-time]");
    const stopsBox = document.querySelector("[data-driver-trip-stops]");
    const seatInfoBox = document.querySelector("[data-driver-trip-seat-info]");
    const summaryBox = document.querySelector("[data-driver-order-summary]");
    const ordersBox = document.querySelector("[data-driver-trip-orders]");
    const refreshButton = document.querySelector("[data-driver-open-verification]");

    const renderTripDetail = (trip) => {
      if (title) {
        title.textContent = `${trip.startCity} -> ${trip.endCity}`;
      }

      if (chips) {
        const sold = Math.max((trip.seatTotal || 0) - (trip.seatAvailable || 0), 0);
        const occupancy = trip.seatTotal ? Math.round((sold / trip.seatTotal) * 100) : 0;
        chips.innerHTML = `
          <span class="mini-chip">${escapeHtml(mapTripStatus(trip.status))}</span>
          <span class="mini-chip">上座?${occupancy}%</span>
        `;
      }

      if (timeBox) {
        timeBox.innerHTML = `
          <span>${escapeHtml(formatFullDateTime(trip.departureTime))}</span>
          <span>预计到达 ${escapeHtml(formatFullDateTime(trip.arrivalTime))}</span>
        `;
      }

      if (stopsBox) {
        const stops = trip.stops || [];
        stopsBox.innerHTML = stops.length
          ? stops.map((stop) => `<span>${escapeHtml(stop.stopName)}</span>`).join("")
          : `<span>${escapeHtml(trip.startCity)}</span><span>${escapeHtml(trip.endCity)}</span>`;
      }

      if (seatInfoBox) {
        const sold = Math.max((trip.seatTotal || 0) - (trip.seatAvailable || 0), 0);
        seatInfoBox.innerHTML = `
          <span>${escapeHtml(formatPriceCent(trip.priceCent))}</span>
          <span>${trip.seatTotal || 0} 座 / 已售 ${sold}</span>
        `;
      }
    };

    const renderSummary = (trip) => {
      const summary = trip?.orderSummary || {};
      if (!summaryBox) {
        return;
      }
      summaryBox.innerHTML = `
        <div class="info-card">
          <strong>待核销 ${summary.pendingVerificationCount || 0} 人</strong>
          <p class="muted">${escapeHtml(summary.pendingVerificationNote || "暂无待核销说明。")}</p>
        </div>
        <div class="info-card">
          <strong>退款申请 ${summary.refundRequestCount || 0} 笔</strong>
          <p class="muted">${escapeHtml(summary.refundRequestNote || "暂无退款申请。")}</p>
        </div>
      `;
    };

    const renderOrders = (trip, reload) => {
      if (!ordersBox) {
        return;
      }

      const orders = Array.isArray(trip?.orders) ? trip.orders : [];
      if (!orders.length) {
        ordersBox.innerHTML = `
          <div class="info-card">
            <strong>当前班次还没有订单</strong>
            <p class="muted">乘客完成下单和支付后，这里会出现真实订单。</p>
          </div>
        `;
        return;
      }

      ordersBox.innerHTML = orders.map((order) => {
        const canVerify = order.payStatus === "paid" && order.orderStatus === "pending_verification" && order.refundStatus === "none";
        return `
          <div class="info-card" data-driver-order-card="${order.id}">
            <div class="row-between">
              <strong>${escapeHtml(order.orderNo || `订单 #${order.id}`)}</strong>
              <span class="tag">${escapeHtml(mapOrderStatus(order.orderStatus, order.payStatus))}</span>
            </div>
            <div class="list-meta">
              <span>${escapeHtml(formatFullDateTime(order.createdAt))}</span>
              <span>${order.ticketCount || 0} 张</span>
              <span>${escapeHtml(mapSeatType(order.seatType))}</span>
              <span>${escapeHtml(formatMoneyFromCent(order.amount || 0))}</span>
            </div>
            <div class="list-meta">
              <span>支付状态：${escapeHtml(order.payStatus || "--")}</span>
              <span>退款状态：${escapeHtml(mapRefundStatus(order.refundStatus))}</span>
            </div>
            <div class="button-row section-block">
              <button class="button button-primary" type="button" data-driver-verify-order="${order.id}" ${canVerify ? "" : "disabled"}>
                ${canVerify ? "核销完成" : "不可核销"}
              </button>
            </div>
          </div>
        `;
      }).join("");

      ordersBox.querySelectorAll("[data-driver-verify-order]").forEach((button) => {
        button.addEventListener("click", async () => {
          const orderId = button.getAttribute("data-driver-verify-order");
          if (!orderId) {
            return;
          }

          const originalText = button.textContent;
          button.disabled = true;
          button.textContent = "核销?..";

          try {
            await api.post(API_ENDPOINTS.driver.verifyOrder, {}, { pathParams: { orderId } });
            showToast("订单已核销完成");
            await reload();
          } catch (error) {
            button.disabled = false;
            button.textContent = originalText;
            showToast(error.message || "核销失败");
          }
        });
      });
    };

    const loadTripDetail = async () => {
      const result = await api.get(API_ENDPOINTS.driver.tripDetail, undefined, { pathParams: { tripId } });
      const trip = result?.data;
      if (!trip) {
        throw new Error("班次不存在");
      }
      renderTripDetail(trip);
      renderSummary(trip);
      renderOrders(trip, loadTripDetail);
    };

    if (refreshButton) {
      refreshButton.addEventListener("click", () => {
        loadTripDetail()
          .then(() => {
            showToast("订单列表已刷新");
          })
          .catch((error) => {
            showToast(error.message || "刷新失败");
          });
      });
    }

    loadTripDetail().catch((error) => {
      if (title) {
        title.textContent = error.message || "班次加载失败";
      }
      if (summaryBox) {
        summaryBox.innerHTML = `
          <div class="info-card">
            <strong>读取失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试。")}</p>
          </div>
        `;
      }
      if (ordersBox) {
        ordersBox.innerHTML = `
          <div class="info-card">
            <strong>无法读取订单列表</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试。")}</p>
          </div>
        `;
      }
    });
  }

  function syncDriverTripOrderStatus() {
    if (getCurrentFileName() !== ROUTES.driver.tripDetail) {
      return;
    }
  }

  function initTicketDetailPage() {
    if (getCurrentFileName() !== ROUTES.passenger.tripDetail) {
      return;
    }

    const ticketId = getQueryParam("ticketId");
    if (!ticketId) {
      return;
    }

    const title = document.querySelector("[data-ticket-detail-title]");
    const routeBox = document.querySelector("[data-ticket-detail-route]");
    const tagsBox = document.querySelector("[data-ticket-detail-tags]");
    const stopsBox = document.querySelector("[data-ticket-detail-stops]");
    const priceBox = document.querySelector("[data-ticket-detail-price]");
    const seatBox = document.querySelector("[data-ticket-detail-seat]");
    const statusBox = document.querySelector("[data-ticket-detail-status]");
    const checkoutLinks = Array.from(document.querySelectorAll("[data-ticket-detail-checkout-link], [data-ticket-detail-checkout-tab]"));

    checkoutLinks.forEach((link) => {
      link.href = `${ROUTES.passenger.checkout}?ticketId=${encodeURIComponent(ticketId)}`;
    });

    api.get(API_ENDPOINTS.passenger.ticketDetail, undefined, { pathParams: { ticketId } })
      .then((result) => {
        const trip = result?.data;
        if (!trip) {
          throw new Error("班次不存?");
        }

        if (title) {
          title.textContent = `${trip.startCity} -> ${trip.endCity}`;
        }

        if (routeBox) {
          routeBox.innerHTML = `
            <span class="route-city">${escapeHtml(formatFullDateTime(trip.departureTime))}</span>
            <span class="route-divider"></span>
            <span class="route-city">${escapeHtml(formatFullDateTime(trip.arrivalTime))}</span>
          `;
        }

        if (tagsBox) {
          tagsBox.innerHTML = `
            <span class="tag">${escapeHtml(formatPriceCent(trip.priceCent))}</span>
            <span class="tag">余票 ${trip.seatAvailable || 0}</span>
            <span class="tag">${escapeHtml(mapTripStatus(trip.status))}</span>
          `;
        }

        if (priceBox) {
          priceBox.textContent = formatPriceCent(trip.priceCent);
        }

        if (seatBox) {
          seatBox.textContent = `${trip.seatAvailable || 0} / ${trip.seatTotal || 0}`;
        }

        if (statusBox) {
          statusBox.textContent = mapTripStatus(trip.status);
        }

        if (stopsBox) {
          const stops = trip.stops || [];
          stopsBox.innerHTML = stops.length
            ? stops.map((stop) => `
                <div class="timeline-item">
                  <span class="timeline-dot"></span>
                  <strong><a class="station-text-link" href="${buildCheckoutHref(ticketId)}">${escapeHtml(stop.stopName)}</a></strong>
                  <p class="muted">${escapeHtml(formatFullDateTime(stop.planArrivalTime || stop.planDepartureTime || trip.departureTime))}</p>
                </div>
              `).join("")
            : `
                <div class="timeline-item">
                  <span class="timeline-dot"></span>
                  <strong><a class="station-text-link" href="${buildCheckoutHref(ticketId)}">${escapeHtml(trip.startCity)}</a></strong>
                  <p class="muted">${escapeHtml(formatFullDateTime(trip.departureTime))}</p>
                </div>
                <div class="timeline-item">
                  <span class="timeline-dot"></span>
                  <strong><a class="station-text-link" href="${buildCheckoutHref(ticketId)}">${escapeHtml(trip.endCity)}</a></strong>
                  <p class="muted">${escapeHtml(formatFullDateTime(trip.arrivalTime))}</p>
                </div>
              `;
        }
      })
      .catch((error) => {
        if (title) {
          title.textContent = error.message || "班次加载失败";
        }
      });
  }

  function initCheckoutPage() {
    if (getCurrentFileName() !== ROUTES.passenger.checkout) {
      return;
    }

    const ticketId = getQueryParam("ticketId");
    const title = document.querySelector("[data-checkout-title]");
    const lede = document.querySelector("[data-checkout-lede]");
    const routeBox = document.querySelector("[data-checkout-route]");
    const basePriceNode = document.querySelector("[data-checkout-base-price]");
    const departureBox = document.querySelector("[data-checkout-departure]");
    const seatBox = document.querySelector("[data-checkout-seat-available]");
    const tripLink = document.querySelector("[data-checkout-trip-link]");
    const backLink = document.querySelector("[data-checkout-back-link]");
    const submitButton = document.querySelector("[data-submit-order]");
    const saveDraftButton = document.querySelector("[data-save-draft]");
    const passengerNameInput = document.querySelector("[name='passengerName']");
    const idCardInput = document.querySelector("[name='idCard']");
    const phoneInput = document.querySelector("[name='phone']");
    const seatTypeInput = document.querySelector("[name='seatType']");
    const ticketCountInput = document.querySelector("[name='ticketCount']");
    const totalOutput = document.querySelector("[data-total-output]");

    function formatCheckoutAmount(cent) {
      return formatPriceCent(cent);
    }

    function renderMissingTicketState() {
      if (title) {
        title.textContent = "缺少班次编号";
      }
      if (lede) {
        lede.textContent = "当前下单确认页没有收到 ticketId，请从班次详情页重新进入。";
      }
      if (routeBox) {
        routeBox.textContent = "无法创建订单";
      }
      if (departureBox) {
        departureBox.textContent = "请返回班次详情页";
      }
      if (seatBox) {
        seatBox.textContent = "--";
      }
      if (totalOutput) {
        totalOutput.textContent = formatPriceCent(0);
      }
      if (tripLink) {
        tripLink.href = ROUTES.passenger.search;
      }
      if (backLink) {
        backLink.href = ROUTES.passenger.search;
        backLink.textContent = "返回票务搜索";
      }
      if (submitButton) {
        submitButton.disabled = false;
        submitButton.textContent = "返回票务搜索";
        submitButton.addEventListener("click", () => {
          window.location.href = ROUTES.passenger.search;
        }, { once: true });
      }
      if (saveDraftButton) {
        saveDraftButton.textContent = "重新选班";
        saveDraftButton.addEventListener("click", () => {
          window.location.href = ROUTES.passenger.search;
        }, { once: true });
      }
    }

    if (!ticketId) {
      renderMissingTicketState();
      return;
    }

    const detailHref = `${ROUTES.passenger.tripDetail}?ticketId=${encodeURIComponent(ticketId)}`;
    if (tripLink) {
      tripLink.href = detailHref;
    }
    if (backLink) {
      backLink.href = detailHref;
    }

    if (saveDraftButton) {
      saveDraftButton.addEventListener("click", () => {
        showToast("草稿只保存在当前页面，刷新后会丢失");
      });
    }

    const auth = readAuth();
    if (passengerNameInput && !passengerNameInput.value && auth?.nickname) {
      passengerNameInput.value = auth.nickname;
    }
    if (phoneInput && !phoneInput.value && auth?.phone) {
      phoneInput.value = auth.phone;
    }

    let currentTrip = null;
    let submitting = false;

    api.get(API_ENDPOINTS.passenger.ticketDetail, undefined, { pathParams: { ticketId } })
      .then((result) => {
        const trip = result?.data;
        if (!trip) {
          throw new Error("班次不存在");
        }

        currentTrip = trip;

        if (title) {
          title.textContent = `${trip.startCity} -> ${trip.endCity}`;
        }
        if (lede) {
          lede.textContent = "确认数量后提交订单，系统会创建真实订单并自动跳转到支付页。";
        }
        if (routeBox) {
          routeBox.textContent = `${trip.startCity} -> ${trip.endCity}`;
        }
        if (basePriceNode) {
          basePriceNode.dataset.basePrice = String(Number(trip.priceCent || 0));
          basePriceNode.textContent = formatCheckoutAmount(trip.priceCent);
        }
        if (departureBox) {
          departureBox.textContent = formatFullDateTime(trip.departureTime);
        }
        if (seatBox) {
          seatBox.textContent = `${trip.seatAvailable || 0} / ${trip.seatTotal || 0}`;
        }

        const updateTotal = () => {
          if (!ticketCountInput || !totalOutput) {
            return;
          }
          const quantity = Math.max(1, Number(ticketCountInput.value || 1));
          totalOutput.textContent = formatCheckoutAmount(Number(trip.priceCent || 0) * quantity);
        };

        if (ticketCountInput) {
          ticketCountInput.addEventListener("input", updateTotal);
          ticketCountInput.addEventListener("change", updateTotal);
        }
        updateTotal();
      })
      .catch((error) => {
        if (title) {
          title.textContent = "班次加载失败";
        }
        if (lede) {
          lede.textContent = error.message || "无法读取班次信息";
        }
        if (submitButton) {
          submitButton.disabled = true;
        }
      });

    if (!submitButton) {
      return;
    }

    submitButton.addEventListener("click", async () => {
      if (submitting) {
        return;
      }

      const passengerName = String(passengerNameInput?.value || "").trim();
      const idCard = String(idCardInput?.value || "").trim();
      const phone = String(phoneInput?.value || "").trim();
      const seatType = String(seatTypeInput?.value || "standard").trim() || "standard";
      const ticketCount = Math.max(1, Number(ticketCountInput?.value || 1));

      if (!currentTrip) {
        showToast("班次信息还没加载完成");
        return;
      }
      if (!passengerName) {
        showToast("请输入乘车人姓名");
        return;
      }
      if (!idCard) {
        showToast("请输入身份证号");
        return;
      }
      if (!phone) {
        showToast("请输入手机号");
        return;
      }
      if (!Number.isInteger(ticketCount) || ticketCount <= 0) {
        showToast("购票数量必须大于 0");
        return;
      }

      submitting = true;
      const originalText = submitButton.textContent;
      submitButton.disabled = true;
      submitButton.textContent = "提交?..";

      try {
        const result = await api.post(API_ENDPOINTS.passenger.createOrder, {
          tripId: Number(ticketId),
          ticketCount,
          seatType,
        });
        const order = result?.data;
        if (!order?.id) {
          throw new Error("订单创建成功，但未返回订单号");
        }
        showToast("订单已创建，正在跳转支付页");
        window.location.href = `${ROUTES.passenger.payment}?orderId=${encodeURIComponent(order.id)}`;
      } catch (error) {
        showToast(error.message || "创建订单失败");
        submitButton.disabled = false;
        submitButton.textContent = originalText;
        submitting = false;
      }
    });
  }

  function mapRefundStatus(refundStatus) {
    if (refundStatus === "requested") {
      return "退款申请中";
    }
    if (refundStatus === "refunded") {
      return "已退款";
    }
    if (refundStatus === "rejected") {
      return "退款已驳回";
    }
    return "无退款";
  }

  function mapUserStatus(status) {
    if (status === "active") {
      return "正常";
    }
    if (status === "frozen") {
      return "冻结";
    }
    if (status === "disabled") {
      return "禁用";
    }
    return status || "--";
  }

  function formatYesNo(value) {
    return value ? "已实名" : "未实名";
  }

  function initAdminUsersPage() {
    if (getCurrentFileName() !== ROUTES.admin.users) {
      return;
    }

    const summaryBox = document.querySelector("[data-admin-user-summary]");
    const listBox = document.querySelector("[data-admin-user-list]");
    const emptyState = document.querySelector("[data-admin-user-empty-state]");
    if (!summaryBox || !listBox || !emptyState) {
      return;
    }

    const maskPhone = (value) => {
      const phone = String(value || "").trim();
      if (phone.length < 7) {
        return phone || "--";
      }
      return `${phone.slice(0, 3)}****${phone.slice(-4)}`;
    };

    const renderSummary = (summary) => {
      summaryBox.innerHTML = `
        <span class="mini-chip">总用户 ${Number(summary?.totalUsers || 0)}</span>
        <span class="mini-chip">乘客 ${Number(summary?.passengerCount || 0)}</span>
        <span class="mini-chip">司机 ${Number(summary?.driverCount || 0)}</span>
        <span class="mini-chip">管理员 ${Number(summary?.adminCount || 0)}</span>
        <span class="mini-chip">活跃 ${Number(summary?.activeCount || 0)}</span>
      `;
    };

    const renderList = (users) => {
      if (!users.length) {
        listBox.innerHTML = `<tr><td colspan="7">暂无用户数据</td></tr>`;
        emptyState.innerHTML = `
          <strong>暂无用户</strong>
          <p class="muted">后端当前还没有可展示的用户记录。</p>
        `;
        return;
      }

      listBox.innerHTML = users.map((user) => `
        <tr>
          <td>${escapeHtml(user.nickname || "--")}</td>
          <td>${escapeHtml(maskPhone(user.phone))}</td>
          <td>${escapeHtml(ROLE_LABELS[user.role] || user.role || "--")}</td>
          <td>${escapeHtml(mapUserStatus(user.status))}</td>
          <td>${escapeHtml(user.email || "--")}</td>
          <td>${escapeHtml(formatYesNo(user.realNameVerified))}</td>
          <td>${escapeHtml(formatFullDateTime(user.createdAt))}</td>
        </tr>
      `).join("");

      emptyState.innerHTML = `
        <strong>用户列表已同步</strong>
        <p class="muted">当前共加载 ${users.length} 条真实用户记录，支持继续扩展搜索、冻结和角色调整功能。</p>
      `;
    };

    Promise.all([
      api.get(API_ENDPOINTS.admin.userSummary),
      api.get(API_ENDPOINTS.admin.users),
    ])
      .then(([summaryResult, listResult]) => {
        renderSummary(summaryResult?.data || {});
        renderList(Array.isArray(listResult?.data) ? listResult.data : []);
      })
      .catch((error) => {
        summaryBox.innerHTML = `<span class="mini-chip">加载失败</span>`;
        listBox.innerHTML = `<tr><td colspan="7">用户列表加载失败</td></tr>`;
        emptyState.innerHTML = `
          <strong>用户数据加载失败</strong>
          <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
        `;
      });
  }

  function initAdminOrdersPage() {
    if (getCurrentFileName() !== ROUTES.admin.orders) {
      return;
    }

    const summaryBox = document.querySelector("[data-admin-order-summary]");
    const filterForm = document.querySelector("[data-admin-order-filter]");
    const listBox = document.querySelector("[data-admin-order-list]");
    const refundStatusInput = filterForm?.querySelector("[name='refundStatus']");
    const reviewNoteInput = filterForm?.querySelector("[name='reviewNote']");
    if (!summaryBox || !filterForm || !listBox || !refundStatusInput) {
      return;
    }

    const renderSummary = (orders) => {
      const requested = orders.filter((order) => order?.refundStatus === "requested").length;
      const refunded = orders.filter((order) => order?.refundStatus === "refunded").length;
      const rejected = orders.filter((order) => order?.refundStatus === "rejected").length;
      summaryBox.innerHTML = `
        <span class="mini-chip">待审核 ${requested}</span>
        <span class="mini-chip">已退款 ${refunded}</span>
        <span class="mini-chip">已驳回 ${rejected}</span>
      `;
    };

    const renderList = (orders, reload) => {
      if (!orders.length) {
        listBox.innerHTML = `
          <div class="info-card">
            <strong>当前没有符合条件的订单</strong>
            <p class="muted">可以切换退款状态筛选，或等待乘客提交新的退款申请。</p>
          </div>
        `;
        return;
      }

      listBox.innerHTML = orders.map((order) => {
        const trip = order?.trip || {};
        const user = order?.user || {};
        const canReview = order?.refundStatus === "requested";
        const reviewNote = order?.refundReviewNote
          ? `<p class="muted">审核备注：${escapeHtml(order.refundReviewNote)}</p>`
          : `<p class="muted">审核备注：暂无</p>`;

        return `
          <div class="order-card" data-admin-order-card="${order.id}">
            <div class="order-top">
              <div>
                <strong>${escapeHtml(order.orderNo || `订单 #${order.id}`)}</strong>
                <div class="list-meta">
                  <span>${escapeHtml(user.nickname || user.phone || "未知用户")}</span>
                  <span>${escapeHtml(`${trip.startCity || "--"} -> ${trip.endCity || "--"}`)}</span>
                </div>
                <div class="list-meta">
                  <span>${escapeHtml(mapOrderStatus(order.orderStatus, order.payStatus))}</span>
                  <span>${escapeHtml(mapRefundStatus(order.refundStatus))}</span>
                  <span>${escapeHtml(formatMoneyFromCent(order.amount || 0))}</span>
                  <span>${escapeHtml(formatFullDateTime(order.createdAt))}</span>
                </div>
                ${reviewNote}
                ${order.refundReviewedAt ? `<p class="muted">审核时间：${escapeHtml(formatFullDateTime(order.refundReviewedAt))}</p>` : ""}
              </div>
              <span class="badge">${escapeHtml(mapRefundStatus(order.refundStatus))}</span>
            </div>
            <div class="button-row">
              <button class="button button-primary" type="button" data-admin-approve-refund="${order.id}" ${canReview ? "" : "disabled"}>
                ${canReview ? "通过退款" : "不可通过"}
              </button>
              <button class="button button-ghost" type="button" data-admin-reject-refund="${order.id}" ${canReview ? "" : "disabled"}>
                ${canReview ? "驳回退款" : "不可驳回"}
              </button>
            </div>
          </div>
        `;
      }).join("");

      listBox.querySelectorAll("[data-admin-approve-refund]").forEach((button) => {
        button.addEventListener("click", async () => {
          const orderId = button.getAttribute("data-admin-approve-refund");
          if (!orderId) {
            return;
          }

          const reviewNote = String(reviewNoteInput?.value || "").trim();
          const originalText = button.textContent;
          button.disabled = true;
          button.textContent = "审核中...";

          try {
            await api.post(API_ENDPOINTS.admin.approveRefund, { reviewNote }, { pathParams: { orderId } });
            showToast("退款已审核通过");
            await reload();
          } catch (error) {
            button.disabled = false;
            button.textContent = originalText;
            showToast(error.message || "审核通过失败");
          }
        });
      });

      listBox.querySelectorAll("[data-admin-reject-refund]").forEach((button) => {
        button.addEventListener("click", async () => {
          const orderId = button.getAttribute("data-admin-reject-refund");
          if (!orderId) {
            return;
          }

          const reviewNote = String(reviewNoteInput?.value || "").trim();
          const originalText = button.textContent;
          button.disabled = true;
          button.textContent = "处理中...";

          try {
            await api.post(API_ENDPOINTS.admin.rejectRefund, { reviewNote }, { pathParams: { orderId } });
            showToast("退款申请已驳回");
            await reload();
          } catch (error) {
            button.disabled = false;
            button.textContent = originalText;
            showToast(error.message || "驳回退款失败");
          }
        });
      });
    };

    const loadOrders = async () => {
      const refundStatus = String(refundStatusInput.value || "").trim();
      const result = await api.get(API_ENDPOINTS.admin.orders, { refundStatus });
      const orders = result?.data || [];
      renderSummary(orders);
      renderList(orders, loadOrders);
    };

    refundStatusInput.addEventListener("change", () => {
      loadOrders().catch((error) => {
        showToast(error.message || "加载订单失败");
      });
    });

    loadOrders().catch((error) => {
      summaryBox.innerHTML = `<span class="mini-chip">加载失败</span>`;
      listBox.innerHTML = `
        <div class="info-card">
          <strong>退款订单加载失败</strong>
          <p class="muted">${escapeHtml(error.message || "请稍后重试。")}</p>
        </div>
      `;
    });
  }

  function initOrdersPage() {
    return;
  }

  function initOrderDetailPage() {
    return;
  }

  function initAdminDashboardPage() {
    if (getCurrentFileName() !== ROUTES.admin.dashboard) {
      return;
    }

    const summaryBox = document.querySelector("[data-admin-dashboard-summary]");
    const pendingBox = document.querySelector("[data-admin-dashboard-pending-refund]");
    const refundedBox = document.querySelector("[data-admin-dashboard-refunded]");
    const rejectedBox = document.querySelector("[data-admin-dashboard-rejected]");
    const pendingNote = document.querySelector("[data-admin-dashboard-pending-note]");
    const emptyState = document.querySelector("[data-admin-dashboard-empty-state]");

    const setLoadError = (message) => {
      if (summaryBox) {
        summaryBox.innerHTML = `<span class="mini-chip">加载失败</span>`;
      }
      if (pendingBox) {
        pendingBox.textContent = "--";
      }
      if (refundedBox) {
        refundedBox.textContent = "--";
      }
      if (rejectedBox) {
        rejectedBox.textContent = "--";
      }
      if (pendingNote) {
        pendingNote.textContent = "请◢后重?";
      }
      if (emptyState) {
        emptyState.innerHTML = `
          <strong>统计加载失败</strong>
          <p class="muted">${escapeHtml(message || "请稍后重试")}</p>
        `;
      }
    };

    api.get(API_ENDPOINTS.admin.dashboard)
      .then((result) => {
        const summary = result?.data || {};
        const pendingRefundCount = Number(summary.pendingRefundCount || 0);
        const refundedCount = Number(summary.refundedCount || 0);
        const rejectedRefundCount = Number(summary.rejectedRefundCount || 0);

        if (summaryBox) {
          summaryBox.innerHTML = `
            <span class="mini-chip">待审核 ${pendingRefundCount}</span>
            <span class="mini-chip">已退款 ${refundedCount}</span>
            <span class="mini-chip">已┏?${rejectedRefundCount}</span>
          `;
        }
        if (pendingBox) {
          pendingBox.textContent = String(pendingRefundCount);
        }
        if (refundedBox) {
          refundedBox.textContent = String(refundedCount);
        }
        if (rejectedBox) {
          rejectedBox.textContent = String(rejectedRefundCount);
        }
        if (pendingNote) {
          pendingNote.textContent = pendingRefundCount > 0 ? `当前还有 ${pendingRefundCount} 笔退款待处理` : "当前没有待审核退款";
        }
        if (emptyState) {
          emptyState.innerHTML = pendingRefundCount > 0
            ? `<strong>退款待办提醒</strong><p class="muted">建议优先进入退款审核页，避免乘客等待过久。</p>`
            : `<strong>暂无退款待办</strong><p class="muted">当前没有待审核的退款申请。</p>`;
        }
      })
      .catch((error) => {
        setLoadError(error.message || "统计加载失败");
      });
  }

  function initAdminDashboardPageV2() {
    if (getCurrentFileName() !== ROUTES.admin.dashboard) {
      return;
    }

    const summaryBox = document.querySelector("[data-admin-dashboard-summary]");
    const totalUsersBox = document.querySelector("[data-admin-dashboard-total-users]");
    const activeUsersBox = document.querySelector("[data-admin-dashboard-active-users]");
    const pendingBox = document.querySelector("[data-admin-dashboard-pending-refund]");
    const refundedBox = document.querySelector("[data-admin-dashboard-refunded]");
    const rejectedBox = document.querySelector("[data-admin-dashboard-rejected]");
    const passengerCountBox = document.querySelector("[data-admin-dashboard-passenger-count]");
    const driverCountBox = document.querySelector("[data-admin-dashboard-driver-count]");
    const adminCountBox = document.querySelector("[data-admin-dashboard-admin-count]");
    const pendingNote = document.querySelector("[data-admin-dashboard-pending-note]");
    const emptyState = document.querySelector("[data-admin-dashboard-empty-state]");

    if (!summaryBox || !pendingBox || !refundedBox || !rejectedBox || !emptyState) {
      return;
    }

    const setLoadError = (message) => {
      summaryBox.innerHTML = `<span class="mini-chip">加载失败</span>`;
      if (totalUsersBox) totalUsersBox.textContent = "--";
      if (activeUsersBox) activeUsersBox.textContent = "--";
      if (pendingBox) pendingBox.textContent = "--";
      if (refundedBox) refundedBox.textContent = "--";
      if (rejectedBox) rejectedBox.textContent = "--";
      if (passengerCountBox) passengerCountBox.textContent = "--";
      if (driverCountBox) driverCountBox.textContent = "--";
      if (adminCountBox) adminCountBox.textContent = "--";
      if (pendingNote) pendingNote.textContent = "请◢后重?";
      emptyState.innerHTML = `
        <strong>后台统计加载失败</strong>
        <p class="muted">${escapeHtml(message || "请稍后重试")}</p>
      `;
    };

    Promise.all([
      api.get(API_ENDPOINTS.admin.dashboard),
      api.get(API_ENDPOINTS.admin.userSummary),
    ])
      .then(([dashboardResult, userSummaryResult]) => {
        const refundSummary = dashboardResult?.data || {};
        const userSummary = userSummaryResult?.data || {};
        const pendingRefundCount = Number(refundSummary.pendingRefundCount || 0);
        const refundedCount = Number(refundSummary.refundedCount || 0);
        const rejectedRefundCount = Number(refundSummary.rejectedRefundCount || 0);
        const totalUsers = Number(userSummary.totalUsers || 0);
        const activeUsers = Number(userSummary.activeCount || 0);
        const passengerCount = Number(userSummary.passengerCount || 0);
        const driverCount = Number(userSummary.driverCount || 0);
        const adminCount = Number(userSummary.adminCount || 0);

        summaryBox.innerHTML = `
          <span class="mini-chip">总用户 ${totalUsers}</span>
          <span class="mini-chip">待审核退款 ${pendingRefundCount}</span>
          <span class="mini-chip">活跃ㄦ埛 ${activeUsers}</span>
        `;
        if (totalUsersBox) totalUsersBox.textContent = String(totalUsers);
        if (activeUsersBox) activeUsersBox.textContent = String(activeUsers);
        if (pendingBox) pendingBox.textContent = String(pendingRefundCount);
        if (refundedBox) refundedBox.textContent = String(refundedCount);
        if (rejectedBox) rejectedBox.textContent = String(rejectedRefundCount);
        if (passengerCountBox) passengerCountBox.textContent = String(passengerCount);
        if (driverCountBox) driverCountBox.textContent = String(driverCount);
        if (adminCountBox) adminCountBox.textContent = String(adminCount);
        if (pendingNote) {
          pendingNote.textContent = pendingRefundCount > 0
            ? `当前还有 ${pendingRefundCount} 笔退款待处理`
            : "当前没有待审核退款";
        }
        emptyState.innerHTML = pendingRefundCount > 0
          ? `<strong>退款待办提醒</strong><p class="muted">建议优先进入退款审核页，避免乘客等待过久。</p>`
          : `<strong>退款队列正常</strong><p class="muted">当前没有待审核退款，可以转去查看用户列表或其他后台模块。</p>`;
      })
      .catch((error) => {
        setLoadError(error.message || "统计加载失败");
      });
  }

  function getOrderPaymentExpireMeta(order) {
    const raw = order?.paymentExpireAt;
    if (!raw) {
      return {
        hasExpireAt: false,
        expired: false,
        remainingMs: 0,
        formattedTime: "--",
      };
    }

    const expireAt = new Date(raw);
    if (Number.isNaN(expireAt.getTime())) {
      return {
        hasExpireAt: false,
        expired: false,
        remainingMs: 0,
        formattedTime: "--",
      };
    }

    const remainingMs = expireAt.getTime() - Date.now();
    return {
      hasExpireAt: true,
      expired: remainingMs <= 0,
      remainingMs,
      formattedTime: formatFullDateTime(raw),
    };
  }

  function formatCountdown(remainingMs) {
    const safeMs = Math.max(0, Number(remainingMs || 0));
    const totalSeconds = Math.floor(safeMs / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`;
  }

  function getDisplayedOrderStatus(order) {
    const expireMeta = getOrderPaymentExpireMeta(order);
    if ((order?.orderStatus === "pending_payment" || order?.payStatus === "unpaid") && expireMeta.expired) {
      return "支付已超?";
    }
    return mapOrderStatus(order?.orderStatus, order?.payStatus);
  }

  function initOrdersPage() {
    if (getCurrentFileName() !== ROUTES.passenger.orders) {
      return;
    }

    const summaryBox = document.querySelector("[data-order-summary]");
    const listBox = document.querySelector("[data-order-list]");
    if (!summaryBox || !listBox) {
      return;
    }

    summaryBox.style.display = "none";
    summaryBox.innerHTML = "";

    api.get(API_ENDPOINTS.passenger.myOrders)
      .then((result) => {
        const orders = Array.isArray(result?.data) ? result.data : [];

        if (!orders.length) {
          listBox.innerHTML = `
            <div class="info-card">
              <strong>暂无订单</strong>
              <p class="muted">当前账号还没有任何订单记录。</p>
            </div>
          `;
          return;
        }

        const pendingPaymentCount = orders.filter((order) => order?.orderStatus === "pending_payment" || order?.payStatus === "unpaid").length;
        const pendingDepartureCount = orders.filter((order) => order?.orderStatus === "pending_verification").length;
        const completedCount = orders.filter((order) => order?.orderStatus === "completed").length;

        summaryBox.style.display = "";
        summaryBox.innerHTML = `
          <span class="mini-chip">待支付 ${pendingPaymentCount}</span>
          <span class="mini-chip">待出行 ${pendingDepartureCount}</span>
          <span class="mini-chip">已完成 ${completedCount}</span>
        `;

        listBox.innerHTML = orders.map((order) => {
          const trip = order?.trip || null;
          const route = trip ? `${trip.startCity} -> ${trip.endCity}` : "";
          const departureTime = trip?.departureTime ? formatFullDateTime(trip.departureTime) : "";
          const orderDetailHref = buildOrderDetailHref(order?.id);
          const expireMeta = getOrderPaymentExpireMeta(order);
          const statusText = getDisplayedOrderStatus(order);
          const routeStations = trip
            ? [
                trip.startCity,
                ...(Array.isArray(trip.stops) ? trip.stops : []).map((stop) => stop?.stopName || stop?.name || ""),
                trip.endCity,
              ]
            : [];
          const orderStations = renderStationLinks(routeStations, "", orderDetailHref, "进入该订单详情");
          const refundTag = order?.refundStatus && order.refundStatus !== "none"
            ? `<span class="tag">${escapeHtml(mapRefundStatus(order.refundStatus))}</span>`
            : "";
          const expireTag = expireMeta.expired && (order?.orderStatus === "pending_payment" || order?.payStatus === "unpaid")
            ? `<span class="tag">支付已超时</span>`
            : "";
          const reviewNote = order?.refundReviewNote
            ? `<p class="muted">审核备注：${escapeHtml(order.refundReviewNote)}</p>`
            : "";
          const reviewedAt = order?.refundReviewedAt
            ? `<p class="muted">审核时间：${escapeHtml(formatFullDateTime(order.refundReviewedAt))}</p>`
            : "";
          const expireNote = expireMeta.hasExpireAt
            ? `<p class="muted">${expireMeta.expired ? "支付截止时间已到：" : "支付截止时间："}${escapeHtml(expireMeta.formattedTime)}</p>`
            : "";

          const canContinuePay = (order?.orderStatus === "pending_payment" || order?.payStatus === "unpaid") && !expireMeta.expired;
          const primaryAction = canContinuePay
            ? `<a class="button button-primary" href="${ROUTES.passenger.payment}?orderId=${order.id}">去支付</a>`
            : `<a class="button button-primary" href="${orderDetailHref}">查看详情</a>`;

          const secondaryAction = canContinuePay
            ? `<a class="button button-secondary" href="${orderDetailHref}">订单详情</a>`
            : `<a class="button button-ghost" href="${orderDetailHref}">查看订单</a>`;

          return `
            <div class="order-card">
              <div class="order-top">
                <div>
                  <strong>${escapeHtml(order.orderNo || `订单 #${order.id}`)}</strong>
                  <div class="list-meta">
                    ${route ? `<span>${escapeHtml(route)}</span>` : ""}
                    ${departureTime ? `<span>${escapeHtml(departureTime)}</span>` : ""}
                    ${refundTag}
                    ${expireTag}
                  </div>
                  ${expireNote}
                  ${orderStations ? `<div class="muted">经停站点可点入对应订单：</div>${orderStations}` : ""}
                  ${reviewNote}
                  ${reviewedAt}
                </div>
                <div class="price-pill">${escapeHtml(formatMoneyFromCent(order.amount))} <small>${escapeHtml(statusText)}</small></div>
              </div>
              <div class="button-row">
                ${primaryAction}
                ${secondaryAction}
              </div>
            </div>
          `;
        }).join("");
      })
      .catch((error) => {
        summaryBox.style.display = "none";
        summaryBox.innerHTML = "";
        listBox.innerHTML = `
          <div class="info-card">
            <strong>订单加载失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后再试。")}</p>
          </div>
        `;
      });
  }

  function initOrderDetailPage() {
    if (getCurrentFileName() !== ROUTES.passenger.orderDetail) {
      return;
    }

    const orderId = getQueryParam("orderId");
    if (!orderId) {
      return;
    }

    const auth = readAuth() || {};
    const title = document.querySelector("[data-order-detail-title]");
    const statusBox = document.querySelector("[data-order-detail-status]");
    const timelineBox = document.querySelector("[data-order-detail-timeline]");
    const metaBox = document.querySelector("[data-order-detail-meta]");
    const pricingBox = document.querySelector("[data-order-detail-pricing]");
    const actionsBox = document.querySelector("[data-order-detail-actions]");

    api.get(API_ENDPOINTS.passenger.orderDetail, undefined, { pathParams: { orderId } })
      .then((result) => {
        const order = result?.data;
        if (!order) {
          throw new Error("订单不存在");
        }

        const trip = order.trip || {};
        const statusText = getDisplayedOrderStatus(order);
        const routeTitle = trip.startCity && trip.endCity ? `${trip.startCity} -> ${trip.endCity}` : "";
        const ticketCount = Number(order.ticketCount || 0);
        const unitPrice = ticketCount > 0 ? Math.round(Number(order.amount || 0) / ticketCount) : Number(order.amount || 0);
        const stops = Array.isArray(trip.stops) ? trip.stops : [];
        const refundReviewNote = order.refundReviewNote || "";
        const refundReviewedAt = order.refundReviewedAt ? formatFullDateTime(order.refundReviewedAt) : "";
        const expireMeta = getOrderPaymentExpireMeta(order);
        const expiredPendingPayment = (order.orderStatus === "pending_payment" || order.payStatus === "unpaid") && expireMeta.expired;

        if (title) {
          title.textContent = order.orderNo || "订单详情";
        }
        if (statusBox) {
          statusBox.textContent = statusText;
        }

        if (timelineBox) {
          const stopMarkup = stops.map((stop) => `
            <div class="timeline-item">
              <span class="timeline-dot"></span>
              <strong>${escapeHtml(stop.stopName || "")}</strong>
              <p class="muted">${escapeHtml(formatFullDateTime(stop.planArrivalTime || stop.planDepartureTime || trip.departureTime))}</p>
            </div>
          `).join("");

          timelineBox.innerHTML = `
            <div class="timeline-item">
              <span class="timeline-dot"></span>
              <strong>${escapeHtml(trip.startCity || "")}</strong>
              <p class="muted">${escapeHtml(formatFullDateTime(trip.departureTime))}</p>
            </div>
            ${stopMarkup}
            <div class="timeline-item">
              <span class="timeline-dot"></span>
              <strong>${escapeHtml(trip.endCity || "")}</strong>
              <p class="muted">${escapeHtml(formatFullDateTime(trip.arrivalTime))}</p>
            </div>
          `;
        }

        if (metaBox) {
          metaBox.innerHTML = `
            <div class="list-item">
              <strong>路线信息</strong>
              <div class="list-meta"><span>${escapeHtml(routeTitle)}</span><span>${escapeHtml(trip.vehicleType || "car")}</span></div>
            </div>
            <div class="list-item">
              <strong>乘车人</strong>
              <div class="list-meta"><span>${escapeHtml(auth.nickname || auth.phone || "当前账号")}</span><span>${escapeHtml(auth.phone || "")}</span></div>
            </div>
            <div class="list-item">
              <strong>座位与张数</strong>
              <div class="list-meta"><span>${escapeHtml(mapSeatType(order.seatType))}</span><span>${ticketCount} 张</span></div>
            </div>
            <div class="list-item">
              <strong>支付截止时间</strong>
              <div class="list-meta"><span>${escapeHtml(expireMeta.formattedTime)}</span><span>${expiredPendingPayment ? "已超时" : "有效"}</span></div>
            </div>
            <div class="list-item">
              <strong>支付与退款</strong>
              <div class="list-meta"><span>${escapeHtml(order.payStatus || "--")}</span><span>${escapeHtml(mapRefundStatus(order.refundStatus))}</span></div>
            </div>
            <div class="list-item">
              <strong>退款审核备注</strong>
              <div class="list-meta"><span>${escapeHtml(refundReviewNote || "暂无")}</span><span>${escapeHtml(refundReviewedAt || "--")}</span></div>
            </div>
          `;
        }

        if (pricingBox) {
          pricingBox.innerHTML = `
            <div class="pricing-row"><span>单价</span><strong>${escapeHtml(formatMoneyFromCent(unitPrice))}</strong></div>
            <div class="pricing-row"><span>数量</span><strong>${ticketCount} 张</strong></div>
            <div class="pricing-row"><span>支付截止</span><strong>${escapeHtml(expireMeta.formattedTime)}</strong></div>
            <div class="pricing-row"><span>退款状态</span><strong>${escapeHtml(mapRefundStatus(order.refundStatus))}</strong></div>
            <div class="pricing-row pricing-total"><span>实付</span><strong>${escapeHtml(formatMoneyFromCent(order.amount))}</strong></div>
          `;
        }

        if (actionsBox) {
          const canContinuePay = (order.orderStatus === "pending_payment" || order.payStatus === "unpaid") && !expireMeta.expired;
          const canCancel = order.orderStatus === "pending_payment";
          const canRequestRefund =
            (order.orderStatus === "pending_verification" || order.orderStatus === "completed") &&
            order.payStatus === "paid" &&
            order.orderStatus !== "cancelled" &&
            (order.refundStatus === "none" || order.refundStatus === "rejected");

          const primaryAction = canContinuePay
            ? `<a class="button button-primary" href="${ROUTES.passenger.payment}?orderId=${order.id}">继续支付</a>`
            : `<button class="button button-primary" type="button" ${expiredPendingPayment ? "disabled" : "data-electronic-ticket"}>${expiredPendingPayment ? "支付已超时" : "查看电子票"}</button>`;

          const cancelAction = canCancel
            ? `<button class="button button-ghost" type="button" data-order-cancel>取消订单</button>`
            : "";

          const refundAction = order.refundStatus === "requested"
            ? `<button class="button button-ghost" type="button" disabled>退款处理中</button>`
            : canRequestRefund
              ? `<button class="button button-ghost" type="button" data-order-refund>${order.refundStatus === "rejected" ? "重新申请退款" : "申请退款"}</button>`
              : "";

          const expireInfo = expireMeta.hasExpireAt
            ? `<div class="info-card"><strong>支付时效</strong><p class="muted">${expiredPendingPayment ? "该订单已超过支付截止时间，不能继续支付。" : `请在 ${escapeHtml(expireMeta.formattedTime)} 前完成支付。`}</p></div>`
            : "";
          const reviewInfo = order.refundStatus !== "none"
            ? `<div class="info-card"><strong>退款结果</strong><p class="muted">${escapeHtml(mapRefundStatus(order.refundStatus))}</p><p class="muted">审核备注：${escapeHtml(refundReviewNote || "暂无")}</p><p class="muted">审核时间：${escapeHtml(refundReviewedAt || "--")}</p></div>`
            : "";

          actionsBox.innerHTML = `
            ${primaryAction}
            ${cancelAction}
            ${refundAction}
            <a class="button button-secondary" href="${ROUTES.passenger.orders}">返回订单列表</a>
            <div class="info-card" data-electronic-ticket-box style="display:none"></div>
            ${expireInfo}
            ${reviewInfo}
          `;
          initToastTriggers();

          const ticketButton = actionsBox.querySelector("[data-electronic-ticket]");
          const ticketBox = actionsBox.querySelector("[data-electronic-ticket-box]");
          if (ticketButton && ticketBox) {
            ticketButton.addEventListener("click", async () => {
              try {
                ticketButton.disabled = true;
                ticketButton.textContent = "出票中...";
                const result = await api.get(API_ENDPOINTS.passenger.orderTicket, undefined, { pathParams: { orderId } });
                const ticket = result?.data || {};
                ticketBox.style.display = "";
                ticketBox.innerHTML = `
                  <strong>电子票 token</strong>
                  <p class="muted">司机端可复制该 token 到核验框完成验票。后续可接入 qrcode 组件渲染为二维码。</p>
                  <textarea class="ticket-token-box" readonly>${escapeHtml(ticket.token || "")}</textarea>
                  <div class="list-meta"><span>状态：${escapeHtml(ticket.status || "--")}</span><span>有效期：${escapeHtml(formatFullDateTime(ticket.expiresAt))}</span></div>
                `;
                ticketButton.textContent = "刷新电子票";
              } catch (error) {
                showToast(error.message || "电子票加载失败");
                ticketButton.textContent = "查看电子票";
              } finally {
                ticketButton.disabled = false;
              }
            });
          }

          const cancelButton = actionsBox.querySelector("[data-order-cancel]");
          if (cancelButton) {
            cancelButton.addEventListener("click", async () => {
              if (!window.confirm("确认取消这个订单吗？")) {
                return;
              }

              try {
                cancelButton.disabled = true;
                cancelButton.textContent = "取消中...";
                await api.post(API_ENDPOINTS.passenger.cancelOrder, {}, { pathParams: { orderId } });
                showToast("订单已取消");
                window.setTimeout(() => window.location.reload(), 300);
              } catch (error) {
                cancelButton.disabled = false;
                cancelButton.textContent = "取消订单";
                showToast(error.message || "取消订单失败");
              }
            });
          }

          const refundButton = actionsBox.querySelector("[data-order-refund]");
          if (refundButton) {
            const defaultRefundText = order.refundStatus === "rejected" ? "重新申请退款" : "申请退款";
            refundButton.addEventListener("click", async () => {
              if (!window.confirm("确认提交退款申请吗？")) {
                return;
              }

              try {
                refundButton.disabled = true;
                refundButton.textContent = "提交?..";
                await api.post(API_ENDPOINTS.passenger.refundOrder, {}, { pathParams: { orderId } });
                showToast("退款申请已提交");
                window.setTimeout(() => window.location.reload(), 300);
              } catch (error) {
                refundButton.disabled = false;
                refundButton.textContent = defaultRefundText;
                showToast(error.message || "退款申请失败");
              }
            });
          }
        }
      })
      .catch((error) => {
        if (title) {
          title.textContent = error.message || "订单加载失败";
        }
        if (timelineBox) {
          timelineBox.innerHTML = `
            <div class="timeline-item">
              <span class="timeline-dot"></span>
              <strong>加载失败</strong>
              <p class="muted">${escapeHtml(error.message || "请稍后再试。")}</p>
            </div>
          `;
        }
      });
  }

  function initPaymentPage() {
    if (getCurrentFileName() !== ROUTES.passenger.payment) {
      return;
    }

    const title = document.querySelector("[data-payment-title]");
    const lede = document.querySelector("[data-payment-lede]");
    const summaryBox = document.querySelector("[data-payment-summary]");
    const orderBox = document.querySelector("[data-payment-order]");
    const amountBox = document.querySelector("[data-payment-amount]");
    const actionsBox = document.querySelector("[data-payment-actions]");
    const orderId = getQueryParam("orderId");
    let countdownTimer = null;
    let pollTimer = null;
    let currentPaymentId = "";

    const stopCountdown = () => {
      if (countdownTimer) {
        window.clearInterval(countdownTimer);
        countdownTimer = null;
      }
    };

    const stopPolling = () => {
      if (pollTimer) {
        window.clearInterval(pollTimer);
        pollTimer = null;
      }
    };

    const renderMissingOrderState = () => {
      if (title) {
        title.textContent = "缺少订单";
      }
      if (lede) {
        lede.textContent = "当前支付页没有收到 orderId，请从订单详情页或订单列表重新进入。";
      }
      if (summaryBox) {
        summaryBox.innerHTML = `
          <div class="list-item">
            <strong>无法发起支付</strong>
            <div class="list-meta"><span>原因</span><span>URL 中缺少 orderId</span></div>
          </div>
        `;
      }
      if (orderBox) {
        orderBox.innerHTML = `
          <div class="list-item">
            <strong>建议操作</strong>
            <div class="list-meta"><span>返回订单详情页</span><span>重新点击支付</span></div>
          </div>
        `;
      }
      if (amountBox) {
        amountBox.innerHTML = `
          <div class="pricing-row"><span>订单金额</span><strong>--</strong></div>
          <div class="pricing-row"><span>支付状态</span><strong>--</strong></div>
          <div class="pricing-row pricing-total"><span>订单状态</span><strong>--</strong></div>
        `;
      }
      if (actionsBox) {
        actionsBox.innerHTML = `
          <a class="button button-primary" href="${ROUTES.passenger.orders}">返回订单列表</a>
        `;
      }
    };

    if (!orderId) {
      renderMissingOrderState();
      return;
    }

    const fetchPaymentStatus = async (paymentId) => {
      const result = await api.get(API_ENDPOINTS.passenger.paymentStatus, undefined, {
        pathParams: { paymentId },
      });
      return result?.data || null;
    };

    const startPolling = (paymentId) => {
      stopPolling();
      pollTimer = window.setInterval(async () => {
        try {
          const payment = await fetchPaymentStatus(paymentId);
          if (payment?.status === "paid") {
            stopPolling();
            stopCountdown();
            showToast("支付成功");
            window.setTimeout(() => {
              redirectTo(`${ROUTES.passenger.orderDetail}?orderId=${orderId}`);
            }, 300);
          }
        } catch (_) {
          stopPolling();
        }
      }, 2000);
    };

    const renderCountdown = (order) => {
      stopCountdown();

      const updateCountdown = () => {
        const expireMeta = getOrderPaymentExpireMeta(order);
        if (title) {
          title.textContent = expireMeta.expired ? "订单支付已超时" : "订单支付处理中";
        }
        if (lede) {
          lede.textContent = expireMeta.expired
            ? "该订单已超过支付截止时间，系统会自动取消并释放座位。"
            : `请在 ${formatCountdown(expireMeta.remainingMs)} 内完成支付。`;
        }
      };

      updateCountdown();

      const expireMeta = getOrderPaymentExpireMeta(order);
      if (!expireMeta.expired && expireMeta.hasExpireAt) {
        countdownTimer = window.setInterval(() => {
          const latestMeta = getOrderPaymentExpireMeta(order);
          if (latestMeta.expired) {
            stopCountdown();
            updateCountdown();
            window.setTimeout(() => window.location.reload(), 800);
            return;
          }
          updateCountdown();
        }, 1000);
      }
    };

    const renderPaymentPage = (order, payment) => {
      const expireMeta = getOrderPaymentExpireMeta(order);
      const expiredPendingPayment = (order?.orderStatus === "pending_payment" || order?.payStatus === "unpaid") && expireMeta.expired;
      const displayedStatus = getDisplayedOrderStatus(order);
      const trip = order?.trip || {};

      renderCountdown(order);

      if (summaryBox) {
        summaryBox.innerHTML = `
          <div class="list-item">
            <strong>支付说明</strong>
            <div class="list-meta"><span>${expiredPendingPayment ? "订单已超时" : "请尽快完成支付"}</span><span>${expireMeta.hasExpireAt ? escapeHtml(expireMeta.formattedTime) : "--"}</span></div>
          </div>
          <div class="list-item">
            <strong>剩余支付时间</strong>
            <div class="list-meta"><span>${expiredPendingPayment ? "00:00" : escapeHtml(formatCountdown(expireMeta.remainingMs))}</span><span>${escapeHtml(payment?.status || "pending")}</span></div>
          </div>
        `;
      }

      if (orderBox) {
        orderBox.innerHTML = `
          <div class="list-item">
            <strong>${escapeHtml(order.orderNo || `订单 #${order.id}`)}</strong>
            <div class="list-meta"><span>${escapeHtml(`${trip.startCity || "--"} -> ${trip.endCity || "--"}`)}</span><span>${escapeHtml(formatFullDateTime(trip.departureTime))}</span></div>
          </div>
        `;
      }

      if (amountBox) {
        amountBox.innerHTML = `
          <div class="pricing-row"><span>订单金额</span><strong>${escapeHtml(formatMoneyFromCent(order.amount))}</strong></div>
          <div class="pricing-row"><span>支付状态</span><strong>${escapeHtml(payment?.status || order.payStatus || "--")}</strong></div>
          <div class="pricing-row"><span>支付截止</span><strong>${escapeHtml(expireMeta.formattedTime)}</strong></div>
          <div class="pricing-row pricing-total"><span>订单状态</span><strong>${escapeHtml(displayedStatus)}</strong></div>
        `;
      }

      if (!actionsBox) {
        return;
      }

      if (order.payStatus === "paid") {
        actionsBox.innerHTML = `
          <a class="button button-primary" href="${ROUTES.passenger.orderDetail}?orderId=${order.id}">查看订单详情</a>
          <a class="button button-secondary" href="${ROUTES.passenger.orders}">返回订单列表</a>
        `;
        return;
      }

      if (expiredPendingPayment || order.orderStatus === "cancelled") {
        actionsBox.innerHTML = `
          <button class="button button-primary" type="button" disabled>订单已超时</button>
          <a class="button button-secondary" href="${ROUTES.passenger.orderDetail}?orderId=${order.id}">查看订单详情</a>
          <a class="button button-ghost" href="${ROUTES.passenger.orders}">返回订单列表</a>
        `;
        return;
      }

      actionsBox.innerHTML = `
        <button class="button button-primary" type="button" data-payment-create>${payment ? "继续支付" : "创建支付单"}</button>
        <button class="button button-ghost" type="button" data-payment-success ${payment ? "" : "disabled"}>模拟支付成功</button>
        <a class="button button-secondary" href="${ROUTES.passenger.orderDetail}?orderId=${order.id}">查看订单详情</a>
      `;

      const createButton = actionsBox.querySelector("[data-payment-create]");
      const successButton = actionsBox.querySelector("[data-payment-success]");

      if (createButton) {
        createButton.addEventListener("click", async () => {
          try {
            createButton.disabled = true;
            createButton.textContent = "创建中...";
            const result = await api.post(API_ENDPOINTS.passenger.createPayment, {
              orderId: Number(order.id),
              channel: "mock",
            });
            const createdPayment = result?.data;
            currentPaymentId = String(createdPayment?.id || "");
            showToast(createdPayment?.id ? "支付单已创建" : "已复用待支付支付单");
            renderPaymentPage(order, createdPayment);
            if (createdPayment?.id) {
              startPolling(String(createdPayment.id));
            }
          } catch (error) {
            createButton.disabled = false;
            createButton.textContent = payment ? "继续支付" : "创建支付单";
            showToast(error.message || "创建支付单失败");
          }
        });
      }

      if (successButton) {
        successButton.addEventListener("click", async () => {
          const targetPaymentId = currentPaymentId || String(payment?.id || "");
          if (!targetPaymentId) {
            showToast("请先创建支付?");
            return;
          }

          try {
            successButton.disabled = true;
            successButton.textContent = "支付?..";
            await api.post(API_ENDPOINTS.passenger.mockPaymentSuccess, {}, {
              pathParams: { paymentId: targetPaymentId },
            });
            stopPolling();
            stopCountdown();
            showToast("支付成功");
            window.setTimeout(() => {
              redirectTo(`${ROUTES.passenger.orderDetail}?orderId=${order.id}`);
            }, 300);
          } catch (error) {
            successButton.disabled = false;
            successButton.textContent = "模拟支付成功";
            showToast(error.message || "支付失败");
          }
        });
      }
    };

    const loadPage = async () => {
      const orderResult = await api.get(API_ENDPOINTS.passenger.orderDetail, undefined, {
        pathParams: { orderId },
      });
      const order = orderResult?.data;
      if (!order) {
        throw new Error("订单不存在");
      }

      let payment = null;
      try {
        const createResult = await api.post(API_ENDPOINTS.passenger.createPayment, {
          orderId: Number(order.id),
          channel: "mock",
        });
        payment = createResult?.data || null;
      } catch (error) {
        if (!String(error.message || "").includes("expired")) {
          throw error;
        }
      }

      currentPaymentId = String(payment?.id || "");
      renderPaymentPage(order, payment);
      if (payment?.id && payment?.status === "pending") {
        startPolling(String(payment.id));
      }
    };

    loadPage().catch((error) => {
      stopCountdown();
      stopPolling();
      if (title) {
        title.textContent = "支付信息加载失败";
      }
      if (lede) {
        lede.textContent = error.message || "请稍后再试";
      }
      if (actionsBox) {
        actionsBox.innerHTML = `
          <a class="button button-primary" href="${ROUTES.passenger.orderDetail}?orderId=${orderId}">返回订单详情</a>
          <a class="button button-secondary" href="${ROUTES.passenger.orders}">返回订单列表</a>
        `;
      }
    });
  }

  function initNotificationCenter() {
    if (getCurrentFileName() !== ROUTES.passenger.profile) {
      return;
    }

    const unreadBox = document.querySelector("[data-notification-unread]");
    const listBox = document.querySelector("[data-notification-list]");
    const readAllButton = document.querySelector("[data-notification-read-all]");

    if (!unreadBox || !listBox || !readAllButton) {
      return;
    }

    const renderList = (notifications) => {
      if (!notifications.length) {
        listBox.innerHTML = `
          <div class="info-card">
            <strong>暂无站内提醒</strong>
            <p class="muted">当前还没有新的订单或退款通知。</p>
          </div>
        `;
        return;
      }

      listBox.innerHTML = notifications.map((item) => {
        const relatedOrderLink = item.relatedOrderId
          ? `<a class="button button-ghost" href="${ROUTES.passenger.orderDetail}?orderId=${item.relatedOrderId}">查看订单</a>`
          : "";

        return `
          <div class="order-card">
            <div class="order-top">
              <div>
                <strong>${escapeHtml(item.title || "系统通知")}</strong>
                <div class="list-meta">
                  <span>${escapeHtml(item.type || "--")}</span>
                  <span>${escapeHtml(formatFullDateTime(item.createdAt))}</span>
                  <span>${item.isRead ? "已读" : "未读"}</span>
                </div>
                <p class="muted">${escapeHtml(item.content || "")}</p>
              </div>
              <span class="badge">${item.isRead ? "已读" : "未读"}</span>
            </div>
            <div class="button-row">
              ${!item.isRead ? `<button class="button button-secondary" type="button" data-notification-read="${item.id}">标记已读</button>` : ""}
              ${relatedOrderLink}
            </div>
          </div>
        `;
      }).join("");

      listBox.querySelectorAll("[data-notification-read]").forEach((button) => {
        button.addEventListener("click", async () => {
          const notificationId = button.getAttribute("data-notification-read");
          if (!notificationId) {
            return;
          }

          try {
            button.disabled = true;
            button.textContent = "处理?..";
            await api.post(API_ENDPOINTS.notification.markRead, {}, {
              pathParams: { notificationId },
            });
            showToast("已标记为已读");
            await loadNotifications();
          } catch (error) {
            button.disabled = false;
            button.textContent = "标记已读";
            showToast(error.message || "操作失败");
          }
        });
      });
    };

    const loadNotifications = async () => {
      const [listResult, unreadResult] = await Promise.all([
        api.get(API_ENDPOINTS.notification.my, { limit: 20 }),
        api.get(API_ENDPOINTS.notification.unreadCount),
      ]);

      const notifications = Array.isArray(listResult?.data) ? listResult.data : [];
      const unreadCount = Number(unreadResult?.data?.unreadCount || 0);

      unreadBox.textContent = `未读 ${unreadCount}`;
      renderList(notifications);
    };

    readAllButton.addEventListener("click", async () => {
      try {
        readAllButton.disabled = true;
        readAllButton.textContent = "处理?..";
        await api.post(API_ENDPOINTS.notification.markAllRead, {});
        showToast("鍏ㄩ通知已标记为已读");
        await loadNotifications();
      } catch (error) {
        showToast(error.message || "操作失败");
      } finally {
        readAllButton.disabled = false;
        readAllButton.textContent = "鍏ㄩ已读";
      }
    });

    loadNotifications().catch((error) => {
      unreadBox.textContent = "未读 --";
      listBox.innerHTML = `
        <div class="info-card">
          <strong>通知加载失败</strong>
          <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
        </div>
      `;
    });
  }

  function initAdminUsersPageV2() {
    if (getCurrentFileName() !== ROUTES.admin.users) {
      return;
    }

    const summaryBox = document.querySelector("[data-admin-user-summary]");
    const filterForm = document.querySelector("[data-admin-user-filter]");
    const resetButton = document.querySelector("[data-admin-user-reset]");
    const listBox = document.querySelector("[data-admin-user-list]");
    const emptyState = document.querySelector("[data-admin-user-empty-state]");
    if (!summaryBox || !filterForm || !listBox || !emptyState) {
      return;
    }

    const maskPhone = (value) => {
      const phone = String(value || "").trim();
      if (phone.length < 7) {
        return phone || "--";
      }
      return `${phone.slice(0, 3)}****${phone.slice(-4)}`;
    };

    const roleOptions = [
      { value: "passenger", label: ROLE_LABELS.passenger || "passenger" },
      { value: "driver", label: ROLE_LABELS.driver || "driver" },
      { value: "admin", label: ROLE_LABELS.admin || "admin" },
    ];

    const statusOptions = [
      { value: "active", label: mapUserStatus("active") },
      { value: "frozen", label: mapUserStatus("frozen") },
      { value: "disabled", label: mapUserStatus("disabled") },
    ];

    const buildSelectOptions = (options, currentValue) => options.map((option) => {
      const selected = option.value === currentValue ? " selected" : "";
      return `<option value="${option.value}"${selected}>${escapeHtml(option.label)}</option>`;
    }).join("");

    const buildQuery = () => {
      const formData = new window.FormData(filterForm);
      return {
        keyword: String(formData.get("keyword") || "").trim(),
        role: String(formData.get("role") || "").trim(),
        status: String(formData.get("status") || "").trim(),
      };
    };

    const renderSummary = (summary) => {
      summaryBox.innerHTML = `
        <span class="mini-chip">用户总数 ${Number(summary?.totalUsers || 0)}</span>
        <span class="mini-chip">${escapeHtml(ROLE_LABELS.passenger || "passenger")} ${Number(summary?.passengerCount || 0)}</span>
        <span class="mini-chip">${escapeHtml(ROLE_LABELS.driver || "driver")} ${Number(summary?.driverCount || 0)}</span>
        <span class="mini-chip">${escapeHtml(ROLE_LABELS.admin || "admin")} ${Number(summary?.adminCount || 0)}</span>
        <span class="mini-chip">正常 ${Number(summary?.activeCount || 0)}</span>
        <span class="mini-chip">冻结 ${Number(summary?.frozenCount || 0)}</span>
        <span class="mini-chip">已实名 ${Number(summary?.verifiedCount || 0)}</span>
      `;
    };

    const renderRealName = (user) => {
      if (user?.realNameVerified && user?.realName) {
        return `${user.realName} / 已实名`;
      }
      if (user?.realNameVerified) {
        return "已实名";
      }
      return "未实名";
    };

    const renderList = (users, reload) => {
      if (!users.length) {
        listBox.innerHTML = `<tr><td colspan="8">暂无用户数据</td></tr>`;
        emptyState.innerHTML = `
          <strong>暂无用户</strong>
          <p class="muted">当前没有符合筛选条件的用户记录。</p>
        `;
        return;
      }

      listBox.innerHTML = users.map((user) => `
        <tr data-admin-user-row="${user.id}" data-current-role="${escapeHtml(user.role || "")}" data-current-status="${escapeHtml(user.status || "")}">
          <td>${escapeHtml(user.nickname || "--")}</td>
          <td>${escapeHtml(maskPhone(user.phone))}</td>
          <td>${escapeHtml(ROLE_LABELS[user.role] || user.role || "--")}</td>
          <td>${escapeHtml(mapUserStatus(user.status))}</td>
          <td>${escapeHtml(user.email || "--")}</td>
          <td>${escapeHtml(renderRealName(user))}</td>
          <td>${escapeHtml(formatFullDateTime(user.createdAt))}</td>
          <td>
            <div class="table-actions">
              <select class="table-inline-select" data-admin-user-role>
                ${buildSelectOptions(roleOptions, user.role)}
              </select>
              <select class="table-inline-select" data-admin-user-status>
                ${buildSelectOptions(statusOptions, user.status)}
              </select>
              <button class="button button-secondary" type="button" data-admin-user-save="${user.id}">保存</button>
            </div>
          </td>
        </tr>
      `).join("");

      emptyState.innerHTML = `
        <strong>用户列表已同步</strong>
        <p class="muted">当前已加载 ${users.length} 条真实用户记录，支持筛选、角色调整和账号状态维护。</p>
      `;

      listBox.querySelectorAll("[data-admin-user-save]").forEach((button) => {
        button.addEventListener("click", async () => {
          const userId = button.getAttribute("data-admin-user-save");
          const row = button.closest("[data-admin-user-row]");
          if (!userId || !row) {
            return;
          }

          const roleSelect = row.querySelector("[data-admin-user-role]");
          const statusSelect = row.querySelector("[data-admin-user-status]");
          const nextRole = String(roleSelect?.value || "").trim();
          const nextStatus = String(statusSelect?.value || "").trim();
          const currentRole = row.getAttribute("data-current-role") || "";
          const currentStatus = row.getAttribute("data-current-status") || "";

          if (nextRole === currentRole && nextStatus === currentStatus) {
            showToast("没有需要保存的修改");
            return;
          }

          try {
            button.disabled = true;
            button.textContent = "保存中...";
            await api.patch(API_ENDPOINTS.admin.updateUser, {
              role: nextRole,
              status: nextStatus,
            }, {
              pathParams: { userId },
            });
            showToast("用户已更新");
            await reload();
          } catch (error) {
            button.disabled = false;
            button.textContent = "保存";
            showToast(error.message || "更新失败");
          }
        });
      });
    };

    const loadUsers = async () => {
      try {
        const [summaryResult, listResult] = await Promise.all([
          api.get(API_ENDPOINTS.admin.userSummary),
          api.get(API_ENDPOINTS.admin.users, buildQuery()),
        ]);
        renderSummary(summaryResult?.data || {});
        renderList(Array.isArray(listResult?.data) ? listResult.data : [], loadUsers);
      } catch (error) {
        summaryBox.innerHTML = `<span class="mini-chip">加载失败</span>`;
        listBox.innerHTML = `<tr><td colspan="8">用户列表加载失败</td></tr>`;
        emptyState.innerHTML = `
          <strong>用户数据加载失败</strong>
          <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
        `;
      }
    };

    filterForm.addEventListener("submit", (event) => {
      event.preventDefault();
      loadUsers();
    });

    if (resetButton) {
      resetButton.addEventListener("click", () => {
        filterForm.reset();
        loadUsers();
      });
    }

    loadUsers();
  }

  function initAdminDashboardPageV3() {
    if (getCurrentFileName() !== ROUTES.admin.dashboard) {
      return;
    }

    const summaryBox = document.querySelector("[data-admin-dashboard-summary]");
    const totalUsersBox = document.querySelector("[data-admin-dashboard-total-users]");
    const activeUsersBox = document.querySelector("[data-admin-dashboard-active-users]");
    const pendingBox = document.querySelector("[data-admin-dashboard-pending-refund]");
    const refundedBox = document.querySelector("[data-admin-dashboard-refunded]");
    const rejectedBox = document.querySelector("[data-admin-dashboard-rejected]");
    const passengerCountBox = document.querySelector("[data-admin-dashboard-passenger-count]");
    const driverCountBox = document.querySelector("[data-admin-dashboard-driver-count]");
    const adminCountBox = document.querySelector("[data-admin-dashboard-admin-count]");
    const pendingNote = document.querySelector("[data-admin-dashboard-pending-note]");
    const emptyState = document.querySelector("[data-admin-dashboard-empty-state]");

    if (!summaryBox || !pendingBox || !refundedBox || !rejectedBox || !emptyState) {
      return;
    }

    const setLoadError = (message) => {
      summaryBox.innerHTML = `<span class="mini-chip">加载失败</span>`;
      if (totalUsersBox) totalUsersBox.textContent = "--";
      if (activeUsersBox) activeUsersBox.textContent = "--";
      if (pendingBox) pendingBox.textContent = "--";
      if (refundedBox) refundedBox.textContent = "--";
      if (rejectedBox) rejectedBox.textContent = "--";
      if (passengerCountBox) passengerCountBox.textContent = "--";
      if (driverCountBox) driverCountBox.textContent = "--";
      if (adminCountBox) adminCountBox.textContent = "--";
      if (pendingNote) pendingNote.textContent = "请稍后重试";
      emptyState.innerHTML = `
        <strong>后台统计加载失败</strong>
        <p class="muted">${escapeHtml(message || "请稍后重试")}</p>
      `;
    };

    api.get(API_ENDPOINTS.admin.dashboard)
      .then((result) => {
        const summary = result?.data || {};
        const pendingRefundCount = Number(summary.pendingRefundCount || 0);
        const refundedCount = Number(summary.refundedCount || 0);
        const rejectedRefundCount = Number(summary.rejectedRefundCount || 0);
        const totalUsers = Number(summary.totalUsers || 0);
        const activeUsers = Number(summary.activeCount || 0);
        const passengerCount = Number(summary.passengerCount || 0);
        const driverCount = Number(summary.driverCount || 0);
        const adminCount = Number(summary.adminCount || 0);
        const frozenCount = Number(summary.frozenCount || 0);
        const disabledCount = Number(summary.disabledCount || 0);
        const verifiedCount = Number(summary.verifiedCount || 0);

        summaryBox.innerHTML = `
          <span class="mini-chip">用户总数 ${totalUsers}</span>
          <span class="mini-chip">正常用户 ${activeUsers}</span>
          <span class="mini-chip">已实名 ${verifiedCount}</span>
          <span class="mini-chip">冻结 ${frozenCount}</span>
          <span class="mini-chip">禁用 ${disabledCount}</span>
          <span class="mini-chip">待退款 ${pendingRefundCount}</span>
        `;
        if (totalUsersBox) totalUsersBox.textContent = String(totalUsers);
        if (activeUsersBox) activeUsersBox.textContent = String(activeUsers);
        if (pendingBox) pendingBox.textContent = String(pendingRefundCount);
        if (refundedBox) refundedBox.textContent = String(refundedCount);
        if (rejectedBox) rejectedBox.textContent = String(rejectedRefundCount);
        if (passengerCountBox) passengerCountBox.textContent = String(passengerCount);
        if (driverCountBox) driverCountBox.textContent = String(driverCount);
        if (adminCountBox) adminCountBox.textContent = String(adminCount);
        if (pendingNote) {
          pendingNote.textContent = pendingRefundCount > 0
            ? `当前还有 ${pendingRefundCount} 笔退款待处理`
            : `当前共有 ${activeUsers} 个正常账号，冻结 ${frozenCount}，禁用 ${disabledCount}`;
        }
        emptyState.innerHTML = pendingRefundCount > 0
          ? `<strong>退款待办提醒</strong><p class="muted">建议优先处理退款审核，避免乘客等待过久。</p>`
          : `<strong>后台运行正常</strong><p class="muted">当前没有待审核退款，已实名用户 ${verifiedCount} 个，可以继续查看用户列表或其他后台模块。</p>`;
      })
      .catch((error) => {
        setLoadError(error.message || "统计加载失败");
    });
  }

  function initDriverDashboardPageV2() {
    if (getCurrentFileName() !== ROUTES.driver.dashboard) {
      return;
    }

    const todayTripsBox = document.querySelector("[data-driver-dashboard-today-trips]");
    const completedTripsBox = document.querySelector("[data-driver-dashboard-completed-trips]");
    const soldTicketsBox = document.querySelector("[data-driver-dashboard-sold-tickets]");
    const seatRateBox = document.querySelector("[data-driver-dashboard-seat-rate]");
    const todayIncomeBox = document.querySelector("[data-driver-dashboard-today-income]");
    const pendingVerifyBox = document.querySelector("[data-driver-dashboard-pending-verify]");
    const refundCountBox = document.querySelector("[data-driver-dashboard-refund-count]");
    const upcomingList = document.querySelector("[data-driver-dashboard-upcoming-list]");
    const alertList = document.querySelector("[data-driver-dashboard-alert-list]");

    if (!todayTripsBox || !upcomingList || !alertList) {
      return;
    }

    api.get(API_ENDPOINTS.driver.dashboard)
      .then((result) => {
        const data = result?.data || {};
        if (todayTripsBox) todayTripsBox.textContent = String(data.todayTripCount || 0);
        if (completedTripsBox) completedTripsBox.textContent = String(data.completedTodayTripCount || 0);
        if (soldTicketsBox) soldTicketsBox.textContent = String(data.soldTicketCount || 0);
        if (seatRateBox) seatRateBox.textContent = `${Math.round(Number(data.seatOccupancyRate || 0) * 100)}%`;
        if (todayIncomeBox) todayIncomeBox.textContent = formatMoneyFromCent(data.todayIncome || 0);
        if (pendingVerifyBox) pendingVerifyBox.textContent = String(data.pendingVerificationCount || 0);
        if (refundCountBox) refundCountBox.textContent = String(data.refundRequestCount || 0);

        const upcomingTrips = Array.isArray(data.upcomingTrips) ? data.upcomingTrips : [];
        upcomingList.innerHTML = upcomingTrips.length
          ? upcomingTrips.map((item) => `
              <div class="list-item">
                <strong>${escapeHtml(item.route || "--")}</strong>
                <div class="list-meta">
                  <span>${escapeHtml(formatFullDateTime(item.departureTime))}</span>
                  <span>已售 ${Number(item.soldTickets || 0)} / ${Number(item.seatTotal || 0)}</span>
                  <span>预计收入 ${escapeHtml(formatMoneyFromCent(item.estimatedIncome || 0))}</span>
                </div>
              </div>
            `).join("")
          : `
              <div class="info-card">
                <strong>未来 4 小时暂无班次</strong>
                <p class="muted">可以前往发布页创建新的班次。</p>
              </div>
            `;

        const alerts = Array.isArray(data.alerts) ? data.alerts : [];
        alertList.innerHTML = alerts.map((text) => `
          <div class="info-card">
            <strong>运营提醒</strong>
            <p class="muted">${escapeHtml(text)}</p>
          </div>
        `).join("");
      })
      .catch((error) => {
        upcomingList.innerHTML = `
          <div class="info-card">
            <strong>班次加载失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
          </div>
        `;
        alertList.innerHTML = `
          <div class="info-card">
            <strong>提醒加载失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
          </div>
        `;
      });
  }

  function initDriverIncomePageV2() {
    if (getCurrentFileName() !== ROUTES.driver.income) {
      return;
    }

    const todayBox = document.querySelector("[data-driver-income-today]");
    const pendingSettleBox = document.querySelector("[data-driver-income-pending-settle]");
    const weekBox = document.querySelector("[data-driver-income-week]");
    const avgOrderBox = document.querySelector("[data-driver-income-avg-order]");
    const refundRateBox = document.querySelector("[data-driver-income-refund-rate]");
    const routeList = document.querySelector("[data-driver-income-route-list]");
    const suggestionList = document.querySelector("[data-driver-income-suggestion-list]");

    if (!todayBox || !routeList || !suggestionList) {
      return;
    }

    api.get(API_ENDPOINTS.driver.income)
      .then((result) => {
        const data = result?.data || {};
        if (todayBox) todayBox.textContent = formatMoneyFromCent(data.todayIncome || 0);
        if (pendingSettleBox) pendingSettleBox.textContent = formatMoneyFromCent(data.pendingSettleAmount || 0);
        if (weekBox) weekBox.textContent = formatMoneyFromCent(data.weeklyIncome || 0);
        if (avgOrderBox) avgOrderBox.textContent = formatMoneyFromCent(data.avgOrderAmount || 0);
        if (refundRateBox) refundRateBox.textContent = `${(Number(data.refundRate || 0) * 100).toFixed(1)}%`;

        const topRoutes = Array.isArray(data.topRoutes) ? data.topRoutes : [];
        routeList.innerHTML = topRoutes.length
          ? topRoutes.map((item) => `
              <div class="list-item">
                <strong>${escapeHtml(item.route || "--")}</strong>
                <div class="list-meta">
                  <span>${escapeHtml(formatMoneyFromCent(item.income || 0))}</span>
                  <span>${Number(item.ticketCount || 0)} 寮?/span>
                  <span>上座率 ${Math.round(Number(item.occupancyRate || 0) * 100)}%</span>
                </div>
              </div>
            `).join("")
          : `
              <div class="info-card">
                <strong>暂无收入数据</strong>
                <p class="muted">当前还没有可统计的司机收入。</p>
              </div>
            `;

        const suggestions = Array.isArray(data.suggestions) ? data.suggestions : [];
        suggestionList.innerHTML = suggestions.map((text) => `
          <div class="info-card">
            <strong>AI 经营建议</strong>
            <p class="muted">${escapeHtml(text)}</p>
          </div>
        `).join("");
      })
      .catch((error) => {
        routeList.innerHTML = `
          <div class="info-card">
            <strong>收入加载失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
          </div>
        `;
        suggestionList.innerHTML = `
          <div class="info-card">
            <strong>建议加载失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
          </div>
        `;
      });
  }

  function initDriverAiPageV2() {
    if (getCurrentFileName() !== ROUTES.driver.ai) {
      return;
    }

    const form = document.querySelector("[data-driver-ai-form]");
    const resultBox = document.querySelector("[data-driver-ai-result]");
    const applyButton = document.querySelector("[data-driver-ai-apply]");
    if (!form || !resultBox || !applyButton) {
      return;
    }

    let currentDraft = null;

    const renderDraft = (draft) => {
      const stops = Array.isArray(draft?.stops) ? draft.stops : [];
      const suggestions = Array.isArray(draft?.suggestions) ? draft.suggestions : [];

      resultBox.innerHTML = `
        <div class="list-item">
          <strong>${escapeHtml(draft.tripName || "--")}</strong>
          <div class="list-meta">
            <span>${escapeHtml(draft.startCity || "--")} -> ${escapeHtml(draft.endCity || "--")}</span>
            <span>${escapeHtml(draft.departureTimeLocal || "--")}</span>
          </div>
        </div>
      <div class="list-item">
        <strong>班次参数</strong>
        <div class="list-meta">
          <span>票价 ${escapeHtml(formatPriceCent(draft.priceCent || 0))}</span>
          <span>座位 ${escapeHtml(String(draft.seatTotal || 0))}</span>
          <span>${escapeHtml(draft.vehicleType || "--")}</span>
        </div>
      </div>
        <div class="list-item">
          <strong>途经站点</strong>
          <div class="list-meta">
            <span>${escapeHtml(stops.length ? stops.join("、") : "无")}</span>
          </div>
        </div>
        <div class="info-card">
          <strong>备注</strong>
          <p class="muted">${escapeHtml(draft.remark || "--")}</p>
        </div>
        ${suggestions.map((item) => `
          <div class="info-card">
            <strong>AI 建议</strong>
            <p class="muted">${escapeHtml(item)}</p>
          </div>
        `).join("")}
      `;
    };

    form.addEventListener("submit", async (event) => {
      event.preventDefault();
      const formData = new window.FormData(form);
      const prompt = String(formData.get("prompt") || "").trim();
      if (!prompt) {
        showToast("请输入班次描述");
        return;
      }

      try {
        applyButton.disabled = true;
        const result = await api.post(API_ENDPOINTS.driver.aiCreateTrip, { prompt });
        currentDraft = result?.data || null;
        if (!currentDraft) {
          throw new Error("AI draft is empty");
        }

        renderDraft(currentDraft);
        applyButton.disabled = false;
        showToast("AI 草稿已生成");
      } catch (error) {
        currentDraft = null;
        resultBox.innerHTML = `
          <div class="info-card">
            <strong>生成失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
          </div>
        `;
        applyButton.disabled = true;
        showToast(error.message || "AI 生成失败");
      }
    });

    applyButton.addEventListener("click", () => {
      if (!currentDraft) {
        return;
      }

      try {
        window.sessionStorage.setItem("tripverse_driver_ai_draft", JSON.stringify(currentDraft));
      } catch (_) {
        // Ignore storage failures.
      }

      redirectTo(ROUTES.driver.publish);
    });
  }

  function normalizeDriverVehicleType(value) {
    const raw = String(value || "").trim();

    if (!raw) {
      return "商务у反";
    }
    if (raw.includes("鎷艰溅")) {
      return "拼车专线";
    }
    if (raw.includes("蹇嚎")) {
      return "鍩庨檯蹇嚎";
    }
    return "商务у反";
  }

  function splitDriverStops(value) {
    const raw = String(value || "").trim();
    if (!raw) {
      return [];
    }

    return raw
      .split(/[銆侊紝,]/)
      .map((item) => item.trim())
      .filter(Boolean);
  }

  function buildDriverTripPayloadFromDraft(draft) {
    const errors = [];
    const warnings = [];

    const startCity = String(draft?.startCity || "").trim();
    const endCity = String(draft?.endCity || "").trim();
    const departureLocal = String(draft?.departureTimeLocal || "").trim();
    let arrivalLocal = String(draft?.arrivalTimeLocal || "").trim();
    const vehicleType = normalizeDriverVehicleType(draft?.vehicleType);
    const seatTotal = Number(draft?.seatTotal || 0);
    const priceCent = Number(draft?.priceCent || 0);

    if (!startCity) {
      errors.push("缺少起点");
    }
    if (!endCity) {
      errors.push("缺少终点");
    }
    if (startCity && endCity && startCity === endCity) {
      errors.push("起点和终点不能相同");
    }
    if (!departureLocal) {
      errors.push("缺少出发时间");
    }

    const departureDate = departureLocal ? new Date(departureLocal) : null;
    if (departureDate && Number.isNaN(departureDate.getTime())) {
      errors.push("出发时间格式不正?");
    }

    let arrivalDate = arrivalLocal ? new Date(arrivalLocal) : null;
    if (arrivalLocal && arrivalDate && Number.isNaN(arrivalDate.getTime())) {
      errors.push("到达时间格式不正?");
    }

    if ((!arrivalLocal || (arrivalDate && Number.isNaN(arrivalDate.getTime()))) && departureDate && !Number.isNaN(departureDate.getTime())) {
      arrivalDate = new Date(departureDate.getTime() + 2 * 60 * 60 * 1000 + 15 * 60 * 1000);
      arrivalLocal = formatLocalInputDateTime(arrivalDate.toISOString());
      warnings.push("未提供有效到达时间，已按出发后 2 小时 15 分自动推算");
    }

    if (departureDate && arrivalDate && !Number.isNaN(departureDate.getTime()) && !Number.isNaN(arrivalDate.getTime()) && arrivalDate.getTime() <= departureDate.getTime()) {
      errors.push("到达时间必须晚于出发时间");
    }

    if (!Number.isFinite(seatTotal) || seatTotal <= 0) {
      errors.push("座位数必须大于 0");
    }
    if (!Number.isFinite(priceCent) || priceCent <= 0) {
      errors.push("票价必须大于 0");
    }

    const stopNames = Array.isArray(draft?.stops)
      ? draft.stops.map((item) => String(item || "").trim()).filter(Boolean)
      : splitDriverStops(draft?.stops);

    const dedupStops = [];
    const seen = new Set();
    for (const stopName of stopNames) {
      if (stopName === startCity || stopName === endCity) {
        warnings.push(`途经站点 ${stopName} 与起终点重复，已自动忽略`);
        continue;
      }
      if (seen.has(stopName)) {
        continue;
      }
      seen.add(stopName);
      dedupStops.push(stopName);
    }

    if (dedupStops.length > 5) {
      warnings.push("途经站点较多，建议控制在 5 个以?");
    }

    return {
      ok: errors.length === 0,
      errors,
      warnings,
      payload: {
        vehicleType,
        startCity,
        endCity,
        departureTime: departureLocal ? toRfc3339FromLocal(departureLocal) : "",
        arrivalTime: arrivalLocal ? toRfc3339FromLocal(arrivalLocal) : "",
        seatTotal,
        priceCent,
        stops: dedupStops.map((stopName, index) => ({
          stopOrder: index + 1,
          stopName,
        })),
      },
    };
  }

  function initDriverAiPageV3() {
    if (getCurrentFileName() !== ROUTES.driver.ai) {
      return;
    }

    const form = document.querySelector("[data-driver-ai-form]");
    const resultBox = document.querySelector("[data-driver-ai-result]");
    const validationBox = document.querySelector("[data-driver-ai-validation]");
    const applyButton = document.querySelector("[data-driver-ai-apply]");
    const publishButton = document.querySelector("[data-driver-ai-publish]");
    if (!form || !resultBox || !validationBox || !applyButton || !publishButton) {
      return;
    }

    let currentDraft = null;

    const renderValidation = (draft) => {
      const check = buildDriverTripPayloadFromDraft(draft);

      if (!check.errors.length && !check.warnings.length) {
        validationBox.innerHTML = `
          <div class="info-card">
            <strong>校验通过</strong>
            <p class="muted">草稿可以直接发布。</p>
          </div>
        `;
      } else {
        const blocks = [];

        if (check.errors.length) {
          blocks.push(`
            <div class="info-card">
              <strong>发布前需修正</strong>
              <p class="muted">${check.errors.map((item) => escapeHtml(item)).join("；")}</p>
            </div>
          `);
        }

        if (check.warnings.length) {
          blocks.push(`
            <div class="info-card">
              <strong>发布提醒</strong>
              <p class="muted">${check.warnings.map((item) => escapeHtml(item)).join("；")}</p>
            </div>
          `);
        }

        validationBox.innerHTML = blocks.join("");
      }

      publishButton.disabled = !check.ok;
    };

    const renderDraft = (draft) => {
      const stops = Array.isArray(draft?.stops) ? draft.stops : splitDriverStops(draft?.stops);
      const suggestions = Array.isArray(draft?.suggestions)
        ? draft.suggestions
        : splitDriverStops(draft?.suggestions);

      resultBox.innerHTML = `
        <div class="list-item">
          <strong>${escapeHtml(draft.tripName || "--")}</strong>
          <div class="list-meta">
            <span>${escapeHtml(draft.startCity || "--")} -> ${escapeHtml(draft.endCity || "--")}</span>
            <span>${escapeHtml(draft.departureTimeLocal || "--")}</span>
          </div>
        </div>
        <div class="list-item">
          <strong>班次参数</strong>
          <div class="list-meta">
            <span>票价 ${escapeHtml(formatPriceCent(draft.priceCent || 0))}</span>
            <span>座位 ${escapeHtml(String(draft.seatTotal || 0))}</span>
            <span>${escapeHtml(draft.vehicleType || "--")}</span>
          </div>
        </div>
        <div class="list-item">
          <strong>途经站点</strong>
          <div class="list-meta">
            <span>${escapeHtml(stops.length ? stops.join("、") : "无")}</span>
          </div>
        </div>
        <div class="info-card">
          <strong>备注</strong>
          <p class="muted">${escapeHtml(draft.remark || "--")}</p>
        </div>
        ${suggestions.map((item) => `
          <div class="info-card">
            <strong>AI 建议</strong>
            <p class="muted">${escapeHtml(item)}</p>
          </div>
        `).join("")}
      `;

      renderValidation(draft);
    };

    form.addEventListener("submit", async (event) => {
      event.preventDefault();
      const formData = new window.FormData(form);
      const prompt = String(formData.get("prompt") || "").trim();
      if (!prompt) {
        showToast("请输入班次描述");
        return;
      }

      try {
        applyButton.disabled = true;
        publishButton.disabled = true;

        const result = await api.post(API_ENDPOINTS.driver.aiCreateTrip, { prompt });
        currentDraft = result?.data || null;
        if (!currentDraft) {
          throw new Error("AI draft is empty");
        }

        renderDraft(currentDraft);
        applyButton.disabled = false;
        showToast("AI 草稿已生成");
      } catch (error) {
        currentDraft = null;
        resultBox.innerHTML = `
          <div class="info-card">
            <strong>生成失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
          </div>
        `;
        validationBox.innerHTML = `
          <div class="info-card">
            <strong>修改建议</strong>
            <p class="muted">你可以换成更标准的描述再试一次，例如：明早 8:30 从杭州东到苏州北，商务大巴，24 座，票价 168，途经嘉兴南、上海虹桥，预计 2 小时 15 分。</p>
          </div>
          <div class="info-card">
            <strong>手工发布建议</strong>
            <p class="muted">如果 AI 持续失败，可以点击“打开发布页”，手工补充起终点、时间、票价和站点。</p>
          </div>
        `;
        applyButton.disabled = true;
        publishButton.disabled = true;
        showToast(error.message || "AI 生成失败");
      }
    });

    applyButton.addEventListener("click", () => {
      if (!currentDraft) {
        return;
      }

      try {
        window.sessionStorage.setItem("tripverse_driver_ai_draft", JSON.stringify(currentDraft));
      } catch (_) {
        // Ignore storage failures.
      }

      redirectTo(ROUTES.driver.publish);
    });

    publishButton.addEventListener("click", async () => {
      if (!currentDraft) {
        return;
      }

      const check = buildDriverTripPayloadFromDraft(currentDraft);
      renderValidation(currentDraft);
      if (!check.ok) {
        showToast(check.errors[0] || "草稿校验未通过");
        return;
      }

      try {
        publishButton.disabled = true;
        publishButton.textContent = "发布?..";

        const result = await api.post(API_ENDPOINTS.driver.trips, check.payload);
        const trip = result?.data || null;

        showToast("班次发布成功");
        window.setTimeout(() => {
          if (trip?.id) {
            redirectTo(`${ROUTES.driver.tripDetail}?tripId=${trip.id}`);
          } else {
            redirectTo(ROUTES.driver.trips);
          }
        }, 300);
      } catch (error) {
        validationBox.innerHTML = `
          <div class="info-card">
            <strong>发布失败</strong>
            <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
          </div>
          <div class="info-card">
            <strong>修改建议</strong>
            <p class="muted">可以先点击“写入发布页”，检查时间、价格和站点后再手工提交。</p>
          </div>
        `;
        showToast(error.message || "发布失败");
      } finally {
        publishButton.textContent = "一键发布";
        renderValidation(currentDraft);
      }
    });
  }

  function hydrateDriverPublishDraftFromAi() {
    if (getCurrentFileName() !== ROUTES.driver.publish) {
      return;
    }

    const form = document.querySelector("[data-driver-publish-form]");
    if (!form) {
      return;
    }

    let draft = null;
    try {
      const raw = window.sessionStorage.getItem("tripverse_driver_ai_draft");
      draft = raw ? JSON.parse(raw) : null;
    } catch (_) {
      draft = null;
    }

    if (!draft) {
      return;
    }

    const setValue = (name, value) => {
      const input = form.querySelector(`[name='${name}']`);
      if (input) {
        input.value = value || "";
      }
    };

    setValue("tripName", draft.tripName || "");
    setValue("depart", formatLocalInputDateTime(draft.departureTimeLocal || ""));
    setValue("arrival", formatLocalInputDateTime(draft.arrivalTimeLocal || ""));
    setValue("start", draft.startCity || "");
    setValue("end", draft.endCity || "");
    setValue("price", draft.priceCent || "");
    setValue("seats", draft.seatTotal || "");
    setValue("vehicleType", draft.vehicleType || "");
    setValue("stops", Array.isArray(draft.stops) ? draft.stops.join("、") : "");
    setValue("remark", draft.remark || "");

    try {
      window.sessionStorage.removeItem("tripverse_driver_ai_draft");
    } catch (_) {
      // Ignore storage failures.
    }

    showToast("AI 草稿已写入发布页");
  }

  function initPassengerAiLegacy() {
    const aiForm = document.querySelector("[data-ai-form]");
    if (!aiForm) {
      return;
    }

    const aiInput = aiForm.querySelector("textarea");
    const aiChat = document.querySelector("[data-ai-chat]");
    if (!aiInput || !aiChat) {
      return;
    }

    const auth = readAuth() || {};
    const aiHistoryKey = `tripverse_passenger_ai_history_${auth.id || auth.userId || auth.phone || "guest"}`;
    const aiDraftKey = `${aiHistoryKey}_draft`;
    const conversation = [];
    const displayHistory = [];

    const readAiHistory = () => {
      try {
        const raw = window.localStorage.getItem(aiHistoryKey);
        const parsed = raw ? JSON.parse(raw) : null;
        return Array.isArray(parsed?.messages) ? parsed.messages : [];
      } catch (_) {
        return [];
      }
    };

    const saveAiHistory = () => {
      try {
        window.localStorage.setItem(aiHistoryKey, JSON.stringify({
          messages: displayHistory.slice(-20),
          savedAt: new Date().toISOString(),
        }));
      } catch (_) {
        // Ignore storage failures.
      }
    };

    const saveAiDraft = () => {
      try {
        window.localStorage.setItem(aiDraftKey, aiInput.value || "");
      } catch (_) {
        // Ignore storage failures.
      }
    };

    const appendMessage = (role, title, content, trustedHtml = false, persist = true) => {
      const node = document.createElement("div");
      node.className = `message ${role}`;
      node.innerHTML = trustedHtml
        ? `<strong>${title}</strong><div>${content}</div>`
        : `<strong>${title}</strong><div>${escapeHtml(content)}</div>`;
      aiChat.appendChild(node);
      aiChat.scrollTop = aiChat.scrollHeight;

      if (persist) {
        displayHistory.push({ role, title, content, trustedHtml: Boolean(trustedHtml) });
        saveAiHistory();
      }
    };

    const restoreAiState = () => {
      const savedMessages = readAiHistory();
      savedMessages.forEach((item) => {
        if (!item || !item.role || !item.title) {
          return;
        }
        displayHistory.push({
          role: item.role,
          title: item.title,
          content: item.content || "",
          trustedHtml: Boolean(item.trustedHtml),
        });
        appendMessage(item.role, item.title, item.content || "", Boolean(item.trustedHtml), false);

        if (item.role === "user") {
          conversation.push({ role: "user", content: String(item.content || "") });
        } else if (item.role === "ai") {
          const plainText = String(item.content || "").replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();
          if (plainText) {
            conversation.push({ role: "assistant", content: plainText });
          }
        }
      });

      try {
        aiInput.value = window.localStorage.getItem(aiDraftKey) || "";
      } catch (_) {
        // Ignore storage failures.
      }
    };

    aiForm.addEventListener("submit", async (event) => {
      event.preventDefault();

      const value = aiInput.value.trim();
      if (!value) {
        return;
      }

      appendMessage("user", "你", value);
      conversation.push({ role: "user", content: value });
      aiInput.value = "";

      try {
        const result = await api.post(API_ENDPOINTS.passenger.aiChat, {
          messages: conversation.slice(-8),
        });

        const data = result?.data || {};
        const reply = String(data.reply || "").trim();
        const suggestions = Array.isArray(data.suggestions) ? data.suggestions : [];

        let html = escapeHtml(reply || "暂时没有生成回复");
        if (suggestions.length) {
          html += `<div class="section-block">${suggestions
            .map((item) => `<div class="muted">${escapeHtml(item)}</div>`)
            .join("")}</div>`;
        }

        appendMessage("ai", "AI 助手", html, true);
        if (reply) {
          conversation.push({ role: "assistant", content: reply });
        }
      } catch (error) {
        appendMessage("ai", "AI 助手", error.message || "请求失败，请稍后重试");
      }
    });
  }

    function renderPassengerAiRouteCards(context) {
    const routeQuery = context?.routeQuery || null;
    const routeResults = Array.isArray(context?.routeResults) ? context.routeResults : [];

    if (!routeQuery && !routeResults.length) {
      return "";
    }

    const header = routeQuery
      ? `<div class="info-card">
          <strong>查询条件</strong>
          <p class="muted">${escapeHtml(routeQuery.startCity || "--")} -> ${escapeHtml(routeQuery.endCity || "--")}，${escapeHtml(routeQuery.date || (routeQuery.dateFrom && routeQuery.dateTo ? `${routeQuery.dateFrom} 至 ${routeQuery.dateTo}` : "--"))}${routeQuery.allowTransfer ? "，允许一次中转" : ""}</p>
        </div>`
      : "";

    const cards = routeResults.length
      ? routeResults.slice(0, 4).map((item) => {
          if (item.kind === "suggestion") {
            const suggestions = Array.isArray(item.suggestions) ? item.suggestions : [];
            return `
              <div class="info-card">
                <strong>可尝试的候选路线</strong>
                <p class="muted">当前没有查到可直接购买的直达或精确中转班次，但当天仍有一些可尝试的中转方向。</p>
                <div class="list-stack section-block">
                  ${suggestions.length
                    ? suggestions.map((suggestion, index) => `
                        <div class="info-card">
                          <strong>方案 ${index + 1}：${escapeHtml(suggestion.route || "--")}</strong>
                          <p class="muted">${escapeHtml(suggestion.reason || "--")}</p>
                          <p class="muted">中转城市：${escapeHtml(suggestion.transferCity || "--")}，前段候选 ${Number(suggestion.firstLegCount || 0)} 班，后段候选 ${Number(suggestion.secondLegCount || 0)} 班</p>
                        </div>
                      `).join("")
                    : `<div class="info-card"><p class="muted">暂无候选路线建议。</p></div>`}
                </div>
              </div>
            `;
          }

          const legs = Array.isArray(item.legs) ? item.legs : [];
          const transferMeta = item.kind === "transfer"
            ? `<p class="muted">经 ${escapeHtml(item.transferCity || "--")} 中转，候车约 ${Number(item.transferWaitMinute || 0)} 分钟</p>`
            : "";
          const legsHtml = item.kind === "transfer" && legs.length
            ? `<div class="list-stack section-block">
                ${legs.map((leg, index) => `
                  <div class="info-card">
                    <strong>第 ${index + 1} 段：${escapeHtml(leg.route || "--")}</strong>
                    <p class="muted">${escapeHtml(leg.departureTime || "--")} - ${escapeHtml(leg.arrivalTime || "--")}</p>
                    <p class="muted">${escapeHtml(leg.vehicleType || "--")}，票价 ${escapeHtml(formatMoneyFromCent(leg.priceCent || 0))}，余座 ${Number(leg.seatAvailable || 0)}</p>
                    <div class="button-row section-block">
                      <a class="button button-secondary" href="${buildTicketDetailHref(leg.tripId)}">查看第 ${index + 1} 段</a>
                      <a class="button button-primary" href="${buildCheckoutHref(leg.tripId)}">购买第 ${index + 1} 段</a>
                    </div>
                  </div>
                `).join("")}
              </div>`
            : "";
          const actions = item.kind === "transfer" && legs.length
            ? legs.map((leg, index) => `
                <a class="button button-secondary" href="${buildTicketDetailHref(leg.tripId)}">查看第 ${index + 1} 段</a>
                <a class="button button-primary" href="${buildCheckoutHref(leg.tripId)}">购买第 ${index + 1} 段</a>
              `).join("")
            : `
                <a class="button button-secondary" href="${buildTicketDetailHref(item.id)}">查看详情</a>
                <a class="button button-primary" href="${buildCheckoutHref(item.id)}">去下单</a>
              `;
          const stationLinks = item.kind !== "transfer" ? renderStationLinks(item.stops, item.id) : "";

          return `
            <div class="info-card">
              <strong>${item.kind === "transfer" ? "一次中转方案" : "直达方案"}：${escapeHtml(item.route || "--")}</strong>
              <p class="muted">${escapeHtml(item.departureTime || "--")} - ${escapeHtml(item.arrivalTime || "--")}</p>
              <p class="muted">余座 ${Number(item.seatAvailable || 0)}，票价 ${escapeHtml(formatMoneyFromCent(item.priceCent || 0))}</p>
              <p class="muted">${escapeHtml(item.vehicleType || "--")}</p>
              ${stationLinks}
              ${transferMeta}
              ${legsHtml}
              <div class="button-row section-block">
                ${actions}
              </div>
            </div>
          `;
        }).join("")
      : `
        <div class="info-card">
          <strong>暂无可售班次</strong>
          <p class="muted">当前条件下没有查到可售班次，建议换日期或确认出发/到达城市。</p>
        </div>
      `;

    return `<div class="section-block">${header}${cards}</div>`;
  }

  function renderPassengerAiOrderSummary(context) {
    const summary = context?.orderSummary || null;
    if (!summary) {
      return "";
    }

    return `
      <div class="info-card">
        <strong>订单摘要</strong>
        <p class="muted">
          总订单 ${Number(summary.totalCount || 0)}，
          待支付 ${Number(summary.pendingPaymentCount || 0)}，
          待核销 ${Number(summary.pendingVerificationCount || 0)}，
          已完成 ${Number(summary.completedCount || 0)}，
          退款申请 ${Number(summary.refundRequestedCount || 0)}
        </p>
      </div>
    `;
  }

  function renderPassengerAiOrderCards(context) {
    const orders = Array.isArray(context?.orderResults) ? context.orderResults : [];
    if (!orders.length) {
      return "";
    }

    return `
      <div class="section-block">
        ${orders.slice(0, 4).map((item) => `
          <div class="info-card">
            <strong>${escapeHtml(item.orderNo || `订单 #${item.id || "--"}`)}</strong>
            <p class="muted">${escapeHtml(item.route || "--")} ${escapeHtml(item.departureTime || "")}</p>
            <p class="muted">订单状态：${escapeHtml(item.orderStatus || "--")}，支付状态：${escapeHtml(item.payStatus || "--")}</p>
            <p class="muted">退款状态：${escapeHtml(item.refundStatus || "--")}，金额：${escapeHtml(formatMoneyFromCent(item.amount || 0))}</p>
            ${item.refundReviewNote ? `<p class="muted">审核备注：${escapeHtml(item.refundReviewNote)}</p>` : ""}
            <div class="button-row section-block">
              <a class="button button-secondary" href="${ROUTES.passenger.orderDetail}?orderId=${encodeURIComponent(item.id || "")}">查看订单</a>
            </div>
          </div>
        `).join("")}
      </div>
    `;
  }

  function renderPassengerAiRefundRules(context) {
    const refundRules = Array.isArray(context?.refundRules) ? context.refundRules : [];
    if (!refundRules.length) {
      return "";
    }

    return `
      <div class="section-block">
        ${refundRules.map((item) => `
          <div class="info-card">
            <strong>退款规则</strong>
            <p class="muted">${escapeHtml(item)}</p>
          </div>
        `).join("")}
      </div>
    `;
  }

  function renderPassengerAiSystemHints(context) {
    const hints = Array.isArray(context?.systemHints) ? context.systemHints : [];
    if (!hints.length) {
      return "";
    }

    return `
      <div class="section-block">
        ${hints.map((item) => `
          <div class="info-card">
            <strong>补充提示</strong>
            <p class="muted">${escapeHtml(item)}</p>
          </div>
        `).join("")}
      </div>
    `;
  }

  function renderPassengerAiSuggestions(suggestions) {
    const items = Array.isArray(suggestions) ? suggestions.filter(Boolean) : [];
    if (!items.length) {
      return "";
    }

    return `
      <div class="section-block">
        <strong>建议追问</strong>
        <div class="button-row section-block">
          ${items.map((item) => `
            <button class="button button-ghost" type="button" data-ai-suggestion="${escapeHtml(item)}">${escapeHtml(item)}</button>
          `).join("")}
        </div>
      </div>
    `;
  }

  function renderPassengerAiActions(intent, context) {
    if (intent === "route") {
      return `
        <div class="button-row section-block">
          <a class="button button-primary" href="${ROUTES.passenger.search}">去购票</a>
        </div>
      `;
    }

    if (intent === "orders") {
      return `
        <div class="button-row section-block">
          <a class="button button-primary" href="${ROUTES.passenger.orders}">查看订单列表</a>
        </div>
      `;
    }

    if (intent === "refund") {
      const firstOrder = Array.isArray(context?.orderResults) ? context.orderResults[0] : null;
      if (firstOrder?.id) {
        return `
          <div class="button-row section-block">
            <a class="button button-primary" href="${ROUTES.passenger.orderDetail}?orderId=${encodeURIComponent(firstOrder.id)}">查看退款相关订单</a>
          </div>
        `;
      }
      return `
        <div class="button-row section-block">
          <a class="button button-primary" href="${ROUTES.passenger.orders}">去订单页查看退款状态</a>
        </div>
      `;
    }

    return "";
  }

  function buildPassengerAiRichHtml(data) {
    const reply = String(data?.reply || "").trim() || "暂时没有生成回复";
    const intent = String(data?.intent || "");
    const context = data?.context || null;
    const suggestions = Array.isArray(data?.suggestions) ? data.suggestions : [];

    return [
      `<div>${escapeHtml(reply)}</div>`,
      renderPassengerAiRouteCards(context),
      renderPassengerAiOrderSummary(context),
      renderPassengerAiOrderCards(context),
      renderPassengerAiRefundRules(context),
      renderPassengerAiSystemHints(context),
      renderPassengerAiSuggestions(suggestions),
      renderPassengerAiActions(intent, context),
    ].join("");
  }

  function buildPassengerAiErrorHtml(error) {
    const message = String(error?.message || "请求失败，请稍后重试").trim();

    if (message.includes("unauthorized") || message.includes("not logged in") || message.includes("Authorization")) {
      return `
        <div>请先登录，再查询个人订单或退款信息。</div>
        <div class="section-block">
          <a class="button button-primary" href="${ROUTES.auth.login}">去登录</a>
        </div>
      `;
    }

    if (message.includes("startCity") || message.includes("endCity") || message.includes("date")) {
      return `
        <div>请补充出发地、目的地和日期，例如：帮我查明天杭州到苏州的票。</div>
        <div class="section-block">
          <button class="button button-ghost" type="button" data-ai-suggestion="帮我查明天杭州到苏州的票">使用示例追问</button>
        </div>
      `;
    }

    return `
      <div>${escapeHtml(message)}</div>
      <div class="section-block">
        <div class="info-card">
          <strong>回退建议</strong>
          <p class="muted">可以换一种更完整的问法，例如：帮我查明天杭州到苏州的票；帮我看下我的订单；我这笔订单能退款吗。</p>
        </div>
      </div>
    `;
  }

  function initPassengerAi() {
    const aiForm = document.querySelector("[data-ai-form]");
    if (!aiForm) {
      return;
    }

    const aiInput = aiForm.querySelector("textarea");
    const aiChat = document.querySelector("[data-ai-chat]");
    if (!aiInput || !aiChat) {
      return;
    }

    const auth = readAuth() || {};
    const aiUserKey = auth.id || auth.userId || auth.phone || "guest";
    const aiSessionsKey = `tripverse_passenger_ai_sessions_${aiUserKey}`;
    const fallbackHistoryKey = `tripverse_passenger_ai_history_${aiUserKey}`;
    const fallbackDraftKey = `${fallbackHistoryKey}_draft`;
    const initialGreeting = aiChat.innerHTML;
    let sessionState = null;

    const createSession = (title) => ({
      id: `ai_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
      title: title || "新窗口",
      messages: [],
      draft: "",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    });

    const readSessions = () => {
      try {
        const raw = window.localStorage.getItem(aiSessionsKey);
        const parsed = raw ? JSON.parse(raw) : null;
        if (parsed && Array.isArray(parsed.sessions) && parsed.sessions.length) {
          const activeId = parsed.sessions.some((session) => session.id === parsed.activeId)
            ? parsed.activeId
            : parsed.sessions[0].id;
          return { activeId, sessions: parsed.sessions };
        }
      } catch (_) {
        // Fall through to migration/default.
      }

      const migratedSession = createSession("历史窗口");
      try {
        const oldRaw = window.localStorage.getItem(fallbackHistoryKey);
        const oldParsed = oldRaw ? JSON.parse(oldRaw) : null;
        migratedSession.messages = Array.isArray(oldParsed?.messages) ? oldParsed.messages : [];
        migratedSession.draft = window.localStorage.getItem(fallbackDraftKey) || "";
      } catch (_) {
        // Ignore old history migration failures.
      }

      const firstSession = migratedSession.messages.length || migratedSession.draft ? migratedSession : createSession("新窗口");
      return {
        activeId: firstSession.id,
        sessions: [firstSession],
      };
    };

    const saveSessions = () => {
      try {
        window.localStorage.setItem(aiSessionsKey, JSON.stringify(sessionState));
      } catch (_) {
        // Ignore storage failures.
      }
    };

    const getActiveSession = () => {
      if (!sessionState.sessions.length) {
        const session = createSession("新窗口");
        sessionState.sessions.push(session);
        sessionState.activeId = session.id;
      }
      const activeSession = sessionState.sessions.find((session) => session.id === sessionState.activeId) || sessionState.sessions[0];
      sessionState.activeId = activeSession.id;
      return activeSession;
    };

    const buildConversation = (session) => {
      const conversation = [];
      (session.messages || []).forEach((item) => {
        if (item.role === "user") {
          conversation.push({ role: "user", content: String(item.content || "") });
        } else if (item.role === "ai") {
          const plainText = String(item.content || "").replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();
          if (plainText) {
            conversation.push({ role: "assistant", content: plainText });
          }
        }
      });
      return conversation;
    };

    const sessionShell = document.createElement("div");
    sessionShell.className = "ai-session-shell";
    sessionShell.innerHTML = `
      <div class="ai-session-tabs" data-ai-session-tabs></div>
      <div class="button-row ai-session-actions">
        <button class="button button-secondary" type="button" data-ai-session-new>新建窗口</button>
        <button class="button button-ghost" type="button" data-ai-session-delete>删除当前窗口</button>
      </div>
    `;
    aiChat.parentNode.insertBefore(sessionShell, aiChat);
    const sessionTabs = sessionShell.querySelector("[data-ai-session-tabs]");

    const renderSessionTabs = () => {
      const activeSession = getActiveSession();
      sessionTabs.innerHTML = sessionState.sessions.map((session, index) => {
        const activeClass = session.id === activeSession.id ? " is-active" : "";
        const title = session.title || `窗口 ${index + 1}`;
        return `<button class="ai-session-tab${activeClass}" type="button" data-ai-session-id="${escapeHtml(session.id)}">${escapeHtml(title)}</button>`;
      }).join("");
    };

    const renderMessage = (role, title, content, trustedHtml = false) => {
      const node = document.createElement("div");
      node.className = `message ${role}`;
      node.innerHTML = trustedHtml
        ? `<strong>${escapeHtml(title)}</strong><div>${content}</div>`
        : `<strong>${escapeHtml(title)}</strong><div>${escapeHtml(content)}</div>`;
      aiChat.appendChild(node);
      aiChat.scrollTop = aiChat.scrollHeight;
    };

    const renderActiveSession = () => {
      const activeSession = getActiveSession();
      aiChat.innerHTML = initialGreeting;
      (activeSession.messages || []).forEach((item) => {
        renderMessage(item.role, item.title, item.content || "", Boolean(item.trustedHtml));
      });
      aiInput.value = activeSession.draft || "";
      renderSessionTabs();
    };

    const updateSessionTitle = (session, firstQuestion) => {
      if (!session || !firstQuestion || session.messages.filter((item) => item.role === "user").length > 1) {
        return;
      }
      const text = String(firstQuestion || "").trim();
      session.title = text.length > 12 ? `${text.slice(0, 12)}...` : text || "新窗口";
    };

    const appendMessageToSession = (session, role, title, content, trustedHtml = false) => {
      const message = { role, title, content, trustedHtml: Boolean(trustedHtml), createdAt: new Date().toISOString() };
      session.messages = [...(session.messages || []), message].slice(-40);
      session.updatedAt = new Date().toISOString();
      if (role === "user") {
        updateSessionTitle(session, content);
      }
    };

    const appendMessage = (role, title, content, trustedHtml = false) => {
      const activeSession = getActiveSession();
      appendMessageToSession(activeSession, role, title, content, trustedHtml);
      saveSessions();
      renderMessage(role, title, content, trustedHtml);
      renderSessionTabs();
    };

    const saveAiDraft = () => {
      const activeSession = getActiveSession();
      activeSession.draft = aiInput.value || "";
      activeSession.updatedAt = new Date().toISOString();
      saveSessions();
    };

    const askAi = async (value) => {
      const activeSession = getActiveSession();
      const activeId = activeSession.id;
      appendMessage("user", "你", value);
      const conversation = buildConversation(getActiveSession());

      const typingNode = document.createElement("div");
      typingNode.className = "message ai";
      typingNode.innerHTML = `<strong>AI 助手</strong><div>正在查询中...</div>`;
      aiChat.appendChild(typingNode);
      aiChat.scrollTop = aiChat.scrollHeight;

      try {
        const result = await api.post(API_ENDPOINTS.passenger.aiChat, {
          messages: conversation.slice(-8),
        });

        const data = result?.data || {};
        if (typingNode.parentNode === aiChat) {
          aiChat.removeChild(typingNode);
        }

          appendMessage("ai", "AI 助手", buildPassengerAiRichHtml(data, value), true);

        const reply = String(data?.reply || "").trim();
        if (reply) {
          conversation.push({ role: "assistant", content: reply });
        }
      } catch (error) {
        if (typingNode.parentNode === aiChat) {
          aiChat.removeChild(typingNode);
        }
        appendMessage("ai", "AI 助手", buildPassengerAiErrorHtml(error), true);
      }
    };

    aiForm.addEventListener("submit", async (event) => {
      event.preventDefault();

      const value = aiInput.value.trim();
      if (!value) {
        return;
      }

      aiInput.value = "";
      saveAiDraft();
      await askAi(value);
    });

    aiInput.addEventListener("input", saveAiDraft);

    aiChat.addEventListener("click", async (event) => {
      const button = event.target.closest("[data-ai-suggestion]");
      if (!button) {
        return;
      }

      const value = String(button.getAttribute("data-ai-suggestion") || "").trim();
      if (!value) {
        return;
      }

      if (aiInput) {
        aiInput.value = value;
        saveAiDraft();
      }
      await askAi(value);
    });

    restoreAiState();
  }

    function renderPassengerAiRouteCards(context) {
    const routeQuery = context?.routeQuery || null;
    const routeResults = Array.isArray(context?.routeResults) ? context.routeResults : [];

    if (!routeQuery && !routeResults.length) {
      return "";
    }

    const header = routeQuery
      ? `<div class="info-card">
          <strong>查询条件</strong>
          <p class="muted">${escapeHtml(routeQuery.startCity || "--")} -> ${escapeHtml(routeQuery.endCity || "--")}，${escapeHtml(routeQuery.date || (routeQuery.dateFrom && routeQuery.dateTo ? `${routeQuery.dateFrom} 至 ${routeQuery.dateTo}` : "--"))}${routeQuery.allowTransfer ? "，允许一次中转" : ""}</p>
        </div>`
      : "";

    const cards = routeResults.length
      ? routeResults.slice(0, 4).map((item) => {
          if (item.kind === "suggestion") {
            const suggestions = Array.isArray(item.suggestions) ? item.suggestions : [];
            return `
              <div class="info-card">
                <strong>可尝试的候选路线</strong>
                <p class="muted">当前没有查到可直接购买的直达或精确中转班次，但当天仍有一些可尝试的中转方向。</p>
                <div class="list-stack section-block">
                  ${suggestions.length
                    ? suggestions.map((suggestion, index) => `
                        <div class="info-card">
                          <strong>方案 ${index + 1}：${escapeHtml(suggestion.route || "--")}</strong>
                          <p class="muted">${escapeHtml(suggestion.reason || "--")}</p>
                          <p class="muted">中转城市：${escapeHtml(suggestion.transferCity || "--")}，前段候选 ${Number(suggestion.firstLegCount || 0)} 班，后段候选 ${Number(suggestion.secondLegCount || 0)} 班</p>
                        </div>
                      `).join("")
                    : `<div class="info-card"><p class="muted">暂无候选路线建议。</p></div>`}
                </div>
              </div>
            `;
          }

          const legs = Array.isArray(item.legs) ? item.legs : [];
          const transferMeta = item.kind === "transfer"
            ? `<p class="muted">经 ${escapeHtml(item.transferCity || "--")} 中转，候车约 ${Number(item.transferWaitMinute || 0)} 分钟</p>`
            : "";
          const legsHtml = item.kind === "transfer" && legs.length
            ? `<div class="list-stack section-block">
                ${legs.map((leg, index) => `
                  <div class="info-card">
                    <strong>第 ${index + 1} 段：${escapeHtml(leg.route || "--")}</strong>
                    <p class="muted">${escapeHtml(leg.departureTime || "--")} - ${escapeHtml(leg.arrivalTime || "--")}</p>
                    <p class="muted">${escapeHtml(leg.vehicleType || "--")}，票价 ${escapeHtml(formatMoneyFromCent(leg.priceCent || 0))}，余座 ${Number(leg.seatAvailable || 0)}</p>
                    <div class="button-row section-block">
                      <a class="button button-secondary" href="${buildTicketDetailHref(leg.tripId)}">查看第 ${index + 1} 段</a>
                      <a class="button button-primary" href="${buildCheckoutHref(leg.tripId)}">购买第 ${index + 1} 段</a>
                    </div>
                  </div>
                `).join("")}
              </div>`
            : "";
          const actions = item.kind === "transfer" && legs.length
            ? legs.map((leg, index) => `
                <a class="button button-secondary" href="${buildTicketDetailHref(leg.tripId)}">查看第 ${index + 1} 段</a>
                <a class="button button-primary" href="${buildCheckoutHref(leg.tripId)}">购买第 ${index + 1} 段</a>
              `).join("")
            : `
                <a class="button button-secondary" href="${buildTicketDetailHref(item.id)}">查看详情</a>
                <a class="button button-primary" href="${buildCheckoutHref(item.id)}">去下单</a>
              `;
          const stationLinks = item.kind !== "transfer" ? renderStationLinks(item.stops, item.id) : "";

          return `
            <div class="info-card">
              <strong>${item.kind === "transfer" ? "一次中转方案" : "直达方案"}：${escapeHtml(item.route || "--")}</strong>
              <p class="muted">${escapeHtml(item.departureTime || "--")} - ${escapeHtml(item.arrivalTime || "--")}</p>
              <p class="muted">余座 ${Number(item.seatAvailable || 0)}，票价 ${escapeHtml(formatMoneyFromCent(item.priceCent || 0))}</p>
              <p class="muted">${escapeHtml(item.vehicleType || "--")}</p>
              ${stationLinks}
              ${transferMeta}
              ${legsHtml}
              <div class="button-row section-block">
                ${actions}
              </div>
            </div>
          `;
        }).join("")
      : `
        <div class="info-card">
          <strong>暂无可售班次</strong>
          <p class="muted">当前条件下没有查到可售班次，建议换日期或确认出发/到达城市。</p>
        </div>
      `;

    return `<div class="section-block">${header}${cards}</div>`;
  }

  function renderPassengerAiOrderSummary(context) {
    const summary = context?.orderSummary || null;
    if (!summary) {
      return "";
    }

    return `
      <div class="info-card">
        <strong>订单摘要</strong>
        <p class="muted">
          总订单 ${Number(summary.totalCount || 0)}，
          待支付 ${Number(summary.pendingPaymentCount || 0)}，
          待核销 ${Number(summary.pendingVerificationCount || 0)}，
          已完成 ${Number(summary.completedCount || 0)}，
          退款申请 ${Number(summary.refundRequestedCount || 0)}
        </p>
      </div>
    `;
  }

  function renderPassengerAiOrderCards(context) {
    const orders = Array.isArray(context?.orderResults) ? context.orderResults : [];
    if (!orders.length) {
      return "";
    }

    return `
      <div class="section-block">
        ${orders.slice(0, 4).map((item) => `
          <div class="info-card">
            <strong>${escapeHtml(item.orderNo || `订单 #${item.id || "--"}`)}</strong>
            <p class="muted">${escapeHtml(item.route || "--")} ${escapeHtml(item.departureTime || "")}</p>
            <p class="muted">订单状态：${escapeHtml(item.orderStatus || "--")}，支付状态：${escapeHtml(item.payStatus || "--")}</p>
            <p class="muted">退款状态：${escapeHtml(item.refundStatus || "--")}，金额：${escapeHtml(formatMoneyFromCent(item.amount || 0))}</p>
            ${item.refundReviewNote ? `<p class="muted">审核备注：${escapeHtml(item.refundReviewNote)}</p>` : ""}
            <div class="button-row section-block">
              <a class="button button-secondary" href="${ROUTES.passenger.orderDetail}?orderId=${encodeURIComponent(item.id || "")}">查看订单</a>
            </div>
          </div>
        `).join("")}
      </div>
    `;
  }

  function renderPassengerAiRefundRules(context) {
    const refundRules = Array.isArray(context?.refundRules) ? context.refundRules : [];
    if (!refundRules.length) {
      return "";
    }

    return `
      <div class="section-block">
        ${refundRules.map((item) => `
          <div class="info-card">
            <strong>退款规则</strong>
            <p class="muted">${escapeHtml(item)}</p>
          </div>
        `).join("")}
      </div>
    `;
  }

  function renderPassengerAiSystemHints(context) {
    const hints = Array.isArray(context?.systemHints) ? context.systemHints : [];
    if (!hints.length) {
      return "";
    }

    return `
      <div class="section-block">
        ${hints.map((item) => `
          <div class="info-card">
            <strong>补充提示</strong>
            <p class="muted">${escapeHtml(item)}</p>
          </div>
        `).join("")}
      </div>
    `;
  }

  function renderPassengerAiSuggestions(suggestions) {
    const items = Array.isArray(suggestions) ? suggestions.filter(Boolean) : [];
    if (!items.length) {
      return "";
    }

    return `
      <div class="section-block">
        <strong>建议追问</strong>
        <div class="button-row section-block">
          ${items.map((item) => `
            <button class="button button-ghost" type="button" data-ai-suggestion="${escapeHtml(item)}">${escapeHtml(item)}</button>
          `).join("")}
        </div>
      </div>
    `;
  }

  function renderPassengerAiActions(intent, context) {
    if (intent === "route") {
      return `
        <div class="button-row section-block">
          <a class="button button-primary" href="${ROUTES.passenger.search}">去购票</a>
        </div>
      `;
    }

    if (intent === "orders") {
      return `
        <div class="button-row section-block">
          <a class="button button-primary" href="${ROUTES.passenger.orders}">查看订单列表</a>
        </div>
      `;
    }

    if (intent === "refund") {
      const firstOrder = Array.isArray(context?.orderResults) ? context.orderResults[0] : null;
      if (firstOrder?.id) {
        return `
          <div class="button-row section-block">
            <a class="button button-primary" href="${ROUTES.passenger.orderDetail}?orderId=${encodeURIComponent(firstOrder.id)}">查看退款相关订单</a>
          </div>
        `;
      }
      return `
        <div class="button-row section-block">
          <a class="button button-primary" href="${ROUTES.passenger.orders}">去订单页查看退款状态</a>
        </div>
      `;
    }

    return "";
  }

  function buildPassengerAiRichHtml(data) {
    const reply = String(data?.reply || "").trim() || "暂时没有生成回复";
    const intent = String(data?.intent || "");
    const context = data?.context || null;
    const suggestions = Array.isArray(data?.suggestions) ? data.suggestions : [];

    return [
      `<div>${escapeHtml(reply)}</div>`,
      renderPassengerAiRouteCards(context),
      renderPassengerAiOrderSummary(context),
      renderPassengerAiOrderCards(context),
      renderPassengerAiRefundRules(context),
      renderPassengerAiSystemHints(context),
      renderPassengerAiSuggestions(suggestions),
      renderPassengerAiActions(intent, context),
    ].join("");
  }

  function buildPassengerAiErrorHtml(error) {
    const message = String(error?.message || "请求失败，请稍后重试").trim();

    if (message.includes("unauthorized") || message.includes("not logged in") || message.includes("Authorization")) {
      return `
        <div>请先登录，再查询个人订单或退款信息。</div>
        <div class="section-block">
          <a class="button button-primary" href="${ROUTES.auth.login}">去登录</a>
        </div>
      `;
    }

    if (message.includes("startCity") || message.includes("endCity") || message.includes("date")) {
      return `
        <div>请补充出发地、目的地和日期，例如：帮我查明天杭州到苏州的票。</div>
        <div class="section-block">
          <button class="button button-ghost" type="button" data-ai-suggestion="帮我查明天杭州到苏州的票">使用示例追问</button>
        </div>
      `;
    }

    return `
      <div>${escapeHtml(message)}</div>
      <div class="section-block">
        <div class="info-card">
          <strong>回退建议</strong>
          <p class="muted">可以换一种更完整的问法，例如：帮我查明天杭州到苏州的票；帮我看下我的订单；我这笔订单能退款吗。</p>
        </div>
      </div>
    `;
  }

  function initPassengerAi() {
    const aiForm = document.querySelector("[data-ai-form]");
    if (!aiForm) {
      return;
    }

    const aiInput = aiForm.querySelector("textarea");
    const aiChat = document.querySelector("[data-ai-chat]");
    if (!aiInput || !aiChat) {
      return;
    }

    const auth = readAuth() || {};
    const aiUserKey = auth.id || auth.userId || auth.phone || "guest";
    const aiSessionsKey = `tripverse_passenger_ai_sessions_${aiUserKey}`;
    const fallbackHistoryKey = `tripverse_passenger_ai_history_${aiUserKey}`;
    const fallbackDraftKey = `${fallbackHistoryKey}_draft`;
    const initialGreeting = aiChat.innerHTML;
    let sessionState = null;

    const createSession = (title) => ({
      id: `ai_${Date.now()}_${Math.random().toString(36).slice(2, 8)}`,
      title: title || "新窗口",
      messages: [],
      draft: "",
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    });

    const readSessions = () => {
      try {
        const raw = window.localStorage.getItem(aiSessionsKey);
        const parsed = raw ? JSON.parse(raw) : null;
        if (parsed && Array.isArray(parsed.sessions) && parsed.sessions.length) {
          const activeId = parsed.sessions.some((session) => session.id === parsed.activeId)
            ? parsed.activeId
            : parsed.sessions[0].id;
          return { activeId, sessions: parsed.sessions };
        }
      } catch (_) {
        // Fall through to migration/default.
      }

      const migratedSession = createSession("历史窗口");
      try {
        const oldRaw = window.localStorage.getItem(fallbackHistoryKey);
        const oldParsed = oldRaw ? JSON.parse(oldRaw) : null;
        migratedSession.messages = Array.isArray(oldParsed?.messages) ? oldParsed.messages : [];
        migratedSession.draft = window.localStorage.getItem(fallbackDraftKey) || "";
      } catch (_) {
        // Ignore old history migration failures.
      }

      const firstSession = migratedSession.messages.length || migratedSession.draft ? migratedSession : createSession("新窗口");
      return { activeId: firstSession.id, sessions: [firstSession] };
    };

    const saveSessions = () => {
      try {
        window.localStorage.setItem(aiSessionsKey, JSON.stringify(sessionState));
      } catch (_) {
        // Ignore storage failures.
      }
    };

    const getActiveSession = () => {
      if (!sessionState.sessions.length) {
        const session = createSession("新窗口");
        sessionState.sessions.push(session);
        sessionState.activeId = session.id;
      }
      const activeSession = sessionState.sessions.find((session) => session.id === sessionState.activeId) || sessionState.sessions[0];
      sessionState.activeId = activeSession.id;
      return activeSession;
    };

    const buildConversation = (session) => {
      const conversation = [];
      (session.messages || []).forEach((item) => {
        if (item.role === "user") {
          conversation.push({ role: "user", content: String(item.content || "") });
        } else if (item.role === "ai") {
          const plainText = String(item.content || "").replace(/<[^>]+>/g, " ").replace(/\s+/g, " ").trim();
          if (plainText) {
            conversation.push({ role: "assistant", content: plainText });
          }
        }
      });
      return conversation;
    };

    const sessionShell = document.createElement("div");
    sessionShell.className = "ai-session-shell";
    sessionShell.innerHTML = `
      <div class="ai-session-tabs" data-ai-session-tabs></div>
      <div class="button-row ai-session-actions">
        <button class="button button-secondary" type="button" data-ai-session-new>新建窗口</button>
        <button class="button button-ghost" type="button" data-ai-session-delete>删除当前窗口</button>
      </div>
    `;
    aiChat.parentNode.insertBefore(sessionShell, aiChat);
    const sessionTabs = sessionShell.querySelector("[data-ai-session-tabs]");

    const renderSessionTabs = () => {
      const activeSession = getActiveSession();
      sessionTabs.innerHTML = sessionState.sessions.map((session, index) => {
        const activeClass = session.id === activeSession.id ? " is-active" : "";
        const title = session.title || `窗口 ${index + 1}`;
        return `<button class="ai-session-tab${activeClass}" type="button" data-ai-session-id="${escapeHtml(session.id)}">${escapeHtml(title)}</button>`;
      }).join("");
    };

    const renderMessage = (role, title, content, trustedHtml = false) => {
      const node = document.createElement("div");
      node.className = `message ${role}`;
      node.innerHTML = trustedHtml
        ? `<strong>${escapeHtml(title)}</strong><div>${content}</div>`
        : `<strong>${escapeHtml(title)}</strong><div>${escapeHtml(content)}</div>`;
      aiChat.appendChild(node);
      aiChat.scrollTop = aiChat.scrollHeight;
    };

    const renderActiveSession = () => {
      const activeSession = getActiveSession();
      aiChat.innerHTML = initialGreeting;
      (activeSession.messages || []).forEach((item) => {
        renderMessage(item.role, item.title, item.content || "", Boolean(item.trustedHtml));
      });
      aiInput.value = activeSession.draft || "";
      renderSessionTabs();
    };

    const updateSessionTitle = (session, firstQuestion) => {
      if (!session || !firstQuestion || session.messages.filter((item) => item.role === "user").length > 1) {
        return;
      }
      const text = String(firstQuestion || "").trim();
      session.title = text.length > 12 ? `${text.slice(0, 12)}...` : text || "新窗口";
    };

    const appendMessageToSession = (session, role, title, content, trustedHtml = false) => {
      const message = { role, title, content, trustedHtml: Boolean(trustedHtml), createdAt: new Date().toISOString() };
      session.messages = [...(session.messages || []), message].slice(-40);
      session.updatedAt = new Date().toISOString();
      if (role === "user") {
        updateSessionTitle(session, content);
      }
    };

    const appendMessage = (role, title, content, trustedHtml = false) => {
      const activeSession = getActiveSession();
      appendMessageToSession(activeSession, role, title, content, trustedHtml);
      saveSessions();
      renderMessage(role, title, content, trustedHtml);
      renderSessionTabs();
    };

    const saveAiDraft = () => {
      const activeSession = getActiveSession();
      activeSession.draft = aiInput.value || "";
      activeSession.updatedAt = new Date().toISOString();
      saveSessions();
    };

    const askAi = async (value) => {
      const activeId = getActiveSession().id;
      appendMessage("user", "你", value);
      const conversation = buildConversation(getActiveSession());

      const typingNode = document.createElement("div");
      typingNode.className = "message ai";
      typingNode.innerHTML = `<strong>AI 助手</strong><div>正在查询中...</div>`;
      aiChat.appendChild(typingNode);
      aiChat.scrollTop = aiChat.scrollHeight;

      try {
        const result = await api.post(API_ENDPOINTS.passenger.aiChat, {
          messages: conversation.slice(-8),
        });

        const data = result?.data || {};
        if (typingNode.parentNode === aiChat) {
          aiChat.removeChild(typingNode);
        }

        if (getActiveSession().id === activeId) {
          appendMessage("ai", "AI 助手", buildPassengerAiRichHtml(data), true);
        } else {
          const targetSession = sessionState.sessions.find((session) => session.id === activeId);
          if (targetSession) {
            appendMessageToSession(targetSession, "ai", "AI 助手", buildPassengerAiRichHtml(data), true);
            saveSessions();
            renderSessionTabs();
          }
        }
      } catch (error) {
        if (typingNode.parentNode === aiChat) {
          aiChat.removeChild(typingNode);
        }
        const errorHtml = buildPassengerAiErrorHtml(error);
        if (getActiveSession().id === activeId) {
          appendMessage("ai", "AI 助手", errorHtml, true);
        } else {
          const targetSession = sessionState.sessions.find((session) => session.id === activeId);
          if (targetSession) {
            appendMessageToSession(targetSession, "ai", "AI 助手", errorHtml, true);
            saveSessions();
            renderSessionTabs();
          }
        }
      }
    };

    aiForm.addEventListener("submit", async (event) => {
      event.preventDefault();
      const value = aiInput.value.trim();
      if (!value) {
        return;
      }
      aiInput.value = "";
      saveAiDraft();
      await askAi(value);
    });

    aiInput.addEventListener("input", saveAiDraft);

    aiChat.addEventListener("click", async (event) => {
      const button = event.target.closest("[data-ai-suggestion]");
      if (!button) {
        return;
      }

      const value = String(button.getAttribute("data-ai-suggestion") || "").trim();
      if (!value) {
        return;
      }

      aiInput.value = value;
      saveAiDraft();
      await askAi(value);
    });

    sessionTabs.addEventListener("click", (event) => {
      const button = event.target.closest("[data-ai-session-id]");
      if (!button) {
        return;
      }
      saveAiDraft();
      sessionState.activeId = button.getAttribute("data-ai-session-id") || sessionState.activeId;
      saveSessions();
      renderActiveSession();
    });

    sessionShell.querySelector("[data-ai-session-new]").addEventListener("click", () => {
      saveAiDraft();
      const session = createSession(`窗口 ${sessionState.sessions.length + 1}`);
      sessionState.sessions.unshift(session);
      sessionState.activeId = session.id;
      saveSessions();
      renderActiveSession();
      aiInput.focus();
    });

    sessionShell.querySelector("[data-ai-session-delete]").addEventListener("click", () => {
      if (sessionState.sessions.length <= 1) {
        const onlySession = getActiveSession();
        onlySession.messages = [];
        onlySession.draft = "";
        onlySession.title = "新窗口";
        onlySession.updatedAt = new Date().toISOString();
      } else {
        const currentId = getActiveSession().id;
        sessionState.sessions = sessionState.sessions.filter((session) => session.id !== currentId);
        sessionState.activeId = sessionState.sessions[0].id;
      }
      saveSessions();
      renderActiveSession();
    });

    sessionState = readSessions();
    renderActiveSession();
  }

  function initDriverDraftGenerator() {
    const fillTripButton = document.querySelector("[data-fill-trip]");
    if (!fillTripButton) {
      return;
    }

    fillTripButton.addEventListener("click", () => {
      redirectTo(ROUTES.driver.ai);
    });
  }

  function initDebugLog() {
    if (!runtimeConfig.debug) {
      return;
    }

    const auth = readAuth();
    const page = document.body.dataset.page || "unknown";
    const role = auth?.role || "guest";
    console.info(`[TripVerse] page=${page}, role=${role}, apiBaseUrl=${runtimeConfig.apiBaseUrl}, useMock=${runtimeConfig.useMock}`);
  }

  function initAdminKnowledgePage() {
    if (getCurrentFileName() !== ROUTES.admin.knowledge) {
      return;
    }

    const uploadForm = document.querySelector("[data-knowledge-upload-form]");
    const documentList = document.querySelector("[data-knowledge-document-list]");
    const searchForm = document.querySelector("[data-knowledge-search-form]");
    const searchResult = document.querySelector("[data-knowledge-search-result]");

    const renderDocuments = (docs) => {
      if (!documentList) {
        return;
      }

      documentList.innerHTML = Array.isArray(docs) && docs.length
        ? docs.map((doc) => `
            <div class="info-card">
              <strong>${escapeHtml(doc.title || "--")}</strong>
              <p class="muted">分类：${escapeHtml(doc.category || "--")}</p>
              <p class="muted">来源：${escapeHtml(doc.sourceName || "--")}</p>
              <p class="muted">状态：${escapeHtml(doc.status || "--")}，切片数：${Number(doc.chunkCount || 0)}</p>
              <p class="muted">创建时间：${escapeHtml(formatFullDateTime(doc.createdAt || ""))}</p>
            </div>
          `).join("")
        : `
          <div class="info-card">
            <strong>暂无文档</strong>
            <p class="muted">可以先上传 Markdown 或 TXT 文档。</p>
          </div>
        `;
    };

    const loadDocuments = async () => {
      try {
        const result = await api.get(API_ENDPOINTS.admin.knowledge);
        renderDocuments(Array.isArray(result?.data) ? result.data : []);
      } catch (error) {
        if (documentList) {
          documentList.innerHTML = `
            <div class="info-card">
              <strong>加载失败</strong>
              <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
            </div>
          `;
        }
      }
    };

    if (uploadForm) {
      uploadForm.addEventListener("submit", async (event) => {
        event.preventDefault();

        const formData = new window.FormData(uploadForm);
        const title = String(formData.get("title") || "").trim();
        const category = String(formData.get("category") || "").trim();
        const file = formData.get("file");

        if (!title || !category || !(file instanceof window.File) || !file.name) {
          showToast("请填写标题、分类并选择文件");
          return;
        }

        const payload = new window.FormData();
        payload.append("title", title);
        payload.append("category", category);
        payload.append("file", file);

        try {
          await api.request(API_ENDPOINTS.admin.knowledgeUpload, {
            method: "POST",
            body: payload,
          });
          showToast("文档上传并切片入库成功");
          uploadForm.reset();
          await loadDocuments();
        } catch (error) {
          showToast(error.message || "文档上传失败");
        }
      });
    }

    if (searchForm && searchResult) {
      searchForm.addEventListener("submit", async (event) => {
        event.preventDefault();

        const formData = new window.FormData(searchForm);
        const query = String(formData.get("query") || "").trim();
        const topK = Number(formData.get("topK") || 6);

        if (!query) {
          showToast("请输入测试问题");
          return;
        }

        try {
          const result = await api.post(API_ENDPOINTS.admin.knowledgeSearch, {
            query,
            topK,
          });

          const items = Array.isArray(result?.data) ? result.data : [];
          searchResult.innerHTML = items.length
            ? items.map((item) => `
                <div class="info-card">
                  <strong>${escapeHtml(item.title || "--")}</strong>
                  <p class="muted">路径：${escapeHtml(item.sectionPath || "--")}</p>
                  <p class="muted">综合分：${Number(item.finalScore || 0).toFixed(4)}，向量分：${Number(item.vectorScore || 0).toFixed(4)}，关键词分：${Number(item.keywordScore || 0).toFixed(4)}</p>
                  <p class="muted">${escapeHtml(item.content || "--")}</p>
                </div>
              `).join("")
            : `
              <div class="info-card">
                <strong>未召回到结果</strong>
                <p class="muted">可以换个问法，或者先上传更多规则文档。</p>
              </div>
            `;
        } catch (error) {
          searchResult.innerHTML = `
            <div class="info-card">
              <strong>检索失败</strong>
              <p class="muted">${escapeHtml(error.message || "请稍后重试")}</p>
            </div>
          `;
        }
      });
    }

    loadDocuments();
  }

  function highlightPassengerAiText(text, query) {
    const source = String(text || "");
    const escaped = escapeHtml(source);
    const terms = String(query || "")
      .split(/[\s,，。、“”"'、:：;；!?？！()（）]+/)
      .map((item) => item.trim())
      .filter((item) => item && item.length >= 2)
      .slice(0, 8);

    let result = escaped;
    terms.forEach((term) => {
      const escapedTerm = escapeHtml(term).replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
      result = result.replace(new RegExp(escapedTerm, "gi"), (matched) => `<mark>${matched}</mark>`);
    });
    return result;
  }

  function renderPassengerAiKnowledgeSources(context, query) {
    const sources = Array.isArray(context?.knowledgeSources) ? context.knowledgeSources : [];
    if (!sources.length) {
      return "";
    }

    return `
      <div class="section-block">
        <strong>答案来源</strong>
        ${sources.map((item, index) => `
          <details class="info-card">
            <summary>来源 ${index + 1}：${highlightPassengerAiText(item.title || "--", query)} / ${highlightPassengerAiText(item.sectionPath || "--", query)} / 分数 ${escapeHtml(item.finalScore || "--")}</summary>
            <p class="muted">${highlightPassengerAiText(item.content || "--", query)}</p>
          </details>
        `).join("")}
      </div>
    `;
  }

  function buildPassengerAiRichHtml(data, query) {
    const reply = String(data?.reply || "").trim() || "暂时没有生成回复";
    const intent = String(data?.intent || "");
    const context = data?.context || null;
    const suggestions = Array.isArray(data?.suggestions) ? data.suggestions : [];

    return [
      `<div>${escapeHtml(reply)}</div>`,
      renderPassengerAiRouteCards(context),
      renderPassengerAiOrderSummary(context),
      renderPassengerAiOrderCards(context),
      renderPassengerAiRefundRules(context),
      renderPassengerAiKnowledgeSources(context, query),
      renderPassengerAiSystemHints(context),
      renderPassengerAiSuggestions(suggestions),
      renderPassengerAiActions(intent, context),
    ].join("");
  }

  function initAdminKnowledgePage() {
    if (getCurrentFileName() !== ROUTES.admin.knowledge) {
      return;
    }

    const uploadForm = document.querySelector("[data-knowledge-upload-form]");
    const filterForm = document.querySelector("[data-knowledge-filter-form]");
    const filterResetButton = document.querySelector("[data-knowledge-filter-reset]");
    const documentList = document.querySelector("[data-knowledge-document-list]");
    const searchForm = document.querySelector("[data-knowledge-search-form]");
    const searchResult = document.querySelector("[data-knowledge-search-result]");
    const detailBox = document.querySelector("[data-knowledge-detail]");
    let currentCategory = "";

    const renderDetail = (doc) => {
      if (!detailBox) {
        return;
      }
      if (!doc) {
        detailBox.innerHTML = `
          <div class="info-card">
            <strong>等待查看</strong>
            <p class="muted">点击左侧文档后，这里会显示原文、切片列表、命中次数和管理操作。</p>
          </div>
        `;
        return;
      }

        const chunks = Array.isArray(doc.chunks) ? doc.chunks : [];
        detailBox.innerHTML = `
          <div class="info-card">
            <strong>${escapeHtml(doc.title || "--")}</strong>
            <p class="muted">分类：${escapeHtml(doc.category || "--")}，状态：${escapeHtml(doc.status || "--")}，切片数：${Number(doc.chunkCount || 0)}</p>
            <p class="muted">来源：${escapeHtml(doc.sourceName || "--")}</p>
            <div class="button-row section-block">
              <button class="button button-dark" type="button" data-knowledge-reindex="${escapeHtml(doc.id || "")}">重新切片</button>
              <button class="button button-secondary" type="button" data-knowledge-toggle="${escapeHtml(doc.id || "")}" data-knowledge-next-status="${doc.status === "active" ? "disabled" : "active"}">${doc.status === "active" ? "禁用文档" : "启用文档"}</button>
              <button class="button button-ghost" type="button" data-knowledge-delete="${escapeHtml(doc.id || "")}">删除文档</button>
            </div>
          </div>
          <form class="search-card section-block" data-knowledge-edit-form data-knowledge-edit-id="${escapeHtml(doc.id || "")}">
            <div class="field span-6">
              <label>文档标题</label>
              <input name="title" type="text" value="${escapeHtml(doc.title || "")}">
            </div>
            <div class="field span-6">
              <label>文档分类</label>
              <input name="category" type="text" value="${escapeHtml(doc.category || "")}">
            </div>
            <div class="field span-12">
              <label>原文内容</label>
              <textarea name="content">${escapeHtml(doc.content || "")}</textarea>
            </div>
            <div class="button-row">
              <button class="button button-primary" type="submit">保存文档并重建切片</button>
            </div>
          </form>
          <details class="info-card" open>
            <summary>查看原文</summary>
            <p class="muted">${escapeHtml(doc.content || "--")}</p>
          </details>
        <div class="section-block">
          <strong>切片列表</strong>
          ${chunks.length ? chunks.map((chunk) => `
            <details class="info-card">
              <summary>#${Number(chunk.chunkIndex || 0) + 1} ${escapeHtml(chunk.title || "--")} / 命中 ${Number(chunk.hitCount || 0)} 次 / 向量维度 ${Number(chunk.embeddingDim || 0)}</summary>
              <p class="muted">路径：${escapeHtml(chunk.sectionPath || "--")}</p>
              <p class="muted">估算 tokens：${Number(chunk.tokenEstimate || 0)}，最后命中：${escapeHtml(chunk.lastHitAt || "--")}</p>
              <p class="muted">${escapeHtml(chunk.content || "--")}</p>
            </details>
          `).join("") : `<div class="info-card"><strong>暂无切片</strong></div>`}
        </div>
      `;
    };

    const loadDocuments = async () => {
      try {
        const result = await api.get(API_ENDPOINTS.admin.knowledge, currentCategory ? { category: currentCategory } : undefined);
        const docs = Array.isArray(result?.data) ? result.data : [];
        if (documentList) {
          documentList.innerHTML = docs.length
            ? docs.map((doc) => `
                <div class="info-card">
                  <strong>${escapeHtml(doc.title || "--")}</strong>
                  <p class="muted">分类：${escapeHtml(doc.category || "--")}</p>
                  <p class="muted">状态：${escapeHtml(doc.status || "--")}，切片数：${Number(doc.chunkCount || 0)}</p>
                  <div class="button-row section-block">
                    <button class="button button-secondary" type="button" data-knowledge-view="${escapeHtml(doc.id || "")}">查看切片</button>
                  </div>
                </div>
              `).join("")
            : `<div class="info-card"><strong>暂无文档</strong><p class="muted">可以先上传 Markdown 或 TXT 文档。</p></div>`;
        }
      } catch (error) {
        if (documentList) {
          documentList.innerHTML = `<div class="info-card"><strong>加载失败</strong><p class="muted">${escapeHtml(error.message || "请稍后重试")}</p></div>`;
        }
      }
    };

    if (filterForm) {
      filterForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new window.FormData(filterForm);
        currentCategory = String(formData.get("category") || "").trim();
        await loadDocuments();
        renderDetail(null);
      });
    }

    if (filterResetButton) {
      filterResetButton.addEventListener("click", async () => {
        currentCategory = "";
        if (filterForm) {
          filterForm.reset();
        }
        await loadDocuments();
        renderDetail(null);
      });
    }

    const loadDetail = async (documentId) => {
      if (!documentId) {
        renderDetail(null);
        return;
      }
      try {
        const result = await api.get(API_ENDPOINTS.admin.knowledgeDetail, undefined, {
          pathParams: { documentId },
        });
        renderDetail(result?.data || null);
      } catch (error) {
        renderDetail(null);
        showToast(error.message || "加载文档详情失败");
      }
    };

    if (uploadForm) {
      uploadForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new window.FormData(uploadForm);
        const title = String(formData.get("title") || "").trim();
        const category = String(formData.get("category") || "").trim();
        const file = formData.get("file");

        if (!title || !category || !(file instanceof window.File) || !file.name) {
          showToast("请填写标题、分类并选择文件");
          return;
        }

        const payload = new window.FormData();
        payload.append("title", title);
        payload.append("category", category);
        payload.append("file", file);

        try {
          const result = await api.request(API_ENDPOINTS.admin.knowledgeUpload, {
            method: "POST",
            body: payload,
          });
          showToast("文档上传并切片入库成功");
          uploadForm.reset();
          await loadDocuments();
          await loadDetail(result?.data?.id || "");
        } catch (error) {
          showToast(error.message || "文档上传失败");
        }
      });
    }

    if (documentList) {
      documentList.addEventListener("click", async (event) => {
        const button = event.target.closest("[data-knowledge-view]");
        if (!button) {
          return;
        }
        await loadDetail(String(button.getAttribute("data-knowledge-view") || "").trim());
      });
    }

    if (detailBox) {
      detailBox.addEventListener("click", async (event) => {
        const reindexButton = event.target.closest("[data-knowledge-reindex]");
        if (reindexButton) {
          const documentId = String(reindexButton.getAttribute("data-knowledge-reindex") || "").trim();
          try {
            await api.post(API_ENDPOINTS.admin.knowledgeReindex, {}, {
              pathParams: { documentId },
            });
            showToast("文档重新切片成功");
            await loadDocuments();
            await loadDetail(documentId);
          } catch (error) {
            showToast(error.message || "重新切片失败");
          }
          return;
        }

        const toggleButton = event.target.closest("[data-knowledge-toggle]");
        if (toggleButton) {
          const documentId = String(toggleButton.getAttribute("data-knowledge-toggle") || "").trim();
          const status = String(toggleButton.getAttribute("data-knowledge-next-status") || "").trim();
          try {
            await api.patch(API_ENDPOINTS.admin.knowledgeDetail, { status }, {
              pathParams: { documentId },
            });
            showToast(status === "disabled" ? "文档已禁用" : "文档已启用");
            await loadDocuments();
            await loadDetail(documentId);
          } catch (error) {
            showToast(error.message || "更新文档状态失败");
          }
          return;
        }

        const deleteButton = event.target.closest("[data-knowledge-delete]");
        if (deleteButton) {
          const documentId = String(deleteButton.getAttribute("data-knowledge-delete") || "").trim();
          if (!window.confirm("确定删除这篇知识文档吗？")) {
            return;
          }
          try {
            await api.delete(API_ENDPOINTS.admin.knowledgeDetail, {
              pathParams: { documentId },
            });
            showToast("文档已删除");
            await loadDocuments();
            renderDetail(null);
          } catch (error) {
            showToast(error.message || "删除文档失败");
          }
        }
      });

      detailBox.addEventListener("submit", async (event) => {
        const form = event.target.closest("[data-knowledge-edit-form]");
        if (!form) {
          return;
        }
        event.preventDefault();

        const documentId = String(form.getAttribute("data-knowledge-edit-id") || "").trim();
        const formData = new window.FormData(form);
        const title = String(formData.get("title") || "").trim();
        const category = String(formData.get("category") || "").trim();
        const content = String(formData.get("content") || "").trim();

        if (!documentId || !title || !category || !content) {
          showToast("请填写标题、分类和原文内容");
          return;
        }

        try {
          await api.patch(API_ENDPOINTS.admin.knowledgeDetail, { title, category, content }, {
            pathParams: { documentId },
          });
          showToast("文档已保存并重建切片");
          await loadDocuments();
          await loadDetail(documentId);
        } catch (error) {
          showToast(error.message || "保存文档失败");
        }
      });
    }

    if (searchForm && searchResult) {
      searchForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const formData = new window.FormData(searchForm);
        const query = String(formData.get("query") || "").trim();
        const topK = Number(formData.get("topK") || 6);
        if (!query) {
          showToast("请输入测试问题");
          return;
        }

        try {
          const result = await api.post(API_ENDPOINTS.admin.knowledgeSearch, { query, topK });
          const items = Array.isArray(result?.data) ? result.data : [];
          searchResult.innerHTML = items.length
            ? items.map((item) => `
                <details class="info-card" open>
                  <summary>${escapeHtml(item.title || "--")} / ${escapeHtml(item.sectionPath || "--")} / 综合分 ${Number(item.finalScore || 0).toFixed(4)}</summary>
                  <p class="muted">向量分：${Number(item.vectorScore || 0).toFixed(4)}，关键词分：${Number(item.keywordScore || 0).toFixed(4)}</p>
                  <p class="muted">命中次数：${Number(item.hitCount || 0)}，最后命中：${escapeHtml(item.lastHitAt || "--")}</p>
                  <p class="muted">${escapeHtml(item.content || "--")}</p>
                </details>
              `).join("")
            : `<div class="info-card"><strong>未召回到结果</strong><p class="muted">可以换个问法，或者先上传更多规则文档。</p></div>`;
        } catch (error) {
          searchResult.innerHTML = `<div class="info-card"><strong>检索失败</strong><p class="muted">${escapeHtml(error.message || "请稍后重试")}</p></div>`;
        }
      });
    }

    renderDetail(null);
    loadDocuments();
  }

  function initAdminTokensPage() {
    if (getCurrentFileName() !== ROUTES.admin.tokens) {
      return;
    }

    const form = document.querySelector("[data-admin-token-filter-form]");
    const userList = document.querySelector("[data-admin-token-user-list]");
    const detailList = document.querySelector("[data-admin-token-detail-list]");

    const render = (data) => {
      const users = Array.isArray(data?.users) ? data.users : [];
      const items = Array.isArray(data?.items) ? data.items : [];

      if (userList) {
        userList.innerHTML = users.length
          ? users.map((item) => `
              <div class="info-card">
                <strong>用户 #${Number(item.userId || 0)} / ${escapeHtml(item.role || "--")}</strong>
                <p class="muted">请求数 ${Number(item.requestCount || 0)}，总 tokens ${Number(item.totalTokens || 0)}</p>
                <p class="muted">输入 ${Number(item.promptTokens || 0)}，输出 ${Number(item.completionTokens || 0)}</p>
                <p class="muted">最后使用：${escapeHtml(formatFullDateTime(item.lastUsedAt || ""))}</p>
              </div>
            `).join("")
          : `<div class="info-card"><strong>暂无数据</strong><p class="muted">当前筛选条件下没有 token 使用记录。</p></div>`;
      }

      if (detailList) {
        detailList.innerHTML = items.length
          ? items.map((item) => `
              <div class="info-card">
                <strong>用户 #${Number(item.userId || 0)} / ${escapeHtml(item.role || "--")}</strong>
                <p class="muted">功能：${escapeHtml(item.feature || "--")}，模型：${escapeHtml(item.model || "--")}</p>
                <p class="muted">请求数 ${Number(item.requestCount || 0)}，总 tokens ${Number(item.totalTokens || 0)}</p>
                <p class="muted">输入 ${Number(item.promptTokens || 0)}，输出 ${Number(item.completionTokens || 0)}</p>
                <p class="muted">最后使用：${escapeHtml(formatFullDateTime(item.lastUsedAt || ""))}</p>
              </div>
            `).join("")
          : `<div class="info-card"><strong>暂无明细</strong><p class="muted">当前筛选条件下没有明细记录。</p></div>`;
      }
    };

    const load = async () => {
      const formData = form ? new window.FormData(form) : new window.FormData();
      const role = String(formData.get("role") || "").trim();
      const feature = String(formData.get("feature") || "").trim();
      const days = String(formData.get("days") || "30").trim();

      try {
        const result = await api.get(API_ENDPOINTS.admin.tokens, { role, feature, days });
        render(result?.data || {});
      } catch (error) {
        showToast(error.message || "加载 token 使用量失败");
      }
    };

    if (form) {
      form.addEventListener("submit", async (event) => {
        event.preventDefault();
        await load();
      });
    }

    load();
  }

  function initAdminRiskPage() {
    if (getCurrentFileName() !== ROUTES.admin.risk) {
      return;
    }

    const form = document.querySelector("[data-admin-risk-filter-form]");
    const summaryBox = document.querySelector("[data-admin-risk-summary]");
    const listBox = document.querySelector("[data-admin-risk-list]");

    const severityLabelMap = {
      high: "高优先级",
      medium: "中优先级",
      low: "低优先级",
    };
    const statusLabelMap = {
      open: "待处理",
      acknowledged: "已确认",
      resolved: "已解决",
    };

    const renderSummary = (summary) => {
      if (!summaryBox) {
        return;
      }

      summaryBox.innerHTML = `
        <div class="info-card">
          <strong>待处理风险</strong>
          <p class="muted">${Number(summary?.openCount || 0)} 条</p>
        </div>
        <div class="info-card">
          <strong>高优先级</strong>
          <p class="muted">${Number(summary?.highCount || 0)} 条</p>
        </div>
        <div class="info-card">
          <strong>中优先级</strong>
          <p class="muted">${Number(summary?.mediumCount || 0)} 条</p>
        </div>
        <div class="info-card">
          <strong>低优先级</strong>
          <p class="muted">${Number(summary?.lowCount || 0)} 条</p>
        </div>
      `;
    };

    const renderMetrics = (metrics) => {
      const entries = Object.entries(metrics || {});
      if (!entries.length) {
        return "";
      }
      return `
        <div class="list-meta">
          ${entries.map(([key, value]) => `<span>${escapeHtml(key)}: ${escapeHtml(String(value))}</span>`).join("")}
        </div>
      `;
    };

    const renderList = (items) => {
      if (!listBox) {
        return;
      }

      listBox.innerHTML = Array.isArray(items) && items.length
        ? items.map((item) => `
            <div class="info-card">
              <strong>${escapeHtml(severityLabelMap[item.severity] || item.severity || "--")} / ${escapeHtml(item.title || "--")}</strong>
              <p class="muted">类型：${escapeHtml(item.eventType || "--")}，主体：${escapeHtml(item.subjectId || "--")}，状态：${escapeHtml(statusLabelMap[item.status] || item.status || "--")}</p>
              <p class="muted">${escapeHtml(item.detail || "--")}</p>
              ${renderMetrics(item.metrics)}
              <p class="muted">发生时间：${escapeHtml(formatFullDateTime(item.createdAt || ""))}</p>
              ${item.status === "open" ? `
                <div class="button-row section-block">
                  <button class="button button-secondary" type="button" data-risk-status="${escapeHtml(String(item.id || ""))}" data-risk-next="acknowledged">标记已确认</button>
                  <button class="button button-dark" type="button" data-risk-status="${escapeHtml(String(item.id || ""))}" data-risk-next="resolved">标记已解决</button>
                </div>
              ` : ""}
            </div>
          `).join("")
        : `
          <div class="info-card">
            <strong>暂无风险事件</strong>
            <p class="muted">当前筛选条件下没有检测到风险告警。</p>
          </div>
        `;
    };

    const load = async () => {
      const formData = form ? new window.FormData(form) : new window.FormData();
      const severity = String(formData.get("severity") || "").trim();
      const status = String(formData.get("status") || "").trim();
      const days = String(formData.get("days") || "30").trim();

      try {
        const result = await api.get(API_ENDPOINTS.admin.riskLogs, { severity, status, days });
        const data = result?.data || {};
        renderSummary(data.summary || {});
        renderList(data.items || []);
      } catch (error) {
        if (summaryBox) {
          summaryBox.innerHTML = `<div class="info-card"><strong>加载失败</strong><p class="muted">${escapeHtml(error.message || "请稍后重试")}</p></div>`;
        }
        if (listBox) {
          listBox.innerHTML = `<div class="info-card"><strong>加载失败</strong><p class="muted">${escapeHtml(error.message || "请稍后重试")}</p></div>`;
        }
      }
    };

    if (form) {
      form.addEventListener("submit", async (event) => {
        event.preventDefault();
        await load();
      });
    }

    if (listBox) {
      listBox.addEventListener("click", async (event) => {
        const button = event.target.closest("[data-risk-status]");
        if (!button) {
          return;
        }

        const eventId = String(button.getAttribute("data-risk-status") || "").trim();
        const status = String(button.getAttribute("data-risk-next") || "").trim();
        if (!eventId || !status) {
          return;
        }

        try {
          await api.patch(API_ENDPOINTS.admin.riskLogDetail, { status }, {
            pathParams: { eventId },
          });
            showToast(status === "resolved" ? "风险已标记为已解决" : "风险已标记为已确认");
          await load();
        } catch (error) {
          showToast(error.message || "更新风险状态失败");
        }
      });
    }

    load();
  }

  function initDriverTicketVerifyForm() {
    const form = document.querySelector("[data-driver-ticket-verify-form]");
    if (!form) {
      return;
    }

    const resultBox = document.querySelector("[data-driver-ticket-verify-result]");
    form.addEventListener("submit", async (event) => {
      event.preventDefault();
      const token = String(new FormData(form).get("token") || "").trim();
      if (!token) {
        showToast("请先粘贴电子票 token");
        return;
      }

      const button = form.querySelector("button[type='submit']");
      try {
        if (button) {
          button.disabled = true;
          button.textContent = "核验中...";
        }
        const result = await api.post(API_ENDPOINTS.driver.verifyTicket, { token });
        const ticket = result?.data || {};
        if (resultBox) {
          resultBox.innerHTML = `
            <div class="info-card">
              <strong>核验成功</strong>
              <p class="muted">订单 #${escapeHtml(ticket.orderId || "")} 已完成核验。</p>
              <div class="list-meta"><span>状态：${escapeHtml(ticket.status || "--")}</span><span>核验时间：${escapeHtml(formatFullDateTime(ticket.verifiedAt))}</span></div>
            </div>
          `;
        }
        showToast("电子票核验成功");
      } catch (error) {
        if (resultBox) {
          resultBox.innerHTML = `
            <div class="info-card">
              <strong>核验失败</strong>
              <p class="muted">${escapeHtml(error.message || "请检查 token 是否正确。")}</p>
            </div>
          `;
        }
        showToast(error.message || "电子票核验失败");
      } finally {
        if (button) {
          button.disabled = false;
          button.textContent = "核验电子票";
        }
      }
    });
  }

  function initApp() {
    if (!requireAccess()) {
      return;
    }

    renderNavigation();
    initNavigation();
    initLogoutActions();
    initToastTriggers();
    initSeatPicker();
    initCheckoutSummary();
    initPassengerAi();
    initDriverDraftGenerator();
    initAuthForms();
    initDriverTripsPage();
    initDriverDashboardPageV2();
    initDriverIncomePageV2();
    initDriverAiPageV3();
    initDriverPublishPage();
    hardenDriverPublishPageV2();
    hydrateDriverPublishDraftFromAi();
    initDriverTripDetailPage();
    initDriverTicketVerifyForm();
    syncDriverTripOrderStatus();
    initTicketSearchPage();
    initTicketDetailPage();
    initCheckoutPage();
    initOrdersPage();
    initOrderDetailPage();
    initPaymentPage();
      initAdminDashboardPageV3();
      initAdminUsersPageV2();
      initAdminOrdersPage();
      initAdminTokensPage();
      initAdminRiskPage();
      initAdminKnowledgePage();
      initProfilePage();
      renderProfileAccountModule();
    initNotificationCenter();
    initDebugLog();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", initApp);
  } else {
    initApp();
  }
})(window, document);





