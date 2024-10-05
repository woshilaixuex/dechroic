// Code generated by goctl. DO NOT EDIT.

package aiUsage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	aiUsageFieldNames          = builder.RawFieldNames(&AiUsage{})
	aiUsageRows                = strings.Join(aiUsageFieldNames, ",")
	aiUsageRowsExpectAutoSet   = strings.Join(stringx.Remove(aiUsageFieldNames, "`usage_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	aiUsageRowsWithPlaceHolder = strings.Join(stringx.Remove(aiUsageFieldNames, "`usage_id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	aiUsageModel interface {
		Insert(ctx context.Context, data *AiUsage) (sql.Result, error)
		FindOne(ctx context.Context, usageId uint64) (*AiUsage, error)
		FindByUserId(ctx context.Context, userId string) ([]AiUsage, error) 
		FindByUserIdAndModelId(ctx context.Context, userId string,modelId uint64) (*AiUsage, error)
		Update(ctx context.Context, data *AiUsage) error
		Delete(ctx context.Context, usageId uint64) error
	}

	defaultAiUsageModel struct {
		conn  sqlx.SqlConn
		table string
	}

	AiUsage struct {
		UsageId    uint64 `db:"usage_id"`    // 自增ID
		UserId     string `db:"user_id"`     // 用户ID
		ModelId    uint64 `db:"model_id"`    // AI模型ID
		ModelName  string `db:"model_name"`  // AI模型名称
		QueryCount int64  `db:"query_count"` // 使用剩余次数
	}
)

func newAiUsageModel(conn sqlx.SqlConn) *defaultAiUsageModel {
	return &defaultAiUsageModel{
		conn:  conn,
		table: "`ai_usage`",
	}
}

func (m *defaultAiUsageModel) Delete(ctx context.Context, usageId uint64) error {
	query := fmt.Sprintf("delete from %s where `usage_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, usageId)
	return err
}

func (m *defaultAiUsageModel) FindOne(ctx context.Context, usageId uint64) (*AiUsage, error) {
	query := fmt.Sprintf("select %s from %s where `usage_id` = ? limit 1", aiUsageRows, m.table)
	var resp AiUsage
	err := m.conn.QueryRowCtx(ctx, &resp, query, usageId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultAiUsageModel) FindByUserId(ctx context.Context, userId string) ([]AiUsage, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ?", aiUsageRows, m.table)
	var resp []AiUsage
	err := m.conn.QueryRowsCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultAiUsageModel) FindByUserIdAndModelId(ctx context.Context, userId string,modelId uint64) (*AiUsage, error){
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `model_id` = ? limit 1", aiUsageRows, m.table)
	var resp AiUsage
	err := m.conn.QueryRowCtx(ctx, &resp, query, userId,modelId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
func (m *defaultAiUsageModel) Insert(ctx context.Context, data *AiUsage) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, aiUsageRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.UserId, data.ModelId, data.ModelName, data.QueryCount)
	return ret, err
}

func (m *defaultAiUsageModel) Update(ctx context.Context, data *AiUsage) error {
	query := fmt.Sprintf("update %s set %s where `usage_id` = ?", m.table, aiUsageRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.ModelId, data.ModelName, data.QueryCount, data.UsageId)
	return err
}

func (m *defaultAiUsageModel) tableName() string {
	return m.table
}