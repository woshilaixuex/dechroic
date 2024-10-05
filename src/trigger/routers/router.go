package routers

import (
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	tr_http "github.com/delyr1c/dechoric/src/trigger/http"
	"github.com/delyr1c/dechoric/src/trigger/ws"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 路由注册
 * @Date: 2024-10-05 10:02
 */
func SetupRouter(sqlConn sqlx.SqlConn, redis redis.RedisService) *gin.Engine {
	router := gin.Default()
	messageRepo := infra_repository.NewMessageRepository(sqlConn, redis)
	userRepo := infra_repository.NewUserRepository(sqlConn, redis)
	// 服务
	userService := tr_http.NewUserService(userRepo)
	MessageService := tr_http.NewMessageService(messageRepo)
	wsHandler := ws.NewWebSocketHandler("your-api-key", messageRepo)
	router.GET("/ws")
	api := router.Group("/api/user")
	{
		api.POST("/register", userService.Register)
		api.POST("/login", userService.Login)
	}
	return router
}
