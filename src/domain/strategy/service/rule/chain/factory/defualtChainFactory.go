package chain_factory

import (
	"context"
	"errors"

	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain"
	chain_ipml "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain/ipml"
	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 责任链工厂
 * @Date: 2024-08-16 23:45
 */

type DefaultChainFactory struct {
	logicChainGroup map[string]chain.ILogiChain
	strategyService repository.StrategyService
}

func NewDefaultLogicFactory(logicChainGroup map[string]chain.ILogiChain, strategyDispatch armory.StrategyDispath, strategyService repository.StrategyService) *DefaultChainFactory {
	logicChainGroup["default"] = chain_ipml.NewDefaultLogicChain(strategyDispatch)
	factory := &DefaultChainFactory{
		logicChainGroup: logicChainGroup,
		strategyService: strategyService,
	}
	factory.registerChainGroup(strategyDispatch, strategyService)
	return factory
}
func (factory *DefaultChainFactory) registerChainGroup(strategyDispatch armory.StrategyDispath, strategyService repository.StrategyService) {
	factory.logicChainGroup["rule_blacklist"] = chain_ipml.NewBackListLogicChain(strategyService)
	factory.logicChainGroup["rule_weight"] = chain_ipml.NewRuleWeightLogicChain(strategyService, strategyDispatch)
}

// 责任链装配/开启
func (factory *DefaultChainFactory) OpenLogicChain(ctx context.Context, strategyId int64) (chain.ILogiChain, error) {
	strategy, err := factory.strategyService.QueryStrategyEntityByStrategyId(ctx, strategyId)
	if err != nil {
		return nil, err
	}
	ruleModels, err := strategy.GetStrsRuleModels()
	if err != nil {
		return nil, err
	}
	if len(ruleModels) == 0 {
		cimpl := factory.logicChainGroup["default"]
		return cimpl, nil
	}
	logicChain := factory.logicChainGroup[ruleModels[0]]
	if logicChain == nil {
		return nil, cerr.LogError(errors.New(ruleModels[0] + "is not exit"))
	}
	current := logicChain
	for i := 1; i < len(ruleModels); i++ {
		nextChain := factory.logicChainGroup[ruleModels[i]]
		current = current.AppendNext(nextChain)
	}
	current.AppendNext(factory.logicChainGroup["default"])
	current = logicChain
	for current != nil {
		logx.Debug(current.ModelType())
		current = current.Next()
	}
	return logicChain, nil
}
