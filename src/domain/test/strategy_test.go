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
	"fmt"
	"math/big"
	"sort"
	"testing"

	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	infra_award "github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/award"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategy"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyRule"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
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

	var strategyMap = map[interface{}]any{
		"1": "1",
		"2": "2",
		"3": "3",
	}
	t.Logf("Key: %v", strategyMap)

	// 尝试存储数据到 Redis
	err := rdb.SetToMap(context.Background(), "testmap", strategyMap)
	if err != nil {
		t.Errorf("Failed to set data to Redis: %v", err)
	}
	// 获取 Redis 中的数据
	getstrategyMap := rdb.GetMap(context.Background(), "testmap")
	if err := getstrategyMap.Err(); err != nil {
		t.Errorf("Failed to get data from Redis: %v", err)
	} else {
		for k, v := range getstrategyMap.Val() {
			t.Logf("Key: %s, Value: %s", k, v)
		}
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
func TestGetFromMap(t *testing.T) {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	rdb := redis.NewRedisService(c.Redis.Host, c.Redis.Password, c.Redis.DB)
	logx.Info(rdb.GetFromMap(context.Background(), "dechoric_strategy_rate_table_key_100002", "1").Int64())
}
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
	sum := common.NewBigFloat()
	sort.Slice(entitis, func(i, j int) bool {
		return entitis[i].AwardRate.Cmp(&entitis[j].AwardRate) < 0
	})
	min := entitis[0].AwardRate
	for _, entity := range entitis {
		t.Log(sum.Float.String())
		sum.Float.Add(sum.Float, entity.AwardRate.Float)
	}
	t.Log(sum.Float.String())
	t.Log(sum.Float.Quo(sum.Float, min.Float).String())
}
func TestAssembleLotteryStrategy(t *testing.T) {
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
	strategyArmory := armory.NewStrategyArmory(*repository.NewStrategyService(strategyRepo))
	strategyArmory.AssembleLotteryStrategy(context.Background(), 100001)
	// RandomAwardId, err := strategyArmory.GetRandomAwardId(context.Background(), 100001)
	// if err != nil {
	// 	t.Fatalf("failed to find all awards: %v", err)
	// }
	// fmt.Print("奖品ID:")
	// fmt.Println(RandomAwardId)
}
func TestGetRandomAwardId(t *testing.T) {
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
	strategyArmory := armory.NewStrategyArmory(*repository.NewStrategyService(strategyRepo))
	for i := 0; i < 10; i++ {
		RandomAwardId, err := strategyArmory.GetRandomAwardId(context.Background(), 100001, 0)
		if err != nil {
			t.Fatalf("failed to find all awards: %v", err)
		}
		fmt.Print("奖品ID:")
		fmt.Println(RandomAwardId)
	}
}
func TestBigFloat(t *testing.T) {
	prec := uint(200)
	rateSum := new(big.Float).SetPrec(prec).SetFloat64(100.2)
	rateMin := new(big.Float).SetPrec(prec).SetFloat64(0.001)

	// 求取范围
	fRateRange := new(big.Float).Quo(rateSum, rateMin)
	iRateRange := new(big.Int)

	// 获取整数部分和舍入情况
	fRateRange.Int(iRateRange)

	fmt.Printf("Rate range: %s\n", fRateRange.String())
	fmt.Printf("Integer Rate range: %s\n", iRateRange.String())
}

// 7.29策略更新增加按量更新功能
func TestAssembleLotteryStrategy729(t *testing.T) {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	rdb := redis.NewRedisService(c.Redis.Host, c.Redis.Password, c.Redis.DB)

	sqlConn := sqlx.NewMysql(c.DB.MySqlDataSource)
	AwardModel := strategyAward.NewStrategyAwardModel(sqlConn)
	RuleMode := strategyRule.NewStrategyRuleModel(sqlConn)
	Model := strategy.NewStrategyModel(sqlConn)
	strategyRepo := &infra_repository.StrategyRepository{
		RedisService:       *rdb,
		StrategyAwardModel: AwardModel,
		StrategyModel:      Model,
		StrategyRuleModel:  RuleMode,
	}
	strategyArmory := armory.NewStrategyArmory(*repository.NewStrategyService(strategyRepo))
	strategyArmory.AssembleLotteryStrategy(context.Background(), 100001)
}
func TestGetRandomAwardId729(t *testing.T) {
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
	strategyArmory := armory.NewStrategyArmory(*repository.NewStrategyService(strategyRepo))
	for i := 0; i < 10; i++ {
		RandomAwardId, err := strategyArmory.GetRandomAwardId(context.Background(), 100001, 6000)
		if err != nil {
			t.Fatalf("failed to find all awards: %v", err)
		}
		fmt.Print("奖品ID:")
		fmt.Println(RandomAwardId)
	}
}
