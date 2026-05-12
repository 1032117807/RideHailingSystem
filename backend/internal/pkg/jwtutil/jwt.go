package jwtutil

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"ridehailing/backend/internal/model"
)

type Claims struct {
	UserID uint   `json:"userId"`
	Role   string `json:"role"`
	Phone  string `json:"phone"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret     []byte
	expireTime time.Duration
	issuer     string
}

func NewManager(secret string, expireTime time.Duration) *Manager {
	return &Manager{
		secret:     []byte(secret),
		expireTime: expireTime,
		issuer:     "ridehailing-backend",
	}
}

func (m *Manager) GenerateToken(user *model.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		Phone:  user.Phone,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expireTime)),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Generate(user *model.User) (string, error) {
	return m.GenerateToken(user)
}

func (m *Manager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
