package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "apguard-secret-change-in-production"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Register(req RegisterRequest) (*User, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("username, email and password are required")
	}
	if req.Role == "" {
		req.Role = "analyst"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	res, err := s.db.Exec(
		`INSERT INTO users (username, email, password_hash, role) VALUES (?,?,?,?)`,
		req.Username, req.Email, string(hash), req.Role,
	)
	if err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}
	id, _ := res.LastInsertId()
	return &User{ID: int(id), Username: req.Username, Email: req.Email, Role: req.Role}, nil
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	var u User
	err := s.db.QueryRow(
		`SELECT id, username, email, password_hash, role FROM users WHERE username=?`, req.Username,
	).Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role)
	if err == sql.ErrNoRows {
		return nil, errors.New("invalid credentials")
	}
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := generateJWT(u)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{Token: token, User: u}, nil
}

func (s *Service) GetUserByID(id int) (*User, error) {
	var u User
	err := s.db.QueryRow(
		`SELECT id, username, email, role FROM users WHERE id=?`, id,
	).Scan(&u.ID, &u.Username, &u.Email, &u.Role)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Service) ListUsers() ([]User, error) {
	rows, err := s.db.Query(`SELECT id, username, email, role, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var u User
		var createdAt string
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &createdAt); err != nil {
			continue
		}
		users = append(users, u)
	}
	return users, nil
}

func generateJWT(u User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"user": u.Username,
		"role": u.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
