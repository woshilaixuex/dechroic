package lotteryHistory

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LotteryHistoryModel = (*customLotteryHistoryModel)(nil)

type (
	// LotteryHistoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLotteryHistoryModel.
	LotteryHistoryModel interface {
		lotteryHistoryModel
		withSession(session sqlx.Session) LotteryHistoryModel
	}

	customLotteryHistoryModel struct {
		*defaultLotteryHistoryModel
	}
)

// NewLotteryHistoryModel returns a model for the database table.
func NewLotteryHistoryModel(conn sqlx.SqlConn) LotteryHistoryModel {
	return &customLotteryHistoryModel{
		defaultLotteryHistoryModel: newLotteryHistoryModel(conn),
	}
}

func (m *customLotteryHistoryModel) withSession(session sqlx.Session) LotteryHistoryModel {
	return NewLotteryHistoryModel(sqlx.NewSqlConnFromSession(session))
}
