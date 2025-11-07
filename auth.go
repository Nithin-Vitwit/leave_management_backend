package main

import (
	
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtSecret string
var jwtExpiryHours int

// Load JWT environment variables
func loadJWTConfig() {
	_ = godotenv.Load()
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret"
	}

	expStr := os.Getenv("JWT_EXP_HOURS")
	if expStr == "" {
		jwtExpiryHours = 2
	} else {
		if h, err := strconv.Atoi(expStr); err == nil && h > 0 {
			jwtExpiryHours = h
		} else {
			jwtExpiryHours = 2
		}
	}
}

// HR Login: returns JWT token
func hrLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Password != hrPassword {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub":  req.ID,
		"role": "hr",
		"exp":  time.Now().Add(time.Duration(jwtExpiryHours) * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "Failed to sign token", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"message": "HR authenticated successfully",
		"token":   signed,
		"expires": time.Now().Add(time.Duration(jwtExpiryHours) * time.Hour).Format(time.RFC3339),
	})
}

// Middleware to require HR JWT
func requireHRAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		tkn, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !tkn.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
