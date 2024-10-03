package chain_ipml

import (
	"context"
	"errors"
	"sort"
	"strconv"
	"strings"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/armory"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain"
	LogicModel "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/factory/model"
	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 权重过滤节点
 * @Date: 2024-08-16 00:27
 */
type RuleWeightLogicChain struct {
	chain.LogicChainNode
	strategyService  repository.StrategyService
	strategyDispatch armory.StrategyDispath
	userScore        int64
}

func NewRuleWeightLogicChain(strategyService repository.StrategyService, strategyDispatch armory.StrategyDispath) *RuleWeightLogicChain {
	ruleWeightLogicChain := &RuleWeightLogicChain{
		LogicChainNode:   *chain.NewLogicChainNode(),
		strategyService:  strategyService,
		strategyDispatch: strategyDispatch,
		userScore:        6000,
	}
	ruleWeightLogicChain.Realize(ruleWeightLogicChain.Logic)
	return ruleWeightLogicChain
}
func (c *RuleWeightLogicChain) Logic(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error) {
	logx.Infof("抽奖责任链-权重校验：userId：%s strategyId：%d ruleModel：%s", userId, strategyId, c.ModelType())
	ruleValue, err := c.strategyService.QueryStrategyRule(ctx, strategyId, c.ModelType())
	if err != nil {
		return &data.StrategyAwardChanVO{}, nil
	}
	logx.Debug(ruleValue.RuleValue)
	analyticalValueGroup, analyticalSortedKeys, err := getAnalyticalValue(ruleValue.RuleValue)
	if err != nil {
		return nil, err
	}
	if len(analyticalValueGroup) == 0 {
		return &data.StrategyAwardChanVO{}, nil
	}
	sort.Slice(analyticalSortedKeys, func(i, j int) bool {
		return analyticalSortedKeys[i] < analyticalSortedKeys[j]
	})
	var nextValue int64 = -1
	for _, analyticalSortedKey := range analyticalSortedKeys {
		if c.userScore < analyticalSortedKey {
			break
		}
		nextValue = analyticalSortedKey
	}
	if nextValue != -1 {
		awardId, err := c.strategyDispatch.GetRandomAwardId(ctx, strategyId, nextValue)
		if err != nil {
			return nil, err
		}
		logx.Infof("抽奖责任链-权重接管：userId：%s strategyId：%d ruleModel：%s", userId, strategyId, c.ModelType())
		return &data.StrategyAwardChanVO{
			AwardId:          int32(awardId),
			AwardRuleValue:   c.ModelType(),
			ChanVOLogicModel: data.RULE_WEIGHT,
		}, nil
	}
	// 抽数不够->直接放行
	return c.Next().Logic(ctx, userId, strategyId)
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
func (c *RuleWeightLogicChain) ModelType() string {
	return LogicModel.RULE_WEIGHT.Code
}
