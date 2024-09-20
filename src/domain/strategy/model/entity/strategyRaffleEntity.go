package entity

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: service层中raffle策略所需实体
 * @Date: 2024-08-05 19:04
 */

type RaffleEntityinterface interface {
}
type RaffleBeforeEntity struct {
	RaffleEntityinterface
	StrategyId         int64
	RuleWeightValueKey string
	AwardId            int32
}
type RaffleCenterEntity struct {
	RaffleEntityinterface
}
type RaffleAfterEntity struct {
	RaffleEntityinterface
}
type RaffleActionEntityInterface interface {
	GetCode() string
	GetInfo() string
}
type RaffleActionEntity[T RaffleEntityinterface] struct {
	RaffleActionEntityInterface
	Code      string
	Info      string
	RuleModel string `db:"rule_model"` // 抽奖规则类型【rule_random - 随机值计算、rule_lock - 抽奖几次后解锁、rule_luck_award - 幸运奖(兜底奖品)】
	Data      T
}

func (entity *RaffleActionEntity[RaffleEntityinterface]) GetCode() string {
	return entity.Code
}
func (entity *RaffleActionEntity[RaffleEntityinterface]) GetInfo() string {
	return entity.Info
}

type RaffleAwardEntity struct {
	Id          int64  `db:"id"`           // 自增ID
	AwardId     int64  `db:"award_id"`     // 抽奖奖品ID - 内部流转使用
	AwardKey    string `db:"award_key"`    // 奖品对接标识 - 每一个都是一个对应的发奖策略
	AwardConfig string `db:"award_config"` // 奖品配置信息
	AwardDesc   string `db:"award_desc"`   // 奖品内容描述
}
type RaffleFactorEntity struct {
	UserId     string
	StrategyId int64
	AwardId    int32
}

// 过滤判断需要的参数实体
type RuleMatterEntity struct {
	UserId     string
	StrategyId int64  `db:"strategy_id"`              // 抽奖策略ID
	AwardId    int32  `db:"award_id" json:"award_id"` // 抽奖奖品ID【规则类型为策略，则不需要奖品ID】
	RuleModel  string `db:"rule_model"`               // 抽奖规则类型【rule_random - 随机值计算、rule_lock - 抽奖几次后解锁、rule_luck_award - 幸运奖(兜底奖品)】
}
