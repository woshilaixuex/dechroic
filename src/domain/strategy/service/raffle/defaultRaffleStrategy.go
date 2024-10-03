package raffle

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/factory"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽奖过滤引擎
 * @Date: 2024-08-07 22:43
 */
type DefaultRaffleStrategy struct {
	service.TemplateRaffleStrategy
	DefaultLogicFactory factory.DefaultLogicFactory
}

func NewDefaultRaffleStrategy(StrategyService repository.StrategyService, strategyDispatch armory.StrategyDispath) *DefaultRaffleStrategy {
	defaultRaffleStrategy := &DefaultRaffleStrategy{
		TemplateRaffleStrategy: *service.NewTemplateRaffleStrategy(StrategyService, strategyDispatch),
		DefaultLogicFactory:    *factory.NewDefaultLogicFactory(StrategyService),
	}
	defaultRaffleStrategy.RaffleLogicChainfunc = defaultRaffleStrategy.raffleLogicChain
	defaultRaffleStrategy.RaffleLogicTreefunc = defaultRaffleStrategy.raffleLogicTree
	return defaultRaffleStrategy
}
func (s *DefaultRaffleStrategy) raffleLogicChain(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error) {
	logicChain, err := s.DefaultChainFactory.OpenLogicChain(ctx, strategyId)
	if err != nil {
		return nil, err
	}
	return logicChain.Logic(ctx, userId, strategyId)
}
func (s *DefaultRaffleStrategy) raffleLogicTree(ctx context.Context, userId string, strategyId int64, awardId int32) (*data.StrategyAwardTreeVO, error) {
	strategyAwardRuleModelVO, err := s.StrategyService.QueryStrategyAwardRuleModelVO(ctx, strategyId, awardId)
	logx.Debug(strategyAwardRuleModelVO)
	if err != nil || strategyAwardRuleModelVO == nil {
		return &data.StrategyAwardTreeVO{
			AwardId: awardId,
		}, err
	}
	ruleTreeVO, err := s.StrategyService.QueryRuleTreeVOByTreeId(ctx, strategyAwardRuleModelVO.RuleModels)
	if err != nil {
		return nil, err
	}
	treeEngine := s.DefaultTreeFactory.OpenLogicTree(ruleTreeVO)
	return treeEngine.Process(userId, strategyId, awardId), nil
}
