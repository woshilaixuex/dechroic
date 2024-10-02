package vo

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 规则树枝（连线）
 * @Date: 2024-09-21 19:22
 */
type RuleTreeNodeLineVO struct {
	TreeId       int32  // 规则数id
	RuleNodeFrom string // 来自节点
	RuleNodeTo   string // 目标节点
	// 限定类型
	RuleLimitTypeVO      *RuleLimitTypeVO
	RuleLogicCheckTypeVO *RuleLogicCheckTypeVO
}
