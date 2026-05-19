package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/repository"
)

const electronicTicketTTL = 30 * 24 * time.Hour

type ElectronicTicketService struct {
	ticketRepo repository.ElectronicTicketRepository
	orderRepo  repository.OrderRepository
	secret     string
}

func NewElectronicTicketService(ticketRepo repository.ElectronicTicketRepository, orderRepo repository.OrderRepository, secret string) *ElectronicTicketService {
	if strings.TrimSpace(secret) == "" {
		secret = "replace-with-your-secret"
	}
	return &ElectronicTicketService{ticketRepo: ticketRepo, orderRepo: orderRepo, secret: secret}
}

func (s *ElectronicTicketService) EnsureForOrder(ctx context.Context, order *model.Order) (*model.ElectronicTicket, error) {
	if order == nil || order.ID == 0 {
		return nil, errors.New("order is required")
	}
	if order.PayStatus != model.PayStatusPaid {
		return nil, errors.New("only paid order can issue electronic ticket")
	}
	if order.OrderStatus != model.OrderStatusPendingVerification && order.OrderStatus != model.OrderStatusCompleted {
		return nil, errors.New("current order status cannot issue electronic ticket")
	}
	if order.RefundStatus == model.RefundStatusRequested || order.RefundStatus == model.RefundStatusRefunded {
		return nil, errors.New("refund processing order cannot issue electronic ticket")
	}

	existing, err := s.ticketRepo.GetByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	expiresAt := time.Now().Add(electronicTicketTTL)
	token, tokenHash, err := s.generateToken(order, expiresAt)
	if err != nil {
		return nil, err
	}

	ticket := &model.ElectronicTicket{
		OrderID:   order.ID,
		UserID:    order.UserID,
		TripID:    order.TripID,
		Token:     token,
		TokenHash: tokenHash,
		Status:    model.ElectronicTicketStatusIssued,
		ExpiresAt: expiresAt,
	}
	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, err
	}
	return s.ticketRepo.GetByOrderID(ctx, order.ID)
}

func (s *ElectronicTicketService) GetMyOrderTicket(ctx context.Context, currentUserID uint, currentUserRole string, orderID uint) (*model.ElectronicTicket, error) {
	if currentUserRole != model.RolePassenger {
		return nil, errors.New("only passenger can view electronic ticket")
	}
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != currentUserID {
		return nil, errors.New("order does not belong to current user")
	}
	return s.EnsureForOrder(ctx, order)
}

func (s *ElectronicTicketService) VerifyByDriver(ctx context.Context, currentUserID uint, currentUserRole string, token string) (*model.ElectronicTicket, error) {
	if currentUserRole != model.RoleDriver {
		return nil, errors.New("only driver can verify electronic ticket")
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("ticket token is required")
	}
	if err := s.validateTokenSignature(token); err != nil {
		return nil, err
	}

	tokenHash := hashToken(token)
	ticket, err := s.ticketRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if ticket == nil {
		return nil, errors.New("electronic ticket not found")
	}

	verifyResult := model.TicketVerificationResultFailed
	verifyMessage := ""
	defer func() {
		_ = s.ticketRepo.CreateVerification(ctx, &model.TicketVerification{
			TicketID: ticket.ID,
			OrderID:  ticket.OrderID,
			DriverID: currentUserID,
			Result:   verifyResult,
			Message:  verifyMessage,
		})
	}()

	if time.Now().After(ticket.ExpiresAt) {
		verifyMessage = "electronic ticket expired"
		return nil, errors.New(verifyMessage)
	}
	if ticket.Trip.DriverID != currentUserID {
		verifyMessage = "ticket does not belong to current driver trip"
		return nil, errors.New(verifyMessage)
	}
	if ticket.Order.PayStatus != model.PayStatusPaid {
		verifyMessage = "order is not paid"
		return nil, errors.New(verifyMessage)
	}
	if ticket.Order.RefundStatus == model.RefundStatusRequested || ticket.Order.RefundStatus == model.RefundStatusRefunded {
		verifyMessage = "refund processing order cannot be verified"
		return nil, errors.New(verifyMessage)
	}
	if ticket.Status == model.ElectronicTicketStatusVerified {
		verifyMessage = "electronic ticket already verified"
		return ticket, nil
	}
	if ticket.Status != model.ElectronicTicketStatusIssued {
		verifyMessage = "electronic ticket is not valid"
		return nil, errors.New(verifyMessage)
	}

	verifiedAt := time.Now()
	updated, err := s.ticketRepo.MarkVerified(ctx, ticket.ID, currentUserID, verifiedAt)
	if err != nil {
		verifyMessage = err.Error()
		return nil, err
	}
	verifyResult = model.TicketVerificationResultSuccess
	verifyMessage = "verified"
	return updated, nil
}

func (s *ElectronicTicketService) VoidByOrderID(ctx context.Context, orderID uint) error {
	if orderID == 0 {
		return nil
	}
	return s.ticketRepo.VoidByOrderID(ctx, orderID)
}

func (s *ElectronicTicketService) generateToken(order *model.Order, expiresAt time.Time) (string, string, error) {
	nonceBytes := make([]byte, 12)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", "", err
	}
	payload := fmt.Sprintf("%d:%d:%d:%d:%s", order.ID, order.UserID, order.TripID, expiresAt.Unix(), hex.EncodeToString(nonceBytes))
	signature := s.sign(payload)
	rawToken := payload + "." + signature
	token := base64.RawURLEncoding.EncodeToString([]byte(rawToken))
	return token, hashToken(token), nil
}

func (s *ElectronicTicketService) validateTokenSignature(token string) error {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return errors.New("invalid electronic ticket token")
	}
	parts := strings.Split(string(decoded), ".")
	if len(parts) != 2 {
		return errors.New("invalid electronic ticket token")
	}
	payload := parts[0]
	signature := parts[1]
	if !hmac.Equal([]byte(signature), []byte(s.sign(payload))) {
		return errors.New("invalid electronic ticket signature")
	}
	payloadParts := strings.Split(payload, ":")
	if len(payloadParts) != 5 {
		return errors.New("invalid electronic ticket payload")
	}
	expiresUnix, err := strconv.ParseInt(payloadParts[3], 10, 64)
	if err != nil {
		return errors.New("invalid electronic ticket expiration")
	}
	if time.Now().After(time.Unix(expiresUnix, 0)) {
		return errors.New("electronic ticket expired")
	}
	return nil
}

func (s *ElectronicTicketService) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
