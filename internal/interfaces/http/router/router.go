package router

import (
	"FLOWGO/internal/interfaces/http/handler"
	"FLOWGO/internal/interfaces/http/handler/devops"
	"FLOWGO/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	projectHandler *handler.ProjectsHandler,
	devopsHandler *devops.DevOpsHandler,
) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 认证相关路由（无需认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// 项目相关路由
		projects := v1.Group("/projects")
		projects.Use(middleware.Auth())
		{
			projects.GET("", projectHandler.ListProjects)
			projects.POST("", projectHandler.CreateProject)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
			projects.GET("/teams/available", projectHandler.ProjectTeams)
			projects.GET("/users/available/:id", projectHandler.ProjectAvailableUsers)
			projects.POST("/:id/users", projectHandler.AddProjectUsers)
			projects.DELETE("/:id/users/:uid", projectHandler.RemoveProjectUser)
		}

		// DevOps 路由 (全局)
		devops := v1.Group("/devops")
		devops.Use(middleware.Auth())
		{
			devops.POST("/config", devopsHandler.ConfigRepo)
			devops.GET("/summary", devopsHandler.GetSummary)
		}

		// Webhook 路由 (一般不需要认证，或者有专门的签名验证)
		webhooks := v1.Group("/webhooks")
		{
			webhooks.POST("/github", devopsHandler.HandleGitHubWebhook)
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
