package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
    DB *sql.DB
}

type UserService struct {
    DB *sql.DB
}

type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    UserID   int    `json:"id"`
    jwt.StandardClaims
}

var jwtKey = []byte("your_secret_key")

func (s *AuthService) Register(username, password, role string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    query := `INSERT INTO users (username, password, role) VALUES ($1, $2, $3)`
    _, err = s.DB.Exec(query, username, string(hashedPassword), role)
    if err != nil {
        return err
    }

    return nil
}

func (s *AuthService) Login(username, password string) (string, error) {
    var id int
    var  hashedPassword, role string // Declare two variables to hold password and role
    query := `SELECT id, password, role FROM users WHERE username = $1`
    
    // Scan both password and role from the result
    err := s.DB.QueryRow(query, username).Scan(&id, &hashedPassword, &role)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", errors.New("user not found")
        }
        return "", err
    }

    // Compare hashed password with the input password
    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    // Generate JWT token with role included
    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &Claims{
        Username: username,
        Role:     role,
        UserID:  id,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func (s *UserService) CountUsersByRole() (map[string]int, error) {
    query := `SELECT role, COUNT(*) AS count FROM users GROUP BY role`
    rows, err := s.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query user roles: %w", err)
    }
    defer rows.Close()

    // Map untuk menyimpan hasil
    roleCounts := make(map[string]int)
    for rows.Next() {
        var role string
        var count int
        if err := rows.Scan(&role, &count); err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }
        roleCounts[role] = count
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return roleCounts, nil
}