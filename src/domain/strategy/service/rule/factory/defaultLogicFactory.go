package factory

import (
	"sync"

	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/factory/model"
	"github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter"
	filter_interface "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter/interface"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 逻辑工厂
 * @Date: 2024-08-06 00:22
 */
type DefaultLogicFactory struct {
	logicFilterMap map[string]filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface]
	models         []model.LogicModel
	mu             sync.RWMutex
}

func NewDefaultLogicFactory(strategyService repository.StrategyService) *DefaultLogicFactory {
	logicFilterFactory := NewLogicFilterFactory(
		RegisterFilter(filter.NewRuleBackListLogicFilter(strategyService)),
		RegisterFilter(filter.NewRuleWeightLogicFilter(strategyService)),
	)
	filters := logicFilterFactory.MakeFilters(strategyService)
	factory := &DefaultLogicFactory{
		logicFilterMap: make(map[string]filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface]),
		models:         make([]model.LogicModel, 0),
	}
	factory.InitFilters(filters)
	return factory
}

// 初始化
func (f *DefaultLogicFactory) InitFilters(filters []filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface]) {
	for _, filter := range filters {
		model := filter.GetLogicModel()
		f.logicFilterMap[model.Code] = filter
		f.models = append(f.models, model)
	}
}
func (f *DefaultLogicFactory) OpenLogicFilter() map[string]filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface] {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.logicFilterMap
}
