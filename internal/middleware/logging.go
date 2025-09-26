package middleware

import (
	"encoding/json"
	"net/http"
	"time"
)

// JSONMiddleware ensures Content-Type header and recovers panics with JSON.
func JSONMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        defer func(start time.Time) {
            if rec := recover(); rec != nil {
                w.WriteHeader(http.StatusInternalServerError)
                json.NewEncoder(w).Encode(map[string]interface{}{
                    "success": false,
                    "message": "Internal server error",
                })
            }
        }(time.Now())
        next.ServeHTTP(w, r)
    })
}