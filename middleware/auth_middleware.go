package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key")

// Claims struct untuk JWT
type Claims struct {
    UserID   int    `json:"id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

// Middleware AuthMiddleware
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            log.Println("Missing Authorization header")
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
            log.Println("Invalid Authorization header format")
            http.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }

        tokenStr := authHeader[7:]
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil {
            log.Printf("Error parsing token: %v\n", err)
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if !token.Valid {
            log.Println("Token is invalid")
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        log.Printf("Valid token: UserID: %d, Username: %s, Role: %s\n", claims.UserID, claims.Username, claims.Role)
        ctx := context.WithValue(r.Context(), "id", claims.UserID)
        ctx = context.WithValue(ctx, "username", claims.Username)
        ctx = context.WithValue(ctx, "role", claims.Role)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Middleware RoleMiddleware
func RoleMiddleware(allowedRoles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("Role Middleware: Checking roles")
			role, ok := r.Context().Value("role").(string)
			if !ok || role == "" {
				http.Error(w, "Unauthorized: Missing or invalid claims", http.StatusUnauthorized)
				return
			}

			log.Println("Role from context:", role) // Debugging untuk memeriksa role

			// Check if the role is allowed
			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: You don't have access to this resource", http.StatusForbidden)
		})
	}
}

// Fungsi GenerateToken untuk login
func GenerateToken(id int, username, role string) (string, error) {
    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        UserID:   id,
        Username: username,
        Role:     role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}
