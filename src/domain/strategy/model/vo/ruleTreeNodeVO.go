package vo

import "fmt"

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 规则树节点
 * @Date: 2024-09-21 19:22
 */
type RuleTreeNodeVO struct {
	TreeId                  string
	RuleKey                 string
	RukeDesc                string
	RuleValue               string
	RuleTreeNodeLineVOSlice []*RuleTreeNodeLineVO
}

func (vo *RuleTreeNodeVO) RuleTreeNodeVoInfo() {
	fmt.Printf("RuleTreeNodeVO Info:\n")
	fmt.Printf("TreeId: %s\n", vo.TreeId)
	fmt.Printf("RuleKey: %s\n", vo.RuleKey)
	fmt.Printf("RukeDesc: %s\n", vo.RukeDesc) // 注意这里的拼写错误，应该是 RuleDesc
	fmt.Printf("RuleValue: %s\n", vo.RuleValue)

	fmt.Println("RuleTreeNodeLineVOSlice:")
	for _, line := range vo.RuleTreeNodeLineVOSlice {
		fmt.Printf("\tRuleNodeFrom: %s, RuleNodeTo: %s\n", line.RuleNodeFrom, line.RuleNodeTo)
	}
}
