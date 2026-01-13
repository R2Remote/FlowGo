package router

import (
	"FLOWGO/internal/application/usecase"
	"FLOWGO/internal/infrastructure/database"
	"FLOWGO/internal/infrastructure/repository"
	"FLOWGO/internal/interfaces/http/handler"
	"FLOWGO/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 依赖注入
	userRepo := repository.NewUserRepository(database.DB)
	createUserUseCase := usecase.NewCreateUserUseCase(userRepo)
	getUserUseCase := usecase.NewGetUserUseCase(userRepo)
	listUsersUseCase := usecase.NewListUsersUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo)

	userHandler := handler.NewUserHandler(createUserUseCase, getUserUseCase, listUsersUseCase)
	authHandler := handler.NewAuthHandler(loginUseCase)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 用户相关路由
		users := v1.Group("/users")
		users.Use(middleware.Auth())
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
		}
	}

	return r
}
