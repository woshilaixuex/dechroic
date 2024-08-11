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
// 工厂需要的仓储
type StrategyRepositoryI interface {
	QueryStrategyAwardList(ctx context.Context, strategyId int64) ([]entity.StrategyAwardEntity, error)
	StoreStrategyAwardSearchRateTable(ctx context.Context, key string, rang int64, shuffleStrategyAwardSearchRateTable map[interface{}]interface{}) error
	GetAssembleRandomVal(ctx context.Context, key string, redisKey int64) (int64, error)
	QueryStrategyEntityByStrategyId(ctx context.Context, strategyId int64) (*entity.StrategyEntity, error)
	QueryStrategyRule(ctx context.Context, strategyId int64, roleModel string) (*entity.StrategyRuleEntity, error)
	QueryStrategyRuleValue(ctx context.Context, strategyId int64, awardId int32, roleModel string) (string, error)
	GetRateRange(ctx context.Context, key string) (int64, error)
	GetRateRangeBeta(ctx context.Context, strategyId int64) (int64, error)
}

type StrategyService struct {
	repo StrategyRepositoryI
}

func NewStrategyService(repo StrategyRepositoryI) *StrategyService {
	return &StrategyService{
		repo: repo,
	}
}

func (s *StrategyService) QueryStrategyAwardList(ctx context.Context, strategyId int64) ([]entity.StrategyAwardEntity, error) {
	return s.repo.QueryStrategyAwardList(ctx, strategyId)
}
func (s *StrategyService) StoreStrategyAwardSearchRateTable(ctx context.Context, key string, rang int64, shuffleStrategyAwardSearchRateTable map[interface{}]interface{}) error {
	return s.repo.StoreStrategyAwardSearchRateTable(ctx, key, rang, shuffleStrategyAwardSearchRateTable)
}
func (s *StrategyService) GetAssembleRandomVal(ctx context.Context, key string, redisKey int64) (int64, error) {
	return s.repo.GetAssembleRandomVal(ctx, key, redisKey)
}
func (s *StrategyService) QueryStrategyEntityByStrategyId(ctx context.Context, strategyId int64) (*entity.StrategyEntity, error) {
	return s.repo.QueryStrategyEntityByStrategyId(ctx, strategyId)
}
func (s *StrategyService) QueryStrategyRule(ctx context.Context, strategyId int64, roleModel string) (*entity.StrategyRuleEntity, error) {
	return s.repo.QueryStrategyRule(ctx, strategyId, roleModel)
}
func (s *StrategyService) QueryStrategyRuleValue(ctx context.Context, strategyId int64, awardId int32, roleModel string) (string, error) {
	return s.repo.QueryStrategyRuleValue(ctx, strategyId, awardId, roleModel)
}

func (s *StrategyService) GetRateRangeBeta(ctx context.Context, strategyId int64) (int64, error) {
	return s.repo.GetRateRangeBeta(ctx, strategyId)
}
func (s *StrategyService) GetRateRange(ctx context.Context, key string) (int64, error) {
	return s.repo.GetRateRange(ctx, key)
}
