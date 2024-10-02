package tree_engine

import (
	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/tree"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 引擎
 * @Date: 2024-09-22 16:52
 */

var LogicNodeGroup = make(map[string]tree.ILogicTreeNode)

type IDecisionTreeEngine interface {
	Process(userId string, strategyId int64, awardId int32) *data.StrategyAwardTreeVO
}
type DecisionTreeEngine struct {
	RuleTreeVO *vo.RuleTreeVO
}

func NewDecisionTreeEngin(ruleTreeVO *vo.RuleTreeVO) *DecisionTreeEngine {
	return &DecisionTreeEngine{
		RuleTreeVO: ruleTreeVO,
	}

}
func (engine *DecisionTreeEngine) Process(userId string, strategyId int64, awardId int32) *data.StrategyAwardTreeVO {
	var strategyAwardData *data.StrategyAwardTreeVO = nil
	nextNode := engine.RuleTreeVO.TreeRootRuleNode
	treeNodeMap := engine.RuleTreeVO.TreeNodeMap
	ruleTreeNode := treeNodeMap[nextNode]
	for ruleTreeNode != nil {
		logicTreeNode := LogicNodeGroup[ruleTreeNode.RuleKey]
		logicEntity := logicTreeNode.Logic(userId, strategyId, awardId)
		ruleLogicCheckTypeVO := logicEntity.RuleLogicCheckType
		strategyAwardData = logicEntity.StrategyAwardData
		logx.Infof("决策树引擎【%v】treeId:%v node:%v code:%v", engine.RuleTreeVO.TreeName, engine.RuleTreeVO.TreeId, nextNode, ruleLogicCheckTypeVO.Code)
		nextNode = engine.NextNode(ruleLogicCheckTypeVO.Code, ruleTreeNode.RuleTreeNodeLineVOSlice)
		ruleTreeNode = treeNodeMap[nextNode]
	}
	return strategyAwardData
}
func (engine *DecisionTreeEngine) NextNode(matterValue string, treeNodeLineVOLSlice []*vo.RuleTreeNodeLineVO) string {
	if len(treeNodeLineVOLSlice) == 0 {
		return ""
	}
	decisionLogic := func(nodeLine *vo.RuleTreeNodeLineVO) bool {
		switch nodeLine.RuleLimitTypeVO {
		case &vo.EQUAL:
			return matterValue == nodeLine.RuleLogicCheckTypeVO.Code
		case &vo.GT:
			return false
		case &vo.LT:
			return false
		case &vo.GE:
			return false
		case &vo.LE:
			return false
		default:
			return false
		}
	}
	for _, nodeLine := range treeNodeLineVOLSlice {

		if decisionLogic(nodeLine) {
			return nodeLine.RuleNodeTo
		}

	}
	panic("nextNode 计算失败，未找到可执行节点")
}
