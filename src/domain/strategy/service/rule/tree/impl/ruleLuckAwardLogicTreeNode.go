package tree_impl

import (
	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-09-21 21:40
 */

type RuleLuckAwardLogicTreeNode struct {
	key string
}

func NewRuleLuckAwardLogicTreeNode() *RuleLuckAwardLogicTreeNode {
	return &RuleLuckAwardLogicTreeNode{key: "rule_luck_award"}
}

func (node *RuleLuckAwardLogicTreeNode) GetKey() string {
	return node.key
}
func (node *RuleLuckAwardLogicTreeNode) Logic(userId string, strategyId int64, awardId int32) *entity.TreeActionEntity {
	return &entity.TreeActionEntity{
		RuleLogicCheckType: vo.TAKE_OVER,
		StrategyAwardData: &data.StrategyAwardTreeVO{
			AwardId:        101,
			AwardRuleValue: "1000",
		},
	}
}
