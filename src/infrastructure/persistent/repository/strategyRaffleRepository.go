package repository

import (
	"context"
	"errors"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyRule"
	"github.com/delyr1c/dechoric/src/types/cerr"
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
func (s *StrategyRepository) QueryStrategyAwardRuleModelVO(ctx context.Context, strategyId int64, awardId int32) (*vo.StrategyAwardRuleModelVO, error) {
	newAwardId := int64(awardId)
	StrategyAwardReq := &strategyAward.FindStrategyAwardReq{
		StrategyId: &strategyId,
		AwardId:    &newAwardId,
	}
	StrategyAwards, err := s.StrategyAwardModel.FindByReq(ctx, StrategyAwardReq)
	if err != nil {
		return nil, err
	}
	StrategyAward := StrategyAwards[0]
	if !StrategyAward.RuleModels.Valid {
		return nil, cerr.LogError(errors.New("Dao RuleModels is null "))
	}
	return &vo.StrategyAwardRuleModelVO{
		RuleModels: StrategyAward.RuleModels.String,
	}, nil
}
