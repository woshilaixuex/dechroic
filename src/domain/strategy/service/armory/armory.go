package armory

import (
	"context"
	"math/big"
	"math/rand"
	"sort"
	"time"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/logx"
)

// 200 位精度，确保足够高的精度
var (
	prec = uint(200)
)

// 策略工厂接口
type Armory interface {
	AssembleLotteryStrategy(strategyId int64) bool
	GetRandomAwardId(ctx context.Context, strategyId int64) (int64, error)
}

// 策略工厂实现
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
		logx.Error(err.Error())
		return false
	}
	logx.Infof("Fetched %d strategy awards", len(entities))

	// 获取最小概率
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].AwardRate.Cmp(&entities[j].AwardRate) < 0
	})
	rateMin := new(big.Float).SetPrec(prec).Set(entities[0].AwardRate.Float)
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
	if err := armory.strategyService.StoreStrategyAwardSearchRateTable(ctx, strategyId, iRateRange.Int64(),
		shuffleStrategyAwardSearchRateTable); err != nil {
		logx.Error(err.Error())
		return false
	}
	logx.Info("Stored rate table in Redis")
	return true
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
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(strategyAwardSearchRateTables), func(i, j int) {
		strategyAwardSearchRateTables[i], strategyAwardSearchRateTables[j] = strategyAwardSearchRateTables[j], strategyAwardSearchRateTables[i]
	})
}

// 获取随机奖品 ID
func (armory *StrategyArmory) GetRandomAwardId(ctx context.Context, strategyId int64) (int64, error) {
	rateRange, err := armory.strategyService.GetRateRange(ctx, strategyId)
	if err != nil {
		logx.Error(err.Error())
		return 0, err
	}
	randomVal := rand.Int63n(rateRange)
	return armory.strategyService.GetAssembleRandomVal(ctx, strategyId, randomVal)
}
