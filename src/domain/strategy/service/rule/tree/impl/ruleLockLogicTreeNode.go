package tree_impl

import (
	"github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-09-21 21:38
 */

type RuleLockLogicTreeNode struct {
	key string
}

func NewRuleLockLogicTreeNode() *RuleLockLogicTreeNode {
	return &RuleLockLogicTreeNode{key: "rule_lock"}
}

func (node *RuleLockLogicTreeNode) GetKey() string {
	return node.key
}
func (node *RuleLockLogicTreeNode) Logic(userId string, strategyId int64, awardId int32) *entity.TreeActionEntity {
	return &entity.TreeActionEntity{
		RuleLogicCheckType: vo.ALLOW,
	}
}
