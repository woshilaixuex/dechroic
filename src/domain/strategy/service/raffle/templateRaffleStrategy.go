package raffle

import (
	"context"
	"errors"
	"strconv"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/factory/model"
	"github.com/delyr1c/dechoric/src/types/cerr"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽象抽奖过滤引擎
 * @Date: 2024-08-07 14:54
 */
type CheckRaffleBeforeLogicFunc func(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error)

type TemplateRaffleStrategyInterface interface {
	PerformRaffle(ctx context.Context, entity StrategyEntity.RaffleFactorEntity) error
	DoCheckRaffleBeforeLogic(ctx context.Context, entity StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error)
}
type TemplateRaffleStrategy struct {
	StrategyService  repository.StrategyService
	StrategyDispatch armory.StrategyDispath
	CheckRaffleBeforeLogicFunc
}

func NewTemplateRaffleStrategy(StrategyService repository.StrategyService, StrategyDispatch armory.StrategyDispath) *TemplateRaffleStrategy {
	return &TemplateRaffleStrategy{
		StrategyService:            StrategyService,
		StrategyDispatch:           StrategyDispatch,
		CheckRaffleBeforeLogicFunc: nil,
	}
}
func (t *TemplateRaffleStrategy) PerformRaffle(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity) (*StrategyEntity.RaffleAwardEntity, error) {
	userId := entity.UserId
	strategyId := entity.StrategyId
	strategy, err := t.StrategyService.QueryStrategyEntityByStrategyId(ctx, strategyId)
	if err != nil {
		return nil, err
	}
	ruleModels, err := strategy.GetStrsRuleModels()
	if err != nil {
		return nil, err
	}
	ruleActionEntityi, err := t.DoCheckRaffleBeforeLogic(ctx, &StrategyEntity.RaffleFactorEntity{
		UserId:     userId,
		StrategyId: strategyId,
	}, ruleModels...)
	if err != nil {
		return nil, err
	}
	ruleActionEntity, ok := ruleActionEntityi.(*StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity])
	if !ok {
		return nil, cerr.LogError(errors.New("ruleActionEntity Conversion failure or ruleActionEntity is null"))
	}
	if vo.TAKE_OVER.Code == ruleActionEntity.Code {
		if model.RULE_BLACKLIST.Code == ruleActionEntity.RuleModel {
			return &StrategyEntity.RaffleAwardEntity{
				AwardId: int64(ruleActionEntity.Data.AwardId),
			}, nil
		} else if model.RULE_WEIGHT.Code == ruleActionEntity.RuleModel {
			raffleBeforeEntity := ruleActionEntity.Data
			ruleWeightValueKey := raffleBeforeEntity.RuleWeightValueKey
			ruleWeightValueKeyNum, err := strconv.ParseInt(ruleWeightValueKey, 10, 64)
			if err != nil {
				return nil, err
			}
			awardId, err := t.StrategyDispatch.GetRandomAwardId(ctx, strategyId, ruleWeightValueKeyNum)
			if err != nil {
				return nil, err
			}
			return &StrategyEntity.RaffleAwardEntity{
				AwardId: awardId,
			}, nil
		}
	}
	awardId, err := t.StrategyDispatch.GetRandomAwardIdBase(ctx, strategyId)
	if err != nil {
		return nil, err
	}
	return &StrategyEntity.RaffleAwardEntity{
		AwardId: awardId,
	}, nil
}

// 抽象方法
func (t *TemplateRaffleStrategy) DoCheckRaffleBeforeLogic(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error) {
	if t.CheckRaffleBeforeLogicFunc == nil {
		panic("CheckRaffleBeforeLogicFunc is not setting")
	}
	return t.CheckRaffleBeforeLogicFunc(ctx, entity, logics...)
}
