package middleware

import (
	"api_gateway/pkg/constant"
	"api_gateway/pkg/utils"
	"encoding/json"
	"net/http"
)

func ErrorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		if err, ok := r.Context().Value(constant.AppErrorContextKey).(utils.AppError); ok {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err)
		}
	})
}
