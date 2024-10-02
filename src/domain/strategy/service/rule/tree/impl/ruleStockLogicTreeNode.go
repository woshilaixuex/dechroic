package tree_impl

import (
	"github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-09-21 21:40
 */
type RuleStockLogicTreeNode struct {
	key string
}

func NewRuleStockLogicTreeNode() *RuleStockLogicTreeNode {
	return &RuleStockLogicTreeNode{key: "rule_stock"}
}

func (node *RuleStockLogicTreeNode) GetKey() string {
	return node.key
}
func (node *RuleStockLogicTreeNode) Logic(userId string, strategyId int64, awardId int32) *entity.TreeActionEntity {
	return &entity.TreeActionEntity{
		RuleLogicCheckType: vo.TAKE_OVER,
	}
}
