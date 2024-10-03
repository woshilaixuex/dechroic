package vo

import "fmt"

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 规则逻辑引擎枚举信息
 * @Date: 2024-08-12 18:32
 */
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

// 创建映射表
var ruleLogicCheckTypeMap = map[string]RuleLogicCheckTypeVO{
	"ALLOW":     ALLOW,
	"TAKE_OVER": TAKE_OVER,
}

// 根据字符串查找匹配的 RuleLogicCheckTypeVO
func GetRuleLogicCheckTypeVOByStr(codeStr string) (*RuleLogicCheckTypeVO, error) {
	// 从 map 中查找是否有匹配的值
	if val, ok := ruleLogicCheckTypeMap[codeStr]; ok {
		return &val, nil
	}
	return &RuleLogicCheckTypeVO{}, fmt.Errorf("no matching RuleLogicCheckTypeVO for code: %s", codeStr)
}
