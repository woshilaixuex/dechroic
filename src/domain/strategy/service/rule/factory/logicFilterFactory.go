package factory

import (
	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/repository"
	filter_interface "github.com/delyr1c/dechoric/src/domain/strategy/service/rule/filter/interface"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 过滤器工厂（逻辑工厂的下层工厂）
 * @Date: 2024-08-08 13:30
 */
type LogicFilterFactoryOption func(*LogicFilterFactory, repository.StrategyService)
type LogicFilterFactory struct {
	options []LogicFilterFactoryOption
	Filters []filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface]
}

func NewLogicFilterFactory(options ...LogicFilterFactoryOption) *LogicFilterFactory {
	return &LogicFilterFactory{
		options: options,
	}
}
func (factory *LogicFilterFactory) MakeFilters(strategyService repository.StrategyService) []filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface] {
	factory.Filters = make([]filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface], 0)
	for _, option := range factory.options {
		option(factory, strategyService)
	}

	return factory.Filters
}

// 注册过滤器->给MakeFilter使用的
func RegisterFilter(filter filter_interface.ILogicFilter[StrategyEntity.RaffleActionEntityInterface]) LogicFilterFactoryOption {
	return func(factory *LogicFilterFactory, strategyService repository.StrategyService) {
		factory.Filters = append(factory.Filters, filter)
	}
}
