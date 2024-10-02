// Code generated by goctl. DO NOT EDIT.

package treeRuleNode

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	ruleTreeNodeFieldNames          = builder.RawFieldNames(&RuleTreeNode{})
	ruleTreeNodeRows                = strings.Join(ruleTreeNodeFieldNames, ",")
	ruleTreeNodeRowsExpectAutoSet   = strings.Join(stringx.Remove(ruleTreeNodeFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	ruleTreeNodeRowsWithPlaceHolder = strings.Join(stringx.Remove(ruleTreeNodeFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	ruleTreeNodeModel interface {
		Insert(ctx context.Context, data *RuleTreeNode) (sql.Result, error)
		FindOne(ctx context.Context, id uint64) (*RuleTreeNode, error)
		FindRuleTreeNodeListByTreeId(ctx context.Context,treeId string) ([]*RuleTreeNode, error)
		Update(ctx context.Context, data *RuleTreeNode) error
		Delete(ctx context.Context, id uint64) error
	}

	defaultRuleTreeNodeModel struct {
		conn  sqlx.SqlConn
		table string
	}

	RuleTreeNode struct {
		Id         uint64         `db:"id"`          // 自增ID
		TreeId     string         `db:"tree_id"`     // 规则树ID
		RuleKey    string         `db:"rule_key"`    // 规则Key
		RuleDesc   string         `db:"rule_desc"`   // 规则描述
		RuleValue  sql.NullString `db:"rule_value"`  // 规则比值
		CreateTime time.Time      `db:"create_time"` // 创建时间
		UpdateTime time.Time      `db:"update_time"` // 更新时间
	}
)

func newRuleTreeNodeModel(conn sqlx.SqlConn) *defaultRuleTreeNodeModel {
	return &defaultRuleTreeNodeModel{
		conn:  conn,
		table: "`rule_tree_node`",
	}
}

func (m *defaultRuleTreeNodeModel) Delete(ctx context.Context, id uint64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultRuleTreeNodeModel) FindOne(ctx context.Context, id uint64) (*RuleTreeNode, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", ruleTreeNodeRows, m.table)
	var resp RuleTreeNode
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultRuleTreeNodeModel) FindRuleTreeNodeListByTreeId(ctx context.Context,treeId string) ([]*RuleTreeNode, error){
	query := fmt.Sprintf("select %s from %s where `tree_id` = ?", ruleTreeNodeRows, m.table)
	var resp []*RuleTreeNode
	err := m.conn.QueryRowsCtx(ctx, &resp, query, treeId)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultRuleTreeNodeModel) Insert(ctx context.Context, data *RuleTreeNode) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, ruleTreeNodeRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.TreeId, data.RuleKey, data.RuleDesc, data.RuleValue)
	return ret, err
}

func (m *defaultRuleTreeNodeModel) Update(ctx context.Context, data *RuleTreeNode) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, ruleTreeNodeRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.TreeId, data.RuleKey, data.RuleDesc, data.RuleValue, data.Id)
	return err
}

func (m *defaultRuleTreeNodeModel) tableName() string {
	return m.table
}
