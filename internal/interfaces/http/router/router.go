package router

import (
	"FLOWGO/internal/interfaces/http/handler"
	"FLOWGO/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	projectHandler *handler.ProjectsHandler,
	statsHandler *handler.StatsHandler,
) *gin.Engine {
	r := gin.New()

	// 中间件
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.CORS())
	r.Use(middleware.VisitLogger())

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

		// 用户相关路由
		users := v1.Group("/users")
		users.Use(middleware.Auth())
		{
			users.POST("", userHandler.CreateUser)
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
		}

		// 统计 API
		stats := v1.Group("/stats")
		// stats.Use(middleware.Auth()) // 根据需求决定是否需要鉴权，这里暂时公开方便查看，或者加 Auth
		{
			stats.GET("/visits", statsHandler.GetVisitStats)
		}
	}

	return r
}
