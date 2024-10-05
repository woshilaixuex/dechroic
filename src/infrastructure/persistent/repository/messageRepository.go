package infra_repository

import (
	"context"
	"strconv"
	"time"

	mess_vo "github.com/delyr1c/dechoric/src/domain/message/model/vo"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/aiUsage"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/redis"
	"github.com/delyr1c/dechoric/src/types/common"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-04 21:28
 */
type MessageRepository struct {
	RedisService redis.RedisService
	AiUsageModel aiUsage.AiUsageModel
}

func NewMessageRepository(sqlConn sqlx.SqlConn, redis redis.RedisService) *MessageRepository {
	return &MessageRepository{
		RedisService: redis,
		AiUsageModel: aiUsage.NewAiUsageModel(sqlConn),
	}
}
func (s *MessageRepository) QueryAIInfoByUserId(ctx context.Context, userId string) (*mess_vo.AIInfosVo, error) {
	aiUsageRecords, err := s.AiUsageModel.FindByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	var aiInfos []mess_vo.AIInfoVO
	for _, record := range aiUsageRecords {
		aiInfos = append(aiInfos, mess_vo.AIInfoVO{
			UsageId:    record.UsageId,
			UserId:     record.UserId,
			ModelId:    record.ModelId,
			ModelName:  record.ModelName,
			QueryCount: record.QueryCount,
		})
	}
	return &mess_vo.AIInfosVo{
		AIInfos: aiInfos,
	}, nil
}
func (s *MessageRepository) UseAIByUserId(ctx context.Context, userId string, modelId uint64) (*mess_vo.AIInfosUsVo, error) {
	// 使用 Redis 锁来防止并发问题
	lockKey := common.RedisKeys.AIUsageLockKey + userId + "_" + strconv.FormatUint(modelId, 10)
	unlock, err := s.RedisService.Lock(ctx, lockKey, 10*time.Second) // 锁定 10 秒
	if err != nil {
		return nil, err
	}
	defer unlock() // 确保解锁

	// 查询用户的 AI 使用记录
	aiUsageRecord, err := s.AiUsageModel.FindByUserIdAndModelId(ctx, userId, modelId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			return &mess_vo.AIInfosUsVo{DeOk: false}, nil // 找不到记录
		}
		return nil, err // 其他查询错误
	}

	// 检查 QueryCount
	if aiUsageRecord.QueryCount <= 0 {
		return &mess_vo.AIInfosUsVo{
			DeOk:        false,
			UserId:      userId,
			DeModelId:   aiUsageRecord.ModelId,
			DeModelName: aiUsageRecord.ModelName,
			QueryCount:  aiUsageRecord.QueryCount,
		}, nil // 使用次数为 0
	}

	// 减少使用次数
	aiUsageRecord.QueryCount--
	err = s.AiUsageModel.Update(ctx, aiUsageRecord)
	if err != nil {
		return nil, err // 更新失败
	}

	return &mess_vo.AIInfosUsVo{
		DeOk:        true,
		UserId:      userId,
		DeModelId:   aiUsageRecord.ModelId,
		DeModelName: aiUsageRecord.ModelName,
		QueryCount:  aiUsageRecord.QueryCount,
	}, nil // 返回使用成功的信息
}
