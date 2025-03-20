package middleware

import (
	"net/http"
	"net/http/httputil"
)

func ProxyMiddleware(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}
