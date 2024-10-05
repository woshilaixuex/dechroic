package main

import (
	"flag"

	"github.com/delyr1c/dechoric/src/app/config"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	"github.com/delyr1c/dechoric/src/trigger/routers"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 启动入口
 * @Date: 2024-10-04 19:35
 */
var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	rdb := redis.NewRedisService(c.Redis.Host, c.Redis.Password, c.Redis.DB)
	sqlConn := sqlx.NewMysql(c.DB.MySqlDataSource)
	router := routers.SetupRouter(sqlConn, *rdb)
	// 启动 HTTP 服务
	router.Run(":8080")
}
