package middlewares

import (
	"KnowledgeSharingPlatform/internal"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

var UserContextKey = "userID"

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the Bearer token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Split the header to get the token part
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Authorization header must be in the format 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]

		// Parse the token
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verify the token signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Return the key for validation
			return internal.GlobalConfig.JWTSecret, nil
		})

		// Handle any errors
		if err != nil {
			var errMessage string
			if err == jwt.ErrSignatureInvalid {
				errMessage = "Invalid token signature"
			} else {
				errMessage = "Invalid token"
			}
			http.Error(w, errMessage, http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Token is valid, proceed with the request
		ctx := context.WithValue(r.Context(), UserContextKey, claims.Subject)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
