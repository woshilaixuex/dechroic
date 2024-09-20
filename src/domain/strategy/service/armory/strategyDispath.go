package armory

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 调度策略
 * @Date: 2024-06-28 19:57
 */

type StrategyDispath interface {
	GetRandomAwardIdBase(ctx context.Context, strategyId int64) (int64, error)
	GetRandomAwardId(ctx context.Context, strategyId, ruleWeightValue int64) (int64, error)
}

// 获取随机奖品 ID
func (armory *StrategyArmory) GetRandomAwardIdBase(ctx context.Context, strategyId int64) (int64, error) {
	rand.NewSource(time.Now().UnixNano())
	rateRange, err := armory.strategyService.GetRateRangeBeta(ctx, strategyId)
	if err != nil {
		logx.Error(err.Error())
		return 0, err
	}
	randomVal := rand.Int63n(rateRange)
	return armory.strategyService.GetAssembleRandomVal(ctx, strconv.FormatInt(strategyId, 10), randomVal)
}
func (armory *StrategyArmory) GetRandomAwardId(ctx context.Context, strategyId, ruleWeightValue int64) (int64, error) {
	key := fmt.Sprintf("%d_%d", strategyId, ruleWeightValue)
	rand.NewSource(time.Now().UnixNano())
	rateRange, err := armory.strategyService.GetRateRange(ctx, key)
	if err != nil {
		logx.Error(err.Error())
		return 0, err
	}
	randomVal := rand.Int63n(rateRange)
	return armory.strategyService.GetAssembleRandomVal(ctx, key, randomVal)
}
