package filter

import (
	"context"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-08-06 00:51
 */
var _ ILogicFilter[StrategyEntity.RaffleActionEntityInterface] = (*RuleWeightLogicFilter)(nil)

type RuleWeightLogicFilter struct {
	ILogicFilter[StrategyEntity.RaffleActionEntityInterface]
	strategyService repository.StrategyService
}

func (filter *RuleWeightLogicFilter) Filter(ctx context.Context, ruleMatter StrategyEntity.RuleMatterEntity) (StrategyEntity.RaffleActionEntityInterface, error) {
	logx.Infof("规则过滤-权重范围 userId:%s strategyId:%d ruleModel:%s", ruleMatter.UserId, ruleMatter.StrategyId, ruleMatter.RuleModel)
}
