package entity

import (
	"github.com/delyr1c/dechoric/src/types/common"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 策略奖品实体
 * @Date: 2024-06-13 22:35
 */

type StrategyAwardEntity struct {
	Id                int64           `db:"id" json:"id"`                                   // 自增ID
	AwardId           int64           `db:"award_id" json:"award_id"`                       // 抽奖策略ID
	AwardCount        int64           `db:"award_count" json:"award_count"`                 // 奖品库存总量
	AwardCountSurplus int64           `db:"award_count_surplus" json:"award_count_surplus"` // 奖品库存剩余
	AwardRate         common.BigFloat `db:"award_rate" json:"award_rate"`                   // 奖品中奖概率
}
