package filter

import (
	"context"
	"strconv"
	"strings"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	LogicModel "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/factory/model"

	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 黑名单过滤引擎
 * @Date: 2024-08-04 21:42
 */

type ILogicFilter[T StrategyEntity.RaffleActionEntityInterface] interface {
	Filter(ctx context.Context, ruleMatter StrategyEntity.RuleMatterEntity) (T, error)
}

var _ ILogicFilter[StrategyEntity.RaffleActionEntityInterface] = (*RuleBackListLogicFilter)(nil)

type RuleBackListLogicFilter struct {
	ILogicFilter[StrategyEntity.RaffleActionEntityInterface]
	strategyService repository.StrategyService
}

// 黑名单过滤
func (filter *RuleBackListLogicFilter) Filter(ctx context.Context, ruleMatter StrategyEntity.RuleMatterEntity) (StrategyEntity.RaffleActionEntityInterface, error) {
	logx.Infof("规则过滤-黑名单 userId:%s strategyId:%d ruleModel:%s", ruleMatter.UserId, ruleMatter.StrategyId, ruleMatter.RuleModel)
	userId := ruleMatter.UserId
	ruleValue, err := filter.strategyService.QueryStrategyRuleValue(ctx, ruleMatter.StrategyId, ruleMatter.AwardId, ruleMatter.RuleModel)
	if err != nil {
		return nil, err
	}
	splitRuleValue := strings.Split(ruleValue, common.COLON)
	awardId, err := strconv.ParseInt(splitRuleValue[0], 10, 32)
	if err != nil {
		return nil, err
	}
	userBlackIds := strings.Split(splitRuleValue[1], common.SPLIT)
	for _, blackId := range userBlackIds {
		if blackId == userId {
			return &StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity]{
				Code:      vo.TAKE_OVER.Code,
				Info:      vo.TAKE_OVER.Info,
				RuleModel: LogicModel.RULE_BLACKLIST.Code,
				Data: StrategyEntity.RaffleBeforeEntity{
					StrategyId: ruleMatter.StrategyId,
					AwardId:    int32(awardId),
				},
			}, nil
		}
	}
	return &StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity]{
		Code: vo.ALLOW.Code,
		Info: vo.ALLOW.Info,
	}, nil
}