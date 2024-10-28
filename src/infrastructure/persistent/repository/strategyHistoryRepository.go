package infra_repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	mess_vo "github.com/delyr1c/dechoric/src/domain/message/model/vo"
	StrategyEntity "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"
	"github.com/delyr1c/dechoric/src/domain/strategy/model/vo"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/aiUsage"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/lotteryHistory"
	"github.com/delyr1c/dechoric/src/infrastructure/persistent/dao/strategyAward"
	"github.com/delyr1c/dechoric/src/types/cerr"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-05 11:03
 */
// 查询奖品清单
func (s *StrategyRepository) QueryPrize(ctx context.Context, strategyId int64) (*vo.PrizesVO, error) {
	// 查询奖品列表
	awards, err := s.StrategyAwardModel.FindListByStrategyId(ctx, strategyId)
	if err != nil {
		return nil, err
	}

	// 创建返回的 PrizesVO 结构体
	prizesVO := &vo.PrizesVO{
		PrizeVOs: make([]vo.PrizeVO, 0, len(*awards)), // 根据奖品列表的长度预先分配空间
	}

	// 遍历查询到的奖品并填充 PrizeVO
	for _, award := range *awards {
		prize := vo.PrizeVO{
			SortId:     award.Sort,
			AwardId:    int32(award.AwardId),
			StrategyId: award.StrategyId,
			AwardTitle: award.AwardTitle,
		}
		prizesVO.PrizeVOs = append(prizesVO.PrizeVOs, prize)
	}
	return prizesVO, nil
}

// QueryStrategyAwardHistory retrieves the paginated lottery award history for a user and returns a slice of HistoryVO and the total number of pages.
func (s *StrategyRepository) QueryStrategyAwardHistory(ctx context.Context, userId string, page int) ([]mess_vo.HistoryVO, int, error) {
	const pageSize = 10
	totalCount, err := s.LotteryHistoryModel.QueryStrategyAwardHistoryCount(ctx, userId)
	if err != nil {
		return nil, 0, cerr.LogError(err)
	}
	histories, err := s.LotteryHistoryModel.QueryStrategyAwardHistory(ctx, userId, page, pageSize)
	if err != nil {
		return nil, 0, cerr.LogError(err)
	}

	totalPages := (totalCount + pageSize - 1) / pageSize

	// Map the result from lotteryHistory.LotteryHistory to HistoryVO
	result := make([]mess_vo.HistoryVO, len(histories))
	for i, history := range histories {
		result[i] = mess_vo.HistoryVO{
			Id:         int32(history.HistoryId), // Assuming history.Id exists and is of int32 type
			Desc:       history.Prize,            // Assuming history.Desc exists
			CreateTime: history.CreateTime,       // Assuming history.CreateTime exists
		}
	}
	return result, totalPages, nil
}

// 增加抽奖记录
func (s *StrategyRepository) AddStrategyAwardHistory(ctx context.Context, userId string, award *StrategyEntity.RaffleAwardEntity) (string, error) {
	// 创建 LotteryHistory 实例
	strategyAwards, err := s.StrategyAwardModel.FindByReq(ctx, &strategyAward.FindStrategyAwardReq{AwardId: &award.AwardId})
	if err != nil {
		return "", err
	}

	history := &lotteryHistory.LotteryHistory{
		UserId: userId,
		TreeDesc: sql.NullString{
			String: strategyAwards[0].AwardTitle,       // 传入的描述
			Valid:  strategyAwards[0].AwardTitle != "", // 仅在描述非空时有效
		},
		Prize:      strategyAwards[0].AwardTitle, // 或使用 award.AwardId 等适合的字段
		CreateTime: time.Now(),                   // 当前时间
	}

	// 将历史记录插入到 LotteryHistoryModel
	if _, err := s.LotteryHistoryModel.Insert(ctx, history); err != nil {
		return "", err // 返回错误信息
	}

	return strategyAwards[0].AwardTitle, nil // 返回 nil 表示成功
}

// 查询并执行点数扣除服务
func (s *StrategyRepository) DoUserCredit(ctx context.Context, userId string) (bool, error) {
	userInfo, err := s.UserModel.FindOne(ctx, userId)
	if err != nil {
		return false, err
	}
	if userInfo.Credit < 100 {
		return false, nil
	}
	userInfo.Credit -= 100
	s.UserModel.Update(ctx, userInfo)
	return true, nil
}

