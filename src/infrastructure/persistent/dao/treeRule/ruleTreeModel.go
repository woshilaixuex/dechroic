package treeRule

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleTreeModel = (*customRuleTreeModel)(nil)

type (
	// RuleTreeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleTreeModel.
	RuleTreeModel interface {
		ruleTreeModel
		withSession(session sqlx.Session) RuleTreeModel
	}

	customRuleTreeModel struct {
		*defaultRuleTreeModel
	}
)

// NewRuleTreeModel returns a model for the database table.
func NewRuleTreeModel(conn sqlx.SqlConn) RuleTreeModel {
	return &customRuleTreeModel{
		defaultRuleTreeModel: newRuleTreeModel(conn),
	}
}

func (m *customRuleTreeModel) withSession(session sqlx.Session) RuleTreeModel {
	return NewRuleTreeModel(sqlx.NewSqlConnFromSession(session))
}
