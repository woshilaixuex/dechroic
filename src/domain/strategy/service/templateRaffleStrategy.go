package service

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain"
	chain_factory "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain/factory"
	tree_factory "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/tree/factory"
	tree_impl "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/tree/impl"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽象抽奖过滤引擎
 * @Date: 2024-08-07 14:54
 */
type RaffleLogicChainfunc func(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error)
type RaffleLogicTreefunc func(ctx context.Context, userId string, strategyId int64, awardId int32) (*data.StrategyAwardTreeVO, error)
type TemplateRaffleStrategyInterface interface {
	PerformRaffle(ctx context.Context, entity StrategyEntity.RaffleFactorEntity) error
	RaffleLogicChainfunc(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error)
	RaffleLogicTreefunc(ctx context.Context, userId string, strategyId int64, awardId int32) (*data.StrategyAwardTreeVO, error)
}
type TemplateRaffleStrategy struct {
	// 策略仓储服务 -> 仓储层提供业务编排所需处理完的数据
	StrategyService repository.StrategyService
	// 策略调度服务 -> 只负责抽奖处理，通过新增接口的方式，隔离职责，不需要使用方关心或者调用抽奖的初始化
	StrategyDispatch armory.StrategyDispath
	// 抽奖的责任链 -> 从抽奖的规则中，解耦出前置规则为责任链处理
	DefaultChainFactory *chain_factory.DefaultChainFactory
	// 抽奖的决策树 -> 负责抽奖中到抽奖后的规则过滤，如抽奖到A奖品ID，之后要做次数的判断和库存的扣减
	DefaultTreeFactory *tree_factory.DefultTreeFactory
	RaffleLogicChainfunc
	RaffleLogicTreefunc
}

func NewTemplateRaffleStrategy(StrategyService repository.StrategyService, StrategyDispatch armory.StrategyDispath) *TemplateRaffleStrategy {
	return &TemplateRaffleStrategy{
		StrategyService:     StrategyService,
		StrategyDispatch:    StrategyDispatch,
		DefaultChainFactory: chain_factory.NewDefaultLogicFactory(make(map[string]chain.ILogiChain), StrategyDispatch, StrategyService),
		// !!未对map进行注入
		DefaultTreeFactory: tree_factory.NewDefultTreeFactory(tree_impl.NewRuleLockLogicTreeNode(), tree_impl.NewRuleLuckAwardLogicTreeNode(),
			tree_impl.NewRuleStockLogicTreeNode()),
		RaffleLogicChainfunc: nil,
		RaffleLogicTreefunc:  nil,
	}
}
func (t *TemplateRaffleStrategy) PerformRaffle(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity) (*StrategyEntity.RaffleAwardEntity, error) {
	// 过滤
	userId := entity.UserId
	strategyId := entity.StrategyId
	// 过滤链
	chainStrategyAwardVO, err := t.raffleLogicChain(ctx, userId, strategyId)
	if err != nil {
		return nil, err
	}
	if !(chainStrategyAwardVO.Code == data.RULE_DEFAULT.Code) {
		return &StrategyEntity.RaffleAwardEntity{
			AwardId: int64(chainStrategyAwardVO.AwardId),
		}, nil
	}
	logx.Infof("抽奖策略计算-责任链 %s %d %d %s", userId, strategyId, chainStrategyAwardVO.AwardId, chainStrategyAwardVO.AwardRuleValue)
	// 策略树
	raffleStrategyAwardVO, err := t.raffleLogicTree(ctx, userId, strategyId, chainStrategyAwardVO.AwardId)
	logx.Debug(raffleStrategyAwardVO)
	if err != nil {
		return nil, err
	}
	logx.Infof("抽奖策略计算-规则树 %s %d %d %s", userId, strategyId, raffleStrategyAwardVO.AwardId, raffleStrategyAwardVO.AwardRuleValue)
	return &StrategyEntity.RaffleAwardEntity{
		AwardId:     int64(raffleStrategyAwardVO.AwardId),
		AwardConfig: raffleStrategyAwardVO.AwardRuleValue,
	}, nil
}
func (t *TemplateRaffleStrategy) raffleLogicChain(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error) {
	if t.RaffleLogicChainfunc == nil {
		panic("RaffleLogicChainfunc is undefined")
	}
	return t.RaffleLogicChainfunc(ctx, userId, strategyId)
}
func (t *TemplateRaffleStrategy) raffleLogicTree(ctx context.Context, userId string, strategyId int64, awardId int32) (*data.StrategyAwardTreeVO, error) {
	if t.RaffleLogicTreefunc == nil {
		panic("RaffleLogicTreefunc is undefined")
	}
	return t.RaffleLogicTreefunc(ctx, userId, strategyId, awardId)
}
