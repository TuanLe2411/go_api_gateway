package middleware

import (
	"api_gateway/pkg/constant"
	"api_gateway/pkg/utils"
	"fmt"
	"net/http"
	"strings"
)

type JwtMiddleware struct {
	*utils.Jwt
}

func NewJwtMiddleware(jwt *utils.Jwt) constant.Middleware {
	return &JwtMiddleware{
		jwt,
	}
}

func (j *JwtMiddleware) Do(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if j.isIgnoredPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		token := r.Header.Get("Authorization")
		if token == "" {
			utils.SetHttpReponseError(r, utils.ErrUnAuthorized)
			return
		}
		token = token[len("Bearer "):]
		isValid, claims := j.ValidateToken(token)
		if !isValid {
			utils.SetHttpReponseError(r, utils.ErrUnAuthorized)
			return
		}

		r.Header.Set("user_id", fmt.Sprintf("%d", claims.UserId))
		r.Header.Set("role", claims.Username)
		next.ServeHTTP(w, r)
	})
}

func (j *JwtMiddleware) isIgnoredPath(path string) bool {
	IgnoredPath := []string{"/auth"}
	for _, p := range IgnoredPath {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
