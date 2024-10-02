package treeRuleNodeLine

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleTreeNodeLineModel = (*customRuleTreeNodeLineModel)(nil)

type (
	// RuleTreeNodeLineModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleTreeNodeLineModel.
	RuleTreeNodeLineModel interface {
		ruleTreeNodeLineModel
		withSession(session sqlx.Session) RuleTreeNodeLineModel
	}

	customRuleTreeNodeLineModel struct {
		*defaultRuleTreeNodeLineModel
	}
)

// NewRuleTreeNodeLineModel returns a model for the database table.
func NewRuleTreeNodeLineModel(conn sqlx.SqlConn) RuleTreeNodeLineModel {
	return &customRuleTreeNodeLineModel{
		defaultRuleTreeNodeLineModel: newRuleTreeNodeLineModel(conn),
	}
}

func (m *customRuleTreeNodeLineModel) withSession(session sqlx.Session) RuleTreeNodeLineModel {
	return NewRuleTreeNodeLineModel(sqlx.NewSqlConnFromSession(session))
}
