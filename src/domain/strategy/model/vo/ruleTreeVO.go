package vo

import (
	"fmt"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 规则树对象
 * @Date: 2024-09-21 19:20
 */
type RuleTreeVO struct {
	TreeId           string
	TreeName         string
	TreeDesc         string
	TreeRootRuleNode string
	TreeNodeMap      map[string]*RuleTreeNodeVO
}

func NewRuleTreeVO() *RuleTreeVO {
	return &RuleTreeVO{
		TreeNodeMap: make(map[string]*RuleTreeNodeVO),
	}
}

// 遍历并打印整颗规则树
func (ruleTree *RuleTreeVO) TraverseRuleTree() {
	if ruleTree == nil || ruleTree.TreeRootRuleNode == "" {
		fmt.Println("空的规则树或根节点不存在")
		return
	}

	for _, node := range ruleTree.TreeNodeMap {
		node.RuleTreeNodeVoInfo()
	}

}