// 执行对应award
func (s *StrategyRepository) PerformAward(ctx context.Context, userId string, awardId int32) (bool, error) {
	// 根据 awardId 执行不同的操作
	switch awardId {
	case 101: // 随机积分
		rand.Seed(time.Now().UnixNano())  // 设定随机数种子
		randomCredit := rand.Int63n(1001) // 生成 0 到 1000 的随机积分
		return s.UpdateUserCredit(ctx, userId, randomCredit)
	case 102: // 全部模型5次使用
		return s.UpdateModelUsage(ctx, userId, 0, 5) // 模型ID为 0 表示全部模型
	case 103: // 全部模型10次使用
		return s.UpdateModelUsage(ctx, userId, 0, 10)
	case 104: // 全部模型20次使用
		return s.UpdateModelUsage(ctx, userId, 0, 20)
	case 105: // 增加10次GPT-4对话模型
		return s.UpdateModelUsage(ctx, userId, 1, 10) // 模型ID为 1 表示 GPT-4
	case 106: // 增加10次DALL-E-2画图模型
		return s.UpdateModelUsage(ctx, userId, 2, 10) // 模型ID为 2 表示 DALL-E-2
	case 107: // 增加10次DALL-E-3画图模型
		return s.UpdateModelUsage(ctx, userId, 3, 10) // 模型ID为 3 表示 DALL-E-3
	case 108: // 增加100次使用
		return s.UpdateModelUsage(ctx, userId, 0, 100) // 增加100次全部模型的使用
	case 109: // 全部模型1000次使用
		return s.UpdateModelUsage(ctx, userId, 0, 1000)
	default:
		return false, fmt.Errorf("未知的奖品 ID: %d", awardId)
	}
}

// 更新用户积分
func (s *StrategyRepository) UpdateUserCredit(ctx context.Context, userId string, creditToAdd int64) (bool, error) {
	// 查找用户信息
	userInfo, err := s.UserModel.FindOne(ctx, userId)
	if err != nil {
		return false, err
	}

	// 更新用户积分
	userInfo.Credit += creditToAdd

	// 保存更新后的用户信息
	if err := s.UserModel.Update(ctx, userInfo); err != nil {
		return false, err
	}

	return true, nil
}

// 更新模型使用次数
func (s *StrategyRepository) UpdateModelUsage(ctx context.Context, userId string, modelId uint64, queryCountToAdd int64) (bool, error) {
	// 如果 modelId 为 0，表示更新所有模型的使用次数
	if modelId == 0 {
		models, err := s.AiUsageModel.FindByUserId(ctx, userId)
		if err != nil {
			return false, err
		}
		for _, model := range models {
			model.QueryCount += queryCountToAdd
			if err := s.AiUsageModel.Update(ctx, &model); err != nil {
				return false, err
			}
		}
		return true, nil
	}

	// 否则更新指定模型的使用次数
	aiUsageData, err := s.AiUsageModel.FindByUserIdAndModelId(ctx, userId, modelId)
	if err != nil {
		if err == sqlx.ErrNotFound {
			// 如果没有找到记录，则插入新的记录
			newUsage := aiUsage.AiUsage{
				UserId:     userId,
				ModelId:    modelId,
				ModelName:  GetModelNameById(modelId), // 根据 modelId 获取模型名称
				QueryCount: queryCountToAdd,           // 初始化使用次数
			}
			_, err := s.AiUsageModel.Insert(ctx, &newUsage)
			return err == nil, err
		}
		return false, err
	}
	// 更新模型使用次数
	aiUsageData.QueryCount += queryCountToAdd
	if err := s.AiUsageModel.Update(ctx, aiUsageData); err != nil {
		return false, err
	}

	return true, nil
}

// 根据 modelId 获取模型名称
func GetModelNameById(modelId uint64) string {
	switch modelId {
	case 1:
		return "gpt-4"
	case 2:
		return "dall-e-2"
	case 3:
		return "dall-e-3"
	default:
		return "未知模型"
	}
}
