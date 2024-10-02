package vo

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 规则树节点
 * @Date: 2024-09-21 19:22
 */
type RuleTreeNodeVO struct {
	TreeId                  int32
	RuleKey                 string
	RukeDesc                string
	RuleValue               string
	RuleTreeNodeLineVOSlice []*RuleTreeNodeLineVO
}
