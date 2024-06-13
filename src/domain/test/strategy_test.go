/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 装配策略测试
 * @Date: 2024-06-13 22:23
 */
package test

import (
	"context"
	"encoding/json"
	"flag"
	"testing"

	"github.com/delyr1c/dechoric/src/domain/DBModel/award"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var configFile = flag.String("f", "etc/test.yaml", "the config file")

type Config struct {
	DB struct {
		MySqlDataSource string
	}
	Redis struct {
		Host     string
		Password string
		DB       int
	}
}

func TestDomainModel(t *testing.T) {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)

	sqlConn := sqlx.NewMysql(c.DB.MySqlDataSource)
	AwardModel := award.NewAwardModel(sqlConn)
	awards, err := AwardModel.FindAll(context.Background())
	if err != nil {
		t.Fatalf("failed to find all awards: %v", err)
	}

	awardsJSON, err := json.Marshal(awards)
	if err != nil {
		t.Fatalf("failed to marshal awards to JSON: %v", err)
	}

	t.Log(string(awardsJSON))
}

func TestInfrastructureRedisMap(t *testing.T) {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	rdb := redis.NewRedisService(c.Redis.Host, c.Redis.Password, c.Redis.DB)
	var strategyMap = map[string]any{
		"1": "1",
		"2": "2",
		"3": "3",
	}
	rdb.AddToMap(context.Background(), "testmap", strategyMap)
	getstrategyMap := rdb.GetMap(context.Background(), "testmap")
	for k, v := range getstrategyMap.Val() {
		t.Logf("Key: %s, Value: %s", k, v)
	}
}
