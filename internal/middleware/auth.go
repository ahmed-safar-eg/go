package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"project/internal/config"
	"project/pkg/utils"
	"strings"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for health check and public routes
        if r.URL.Path == "/health" || strings.HasPrefix(r.URL.Path, "/public") || strings.HasPrefix(r.URL.Path, "/auth/") {
			next.ServeHTTP(w, r)
			return
		}
		
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": false,
                "message": "Authorization header required",
            })
			return
		}
		
		// Remove "Bearer " prefix
		if len(tokenString) > 7 && strings.ToUpper(tokenString[0:6]) == "BEARER" {
			tokenString = tokenString[7:]
		}
		
        cfg := config.LoadConfig()
        claims, err := utils.ValidateJWTWithSecret(tokenString, cfg.JWT.Secret)
		if err != nil {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]interface{}{
                "success": false,
                "message": "Invalid token",
            })
			return
		}
		
		// Add claims to context
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}