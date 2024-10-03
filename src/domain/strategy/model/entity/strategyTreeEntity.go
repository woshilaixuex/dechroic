package entity

import (
	"github.com/delyr1c/dechoric/src/domain/strategy/model/data"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-09-23 21:39
 */
type TreeActionEntity struct {
	RuleLogicCheckType vo.RuleLogicCheckTypeVO
	StrategyAwardData  *data.StrategyAwardTreeVO
}
