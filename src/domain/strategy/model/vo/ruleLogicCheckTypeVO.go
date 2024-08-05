package vo

type RuleLogicCheckTypeVO struct {
	Code string
	Info string
}

var (
	ALLOW = RuleLogicCheckTypeVO{
		Code: "0000",
		Info: "放行；执行后续的流程，不受规则引擎影响",
	}
	TAKE_OVER = RuleLogicCheckTypeVO{
		Code: "0001",
		Info: "接管；后续的流程，受规则引擎执行结果影响",
	}
)
