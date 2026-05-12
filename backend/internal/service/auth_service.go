package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"ridehailing/backend/internal/model"
	"ridehailing/backend/internal/pkg/jwtutil"
	"ridehailing/backend/internal/pkg/mailer"
	"ridehailing/backend/internal/repository"
)

var phoneRegexp = regexp.MustCompile(`^\d{11}$`)
var emailRegexp = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)

type AuthResult struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

type AuthService struct {
	userRepo repository.UserRepository
	codeRepo repository.CodeRepository
	jwt      *jwtutil.Manager
	codeTTL  time.Duration
	mailer   mailer.Mailer
}

func NewAuthService(
	userRepo repository.UserRepository,
	codeRepo repository.CodeRepository,
	jwt *jwtutil.Manager,
	codeTTL time.Duration,
	mailer mailer.Mailer,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		codeRepo: codeRepo,
		jwt:      jwt,
		codeTTL:  codeTTL,
		mailer:   mailer,
	}
}

func (s *AuthService) SendEmailCode(ctx context.Context, email, scene string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	scene = strings.TrimSpace(scene)

	if err := validateEmail(email); err != nil {
		return "", err
	}
	if err := validateScene(scene); err != nil {
		return "", err
	}

	code, err := generateCode(6)
	if err != nil {
		return "", err
	}

	if err := s.codeRepo.Set(ctx, email, scene, code, s.codeTTL); err != nil {
		return "", err
	}
	if err := s.mailer.SendVerificationCode(ctx, email, scene, code); err != nil {
		return "", err
	}
	return code, nil
}

func (s *AuthService) Register(ctx context.Context, role, phone, email, emailCode, password, nickname string) (*AuthResult, error) {
	role = normalizeRole(role)
	phone = strings.TrimSpace(phone)
	email = strings.ToLower(strings.TrimSpace(email))
	emailCode = strings.TrimSpace(emailCode)
	password = strings.TrimSpace(password)
	nickname = strings.TrimSpace(nickname)

	if role != model.RolePassenger && role != model.RoleDriver {
		return nil, errors.New("register role must be passenger or driver")
	}
	if err := validatePhone(phone); err != nil {
		return nil, err
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}
	if len(password) < 6 {
		return nil, errors.New("password length must be at least 6")
	}

	ok, err := s.codeRepo.Verify(ctx, email, "register", emailCode)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid email code")
	}

	exist, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if exist != nil {
		return nil, errors.New("phone already registered")
	}

	existByEmail, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existByEmail != nil {
		return nil, errors.New("email already registered")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if nickname == "" {
		if role == model.RoleDriver {
			nickname = "新司机"
		} else {
			nickname = "新乘客"
		}
	}

	user := &model.User{
		Phone:        phone,
		Email:        email,
		PasswordHash: string(passwordHash),
		Nickname:     nickname,
		Role:         role,
		DefaultRole:  role,
		Status:       model.UserStatusActive,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.jwt.Generate(user)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) LoginByPassword(ctx context.Context, role, phone, password string) (*AuthResult, error) {
	role = normalizeRole(role)
	phone = strings.TrimSpace(phone)
	password = strings.TrimSpace(password)

	if role == "" {
		return nil, errors.New("role is required")
	}
	if err := validatePhone(phone); err != nil {
		return nil, err
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	user, err := s.userRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid phone or password")
	}
	if user.Role != role {
		return nil, errors.New("current account role does not match")
	}
	if user.Status != model.UserStatusActive {
		return nil, errors.New("account status is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid phone or password")
	}

	token, err := s.jwt.Generate(user)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) LoginByCode(ctx context.Context, role, email, emailCode string) (*AuthResult, error) {
	role = normalizeRole(role)
	email = strings.ToLower(strings.TrimSpace(email))
	emailCode = strings.TrimSpace(emailCode)

	if role == "" {
		return nil, errors.New("role is required")
	}
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	ok, err := s.codeRepo.Verify(ctx, email, "login", emailCode)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("invalid email code")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("account not found")
	}
	if user.Role != role {
		return nil, errors.New("current account role does not match")
	}
	if user.Status != model.UserStatusActive {
		return nil, errors.New("account status is not active")
	}

	token, err := s.jwt.Generate(user)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Logout(_ context.Context, _ uint) error {
	return nil
}

func (s *AuthService) Me(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func validatePhone(phone string) error {
	if !phoneRegexp.MatchString(phone) {
		return errors.New("phone must be 11 digits")
	}
	return nil
}

func validateEmail(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func validateScene(scene string) error {
	switch scene {
	case "register", "login", "reset_password":
		return nil
	default:
		return errors.New("invalid email code scene")
	}
}

func normalizeRole(role string) string {
	return strings.ToLower(strings.TrimSpace(role))
}

func generateCode(length int) (string, error) {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = byte('0' + n.Int64())
	}
	return string(result), nil
}
