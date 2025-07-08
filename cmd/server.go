package main

import (
	"context"
	"log"
	"net/http"

	"simple_gin_server/configs"
	"simple_gin_server/internal/auth"
	"simple_gin_server/internal/users"
	"simple_gin_server/pkg/db"
	"simple_gin_server/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer  *http.Server
	router      *gin.Engine
	authHandler *auth.AuthHandler
	config      *configs.Config
	db          db.PgRepoInterface
}

// Конструктор для сервера
func NewServer(ctx context.Context) *Server {
	conf := configs.LoadConfig()

	// создаём экземпляр пула соединений на базе конфига и контекста
	db_pg, err := db.NewPgRepo(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}

	// создаём экземпляр reddis, используя config
	redisRepo := db.NewRedisRepo(ctx, conf)

	// Инициализация слоёв приложения
	userRepository := users.NewUserRepository(db_pg)
	authService := auth.NewAuthService(userRepository, redisRepo)
	authHandler := auth.NewAuthHandler(authService, conf)

	// создаём экземпляр роутера
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Добавляем middleware для проброса контекста
	router.Use(func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), "request_id", c.GetHeader("X-Request-ID"))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	return &Server{
		router:      router,
		authHandler: authHandler,
		config:      conf,
		db:          db_pg,
	}
}

// Метод для маршрутизации сервера
func (s *Server) SetUpRoutes() {
	s.router.POST("/register", middleware.ValidateAuthMiddleware(&auth.RegisterRequest{}), s.authHandler.RegisterHandler)
	s.router.POST("/login", middleware.ValidateAuthMiddleware(&auth.LoginRequest{}), s.authHandler.LoginHandler)
	s.router.POST("auth/refresh", s.authHandler.ProcessRefreshTokenHandler)

	authGroup := s.router.Group("/")
	authGroup.Use(middleware.AuthMiddleware(s.config))
	{
		authGroup.GET("/health", s.authHandler.Check)
		authGroup.GET("/list", s.authHandler.ListHandler)
		authGroup.POST("/logout", s.authHandler.LogoutHandler)
	}
}

// Метод для запуска сервера
func (s *Server) Run() error {
	s.SetUpRoutes()

	s.httpServer = &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}
	log.Println("Server is running on port 8080")
	return s.httpServer.ListenAndServe()
}

// Метод для graceful shutdown
func (s *Server) Shutdown(ctx context.Context) error {
	// Закрываем соединение с БД
	s.db.Close()

	// Останавливаем HTTP сервер
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server shutdown completed")
	return nil
}
