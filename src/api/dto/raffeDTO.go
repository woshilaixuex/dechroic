package dto

import mess_vo "github.com/delyr1c/dechoric/src/domain/message/model/vo"

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 抽奖相关dto
 * @Date: 2024-10-04 17:30
 */
type RaffleAwardListRequestDTO struct {
	StrategyId int64
}
type RaffleAwardListResponceDTO struct {
	AwardId       int32
	AwardTitle    string
	AwardSubtitle string
	Sort          int32
}
type RaffleRequestDTO struct {
	StrategyId int64
}
type RaffleResponceDTO struct {
	AwardId    int32
	AwardIndex int32
}
type PrizeRequsetDTO struct {
	StrategyId int64 `json:"strategy_id"`
}
type PrizeResponseDTO struct {
	SortId     int64  `json:"sort_id"`     // 奖品的排序ID
	AwardId    int32  `json:"award_id"`    //策略
	StrategyId int64  `json:"strategy_id"` // 策略ID
	AwardTitle string `json:"award_title"` // 奖品标题
}

type PrizesResponseDTO struct {
	PrizeList []PrizeResponseDTO `json:"prize_list"` // 奖品列表
}

type AwardPrizeRequsetDTO struct {
	UserId     string `json:"user_id"`
	StrategyId int64  `json:"strategy_id"`
}
type AwardPrizeResponseDTO struct {
	AwardId    int32  `json:"award_id"`
	AwardTitle string `json:"award_title"` // 奖品标题
}
type GetHistoryRequsetDTO struct {
	UserId string `json:"user_id"`
	Page   int32  `json:"page"`
}
type GetHistoryResponseDTO struct {
	Page       int32               `json:"page"`
	TotalPages int32               `json:"total_pages"`
	Historys   []mess_vo.HistoryVO `json:"historys"`
}
