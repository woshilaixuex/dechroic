package test

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 装配策略测试
 * @Date: 2024-06-13 22:23
 */
import (
	"context"
	"encoding/json"
	"flag"
	"testing"

	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	infra_award "github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/award"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// orm测试

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
	AwardModel := infra_award.NewAwardModel(sqlConn)
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
	rdb.SetToMap(context.Background(), "testmap", strategyMap)
	getstrategyMap := rdb.GetMap(context.Background(), "testmap")
	for k, v := range getstrategyMap.Val() {
		t.Logf("Key: %s, Value: %s", k, v)
	}
}

type TestEntity struct {
	Test1 string `json:"test1"`
	Test2 string `json:"test2"`
}

func TestInfrastructureRedisArray(t *testing.T) {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	rdb := redis.NewRedisService(c.Redis.Host, c.Redis.Password, c.Redis.DB)

	setarr := []TestEntity{
		{Test1: "test1", Test2: "test2"},
		{Test1: "test3", Test2: "test4"},
		{Test1: "test5", Test2: "test6"},
	}

	// 存储数据到 Redis
	if err := rdb.SetArray(context.Background(), "testarray", setarr); err != nil {
		t.Fatalf("Error setting array to Redis: %v", err)
	}

	// 从 Redis 获取数据
	var getarr []TestEntity
	if err := rdb.GetArray(context.Background(), "testarray", &getarr); err != nil {
		t.Fatalf("Error getting array from Redis: %v", err)
	}

	// 打印获取的数据
	for _, v := range getarr {
		t.Logf("Value: %+v", v)
	}
}

// strategy业务测试

func TestGetStrategy(t *testing.T) {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	rdb := redis.NewRedisService(c.Redis.Host, c.Redis.Password, c.Redis.DB)

	sqlConn := sqlx.NewMysql(c.DB.MySqlDataSource)
	AwardModel := strategyAward.NewStrategyAwardModel(sqlConn)
	strategyRepo := &infra_repository.StrategyRepository{
		RedisService:       *rdb,
		StrategyAwardModel: AwardModel,
	}
	strategyService := repository.NewStrategyService(strategyRepo)
	entitis, err := strategyService.QueryStrategyAwardList(context.Background(), 100001)
	if err != nil {
		t.Fatalf("failed to find all awards: %v", err)
	}
	for _, entity := range entitis {
		t.Logf("Value: %+v %v", entity, entity.AwardRate.String())
	}
}
