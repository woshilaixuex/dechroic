package model

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 这个model只是用来切断依赖环的
 * @Date: 2024-08-06 00:54
 */
type LogicModel struct {
	Code string
	Info string
}

var (
	RULE_WEIGHT = LogicModel{
		Code: "rule_weight",
		Info: "【抽奖前规则】根据抽奖权重返回可抽奖范围KEY",
	}
	RULE_BLACKLIST = LogicModel{
		Code: "rule_blacklist",
		Info: "【抽奖前规则】黑名单规则过滤，命中黑名单则直接返回",
	}
)
