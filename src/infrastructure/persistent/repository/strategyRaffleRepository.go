package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyRule"
	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/delyr1c/dechoric/src/types/common"
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
		return nil, cerr.LogError(err)
	}
	if ruleTreeVOCache != nil {
		ruleTreeVO := new(vo.RuleTreeVO)
		err := json.Unmarshal([]byte(ruleTreeVOCache.(string)), ruleTreeVO)
		if err != nil {
			return nil, cerr.LogError(err) // 处理 JSON 解码错误
		}
		return ruleTreeVO, nil
	}
	ruleTree, err := s.TreeRuleModel.FindOneByTreeId(ctx, treeId)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	ruleTreeNodes, err := s.TreeRuleNodeModel.FindRuleTreeNodeListByTreeId(ctx, treeId)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	ruleTreeNodeLines, err := s.TreeRuleNodeLineModel.FindRuleTreeNodeLineListByTreeId(ctx, treeId)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	ruleTreeNodeLineMap := make(map[string][]*vo.RuleTreeNodeLineVO)
	for _, ruleTreeNodeLine := range ruleTreeNodeLines {
		ruleTreeNodeLineEle := &vo.RuleTreeNodeLineVO{
			TreeId:       ruleTreeNodeLine.TreeId,
			RuleNodeFrom: ruleTreeNodeLine.RuleNodeFrom,
			RuleNodeTo:   ruleTreeNodeLine.RuleNodeTo,
		}
		ruleTreeNodeLineEle.RuleLimitTypeVO, err = vo.GetRuleLimitTypeVOByStr(ruleTreeNodeLine.RuleLimitType)
		if err != nil {
			return nil, cerr.LogError(err)
		}
		ruleTreeNodeLineEle.RuleLogicCheckTypeVO, err = vo.GetRuleLogicCheckTypeVOByStr(ruleTreeNodeLine.RuleLimitValue)
		if err != nil {
			return nil, cerr.LogError(err)
		}
		ruleTreeNodeLineVOList, exists := ruleTreeNodeLineMap[ruleTreeNodeLine.RuleNodeFrom]
		if !exists {
			ruleTreeNodeLineVOList = make([]*vo.RuleTreeNodeLineVO, 0, 10)
		}
		// 将新的元素添加到列表中
		ruleTreeNodeLineVOList = append(ruleTreeNodeLineVOList, ruleTreeNodeLineEle)
		// 更新 map
		ruleTreeNodeLineMap[ruleTreeNodeLine.RuleNodeFrom] = ruleTreeNodeLineVOList
	}
	treeNodeMap := make(map[string]*vo.RuleTreeNodeVO)
	for _, ruleTreeNode := range ruleTreeNodes {
		ruleTreeNodeEle := &vo.RuleTreeNodeVO{
			TreeId:                  ruleTreeNode.TreeId,
			RuleKey:                 ruleTreeNode.RuleKey,
			RukeDesc:                ruleTreeNode.RuleDesc,
			RuleValue:               ruleTreeNode.RuleValue.String,
			RuleTreeNodeLineVOSlice: make([]*vo.RuleTreeNodeLineVO, 0),
		}
		ruleTreeNodeEle.RuleTreeNodeLineVOSlice = ruleTreeNodeLineMap[ruleTreeNode.RuleKey]
		treeNodeMap[ruleTreeNodeEle.RuleKey] = ruleTreeNodeEle
	}
	ruleTreeDB := &vo.RuleTreeVO{
		TreeId:           ruleTree.TreeId,
		TreeName:         ruleTree.TreeName,
		TreeDesc:         ruleTree.TreeDesc.String,
		TreeRootRuleNode: ruleTree.TreeNodeRuleKey,
		TreeNodeMap:      treeNodeMap,
	}

	s.RedisService.SetValue(ctx, cacheKey, ruleTreeDB)
	return ruleTreeDB, nil
}
