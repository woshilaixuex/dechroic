package aiUsage

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiUsageModel = (*customAiUsageModel)(nil)

type (
	// AiUsageModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiUsageModel.
	AiUsageModel interface {
		aiUsageModel
		withSession(session sqlx.Session) AiUsageModel
	}

	customAiUsageModel struct {
		*defaultAiUsageModel
	}
)

// NewAiUsageModel returns a model for the database table.
func NewAiUsageModel(conn sqlx.SqlConn) AiUsageModel {
	return &customAiUsageModel{
		defaultAiUsageModel: newAiUsageModel(conn),
	}
}

func (m *customAiUsageModel) withSession(session sqlx.Session) AiUsageModel {
	return NewAiUsageModel(sqlx.NewSqlConnFromSession(session))
}
