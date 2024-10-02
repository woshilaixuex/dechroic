package chain_ipml

import (
	"context"
	"strconv"
	"strings"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/chain"
	LogicModel "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/factory/model"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 黑名单过滤节点
 * @Date: 2024-08-16 00:28
 */

type BackListLogicChain struct {
	chain.LogicChainNode
	strategyService repository.StrategyService
}

func NewBackListLogicChain(strategyService repository.StrategyService) *BackListLogicChain {
	backListLogicChain := &BackListLogicChain{
		LogicChainNode:  *chain.NewLogicChainNode(),
		strategyService: strategyService,
	}
	backListLogicChain.Realize(backListLogicChain.Logic)
	return backListLogicChain
}
func (c *BackListLogicChain) Logic(ctx context.Context, userId string, strategyId int64) (*data.StrategyAwardChanVO, error) {
	logx.Infof("抽奖责任链-黑名单：userId：%s strategyId：%d ruleModel：%s", userId, strategyId, c.ModelType())
	ruleValue, err := c.strategyService.QueryStrategyRule(ctx, strategyId, c.ModelType())
	if err != nil {
		return nil, nil
	}
	splitRuleValue := strings.Split(ruleValue.RuleValue, common.COLON)
	awardId, err := strconv.ParseInt(splitRuleValue[0], 10, 32)
	if err != nil {
		return nil, err
	}
	userBlackIds := strings.Split(splitRuleValue[1], common.SPLIT)
	for _, blackId := range userBlackIds {
		if blackId == userId {
			return &data.StrategyAwardChanVO{
				AwardId:          int32(awardId),
				AwardRuleValue:   c.ModelType(),
				ChanVOLogicModel: data.RULE_BLACKLIST,
			}, nil
		}
	}
	logx.Infof("抽奖责任链-黑名单放行：userId：%s strategyId：%d ruleModel：%s", userId, strategyId, c.ModelType())
	return c.Next().Logic(ctx, userId, strategyId)
}
func (c *BackListLogicChain) ModelType() string {
	return LogicModel.RULE_BLACKLIST.Code
}
