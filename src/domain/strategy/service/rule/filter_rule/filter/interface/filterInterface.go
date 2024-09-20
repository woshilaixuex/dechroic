package filter_interface

import (
	"context"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	LogicModel "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter_rule/factory/model"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 工厂接口（用来约束和切断依赖环的）
 * @Date: 2024-08-08 15:22
 */

// 内部除了规则模型一律都不暴露
type ILogicFilter[T StrategyEntity.RaffleActionEntityInterface] interface {
	Filter(ctx context.Context, ruleMatter StrategyEntity.RuleMatterEntity) (T, error)
	GetLogicModel() LogicModel.LogicModel
}
