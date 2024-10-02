package vo

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 规则树对象
 * @Date: 2024-09-21 19:20
 */
type RuleTreeVO struct {
	TreeId           int32
	TreeName         string
	TreeDesc         string
	TreeRootRuleNode string
	TreeNodeMap      map[string]*RuleTreeNodeVO
}

func NewRuleTreeVO() *RuleTreeVO {
	return &RuleTreeVO{
		TreeNodeMap: make(map[string]*RuleTreeNodeVO),
	}
}
