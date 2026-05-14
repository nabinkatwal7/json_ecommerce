package service

import (
	"os"
	"strings"
	"time"

	"go-ecommerce-json/internal/auth"
	"go-ecommerce-json/internal/models"
	"go-ecommerce-json/internal/repository"

	"github.com/google/uuid"
)

type UserService struct {
	Store     *repository.Store
	JWTSecret []byte
	JWTTTL    time.Duration
}

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResult struct {
	Token string       `json:"token"`
	User  models.User  `json:"user"`
}

func (s *UserService) Register(in RegisterInput) (*AuthResult, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if in.Name == "" || email == "" || len(in.Password) < 8 {
		return nil, ErrValidation
	}
	existing, err := s.Store.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrConflict
	}
	hash, err := auth.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	role := "customer"
	if want := strings.TrimSpace(os.Getenv("ADMIN_EMAIL")); want != "" && email == strings.ToLower(strings.TrimSpace(want)) {
		role = "admin"
	}
	u := models.User{
		ID:           uuid.NewString(),
		Name:         strings.TrimSpace(in.Name),
		Email:        email,
		PasswordHash: hash,
		Role:         role,
		Addresses:    []models.Address{},
		CreatedAt:    now,
	}
	if err := s.Store.UpsertUser(u); err != nil {
		return nil, err
	}
	token, err := auth.SignJWT(s.JWTSecret, u.ID, u.Role, s.JWTTTL)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return &AuthResult{Token: token, User: u}, nil
}

func (s *UserService) Login(in LoginInput) (*AuthResult, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	if email == "" || in.Password == "" {
		return nil, ErrValidation
	}
	u, err := s.Store.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if u == nil || !auth.CheckPassword(u.PasswordHash, in.Password) {
		return nil, ErrUnauthorized
	}
	token, err := auth.SignJWT(s.JWTSecret, u.ID, u.Role, s.JWTTTL)
	if err != nil {
		return nil, err
	}
	out := *u
	out.PasswordHash = ""
	return &AuthResult{Token: token, User: out}, nil
}

func (s *UserService) GetProfile(userID string) (*models.User, error) {
	u, err := s.Store.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	u.PasswordHash = ""
	return u, nil
}

type PatchProfileInput struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

// UpdateProfile updates the caller's display name and/or email (partial JSON).
func (s *UserService) UpdateProfile(userID string, in PatchProfileInput) (*models.User, error) {
	if in.Name == nil && in.Email == nil {
		return nil, ErrValidation
	}
	u, err := s.Store.FindUserByID(userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrNotFound
	}
	if in.Name != nil {
		n := strings.TrimSpace(*in.Name)
		if n == "" {
			return nil, ErrValidation
		}
		u.Name = n
	}
	if in.Email != nil {
		e := strings.TrimSpace(strings.ToLower(*in.Email))
		if e == "" {
			return nil, ErrValidation
		}
		other, err := s.Store.FindUserByEmail(e)
		if err != nil {
			return nil, err
		}
		if other != nil && other.ID != u.ID {
			return nil, ErrConflict
		}
		u.Email = e
	}
	if err := s.Store.UpsertUser(*u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}
