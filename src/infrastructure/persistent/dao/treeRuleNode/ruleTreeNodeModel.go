package treeRuleNode

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RuleTreeNodeModel = (*customRuleTreeNodeModel)(nil)

type (
	// RuleTreeNodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRuleTreeNodeModel.
	RuleTreeNodeModel interface {
		ruleTreeNodeModel
		withSession(session sqlx.Session) RuleTreeNodeModel
	}

	customRuleTreeNodeModel struct {
		*defaultRuleTreeNodeModel
	}
)

// NewRuleTreeNodeModel returns a model for the database table.
func NewRuleTreeNodeModel(conn sqlx.SqlConn) RuleTreeNodeModel {
	return &customRuleTreeNodeModel{
		defaultRuleTreeNodeModel: newRuleTreeNodeModel(conn),
	}
}

func (m *customRuleTreeNodeModel) withSession(session sqlx.Session) RuleTreeNodeModel {
	return NewRuleTreeNodeModel(sqlx.NewSqlConnFromSession(session))
}
