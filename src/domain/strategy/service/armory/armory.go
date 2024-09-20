package armory

import (
	"context"
	"fmt"
	"math/big"
	"math/rand"
	"sort"
	"strconv"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 策略装配工厂的实现
 * @Date: 2024-07-24 12:40
 */

// 40 位精度，确保足够高的精度
var (
	prec = uint(40)
)

// ------------------------------------------------------------
// 策略工厂接口
type StrategyAssemble interface {
	AssembleLotteryStrategy(ctx context.Context, strategyId int64) bool
}

type Armory interface {
	StrategyAssemble
	StrategyDispath
}

var _ Armory = (*StrategyArmory)(nil)

// ------------------------------------------------------------
// 策略工厂实体
type StrategyArmory struct {
	strategyService repository.StrategyService
}

// 创建策略工厂
func NewStrategyArmory(strategyService repository.StrategyService) *StrategyArmory {
	return &StrategyArmory{
		strategyService: strategyService,
	}
}

// 实现装配抽奖策略
func (armory *StrategyArmory) AssembleLotteryStrategy(ctx context.Context, strategyId int64) bool {

	// 获取策略列表
	entities, err := armory.strategyService.QueryStrategyAwardList(ctx, strategyId)
	if err != nil {
		return false
	}
	logx.Infof("Fetched %d strategy awards", len(entities))
	armory.assembleLotteryStrategy(ctx, strconv.FormatInt(strategyId, 10), entities)
	// 权重规则配置
	strategyEntity, err := armory.strategyService.QueryStrategyEntityByStrategyId(ctx, strategyId)
	if err != nil {
		return false
	}
	ruleWeight, err := strategyEntity.GetRuleWeight()
	if err != nil {
		return false
	}
	if ruleWeight == "" {
		return true
	}
	// 查询规则
	strategyRule, err := armory.strategyService.QueryStrategyRule(ctx, strategyId, ruleWeight)
	if err != nil {
		return false
	}
	ruleMap, err := strategyRule.GetRule()
	if err != nil {
		return false
	}
	logx.Debug(ruleMap)
	for key, set := range ruleMap {
		var entitiesClone []entity.StrategyAwardEntity
		for _, etty := range entities {
			if set[etty.AwardId] {
				entitiesClone = append(entitiesClone, etty)
			}
		}
		strategyKey := fmt.Sprintf("%d_%d", strategyId, key)
		armory.assembleLotteryStrategy(ctx, strategyKey, entitiesClone)
	}
	return true
}

// 装配工厂核心实现
func (armory *StrategyArmory) assembleLotteryStrategy(ctx context.Context, key string, entities []entity.StrategyAwardEntity) {
	// 获取最小概率
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].AwardRate.Cmp(&entities[j].AwardRate) < 0
	})
	rateMin := new(big.Float).SetPrec(prec).Set(entities[0].AwardRate.Float)
	rateMin, err := common.GetSmallestUnitIncrementByStr(rateMin)
	if err != nil {
		return
	}
	rateSum := common.NewBigFloat().SetPrec(prec)
	for _, entity := range entities {
		rateSum.Add(rateSum, entity.AwardRate.Float)
	}
	logx.Infof("Minimum rate: %s Total rate sum: %s", rateMin.String(), rateSum.String())

	// 求取范围
	fRateRange := new(big.Float).Quo(rateSum, rateMin)
	iRateRange := new(big.Int)
	fRateRange.Int(iRateRange)
	logx.Infof("Rate range: %s", iRateRange.String())

	// 生成概率表
	strategyAwardSearchRateTables := generateRateTable(entities, rateMin)

	// 乱序搜索表
	shuffleStrategyAwardSearchRateTable(strategyAwardSearchRateTables)
	shuffleStrategyAwardSearchRateTable := make(map[interface{}]interface{}, len(strategyAwardSearchRateTables))
	for i, awardId := range strategyAwardSearchRateTables {
		shuffleStrategyAwardSearchRateTable[i] = awardId
	}
	logx.Infof("Rate map: %v", len(shuffleStrategyAwardSearchRateTable))
	// 存放到 Redis
	if err := armory.strategyService.StoreStrategyAwardSearchRateTable(ctx, key, iRateRange.Int64(),
		shuffleStrategyAwardSearchRateTable); err != nil {
		logx.Error(err.Error())
		return
	}
	logx.Info("Stored rate table in Redis")
}

// 生成概率表
func generateRateTable(entities []entity.StrategyAwardEntity, rateMin *big.Float) []int64 {
	var strategyAwardSearchRateTables []int64
	for _, entity := range entities {
		awardRate := entity.AwardRate.Float
		fCount := new(big.Float).SetPrec(prec).Quo(awardRate, rateMin)
		count := new(big.Int)
		fCount.Int(count)
		for i := int64(0); i < count.Int64(); i++ {
			strategyAwardSearchRateTables = append(strategyAwardSearchRateTables, entity.AwardId)
		}
	}
	logx.Infof("Generated rate table with %d entries", len(strategyAwardSearchRateTables))
	return strategyAwardSearchRateTables
}

// 乱序搜索表
func shuffleStrategyAwardSearchRateTable(strategyAwardSearchRateTables []int64) {
	rand.Shuffle(len(strategyAwardSearchRateTables), func(i, j int) {
		strategyAwardSearchRateTables[i], strategyAwardSearchRateTables[j] = strategyAwardSearchRateTables[j], strategyAwardSearchRateTables[i]
	})
}
