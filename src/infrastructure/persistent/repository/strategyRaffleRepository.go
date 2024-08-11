package repository

import (
	"context"

	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyRule"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: raffle的repository实现
 * @Date: 2024-08-05 22:14
 */
func (s *StrategyRepository) QueryStrategyRuleValue(ctx context.Context, strategyId int64, awardId int32, roleModel string) (string, error) {
	strategyRuleReq := &strategyRule.FindStrategyRuleReq{
		StrategyId: &strategyId,
		RuleModel:  &roleModel,
	}
	if awardId == 0 {
		strategyRuleReq.AwardId = nil
	}
	return s.StrategyRuleModel.FindRuleValueByReq(ctx, strategyRuleReq)
}
