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

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/raffle"
	infra_award "github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/award"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategy"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyRule"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/treeRule"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/treeRuleNode"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/treeRuleNodeLine"
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

// strategy armory业务测试
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
	RuleMode := strategyRule.NewStrategyRuleModel(sqlConn)
	Model := strategy.NewStrategyModel(sqlConn)
	strategyRepo := &infra_repository.StrategyRepository{
		RedisService:       *rdb,
		StrategyAwardModel: AwardModel,
		StrategyModel:      Model,
		StrategyRuleModel:  RuleMode,
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

// strategy 过滤器测试
func TestPerformRaffleBlacklist(t *testing.T) { // 黑名单测试
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
	defaultRaffleStrategy := raffle.NewDefaultRaffleStrategy(*repository.NewStrategyService(strategyRepo), strategyArmory)
	entity := new(StrategyEntity.RaffleFactorEntity)
	entity.StrategyId = 100001
	entity.UserId = "user005"
	awradEntity, err := defaultRaffleStrategy.PerformRaffle(context.Background(), entity)
	if err != nil {
		t.Fatalf("PerformRaffle get err: %v", err)
	}
	t.Logf("奖品策略Id:%d", awradEntity.AwardId)
}

// 8.19策略更新增加按量更新功能
func TestPerformRaffleBlacklist819(t *testing.T) { // 黑名单测试
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	// 初始化依赖
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
	// defaultRaffleStrategy := raffle.NewDefaultRaffleStrategy(*repository.NewStrategyService(strategyRepo), strategyArmory)
	// entity := new(StrategyEntity.RaffleFactorEntity)
	// 黑名单100:user001,user002,user003
	// 4000:102,103,104,105 5000:102,103,104,105,106,107 6000:102,103,104,105,106,107,108,109
	// entity.StrategyId = 100001
	// entity.UserId = "user001"
	// // 查询流程
	// awradEntity, err := defaultRaffleStrategy.PerformRaffle(context.Background(), entity)
	// if err != nil {
	// 	t.Fatalf("PerformRaffle get err: %v", err)
	// }
	// t.Logf("奖品策略Id:%d", awradEntity.AwardId)
}

// 9.20过滤器责任链测试
func TestPerformRaffleBlacklist920(t *testing.T) {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	// 初始化依赖
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
	// strategyArmory.AssembleLotteryStrategy(context.Background(), 100001)
	defaultRaffleStrategy := raffle.NewDefaultRaffleStrategy(*repository.NewStrategyService(strategyRepo), strategyArmory)
	entity := new(StrategyEntity.RaffleFactorEntity)
	// 黑名单100:user001,user002,user003
	// 4000:102,103,104,105 5000:102,103,104,105,106,107 6000:102,103,104,105,106,107,108,109
	entity.StrategyId = 100001
	entity.UserId = "user005"
	// 查询流程
	awradEntity, err := defaultRaffleStrategy.PerformRaffle(context.Background(), entity)
	if err != nil {
		t.Fatalf("PerformRaffle get err: %v", err)
	}
	t.Logf("奖品策略Id:%d", awradEntity.AwardId)
}

// 初始化仓储
func initStrategyRepo() *infra_repository.StrategyRepository {
	flag.Parse()
	var c Config
	conf.MustLoad(*configFile, &c)
	// 初始化依赖
	rdb := redis.NewRedisService(c.Redis.Host, c.Redis.Password, c.Redis.DB)
	sqlConn := sqlx.NewMysql(c.DB.MySqlDataSource)
	AwardModel := strategyAward.NewStrategyAwardModel(sqlConn)
	RuleMode := strategyRule.NewStrategyRuleModel(sqlConn)
	Model := strategy.NewStrategyModel(sqlConn)
	TreeModel := treeRule.NewRuleTreeModel(sqlConn)
	TreeNodeModel := treeRuleNode.NewRuleTreeNodeModel(sqlConn)
	TreeNodeLineModel := treeRuleNodeLine.NewRuleTreeNodeLineModel(sqlConn)
	return &infra_repository.StrategyRepository{
		RedisService:          *rdb,
		StrategyAwardModel:    AwardModel,
		StrategyModel:         Model,
		StrategyRuleModel:     RuleMode,
		TreeRuleModel:         TreeModel,
		TreeRuleNodeModel:     TreeNodeModel,
		TreeRuleNodeLineModel: TreeNodeLineModel,
	}
}

// 9.23策略树测试
// func TestLogicTree(t *testing.T) {
// 	// 带切片的VO
// 	rule_lock := &vo.RuleTreeNodeVO{
// 		TreeId:                  "100000001",
// 		RuleKey:                 "rule_lock",
// 		RukeDesc:                "限定用户已完成N次抽奖后解锁",
// 		RuleValue:               "1",
// 		RuleTreeNodeLineVOSlice: make([]*vo.RuleTreeNodeLineVO, 0),
// 	}
// 	ruleTreeNodeLineVO1 := &vo.RuleTreeNodeLineVO{
// 		TreeId:               "100000001",
// 		RuleNodeFrom:         "rule_lock",
// 		RuleNodeTo:           "rule_luck_award",
// 		RuleLimitTypeVO:      &vo.EQUAL,
// 		RuleLogicCheckTypeVO: &vo.TAKE_OVER,
// 	}
// 	ruleTreeNodeLineVO2 := &vo.RuleTreeNodeLineVO{
// 		TreeId:               "100000001",
// 		RuleNodeFrom:         "rule_lock",
// 		RuleNodeTo:           "rule_stock",
// 		RuleLimitTypeVO:      &vo.EQUAL,
// 		RuleLogicCheckTypeVO: &vo.ALLOW,
// 	}
// 	rule_lock.RuleTreeNodeLineVOSlice = append(rule_lock.RuleTreeNodeLineVOSlice, ruleTreeNodeLineVO1)
// 	rule_lock.RuleTreeNodeLineVOSlice = append(rule_lock.RuleTreeNodeLineVOSlice, ruleTreeNodeLineVO2)
// 	// 无值VO
// 	rule_luck_award := &vo.RuleTreeNodeVO{
// 		TreeId:                  "100000001",
// 		RuleKey:                 "rule_luck_award",
// 		RukeDesc:                "限定用户已完成N次抽奖后解锁",
// 		RuleValue:               "1",
// 		RuleTreeNodeLineVOSlice: nil,
// 	}
// 	rule_stock := &vo.RuleTreeNodeVO{
// 		TreeId:                  "100000001",
// 		RuleKey:                 "rule_stock",
// 		RukeDesc:                "库存处理规则",
// 		RuleValue:               "",
// 		RuleTreeNodeLineVOSlice: make([]*vo.RuleTreeNodeLineVO, 0),
// 	}
// 	rule_stock.RuleTreeNodeLineVOSlice = append(rule_stock.RuleTreeNodeLineVOSlice, ruleTreeNodeLineVO1)
// 	ruleTreeVO := vo.NewRuleTreeVO()
// 	ruleTreeVO.TreeId = "100000001"
// 	ruleTreeVO.TreeName = "决策树规则；增加dall-e-3画图模型"
// 	ruleTreeVO.TreeDesc = "决策树规则；增加dall-e-3画图模型"
// 	ruleTreeVO.TreeRootRuleNode = "rule_lock"
// 	ruleTreeVO.TreeNodeMap["rule_lock"] = rule_lock
// 	ruleTreeVO.TreeNodeMap["rule_stock"] = rule_stock
// 	ruleTreeVO.TreeNodeMap["rule_luck_award"] = rule_luck_award
// 	treeFactory := tree_factory.NewDefultTreeFactory(tree_impl.NewRuleLockLogicTreeNode(),
// 		tree_impl.NewRuleLuckAwardLogicTreeNode(),
// 		tree_impl.NewRuleStockLogicTreeNode())
// 	treeComposite := treeFactory.OpenLogicTree(ruleTreeVO)
// 	data := treeComposite.Process("delyr1c", 100001, 100)
// 	t.Log(data.AwardRuleValue)
// }

// 10.3策略数链接
func TestLogicTreeLink(t *testing.T) {
	strategyRepo := initStrategyRepo()
	strategyArmory := armory.NewStrategyArmory(*repository.NewStrategyService(strategyRepo))
	// strategyArmory.AssembleLotteryStrategy(context.Background(), 100001)
	defaultRaffleStrategy := raffle.NewDefaultRaffleStrategy(*repository.NewStrategyService(strategyRepo), strategyArmory)
	raffleAwardEntity := &StrategyEntity.RaffleFactorEntity{
		UserId:     "delyr1c",
		StrategyId: 100001,
	}
	award, err := defaultRaffleStrategy.PerformRaffle(context.Background(), raffleAwardEntity)
	if err != nil {
		t.Fail()
	}
	// data, _ := repository.NewStrategyService(strategyRepo).QueryRuleTreeVOByTreeId(context.Background(), "tree_lock")
	// data.TraverseRuleTree()
	t.Log(award)
}
