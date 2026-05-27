// Package router 定義了 Gin 應用程式的路由設置
package router

import (
	"os"
	"path/filepath"
	"self_go_gin/container"
	"self_go_gin/domains/user/events"
	v1_admin "self_go_gin/gin_application/api/v1/admin"
	v1_user "self_go_gin/gin_application/api/v1/user"
	middleware "self_go_gin/gin_application/middleware"
	"self_go_gin/infra/event"
	"self_go_gin/infra/log/zaplog"
	"time"

	_ "self_go_gin/docs" // Swagger 文档生成需要导入

	"github.com/gin-contrib/cors"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func setDefaultMiddlewares(router *gin.Engine) {
	// Trace 中間件（確保所有請求都有 ID）
	router.Use(middleware.TraceMiddleware())
	serverConfig := container.GetContainer().GetConfig()
	if serverConfig.IsMetrics {
		// Metrics 中間件（記錄 HTTP 請求指標）
		router.Use(middleware.MetricsMiddleware())
	}
	// 獲取當前工作目錄並構建 log 路徑
	cwd, err := os.Getwd()
	if err != nil {
		panic("無法獲取工作目錄: " + err.Error())
	}
	logPath := filepath.Join(cwd, "log") + string(filepath.Separator)
	zapLogger := zaplog.GetZapLogger(logPath)
	// Add a ginzap middleware
	router.Use(ginzap.Ginzap(zapLogger, "", true))

	// Logs all panic to error log
	router.Use(ginzap.RecoveryWithZap(zapLogger, true))
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Content-Type", "Authorization", "Access-Control-Allow-Origin"},
	})) //跨域請求的中間件
}

// Router 路由
func Router() *gin.Engine {
	router := gin.New()
	// pprof.Register(router)
	setDefaultMiddlewares(router)
	registerSwagger(router)
	apiV1Group := router.Group("/api/v1")
	setNoAuthRoutes(apiV1Group)
	setAuthRoutes(apiV1Group)
	return router
}

func registerSwagger(router *gin.Engine) {
	if gin.Mode() != gin.ReleaseMode {
		router.GET("/swagger-test/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func setNoAuthRoutes(apiV1Group *gin.RouterGroup) {
	apiV1UsersGroup := apiV1Group.Group("/users")
	apiV1AdminsGroup := apiV1Group.Group("/admins")

	// Health check endpoint
	apiV1Group.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	serverConfig := container.GetContainer().GetConfig()
	if serverConfig.IsMetrics {
		apiV1Group.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}
	// apiV1Group.Use(middleware.RateLimit("test-limit")).GET("/limit_ping", func(c *gin.Context) {
	// 	c.String(200, "pong "+fmt.Sprint(time.Now().Unix()))
	// })

	apiV1Group.GET("/logtest", func(_ *gin.Context) {
		test := true
		if test {
			zap.L().Info("Logger  Success..",
				zap.String("GgGGGG", "200"))
		} else {
			zap.L().Error(
				"Logger  Error ..",
				zap.String("test", "just for test"))
		}
	})
	Login(apiV1UsersGroup, apiV1AdminsGroup)
	apiV1UsersGroup.POST("createUser", v1_user.CreateUser)

}

func setAuthRoutes(apiV1Group *gin.RouterGroup) {
	// apiV1AuthGroup := apiV1Group.Group("/auth")
	apiV1Group.Use(middleware.JwtAuthMiddleware())

	// Users
	apiV1AuthUsersGroup := apiV1Group.Group("/users")
	Users(apiV1AuthUsersGroup)

	// Admins
	apiV1AuthAdminsGroup := apiV1Group.Group("/admins")
	Admins(apiV1AuthAdminsGroup)

}

// =================================   no auth group   =====================================

// Login 登入
func Login(userRouter, adminRouter *gin.RouterGroup) {
	userRouter.POST("/login", v1_user.UserLogin)
	adminRouter.POST("/login", v1_admin.AdminLogin)
}

// =================================   auth group   =====================================

// Users 用戶
func Users(router *gin.RouterGroup) {
	router.Use(middleware.OpaMiddleware())
	router.GET("/:filterUsersID", v1_user.GetUsersByID)
}

// Admins 管理員
func Admins(router *gin.RouterGroup) {
	router.Use(middleware.OpaMiddleware())
	router.GET("/:filterAdminsID", v1_admin.GetAdminsByID)
	Shutdown(router)
}

// Shutdown 優雅關閉服務測試
func Shutdown(router *gin.RouterGroup) {
	router.GET("/slow_test", func(c *gin.Context) {
		time.Sleep(5 * time.Second) // 模擬慢速API
		app := container.GetContainer()
		if app.GetConfig().IsEventBroker {
			broker := app.GetEventBroker()
			if broker.Publisher() != nil {
				payload := events.UserCreatedEventPayload{
					UserID:   1234,
					Account:  "slow_test_user",
					CreateAt: time.Now().Format(time.RFC3339),
				}
				traceID, ok := c.Get("traceID")
				if !ok {
					traceID = "no-trace-id"
				}
				evt, _ := event.NewEvent(events.UserCreatedEventType, payload, "slow-test", traceID.(string), time.Now())

				err := broker.Publisher().Publish(c, evt)
				if err != nil {
					zap.L().Error("Failed to publish event", zap.Error(err))
				}
			}
		}
		zap.L().Info("慢速API完成")
		c.String(200, "shutdown slow test")
	})
}
