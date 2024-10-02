package data

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-02 22:34
 */
type ChanVOLogicModel struct {
	Code string
	Info string
}
type StrategyAwardChanVO struct {
	AwardId        int32
	AwardRuleValue string
	ChanVOLogicModel
}

var (
	RULE_DEFAULT = ChanVOLogicModel{
		Code: "rule_default",
		Info: "默认抽奖",
	}
	RULE_BLACKLIST = ChanVOLogicModel{
		Code: "rule_blacklist",
		Info: "黑名单抽奖",
	}
	RULE_WEIGHT = ChanVOLogicModel{
		Code: "rule_weight",
		Info: "权重规则",
	}
)
