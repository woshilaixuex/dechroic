package tree_factory

import (
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/tree"
	tree_engine "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/tree/factory/engine"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 规则树工厂
 * @Date: 2024-09-21 21:41
 */
type DefultTreeFactory struct {
	RuleLogicNodeGroup map[string]tree.ILogicTreeNode
}

func NewDefultTreeFactory(group ...tree.ILogicTreeNode) *DefultTreeFactory {
	nodeGroup := make(map[string]tree.ILogicTreeNode)
	for _, node := range group {
		nodeGroup[node.GetKey()] = node
	}
	tree_engine.LogicNodeGroup = nodeGroup
	logx.Debug(nodeGroup)
	return &DefultTreeFactory{RuleLogicNodeGroup: nodeGroup}
}
func (factory *DefultTreeFactory) OpenLogicTree(ruleTreeVO *vo.RuleTreeVO) tree_engine.IDecisionTreeEngine {
	return tree_engine.NewDecisionTreeEngin(ruleTreeVO)
}
