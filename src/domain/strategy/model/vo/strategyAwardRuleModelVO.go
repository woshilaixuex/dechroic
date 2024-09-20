package vo

import (
	"strings"

	LogicModel "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/factory/model"
	"github.com/delyr1c/dechoric/src/types/common"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-08-12 18:31
 */
type StrategyAwardRuleModelVO struct {
	RuleModels string
}

func (vo *StrategyAwardRuleModelVO) RaffleCenterModelList() []string {
	ruleModelsSlice := make([]string, 0)
	ruleModelValues := strings.Split(vo.RuleModels, common.SPLIT)
	for _, ruleModelValue := range ruleModelValues {
		if LogicModel.IsCenter(ruleModelValue) {
			ruleModelsSlice = append(ruleModelsSlice, ruleModelValue)
		}
	}
	return ruleModelsSlice
}
