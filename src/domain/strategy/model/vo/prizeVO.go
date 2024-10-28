package vo

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-06 00:11
 */

type PrizeVO struct {
	SortId     int64
	AwardId    int32
	StrategyId int64
	AwardTitle string
}
type PrizesVO struct {
	PrizeVOs []PrizeVO
}

type RaffeAwardVO struct {
	AwardDesc string
}
