package repository

import (
	"context"
	"fmt"
	"strconv"

	strategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	strategyRepository "github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/jinzhu/copier"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽奖策略仓储具体实现
 * @Date: 2024-06-16 23:42
 */

// 确保 StrategyRepository 实现 StrategyRepository1 接口
var _ strategyRepository.StrategyRepository1 = (*StrategyRepository)(nil)

type StrategyRepository struct {
	RedisService       redis.RedisService
	StrategyAwardModel strategyAward.StrategyAwardModel
}

func (s *StrategyRepository) QueryStrategyAwardList(ctx context.Context, strategyId int64) ([]strategyEntity.StrategyAwardEntity, error) {
	cacheKey := common.RedisKeys.StrategyAwardListKey + strconv.FormatInt(strategyId, 10)

	var strategyAwardEntityList []strategyEntity.StrategyAwardEntity

	if err := s.RedisService.GetArray(ctx, cacheKey, &strategyAwardEntityList); err != nil {
		return nil, err
	}

	if len(strategyAwardEntityList) > 0 {
		return strategyAwardEntityList, nil
	}

	queryStrategyAwardEntityList, err := s.StrategyAwardModel.FindListByStrategyId(ctx, strategyId)
	if err != nil {
		return nil, err
	}
	copier.Copy(&strategyAwardEntityList, queryStrategyAwardEntityList)

	if err := s.RedisService.SetArray(ctx, cacheKey, strategyAwardEntityList); err != nil {
		fmt.Println("Failed to set data to Redis:", err)
	}

	return strategyAwardEntityList, nil
}
