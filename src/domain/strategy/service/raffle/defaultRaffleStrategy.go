package raffle

import (
	"context"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/factory"
	LogicModel "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/factory/model"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽奖过滤引擎
 * @Date: 2024-08-07 22:43
 */
type DefaultRaffleStrategy struct {
	TemplateRaffleStrategy
	DefaultLogicFactory factory.DefaultLogicFactory
}

func NewDefaultRaffleStrategy(StrategyService repository.StrategyService, strategyDispatch armory.StrategyDispath) *DefaultRaffleStrategy {
	defaultRaffleStrategy := &DefaultRaffleStrategy{
		TemplateRaffleStrategy: *NewTemplateRaffleStrategy(StrategyService, strategyDispatch),
		DefaultLogicFactory:    *factory.NewDefaultLogicFactory(StrategyService),
	}

	defaultRaffleStrategy.TemplateRaffleStrategy.CheckRaffleBeforeLogicFunc = defaultRaffleStrategy.DoCheckRaffleBeforeLogic
	return defaultRaffleStrategy
}
func (t *DefaultRaffleStrategy) DoCheckRaffleBeforeLogic(ctx context.Context, entity *StrategyEntity.RaffleFactorEntity, logics ...string) (StrategyEntity.RaffleActionEntityInterface, error) {
	logicFilterGroup := t.DefaultLogicFactory.OpenLogicFilter()
	// 将黑名单策略提前
	var isNotNull bool = false
	for i, logic := range logics {
		if logic == LogicModel.RULE_BLACKLIST.Code {
			temp := logics[0]
			logics[0] = logics[i]
			logics[i] = temp
			isNotNull = true
			break
		}
	}
	if isNotNull {
		logicFilter := logicFilterGroup[LogicModel.RULE_BLACKLIST.Code]
		ruleMatterEntity := &StrategyEntity.RuleMatterEntity{
			UserId:     entity.UserId,
			StrategyId: entity.StrategyId,
			RuleModel:  LogicModel.RULE_BLACKLIST.Code,
		}
		ruleActionEntity, err := logicFilter.Filter(ctx, *ruleMatterEntity)
		if err != nil {
			return nil, err
		}
		if !(vo.ALLOW.Code == ruleActionEntity.GetCode()) {
			return ruleActionEntity, nil
		}
	}
	// 去除黑名单
	ruleList := logics
	if isNotNull {
		ruleList = ruleList[1:]
	}
	ruleActionEntity := new(StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity])
	for _, ruleModel := range ruleList {
		logicFilter := logicFilterGroup[ruleModel]
		if logicFilter == nil {
			continue
		}
		ruleMatterEntity := &StrategyEntity.RuleMatterEntity{
			UserId:     entity.UserId,
			StrategyId: entity.StrategyId,
			RuleModel:  ruleModel,
		}
		ruleActionEntity, err := logicFilter.Filter(ctx, *ruleMatterEntity)
		if err != nil {
			return nil, err
		}
		logx.Infof("抽奖前规则过滤 userId: %s ruleModel: %s code: %s info: %s", entity.UserId, ruleModel, ruleActionEntity.GetCode(), ruleActionEntity.GetInfo())
		if !(vo.ALLOW.Code == ruleActionEntity.GetCode()) {
			// ruleActionEntity1, _ := ruleActionEntity.(*StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity])
			// logx.Debug(ruleActionEntity1.RuleModel)
			return ruleActionEntity, nil
		}
	}
	ruleActionEntity = nil
	return ruleActionEntity, nil
}
