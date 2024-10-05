package dto

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
	StrategyId int64
}
