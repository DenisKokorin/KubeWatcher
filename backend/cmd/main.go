package main

import (
	"errors"
	"flag"
	"k8s-mon/config"
	externalapi "k8s-mon/internal/external"
	"k8s-mon/internal/handler"
	"k8s-mon/internal/middleware"
	"k8s-mon/internal/monitor"
	"k8s-mon/internal/repository"
	reviewservice "k8s-mon/internal/review_service"
	"k8s-mon/internal/roles"
	service "k8s-mon/internal/user_service"
	"k8s-mon/pkg/db"
	"k8s-mon/pkg/jwt"
	"k8s-mon/pkg/logger"
	"k8s-mon/pkg/redis"
	"k8s-mon/pkg/storage"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
	"k8s.io/client-go/kubernetes"
)

func main() {
	cfg := config.InitConfig("env/.env")
	db := db.InitDB(cfg.PostgresURL)
	redisDB := redis.InitRedisDB(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)

	log := logger.SetUp(cfg.LogLevel)

	JWT := jwt.NewUserJWTpkg(cfg.JWTSecret, cfg.AccessTTL)
	jwtMiddleware := middleware.NewJWTMiddleware(JWT)
	roleMiddleware := middleware.NewRoleMiddleware(JWT)

	userRepo := repository.NewUserRepository(db, redisDB)

	// Initialize storage service
	storageConfig := storage.StorageConfig{
		Endpoint:        cfg.StorageEndpoint,
		AccessKeyID:     cfg.StorageAccessKey,
		SecretAccessKey: cfg.StorageSecretKey,
		BucketName:      cfg.StorageBucketName,
		UseSSL:          cfg.StorageUseSSL,
	}
	storageSvc, err := storage.NewMinIOStorage(storageConfig, log)
	if err != nil {
		log.Error("Failed to initialize storage service", "error", err)
		panic(err)
	}

	userService := service.NewUserService(userRepo, storageSvc, log, JWT, cfg.RefreshTTL)
	userHandler := handler.NewUserHandler(userService)

	externalSvc, err := externalapi.NewExternalAPI(cfg.ExternalAPIURL, cfg.ExternalAPITimeout, cfg.ExternalAPIRetry, cfg.ExternalAPIRate)
	if err != nil {
		log.Error("Failed to initialize external API service", "error", err)
		panic(err)
	}
	externalHandler := handler.NewExternalHandler(externalSvc)

	revRepo := repository.NewReviewRepository(db)
	revService := reviewservice.NewReviewService(revRepo, log)
	revHandler := handler.NewReviewHandler(revService)

	var kubeconfig string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	flag.Parse()

	k8s_cfg, err := config.LoadKubeconfig(kubeconfig)
	if err != nil {
		log.Warn("Kubernetes config unavailable, continuing without cluster monitoring", "error", err)
	}

	var monitorService *monitor.Monitor
	if err == nil {
		clientset, err := kubernetes.NewForConfig(k8s_cfg)
		if err != nil {
			log.Warn("Failed to initialize Kubernetes client, continuing without cluster monitoring", "error", err)
		} else {
			monitorService = monitor.NewMonitor(clientset)
		}
	}

	if monitorService == nil {
		monitorService = monitor.NewMonitor(nil)
	}

	monHandler := handler.NewMonitorHandler(*monitorService)

	seoHandler := handler.NewSEOHandler(cfg.SiteBaseURL)
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var httpErr *echo.HTTPError
		code := http.StatusInternalServerError
		if errors.As(err, &httpErr) {
			code = httpErr.Code
		}

		path := c.Request().URL.Path
		if path == "/robots.txt" || path == "/sitemap.xml" {
			c.String(code, http.StatusText(code))
			return
		}

		if strings.HasPrefix(c.Request().URL.Path, "/api/") {
			c.JSON(code, map[string]interface{}{
				"error": http.StatusText(code),
			})
			return
		}

		c.String(code, http.StatusText(code))
	}
	e.Use(mid.CORSWithConfig(mid.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:8080"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	}))
	e.Static("/uploads", "uploads")

	auth := e.Group("/api/v1/auth")
	{
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)
		auth.POST("/refresh", userHandler.Refresh)
		auth.POST("/logout", userHandler.Logout)
		auth.GET("/me", userHandler.GetMe, jwtMiddleware.Auth)
	}

	mon := e.Group("/api/v1/monitor", roleMiddleware.RequirePermission(roles.PermMonitorView))
	{
		mon.GET("/cluster", monHandler.Overview)
		mon.GET("/nodes", monHandler.Nodes)
		mon.GET("/pods", monHandler.Pods)
	}

	e.GET("/robots.txt", seoHandler.Robots)
	e.GET("/sitemap.xml", seoHandler.Sitemap)

	rev := e.Group("/api/v1/reviews", roleMiddleware.RequirePermission(roles.PermReviewView))
	{
		rev.GET("/employees", revHandler.GetAllEmployees)
		rev.GET("/employees/:id", revHandler.GetEmployee)
		rev.GET("/employees/:id/reviews", revHandler.GetEmployeeWithReviews)
		rev.POST("", roleMiddleware.RequirePermission(roles.PermReviewCreate)(revHandler.CreateReview))
		rev.GET("/employees/:team", revHandler.GetEmployeesByTeam)
		rev.GET("/team/:team", revHandler.GetTeamStats)
	}

	admin := e.Group("/api/v1/admin", roleMiddleware.RequireRole(roles.RoleAdmin))
	{
		admin.GET("/users", userHandler.GetAllUsers)
		admin.POST("/users", userHandler.CreateUser)
		admin.PUT("/users/:id/role", userHandler.UpdateUserRole)
		admin.POST("/users/:id/document", userHandler.UploadUserDocument)
		admin.GET("/users/:id/document", userHandler.GetUserDocumentURL)
		admin.DELETE("/users/:id/document", userHandler.DeleteUserDocument)
		admin.DELETE("/users/:id", userHandler.DeleteUser)
	}

	external := e.Group("/api/v1/external", jwtMiddleware.Auth)
	{
		external.GET("/posts", externalHandler.GetPosts)
	}

	e.GET("healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status": "ok",
		})
	})

	e.Logger.Fatal(e.Start(":" + cfg.HTTPPort))
}
