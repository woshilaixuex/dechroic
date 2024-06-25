package repository

import (
	"context"

	"github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽奖策略仓储实现方法（接口）
 * @Date: 2024-06-16 11:29
 */

type StrategyRepository1 interface {
	QueryStrategyAwardList(ctx context.Context, strategyId int64) ([]entity.StrategyAwardEntity, error)
}

type StrategyService struct {
	repo StrategyRepository1
}

func NewStrategyService(repo StrategyRepository1) *StrategyService {
	return &StrategyService{
		repo: repo,
	}
}

func (s *StrategyService) QueryStrategyAwardList(ctx context.Context, strategyId int64) ([]entity.StrategyAwardEntity, error) {
	return s.repo.QueryStrategyAwardList(ctx, strategyId)
}
