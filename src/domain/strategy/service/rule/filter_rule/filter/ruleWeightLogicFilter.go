package filter

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	LogicModel "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/factory/model"
	filter_interface "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/filter/interface"
	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 权重过滤器引擎(产品)
 * @Date: 2024-08-06 00:51
 */
var _ filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface] = (*RuleWeightLogicFilter)(nil)

type RuleWeightLogicFilter struct {
	strategyService repository.StrategyService
	logicModel      LogicModel.LogicModel
	userScore       int64
}

func NewRuleWeightLogicFilter(strategyService repository.StrategyService) *RuleWeightLogicFilter {
	return &RuleWeightLogicFilter{
		strategyService: strategyService,
		logicModel:      LogicModel.RULE_WEIGHT,
		userScore:       -1, // 默认负一
	}
}
func (filter *RuleWeightLogicFilter) Filter(ctx context.Context, ruleMatter StrategyEntity.RuleMatterEntity) (StrategyEntity.RaffleActionEntityInterface, error) {
	filter.userScore = 6000
	logx.Infof("规则过滤-权重范围 userId:%s strategyId:%d ruleModel:%s", ruleMatter.UserId, ruleMatter.StrategyId, ruleMatter.RuleModel)
	// userId := ruleMatter.UserId
	strategyId := ruleMatter.StrategyId
	ruleValue, err := filter.strategyService.QueryStrategyRuleValue(ctx, ruleMatter.StrategyId, ruleMatter.AwardId, ruleMatter.RuleModel)
	if err != nil {
		return nil, err
	}
	analyticalValueGroup, analyticalSortedKeys, err := getAnalyticalValue(ruleValue)
	if err != nil {
		return nil, err
	}
	// 未设置权重或者查不到->直接放行
	if len(analyticalValueGroup) == 0 {
		return &StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity]{
			Code: vo.ALLOW.Code,
			Info: vo.ALLOW.Info,
		}, nil
	}
	sort.Slice(analyticalSortedKeys, func(i, j int) bool {
		return analyticalSortedKeys[i] < analyticalSortedKeys[j]
	})
	var nextValue int64 = -1
	for _, analyticalSortedKey := range analyticalSortedKeys {
		if filter.userScore < analyticalSortedKey {
			break
		}
		nextValue = analyticalSortedKey
	}
	if nextValue != -1 {
		return &StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity]{
			Code:      vo.TAKE_OVER.Code,
			Info:      vo.TAKE_OVER.Info,
			RuleModel: LogicModel.RULE_WEIGHT.Code,
			Data: StrategyEntity.RaffleBeforeEntity{
				StrategyId:         strategyId,
				RuleWeightValueKey: strconv.FormatInt(nextValue, 10),
			},
		}, nil
	}
	// 抽数不够->直接放行
	return &StrategyEntity.RaffleActionEntity[StrategyEntity.RaffleBeforeEntity]{
		Code: vo.ALLOW.Code,
		Info: vo.ALLOW.Info,
	}, nil
}
func getAnalyticalValue(ruleValue string) (map[int64]string, []int64, error) {
	ruleValueGroups := strings.Split(ruleValue, common.SPACE)
	ruleValueMap := make(map[int64]string)
	analyticalSortedKeys := make([]int64, 0)
	for _, ruleValueKey := range ruleValueGroups {
		parts := strings.Split(ruleValueKey, common.COLON)
		if len(parts) != 2 {
			return nil, nil, cerr.LogError(errors.New("rule_weight rule_rule invalid input format:" + ruleValueKey))
		}
		key, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, nil, err
		}
		analyticalSortedKeys = append(analyticalSortedKeys, key)
		ruleValueMap[key] = parts[1]
	}
	return ruleValueMap, analyticalSortedKeys, nil
}
func (filter *RuleWeightLogicFilter) GetLogicModel() LogicModel.LogicModel {
	return filter.logicModel
}
