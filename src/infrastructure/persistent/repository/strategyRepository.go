package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	strategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	strategyRepository "github.com/delyr1c/dechoric/src/domain/strategy/repository"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategy"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyRule"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/treeRule"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/treeRuleNode"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/treeRuleNodeLine"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽奖策略仓储具体实现
 * @Date: 2024-06-16 23:42
 */

// 确保 StrategyRepository 实现 IStrategyRepository 接口
var _ strategyRepository.StrategyRepositoryI = (*StrategyRepository)(nil)

type StrategyRepository struct {
	RedisService          redis.RedisService
	StrategyAwardModel    strategyAward.StrategyAwardModel
	StrategyModel         strategy.StrategyModel
	StrategyRuleModel     strategyRule.StrategyRuleModel
	TreeRuleModel         treeRule.RuleTreeModel
	TreeRuleNodeModel     treeRuleNode.RuleTreeNodeModel
	TreeRuleNodeLineModel treeRuleNodeLine.RuleTreeNodeLineModel
}

// 返回strategyAwardEntityList（redis->db)
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
	if err := copier.Copy(&strategyAwardEntityList, queryStrategyAwardEntityList); err != nil {
		return nil, err
	}

	if err := s.RedisService.SetArray(ctx, cacheKey, strategyAwardEntityList); err != nil {
		fmt.Println("Failed to set data to Redis:", err)
	}

	return strategyAwardEntityList, nil
}
func (s *StrategyRepository) StoreStrategyAwardSearchRateTable(ctx context.Context, key string, rang int64, shuffleStrategyAwardSearchRateTable map[interface{}]interface{}) error {
	err := s.RedisService.SetValue(ctx, common.RedisKeys.StrategyRateRangeKey+key, rang)
	if err != nil {
		logx.Info(err.Error())
		return err
	}
	// 把SetToMap改成SetToMapString是序列化方法
	err = s.RedisService.SetToMap(ctx, common.RedisKeys.StrategyRateTableKey+key, shuffleStrategyAwardSearchRateTable)
	if err != nil {
		logx.Info(err.Error())
		return err
	}
	return nil
}
func (s *StrategyRepository) GetAssembleRandomVal(ctx context.Context, key string, redisKey int64) (int64, error) {
	// 这是序列化方法
	// string
	// randomVal, err := s.RedisService.GetFromMapString(ctx, common.RedisKeys.StrategyRateTableKey+strconv.FormatInt(strategyId, 10), strconv.FormatInt(redisKey, 10))
	// if err != nil {
	// 	return -1, err
	// }
	// if randomVal == nil {
	// 	return -1, fmt.Errorf("key %s not found in Redis", common.RedisKeys.StrategyRateTableKey+strconv.FormatInt(strategyId, 10))
	// }

	// // 尝试将 randomVal 转换为 float64
	// val, ok := randomVal.(float64)
	// if !ok {
	// 	return -1, fmt.Errorf("unexpected value type for randomVal: %T", randomVal)
	// }
	// return int64(val), nil
	val := s.RedisService.GetFromMap(ctx, common.RedisKeys.StrategyRateTableKey+key, strconv.FormatInt(redisKey, 10))
	return strconv.ParseInt(val.Val(), 10, 64)
}

// 查询策略权重
func (s *StrategyRepository) QueryStrategyEntityByStrategyId(ctx context.Context, strategyId int64) (*strategyEntity.StrategyEntity, error) {
	cacheKey := common.RedisKeys.StrategyKey + strconv.FormatInt(strategyId, 10)
	val, err := s.RedisService.GetValue(ctx, cacheKey)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	entity := new(strategyEntity.StrategyEntity)
	if val != "" && val != nil {
		if err := json.Unmarshal([]byte(val.(string)), entity); err != nil {
			return nil, cerr.LogError(errors.New("failed to unmarshal value to StrategyEntity"))
		}
		return entity, nil
	}
	strategyModel, err := s.StrategyModel.FindByStrategyId(ctx, strategyId)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	err = copier.Copy(entity, strategyModel)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	// redis存储
	err = s.RedisService.SetValue(ctx, cacheKey, entity)
	if err != nil {
		return nil, cerr.LogError(err)
	}
	return entity, nil
}

// 查询规则
func (s *StrategyRepository) QueryStrategyRule(ctx context.Context, strategyId int64, roleModel string) (*strategyEntity.StrategyRuleEntity, error) {
	strategyRuleReq := &strategyRule.FindStrategyRuleReq{
		StrategyId: &strategyId,
		RuleModel:  &roleModel,
	}
	strategyRules, err := s.StrategyRuleModel.FindByReq(ctx, strategyRuleReq)
	if err != nil {
		return nil, err
	}
	strategyRule := new(strategyEntity.StrategyRuleEntity)
	copier.Copy(strategyRule, strategyRules[0])
	return strategyRule, nil
}
func (s *StrategyRepository) GetRateRangeBeta(ctx context.Context, strategyId int64) (int64, error) {
	key := strconv.FormatInt(strategyId, 10)
	return s.GetRateRange(ctx, key)
}
func (s *StrategyRepository) GetRateRange(ctx context.Context, key string) (int64, error) {
	val, err := s.RedisService.GetValue(ctx, common.RedisKeys.StrategyRateRangeKey+key)
	if err != nil {
		return 0, cerr.LogError(err)
	}
	rateRangStr, ok := val.(string)
	if !ok {
		err := errors.New("failed to convert rateRangStr to string")
		return 0, cerr.LogError(err)
	}
	rateRang, err := strconv.ParseInt(rateRangStr, 10, 64)
	if err != nil {
		err := errors.New("failed to convert rateRangStr to int64")
		return 0, cerr.LogError(err)
	}
	return rateRang, nil
}
