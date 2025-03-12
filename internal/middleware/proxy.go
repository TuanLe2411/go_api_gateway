package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func ProxyMiddleware(proxy *httputil.ReverseProxy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Proxying request to:", r.URL.Path)
		proxy.ServeHTTP(w, r)
	})
}
