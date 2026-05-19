import http from "k6/http";
import { check, group, sleep } from "k6";
import { Trend, Rate } from "k6/metrics";

const BASE_URL = __ENV.BASE_URL || "http://localhost:8080/api";
const PEAK_VUS = Number(__ENV.PEAK_VUS || 20);
const RAMP_UP = __ENV.RAMP_UP || "15s";
const HOLD = __ENV.HOLD || "30s";
const RAMP_DOWN = __ENV.RAMP_DOWN || "15s";

const reqDuration = new Trend("tripverse_req_duration", true);
const reqFailRate = new Rate("tripverse_req_failed");
const endpointMetrics = {
  publicSearchMain: new Trend("endpoint_public_search_main", true),
  publicSearchTransfer: new Trend("endpoint_public_search_transfer", true),
  publicTicketDetail: new Trend("endpoint_public_ticket_detail", true),
  passengerMe: new Trend("endpoint_passenger_auth_me", true),
  passengerProfile: new Trend("endpoint_passenger_profile", true),
  passengerAccount: new Trend("endpoint_passenger_account_status", true),
  passengerOrders: new Trend("endpoint_passenger_orders_my", true),
  passengerNotifications: new Trend("endpoint_passenger_notifications_my", true),
  passengerUnreadCount: new Trend("endpoint_passenger_unread_count", true),
  driverTrips: new Trend("endpoint_driver_trips", true),
  driverDashboard: new Trend("endpoint_driver_dashboard", true),
  driverIncome: new Trend("endpoint_driver_income", true),
  adminDashboard: new Trend("endpoint_admin_dashboard", true),
  adminUsers: new Trend("endpoint_admin_users", true),
  adminUsersSummary: new Trend("endpoint_admin_users_summary", true),
  adminOrders: new Trend("endpoint_admin_orders", true),
  adminTokens: new Trend("endpoint_admin_tokens", true),
  adminRiskLogs: new Trend("endpoint_admin_risk_logs", true),
  adminKnowledge: new Trend("endpoint_admin_knowledge", true),
};

export const options = {
  scenarios: {
    core_read_flow: {
      executor: "ramping-vus",
      stages: [
        { duration: "15s", target: 10 },
        { duration: RAMP_UP, target: PEAK_VUS },
        { duration: HOLD, target: PEAK_VUS },
        { duration: RAMP_DOWN, target: 0 },
      ],
      gracefulRampDown: "5s",
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.05"],
    http_req_duration: ["p(95)<1000"],
    tripverse_req_failed: ["rate<0.05"],
    tripverse_req_duration: ["p(95)<1000"],
  },
};

function jsonHeaders(token) {
  const headers = { "Content-Type": "application/json" };
  if (token) {
    headers.Authorization = `Bearer ${token}`;
  }
  return { headers };
}

function record(res, expectedStatuses, metric) {
  reqDuration.add(res.timings.duration);
  if (metric) {
    metric.add(res.timings.duration);
  }
  const allowed = expectedStatuses || [200];
  const ok = allowed.indexOf(res.status) !== -1;
  reqFailRate.add(!ok);
  check(res, {
    [`status in ${allowed.join("/")}`]: () => ok,
  });
  return ok;
}

function login(role, phone) {
  const res = http.post(
    `${BASE_URL}/auth/login/password`,
    JSON.stringify({ role, phone, password: "123456" }),
    jsonHeaders()
  );
  record(res, [200]);
  if (res.status !== 200) {
    return "";
  }
  const body = res.json();
  return body && body.data && body.data.token ? body.data.token : "";
}

function firstTicketId() {
  const res = http.get(
    `${BASE_URL}/tickets/search?startCity=${encodeURIComponent("北京")}&endCity=${encodeURIComponent("上海")}&date=2026-05-12&allowTransfer=true`
  );
  record(res, [200]);
  if (res.status !== 200) {
    return 1;
  }
  const body = res.json();
  const list = body && body.data && body.data.length ? body.data : [];
  for (let i = 0; i < list.length; i += 1) {
    if (list[i].id) {
      return list[i].id;
    }
    if (list[i].legs && list[i].legs.length && list[i].legs[0].tripId) {
      return list[i].legs[0].tripId;
    }
  }
  return 1;
}

export function setup() {
  return {
    passengerToken: login("passenger", "13000000001"),
    driverToken: login("driver", "15000000001"),
    adminToken: login("admin", "18800000000"),
    ticketId: firstTicketId(),
  };
}

export default function (data) {
  group("public ticket APIs", () => {
    record(http.get(`${BASE_URL}/tickets/search?startCity=${encodeURIComponent("北京")}&endCity=${encodeURIComponent("上海")}&date=2026-05-12&allowTransfer=true`), [200], endpointMetrics.publicSearchMain);
    record(http.get(`${BASE_URL}/tickets/search?startCity=${encodeURIComponent("杭州")}&endCity=${encodeURIComponent("苏州")}&date=2026-05-12&allowTransfer=true`), [200], endpointMetrics.publicSearchTransfer);
    record(http.get(`${BASE_URL}/tickets/${data.ticketId}`), [200, 404], endpointMetrics.publicTicketDetail);
  });

  group("passenger read APIs", () => {
    const auth = jsonHeaders(data.passengerToken);
    record(http.get(`${BASE_URL}/auth/me`, auth), [200], endpointMetrics.passengerMe);
    record(http.get(`${BASE_URL}/users/profile`, auth), [200], endpointMetrics.passengerProfile);
    record(http.get(`${BASE_URL}/users/account/status`, auth), [200], endpointMetrics.passengerAccount);
    record(http.get(`${BASE_URL}/orders/my`, auth), [200], endpointMetrics.passengerOrders);
    record(http.get(`${BASE_URL}/notifications/my`, auth), [200], endpointMetrics.passengerNotifications);
    record(http.get(`${BASE_URL}/notifications/unread-count`, auth), [200], endpointMetrics.passengerUnreadCount);
  });

  group("driver read APIs", () => {
    const auth = jsonHeaders(data.driverToken);
    record(http.get(`${BASE_URL}/driver/trips`, auth), [200], endpointMetrics.driverTrips);
    record(http.get(`${BASE_URL}/driver/dashboard`, auth), [200], endpointMetrics.driverDashboard);
    record(http.get(`${BASE_URL}/driver/income`, auth), [200], endpointMetrics.driverIncome);
  });

  group("admin read APIs", () => {
    const auth = jsonHeaders(data.adminToken);
    record(http.get(`${BASE_URL}/admin/dashboard`, auth), [200], endpointMetrics.adminDashboard);
    record(http.get(`${BASE_URL}/admin/users`, auth), [200], endpointMetrics.adminUsers);
    record(http.get(`${BASE_URL}/admin/users/summary`, auth), [200], endpointMetrics.adminUsersSummary);
    record(http.get(`${BASE_URL}/admin/orders`, auth), [200], endpointMetrics.adminOrders);
    record(http.get(`${BASE_URL}/admin/tokens`, auth), [200], endpointMetrics.adminTokens);
    record(http.get(`${BASE_URL}/admin/risk/logs`, auth), [200], endpointMetrics.adminRiskLogs);
    record(http.get(`${BASE_URL}/admin/knowledge`, auth), [200], endpointMetrics.adminKnowledge);
  });

  sleep(1);
}
