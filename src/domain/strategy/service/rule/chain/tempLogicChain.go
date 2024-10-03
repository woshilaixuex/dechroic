package chain

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 过滤链模版（impl为模版实现）
 * @Date: 2024-08-15 23:19
 */
type ChainLogic func(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error)

var _ ILogiChain = (*LogicChainNode)(nil)

// 过滤链节点
type LogicChainNode struct {
	next       ILogiChain
	chainLogic ChainLogic
}

// 初始化过滤链节点
func NewLogicChainNode() *LogicChainNode {
	return &LogicChainNode{
		next: nil,
	}
}
func (chain *LogicChainNode) Realize(chainLogic ChainLogic) {
	chain.chainLogic = chainLogic
}
func (chain *LogicChainNode) Logic(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error) {
	return chain.chainLogic(ctx, userId, strategyId)
}
func (chain *LogicChainNode) AppendNext(next ILogiChain) ILogiChain {
	chain.next = next
	return next
}
func (chain *LogicChainNode) Next() ILogiChain {
	return chain.next
}
func (chain *LogicChainNode) ModelType() string {
	return "start"
}
