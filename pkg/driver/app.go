package driver

import (
	"api_gateway/internal/drivers/app_controller"
	"api_gateway/internal/middleware"
	"api_gateway/pkg"
	"api_gateway/pkg/constant"
	"api_gateway/pkg/utils"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func Run() {
	routes, err := pkg.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	jwtAccessTokenTtl, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_TTL_S"))
	jwtRefreshTokenTtl, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_TTL_S"))
	jwt := utils.NewJwt(
		os.Getenv("JWT_ACCESS_TOKEN_SECRET"),
		os.Getenv("JWT_REFRESH_TOKEN_SECRET"),
		jwtAccessTokenTtl,
		jwtRefreshTokenTtl,
	)

	router := mux.NewRouter()

	// Middleware cho toàn bộ router
	router.Use(
		middleware.XssProtectionMiddleware,
		middleware.CorsMiddleware,
		middleware.NewJwtMiddleware(jwt).Do,
		middleware.ErrorHandlerMiddleware,
	)

	for _, route := range routes.Routes {
		proxy, err := utils.NewProxy(route.Target)
		if err != nil {
			log.Fatalf("Failed to create proxy for %s: %v", route.Name, err)
		}
		handler := middleware.ProxyMiddleware(proxy)
		router.Handle(route.Context+"/{.*}", handler).Methods("GET", "POST", "PUT", "DELETE")
		router.Handle(route.Context, handler).Methods("GET", "POST", "PUT", "DELETE")
	}

	appController := app_controller.AppController{}
	router.HandleFunc("/health", appController.HealthCheck).Methods(constant.GetMethod)

	log.Println("Server is running on port: " + os.Getenv("SERVER_PORT"))
	http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), router)
}
