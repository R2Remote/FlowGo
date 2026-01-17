package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"FLOWGO/internal/application/service"
	"FLOWGO/internal/application/service/devops"
	"FLOWGO/internal/infrastructure/config"
	"FLOWGO/internal/infrastructure/database"
	"FLOWGO/internal/infrastructure/redis"
	"FLOWGO/internal/infrastructure/repository"
	devopsRepo "FLOWGO/internal/infrastructure/repository/devops"
	"FLOWGO/internal/interfaces/http/handler"
	devopsHandler "FLOWGO/internal/interfaces/http/handler/devops"
	"FLOWGO/internal/interfaces/http/router"
	"FLOWGO/pkg/jwt"
)

func main() {
	// 加载配置文件（可通过环境变量指定配置文件路径）
	configPath := os.Getenv("CONFIG_PATH")
	if err := config.LoadConfig(configPath); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// 初始化Redis
	if err := redis.InitRedis(); err != nil {
		log.Printf("Warning: Failed to initialize redis: %v", err)
	}
	defer redis.CloseRedis()

	// 初始化JWT配置
	if config.AppConfig.JWT.PrivateKeyLocation != "" && config.AppConfig.JWT.PublicKeyLocation != "" {
		// 使用RSA密钥对（RS256）
		if err := jwt.SetRSAPrivateKey(config.AppConfig.JWT.PrivateKeyLocation); err != nil {
			log.Fatalf("Failed to load RSA private key: %v", err)
		}
		if err := jwt.SetRSAPublicKey(config.AppConfig.JWT.PublicKeyLocation); err != nil {
			log.Fatalf("Failed to load RSA public key: %v", err)
		}
		log.Println("JWT initialized with RSA key pair (RS256)")
	} else if config.AppConfig.JWT.SecretKey != "" {
		// 使用对称密钥（HS256）
		jwt.SetSecretKey(config.AppConfig.JWT.SecretKey)
		log.Println("JWT initialized with symmetric key (HS256)")
	} else {
		log.Fatal("JWT configuration is missing: either provide RSA key pair or secret key")
	}
	jwt.SetTokenExpiration(time.Duration(config.AppConfig.JWT.Expiration) * time.Hour)

	// 设置Gin模式
	ginMode := config.AppConfig.Server.Mode
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	// 依赖注入
	// 基础设施
	userRepo := repository.NewUserRepository(database.DB)
	projectRepo := repository.NewProjectsRepository(database.DB)
	devopsRepository := devopsRepo.NewDevOpsRepository(database.DB)

	// 应用服务
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	devopsService := devops.NewDevOpsService(devopsRepository)

	// 控制器
	userHandler := handler.NewUserHandler(userService)
	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectsHandler(projectService)
	doHandler := devopsHandler.NewDevOpsHandler(devopsService)

	// 设置路由
	r := router.SetupRouter(authHandler, userHandler, projectHandler, doHandler)

	// 启动服务器
	port := config.AppConfig.Server.Port
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  time.Duration(config.AppConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.AppConfig.Server.WriteTimeout) * time.Second,
	}

	// 优雅关闭
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
