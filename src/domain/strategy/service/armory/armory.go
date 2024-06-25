package armory

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 装配工厂
 * @Date: 2024-06-16 11:18
 */

// 抽象抽奖策略工厂
type Armory interface {
	AssembleLotteryStrategy(strategyId int64)
}

type StrategyArmory struct {
	strategyService repository.StrategyService
}

func (armory *StrategyArmory) AssembleLotteryStrategy(ctx context.Context, strategyId int64) bool {
	if entitis, err := armory.strategyService.QueryStrategyAwardList(ctx, strategyId); err != nil {
		logx.Error(entitis)
		return false
	}
	return false
}
