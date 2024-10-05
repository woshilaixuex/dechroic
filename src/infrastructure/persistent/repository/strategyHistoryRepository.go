package infra_repository

import (
	"context"

	mess_vo "github.com/delyr1c/dechoric/src/domain/message/model/vo"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-05 11:03
 */
func (s *StrategyRepository) QueryStrategyAwardHistory() {

}
func (s *StrategyRepository) QueryPrize(ctx context.Context, strategyId int64) (*mess_vo.PrizesVO, error) {
	// 查询奖品列表
	awards, err := s.StrategyAwardModel.FindListByStrategyId(ctx, strategyId)
	if err != nil {
		return nil, err
	}

	// 创建返回的 PrizesVO 结构体
	prizesVO := &mess_vo.PrizesVO{
		PrizeVOs: make([]mess_vo.PrizeVO, 0, len(*awards)), // 根据奖品列表的长度预先分配空间
	}

	// 遍历查询到的奖品并填充 PrizeVO
	for _, award := range *awards {
		prize := mess_vo.PrizeVO{
			Id:         award.Id,
			StrategyId: award.StrategyId,
			AwardId:    award.AwardId,
			AwardTitle: award.AwardTitle,
			Sort:       award.Sort,
		}
		prizesVO.PrizeVOs = append(prizesVO.PrizeVOs, prize)
	}
	return prizesVO, nil
}
