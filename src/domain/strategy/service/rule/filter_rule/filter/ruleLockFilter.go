package filter

import (
	"context"
	"strconv"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	LogicModel "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/factory/model"
	filter_interface "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/filter/interface"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽奖次数累计解锁
 * @Date: 2024-08-13 23:07
 */
var _ filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface] = (*RuleLockLogicFilter)(nil)

type RuleLockLogicFilter struct {
	strategyService repository.StrategyService
	logicModel      LogicModel.LogicModel
	userRaffleCount int64
}

func NewRuleLockLogicFilter(strategyService repository.StrategyService) *RuleLockLogicFilter {
	return &RuleLockLogicFilter{
		strategyService: strategyService,
		logicModel:      LogicModel.RULE_LOCK,
	}
}
func (filter *RuleLockLogicFilter) Filter(ctx context.Context, ruleMatter StrategyEntity.RuleMatterEntity) (StrategyEntity.RaffleActionEntityInterface, error) {

	logx.Infof("规则过滤-次数 userId:%s strategyId:%d ruleModel:%s", ruleMatter.UserId, ruleMatter.StrategyId, ruleMatter.RuleModel)
	ruleValue, err := filter.strategyService.QueryStrategyRuleValue(ctx, ruleMatter.StrategyId, ruleMatter.AwardId, ruleMatter.RuleModel)
	if err != nil {
		return nil, err
	}
	raffleCount, err := strconv.ParseInt(ruleValue, 10, 64)
	if err != nil {
		return nil, err
	}
	if filter.userRaffleCount >= raffleCount {
		// 通过次规则
		return &StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleCenterEntity]{
			Code: vo.ALLOW.Code,
			Info: vo.ALLOW.Info,
		}, nil
	}
	return &StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleCenterEntity]{
		Code: vo.TAKE_OVER.Code,
		Info: vo.TAKE_OVER.Info,
	}, nil
}
func (filter *RuleLockLogicFilter) GetLogicModel() LogicModel.LogicModel {
	return filter.logicModel
}
