package vo

import "fmt"

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-09-21 19:26
 */
type RuleLimitTypeVO struct {
	Code int
	Info string
}

var (
	EQUAL = RuleLimitTypeVO{
		Code: 1,
		Info: "等于",
	}
	GT = RuleLimitTypeVO{
		Code: 2,
		Info: "大于",
	}
	LT = RuleLimitTypeVO{
		Code: 3,
		Info: "小于",
	}
	GE = RuleLimitTypeVO{
		Code: 4,
		Info: "大于&等于",
	}
	LE = RuleLimitTypeVO{
		Code: 5,
		Info: "小于&等于",
	}
	ENUM = RuleLimitTypeVO{
		Code: 6,
		Info: "枚举",
	}
)

// 创建映射表
var ruleLimitTypeMap = map[string]RuleLimitTypeVO{
	"EQUAL": EQUAL,
	"GT":    GT,
	"LT":    LT,
	"GE":    GE,
	"LE":    LE,
	"ENUM":  ENUM,
}

// 根据字符串查找匹配的 RuleLimitTypeVO
func GetRuleLimitTypeVOByStr(typeStr string) (*RuleLimitTypeVO, error) {
	if val, ok := ruleLimitTypeMap[typeStr]; ok {
		return &val, nil
	}
	return &RuleLimitTypeVO{}, fmt.Errorf("no matching RuleLimitTypeVO for string: %s", typeStr)
}
