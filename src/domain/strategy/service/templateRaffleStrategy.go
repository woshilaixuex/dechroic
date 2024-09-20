package service

import (
	"context"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain"
	chain_factory "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain/factory"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽象抽奖过滤引擎
 * @Date: 2024-08-07 14:54
 */
type CheckRaffleBeforeLogicFunc func(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error)
type CheckRaffleCenterLogicFunc func(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error)
type TemplateRaffleStrategyInterface interface {
	PerformRaffle(ctx context.Context, entity StrategyEntity.RaffleFactorEntity) error
	DoCheckRaffleBeforeLogic(ctx context.Context, entity StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error)
	DoCheckRaffleCenterLogic(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error)
}
type TemplateRaffleStrategy struct {
	StrategyService     repository.StrategyService
	StrategyDispatch    armory.StrategyDispath
	DefaultChainFactory *chain_factory.DefaultChainFactory
	CheckRaffleBeforeLogicFunc
	CheckRaffleCenterLogicFunc
}

func NewTemplateRaffleStrategy(StrategyService repository.StrategyService, StrategyDispatch armory.StrategyDispath) *TemplateRaffleStrategy {
	return &TemplateRaffleStrategy{
		StrategyService:            StrategyService,
		StrategyDispatch:           StrategyDispatch,
		DefaultChainFactory:        chain_factory.NewDefaultLogicFactory(make(map[string]chain.ILogiChain), StrategyDispatch, StrategyService),
		CheckRaffleBeforeLogicFunc: nil,
		CheckRaffleCenterLogicFunc: nil,
	}
}
func (t *TemplateRaffleStrategy) PerformRaffle(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity) (*StrategyEntity.RaffleAwardEntity, error) {
	// 过滤
	userId := entity.UserId
	strategyId := entity.StrategyId
	// // 查找策略实体
	// strategy, err := t.StrategyService.QueryStrategyEntityByStrategyId(ctx, strategyId)
	// if err != nil {
	// 	return nil, err
	// }
	// ruleModels, err := strategy.GetStrsRuleModels()
	// if err != nil {
	// 	return nil, err
	// }
	logicChain, err := t.DefaultChainFactory.OpenLogicChain(ctx, strategyId)
	if err != nil {
		return nil, err
	}
	logx.Debug("factory test")
	awardId, err := logicChain.Logic(ctx, userId, strategyId)
	if err != nil {
		return nil, err
	}

	// // 过滤前置
	// ruleActionEntityBefore, err := t.DoCheckRaffleBeforeLogic(ctx, &StrategyEntity.RaffleFactorEntity{
	// 	UserId:     userId,
	// 	StrategyId: strategyId,
	// }, ruleModels...)
	// if err != nil {
	// 	return nil, err
	// }
	// ruleActionEntity, ok := ruleActionEntityBefore.(*StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity])
	// if !ok {
	// 	return nil, cerr.LogError(errors.New("ruleActionEntity Conversion failure or ruleActionEntity is null"))
	// }
	// if vo.TAKE_OVER.Code == ruleActionEntity.Code {
	// 	if model.RULE_BLACKLIST.Code == ruleActionEntity.RuleModel {
	// 		return &StrategyEntity.RaffleAwardEntity{
	// 			AwardId: int64(ruleActionEntity.Data.AwardId),
	// 		}, nil
	// 	} else if model.RULE_WEIGHT.Code == ruleActionEntity.RuleModel {
	// 		raffleBeforeEntity := ruleActionEntity.Data
	// 		ruleWeightValueKey := raffleBeforeEntity.RuleWeightValueKey
	// 		ruleWeightValueKeyNum, err := strconv.ParseInt(ruleWeightValueKey, 10, 64)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		awardId, err := t.StrategyDispatch.GetRandomAwardId(ctx, strategyId, ruleWeightValueKeyNum)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		return &StrategyEntity.RaffleAwardEntity{
	// 			AwardId: awardId,
	// 		}, nil
	// 	}
	// }
	// // 过滤
	// awardId, err := t.StrategyDispatch.GetRandomAwardIdBase(ctx, strategyId)
	// if err != nil {
	// 	return nil, err
	// }

	strategyAwardRuleModelVO, err := t.StrategyService.QueryStrategyAwardRuleModelVO(ctx, strategyId, int32(awardId))
	if err != nil {
		return nil, err
	}
	ruleActionEntityCenter, err := t.DoCheckRaffleCenterLogic(ctx, entity, strategyAwardRuleModelVO.RaffleCenterModelList()...)
	if err != nil {
		return nil, err
	}
	if vo.TAKE_OVER.Code == ruleActionEntityCenter.GetCode() {
		logx.Info("【临时日志】中奖中规则拦截，通过抽奖后规则 rule_luck_awrad 走兜底奖励")
		return &StrategyEntity.RaffleAwardEntity{
			AwardDesc: "中奖中规则拦截，通过抽奖后规则 rule_luck_awrad 走兜底奖励",
		}, nil
	}
	return &StrategyEntity.RaffleAwardEntity{
		AwardId: int64(awardId),
	}, nil
}

// 抽象方法 过滤前执行
func (t *TemplateRaffleStrategy) DoCheckRaffleBeforeLogic(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error) {
	if t.CheckRaffleBeforeLogicFunc == nil {
		panic("CheckRaffleBeforeLogicFunc is not setting")
	}
	return t.CheckRaffleBeforeLogicFunc(ctx, entity, logics...)
}

// 抽象方法 过滤时执行
func (t *TemplateRaffleStrategy) DoCheckRaffleCenterLogic(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error) {
	if t.CheckRaffleCenterLogicFunc == nil {
		panic("CheckRaffleLogicFunc is not setting")
	}
	return t.CheckRaffleCenterLogicFunc(ctx, entity, logics...)
}
