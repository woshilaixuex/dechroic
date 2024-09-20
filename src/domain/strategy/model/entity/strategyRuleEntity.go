package entity

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: service层中rule策略所需实体
 * @Date: 2024-08-05 19:20
 */
// 策略规则实体
type StrategyRuleEntity struct {
	StrategyId int64         `db:"strategy_id"` // 抽奖策略ID
	AwardId    sql.NullInt64 `db:"award_id"`    // 抽奖奖品ID【规则类型为策略，则不需要奖品ID】
	RuleType   int64         `db:"rule_type"`   // 抽象规则类型；1-策略规则、2-奖品规则
	RuleModel  string        `db:"rule_model"`  // 抽奖规则类型【rule_random - 随机值计算、rule_lock - 抽奖几次后解锁、rule_luck_award - 幸运奖(兜底奖品)】
	RuleValue  string        `db:"rule_value"`  // 抽奖规则比值
	RuleDesc   string        `db:"rule_desc"`   // 抽奖规则描述
}

func (e *StrategyRuleEntity) GetRule() (map[int64]map[int64]bool, error) {
	if e.RuleModel != "rule_weight" {
		logx.Info("RuleModel is not rule_weight", 4)
		return nil, nil
	}
	ruleMap := make(map[int64]map[int64]bool)
	ruleValues := strings.Split(e.RuleValue, common.SPACE)
	if ruleValues == nil {
		return nil, cerr.LogError(errors.New("RuleValue is null"))
	}
	for _, ruleValue := range ruleValues {
		keyVal := strings.Split(ruleValue, common.COLON)
		if len(keyVal) != 2 {
			return nil, cerr.LogError(errors.New("invalid ruleValue format: key val"))
		}
		key, err := strconv.ParseInt(keyVal[0], 10, 64)
		if err != nil {
			return nil, cerr.LogError(errors.New("failed to convert key[0] to int64"))
		}
		vals := strings.Split(keyVal[1], common.SPLIT)
		set := make(map[int64]bool)
		for _, val := range vals {
			numval, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, cerr.LogError(errors.New("failed to convert val to int64"))
			}
			set[numval] = true
		}
		ruleMap[key] = set
	}
	return ruleMap, nil
}
