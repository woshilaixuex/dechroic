package entity

import (
	"strings"

	"github.com/delyr1c/dechoric/src/types/common"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: service层中armoy策略所需实体
 * @Date: 2024-06-13 22:35
 */

// 具体策略实体
type StrategyAwardEntity struct {
	Id                int64           `db:"id" json:"id"`                                   // 自增ID
	AwardId           int64           `db:"award_id" json:"award_id"`                       // 抽奖策略ID
	AwardCount        int64           `db:"award_count" json:"award_count"`                 // 奖品库存总量
	AwardCountSurplus int64           `db:"award_count_surplus" json:"award_count_surplus"` // 奖品库存剩余
	AwardRate         common.BigFloat `db:"award_rate" json:"award_rate"`                   // 奖品中奖概率
}

// 策略实体
type StrategyEntity struct {
	StrategyId   int64  `db:"strategy_id"`   // 抽奖策略ID
	StrategyDesc string `db:"strategy_desc"` // 抽奖策略描述
	RuleModels   string `db:"rule_model"`    // 抽奖规则的模型
}

func (e *StrategyEntity) GetStrsRuleModels() ([]string, error) {
	if strings.TrimSpace(e.RuleModels) == "" {
		return nil, nil
	}
	ruleModels := strings.Split(e.RuleModels, common.SPLIT)
	return ruleModels, nil
}

func (e *StrategyEntity) GetRuleWeight() (string, error) {
	ruleModels, err := e.GetStrsRuleModels()
	if err != nil {
		return "", err
	}
	for _, ruleModel := range ruleModels {
		if ruleModel == "rule_weight" {
			return ruleModel, nil
		}
	}
	return "", nil
}
