package strategyAward

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ StrategyAwardModel = (*customStrategyAwardModel)(nil)

type (
	// StrategyAwardModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStrategyAwardModel.
	StrategyAwardModel interface {
		strategyAwardModel
	}

	customStrategyAwardModel struct {
		*defaultStrategyAwardModel
	}
)

// NewStrategyAwardModel returns a model for the database table.
func NewStrategyAwardModel(conn sqlx.SqlConn) StrategyAwardModel {
	return &customStrategyAwardModel{
		defaultStrategyAwardModel: newStrategyAwardModel(conn),
	}
}
