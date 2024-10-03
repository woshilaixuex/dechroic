package tree

import "github.com/delyr1c/dechoric/src/domain/strategy/model/entity"

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 策略树节点接口
 * @Date: 2024-09-21 21:21
 */
type ILogicTreeNode interface {
	GetKey() string
	Logic(userId string, strategyId int64, awardId int32) *entity.TreeActionEntity
}
