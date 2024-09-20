package model

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: (过滤器的逻辑叙述模型)这个model只是用来切断依赖环的
 * @Date: 2024-08-06 00:54
 */
type LogicModel struct {
	Code string
	Info string
	Type string
}

var (
	LogicModelMap = make(map[string]LogicModel, 0)
	RULE_WEIGHT   = LogicModel{
		Code: "rule_weight",
		Info: "【抽奖前规则】根据抽奖权重返回可抽奖范围KEY",
		Type: "before",
	}
	RULE_BLACKLIST = LogicModel{
		Code: "rule_blacklist",
		Info: "【抽奖前规则】黑名单规则过滤，命中黑名单则直接返回",
		Type: "before",
	}
	RULE_LOCK = LogicModel{
		Code: "rule_lock",
		Info: "【抽奖中规则】抽N次后，对应可解锁抽奖",
		Type: "center",
	}
	RULE_LUCK_AWARD = LogicModel{
		Code: "rule_luck_award",
		Info: "【抽奖后规则】幸运奖兜底",
		Type: "after",
	}
)

func init() {
	registerLogicModels(
		RULE_WEIGHT,
		RULE_BLACKLIST,
		RULE_LOCK,
		RULE_LUCK_AWARD,
	)
}

func registerLogicModels(models ...LogicModel) {
	for _, model := range models {
		LogicModelMap[model.Code] = model
	}
}
func IsBefore(code string) bool {
	return LogicModelMap[code].Type == "before"
}
func IsCenter(code string) bool {
	return LogicModelMap[code].Type == "center"
}
func IsAfter(code string) bool {
	return LogicModelMap[code].Type == "after"
}
