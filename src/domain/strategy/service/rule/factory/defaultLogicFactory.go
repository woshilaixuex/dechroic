package factory

import (
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/factory/model"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-08-06 00:22
 */
type DefaultLogicFactory struct {
	logicFilterMap map[string]filter.ILogicFilter[any] // 也可以不是any
	model.LogicModel
}
