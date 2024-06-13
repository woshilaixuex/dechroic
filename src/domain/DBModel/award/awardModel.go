package award

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AwardModel = (*customAwardModel)(nil)

type (
	// AwardModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAwardModel.
	AwardModel interface {
		awardModel
	}

	customAwardModel struct {
		*defaultAwardModel
	}
)

// NewAwardModel returns a model for the database table.
func NewAwardModel(conn sqlx.SqlConn) AwardModel {
	return &customAwardModel{
		defaultAwardModel: newAwardModel(conn),
	}
}
