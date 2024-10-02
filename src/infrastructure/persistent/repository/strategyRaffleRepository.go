package repository

import (
	"context"
	"errors"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyRule"
	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/jinzhu/copier"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: raffle的repository实现
 * @Date: 2024-08-05 22:14
 */
func (s *StrategyRepository) QueryStrategyRuleValue(ctx context.Context, strategyId int64, awardId int32, roleModel string) (string, error) {
	strategyRuleReq := &strategyRule.FindStrategyRuleReq{
		StrategyId: &strategyId,
		RuleModel:  &roleModel,
	}
	if awardId == 0 {
		strategyRuleReq.AwardId = nil
	}
	return s.StrategyRuleModel.FindRuleValueByReq(ctx, strategyRuleReq)
}
func (s *StrategyRepository) QueryStrategyAwardRuleModelVO(ctx context.Context, strategyId int64, awardId int32) (*vo.StrategyAwardRuleModelVO, error) {
	newAwardId := int64(awardId)
	StrategyAwardReq := &strategyAward.FindStrategyAwardReq{
		StrategyId: &strategyId,
		AwardId:    &newAwardId,
	}
	StrategyAwards, err := s.StrategyAwardModel.FindByReq(ctx, StrategyAwardReq)
	if err != nil {
		return nil, err
	}
	StrategyAward := StrategyAwards[0]
	if !StrategyAward.RuleModels.Valid {
		return nil, cerr.LogError(errors.New("Dao RuleModels is null "))
	}
	return &vo.StrategyAwardRuleModelVO{
		RuleModels: StrategyAward.RuleModels.String,
	}, nil
}

// 查询规则树
func (s *StrategyRepository) QueryRuleTreeVOByTreeId(ctx context.Context, treeId string) (*vo.RuleTreeVO, error) {
	cacheKey := common.RedisKeys.RuleTreeVOKey + treeId
	ruleTreeVOCache, err := s.RedisService.GetValue(ctx, cacheKey)
	if err != nil {
		return nil, err
	}
	if ruleTreeVOCache != nil {
		ruleTreeVO, ok := ruleTreeVOCache.(vo.RuleTreeVO)
		if ok {
			return &ruleTreeVO, nil
		}
	}
	ruleTree, err := s.TreeRuleModel.FindOneByTreeId(ctx, treeId)
	if err != nil {
		return nil, err
	}
	ruleTreeNodes, err := s.TreeRuleNodeModel.FindRuleTreeNodeListByTreeId(ctx, treeId)
	if err != nil {
		return nil, err
	}
	ruleTreeNodeLines, err := s.TreeRuleNodeLineModel.FindRuleTreeNodeLineListByTreeId(ctx, treeId)
	if err != nil {
		return nil, err
	}
	ruleTreeNodeLineMap := make(map[string][]*vo.RuleTreeNodeLineVO)
	for _, ruleTreeNodeLine := range ruleTreeNodeLines {
		ruleTreeNodeLineEle := new(vo.RuleTreeNodeLineVO)
		err := copier.Copy(ruleTreeNodeLineEle, ruleTreeNodeLine)
		if err != nil {
			return nil, err
		}
		ruleTreeNodeLineEle.RuleLimitTypeVO, err = vo.GetRuleLimitTypeVOByStr(ruleTreeNodeLine.RuleLimitType)
		if err != nil {
			return nil, err
		}
		ruleTreeNodeLineEle.RuleLogicCheckTypeVO, err = vo.GetRuleLogicCheckTypeVOByStr(ruleTreeNodeLine.RuleLimitValue)
		ruleTreeNodeLineVOList, exists := ruleTreeNodeLineMap[ruleTreeNodeLine.RuleNodeFrom]
		if !exists {
			ruleTreeNodeLineVOList = []*vo.RuleTreeNodeLineVO{}
		}
		// 将新的元素添加到列表中
		ruleTreeNodeLineVOList = append(ruleTreeNodeLineVOList, ruleTreeNodeLineEle)
		// 更新 map
		ruleTreeNodeLineMap[ruleTreeNodeLine.RuleNodeFrom] = ruleTreeNodeLineVOList
	}
	treeNodeMap := make(map[string]*vo.RuleTreeNodeVO)
	for _, ruleTreeNode := range ruleTreeNodes {
		ruleTreeNodeEle := new(vo.RuleTreeNodeVO)
		err := copier.Copy(ruleTreeNodeEle, ruleTreeNode)
		if err != nil {
			return nil, err
		}
		ruleTreeNodeEle.RuleTreeNodeLineVOSlice = ruleTreeNodeLineMap[ruleTreeNode.RuleKey]
		treeNodeMap[ruleTreeNodeEle.RuleKey] = ruleTreeNodeEle
	}
	ruleTreeDB := new(vo.RuleTreeVO)
	err = copier.Copy(ruleTreeDB, ruleTree)
	if err != nil {
		return nil, err
	}
	s.RedisService.SetValue(ctx, cacheKey, ruleTree)
	return ruleTreeDB, err
}
