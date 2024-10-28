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
	userService := tr_http.NewUserService(userRepo)

	raffleRepo := infra_repository.NewStrategyRepository(sqlConn, redis)
	raffleService := tr_http.NewRaffeService(raffleRepo)
	MessageService := tr_http.NewMessageService(messageRepo)
	// websockt
	wsHandler := ws.NewWebSocketHandler("sk-97c2bf6ddc154c239a29369491429041", messageRepo)
	router.GET("", wsHandler.HandleWebSocket)
	apiUser := router.Group("/api/user")
	{
		apiUser.POST("/register", userService.Register)
		apiUser.POST("/login", userService.Login)
		apiUser.GET("/info", userService.GetUserInfo)
	}
	apiRaffle := router.Group("/api/raffle")
	{
		apiRaffle.POST("/prizes", raffleService.GetPrize)
		apiRaffle.POST("/award", raffleService.GetAwardPrize)
		apiRaffle.POST("/history", raffleService.GetHistory)
	}
	apiMess := router.Group("/api/message")
	{
		apiMess.POST("/ai_models", MessageService.GetUserAIInfo)
	}
	return router
}
