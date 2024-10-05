package aiModel

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ AiModelsModel = (*customAiModelsModel)(nil)

type (
	// AiModelsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAiModelsModel.
	AiModelsModel interface {
		aiModelsModel
		withSession(session sqlx.Session) AiModelsModel
	}

	customAiModelsModel struct {
		*defaultAiModelsModel
	}
)

// NewAiModelsModel returns a model for the database table.
func NewAiModelsModel(conn sqlx.SqlConn) AiModelsModel {
	return &customAiModelsModel{
		defaultAiModelsModel: newAiModelsModel(conn),
	}
}

func (m *customAiModelsModel) withSession(session sqlx.Session) AiModelsModel {
	return NewAiModelsModel(sqlx.NewSqlConnFromSession(session))
}
