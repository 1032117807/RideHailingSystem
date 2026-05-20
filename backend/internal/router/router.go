package router

import (
	"log"
	"net/http"

	"ridehailing/backend/internal/handler"
	"ridehailing/backend/internal/pkg/response"
)

type Dependencies struct {
	AuthHandler             *handler.AuthHandler
	UserHandler             *handler.UserHandler
	DriverTripHandler       *handler.DriverTripHandler
	AIHandler               *handler.AIHandler
	KnowledgeHandler        *handler.KnowledgeHandler
	RiskHandler             *handler.RiskHandler
	TicketHandler           *handler.TicketHandler
	OrderHandler            *handler.OrderHandler
	PaymentHandler          *handler.PaymentHandler
	NotificationHandler     *handler.NotificationHandler
	ElectronicTicketHandler *handler.ElectronicTicketHandler
	AuditHandler            *handler.AuditHandler
	PassengerHandler        *handler.PassengerHandler
	PriceAlertHandler       *handler.PriceAlertHandler
	DriverProfileHandler    *handler.DriverProfileHandler
	AuthMiddleware          func(http.Handler) http.Handler

	EnableMockPayment bool
	TokenUsageHandler *handler.TokenUsageHandler
}

func NewRouter(dep Dependencies) http.Handler {
	mux := http.NewServeMux()

	mount(mux, "/api/health", map[string]http.Handler{
		http.MethodGet: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response.Success(w, map[string]string{"status": "ok"})
		}),
	})

	mount(mux, "/api/auth/email/send", map[string]http.Handler{
		http.MethodPost: http.HandlerFunc(dep.AuthHandler.SendEmailCode),
	})
	mount(mux, "/api/auth/register", map[string]http.Handler{
		http.MethodPost: http.HandlerFunc(dep.AuthHandler.Register),
	})
	mount(mux, "/api/auth/login/password", map[string]http.Handler{
		http.MethodPost: http.HandlerFunc(dep.AuthHandler.LoginByPassword),
	})
	mount(mux, "/api/auth/login/code", map[string]http.Handler{
		http.MethodPost: http.HandlerFunc(dep.AuthHandler.LoginByCode),
	})

	mount(mux, "/api/auth/logout", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.AuthHandler.Logout)),
	})
	mount(mux, "/api/auth/me", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.AuthHandler.Me)),
	})

	mount(mux, "/api/users/profile", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.UserHandler.GetProfile)),
		http.MethodPut: dep.AuthMiddleware(http.HandlerFunc(dep.UserHandler.UpdateProfile)),
	})
	mount(mux, "/api/users/verify-real-name", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.UserHandler.VerifyRealName)),
	})
	mount(mux, "/api/users/account/status", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.UserHandler.GetAccountStatus)),
	})
	mount(mux, "/api/admin/users", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.UserHandler.ListUsersForAdmin)),
	})
	mount(mux, "/api/admin/users/{userId}", map[string]http.Handler{
		http.MethodPatch: dep.AuthMiddleware(http.HandlerFunc(dep.UserHandler.UpdateUserByAdmin)),
	})
	mount(mux, "/api/admin/users/summary", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.UserHandler.GetAdminUserSummary)),
	})
	mount(mux, "/api/notifications/my", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.NotificationHandler.ListMyNotifications)),
	})
	mount(mux, "/api/notifications/unread-count", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.NotificationHandler.CountMyUnreadNotifications)),
	})
	mount(mux, "/api/notifications/{notificationId}/read", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.NotificationHandler.MarkMyNotificationRead)),
	})
	mount(mux, "/api/notifications/read-all", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.NotificationHandler.MarkAllMyNotificationsRead)),
	})
	mount(mux, "/api/passengers", map[string]http.Handler{
		http.MethodGet:  dep.AuthMiddleware(http.HandlerFunc(dep.PassengerHandler.List)),
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.PassengerHandler.Create)),
	})
	mount(mux, "/api/passengers/{passengerId}", map[string]http.Handler{
		http.MethodDelete: dep.AuthMiddleware(http.HandlerFunc(dep.PassengerHandler.Delete)),
	})
	mount(mux, "/api/price-alerts", map[string]http.Handler{
		http.MethodGet:  dep.AuthMiddleware(http.HandlerFunc(dep.PriceAlertHandler.List)),
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.PriceAlertHandler.Create)),
	})
	mount(mux, "/api/price-alerts/{alertId}/disable", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.PriceAlertHandler.Disable)),
	})

	mount(mux, "/api/driver/trips", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.DriverTripHandler.CreateTrip)),
		http.MethodGet:  dep.AuthMiddleware(http.HandlerFunc(dep.DriverTripHandler.ListTrips)),
	})
	mount(mux, "/api/driver/dashboard", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.DriverTripHandler.GetDashboard)),
	})
	mount(mux, "/api/driver/income", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.DriverTripHandler.GetIncome)),
	})
	mount(mux, "/api/driver/profile", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.DriverProfileHandler.GetProfile)),
		http.MethodPut: dep.AuthMiddleware(http.HandlerFunc(dep.DriverProfileHandler.UpsertProfile)),
	})
	mount(mux, "/api/driver/vehicles", map[string]http.Handler{
		http.MethodGet:  dep.AuthMiddleware(http.HandlerFunc(dep.DriverProfileHandler.ListVehicles)),
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.DriverProfileHandler.CreateVehicle)),
	})

	mount(mux, "/api/driver/trips/{tripId}", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.DriverTripHandler.GetTripDetail)),
	})
	mount(mux, "/api/driver/orders/{orderId}/verify", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.VerifyOrderByDriver)),
	})
	mount(mux, "/api/driver/tickets/verify", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.ElectronicTicketHandler.VerifyByDriver)),
	})

	mount(mux, "/api/tickets/search", map[string]http.Handler{
		http.MethodGet: http.HandlerFunc(dep.TicketHandler.Search),
	})
	mount(mux, "/api/tickets/{ticketId}", map[string]http.Handler{
		http.MethodGet: http.HandlerFunc(dep.TicketHandler.Detail),
	})

	mount(mux, "/api/orders", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.CreateOrder)),
	})
	mount(mux, "/api/orders/my", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.ListMyOrders)),
	})
	mount(mux, "/api/orders/{orderId}", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.GetMyOrderDetail)),
	})
	mount(mux, "/api/orders/{orderId}/cancel", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.CancelMyOrder)),
	})
	mount(mux, "/api/orders/{orderId}/refund", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.RequestRefund)),
	})
	mount(mux, "/api/orders/{orderId}/ticket", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.ElectronicTicketHandler.GetMyOrderTicket)),
	})
	mount(mux, "/api/admin/orders", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.ListRefundOrdersForAdmin)),
	})
	mount(mux, "/api/admin/dashboard", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.GetAdminDashboardSummary)),
	})
	mount(mux, "/api/admin/knowledge", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.ListDocuments)),
	})
	mount(mux, "/api/admin/knowledge/upload", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.UploadDocument)),
	})
	mount(mux, "/api/admin/knowledge/search", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.SearchDocuments)),
	})
	mount(mux, "/api/admin/knowledge/{documentId}", map[string]http.Handler{
		http.MethodGet:    dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.GetDocument)),
		http.MethodPatch:  dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.UpdateDocumentStatus)),
		http.MethodDelete: dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.DeleteDocument)),
	})
	mount(mux, "/api/admin/knowledge/{documentId}/reindex", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.KnowledgeHandler.ReindexDocument)),
	})
	mount(mux, "/api/admin/orders/{orderId}/refund/approve", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.ApproveRefundByAdmin)),
	})
	mount(mux, "/api/admin/orders/{orderId}/refund/reject", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.OrderHandler.RejectRefundByAdmin)),
	})
	mount(mux, "/api/admin/orders/{orderId}/refund/audit-logs", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.AuditHandler.ListRefundAuditLogs)),
	})
	mount(mux, "/api/admin/audit-logs", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.AuditHandler.ListAuditLogs)),
	})
	mount(mux, "/api/payments/create", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.PaymentHandler.CreatePayment)),
	})
	mount(mux, "/api/payments/{paymentId}/status", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.PaymentHandler.GetPaymentStatus)),
	})

	mount(mux, "/api/ai/driver/create-trip", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.AIHandler.CreateDriverTripDraft)),
	})

	mount(mux, "/api/ai/chat", map[string]http.Handler{
		http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.AIHandler.ChatPassenger)),
	})

	mount(mux, "/api/admin/tokens", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.TokenUsageHandler.ListAdminUsage)),
	})
	mount(mux, "/api/admin/risk/logs", map[string]http.Handler{
		http.MethodGet: dep.AuthMiddleware(http.HandlerFunc(dep.RiskHandler.ListAdminRisks)),
	})
	mount(mux, "/api/admin/risk/logs/{eventId}", map[string]http.Handler{
		http.MethodPatch: dep.AuthMiddleware(http.HandlerFunc(dep.RiskHandler.UpdateRiskStatus)),
	})

	if dep.EnableMockPayment {
		mount(mux, "/api/payments/{paymentId}/mock-success", map[string]http.Handler{
			http.MethodPost: dep.AuthMiddleware(http.HandlerFunc(dep.PaymentHandler.MockPaySuccess)),
		})
	}

	return withRecovery(withCORS(mux))
}

func mount(mux *http.ServeMux, path string, handlers map[string]http.Handler) {
	mux.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h, ok := handlers[r.Method]
		if !ok {
			response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h.ServeHTTP(w, r)
	}))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v", rec)
				response.Error(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
