package chain_ipml

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 默认过滤节点（兜底）
 * @Date: 2024-08-16 00:28
 */

type DefaultLogicChain struct {
	chain.LogicChainNode
	strategyDispatch armory.StrategyDispath
}

func NewDefaultLogicChain(strategyDispatch armory.StrategyDispath) *DefaultLogicChain {
	defaultLogicChain := &DefaultLogicChain{
		LogicChainNode:   *chain.NewLogicChainNode(),
		strategyDispatch: strategyDispatch,
	}
	defaultLogicChain.Realize(defaultLogicChain.Logic)
	return defaultLogicChain
}
func (c *DefaultLogicChain) Logic(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error) {
	awardId, err := c.strategyDispatch.GetRandomAwardIdBase(ctx, strategyId)
	if err != nil {
		return nil, err
	}
	return &data.StrategyAwardChanVO{
		AwardId:          int32(awardId),
		AwardRuleValue:   c.ModelType(),
		ChanVOLogicModel: data.RULE_DEFAULT,
	}, nil
}
func (c *DefaultLogicChain) ModelType() string {
	return "default"
}
