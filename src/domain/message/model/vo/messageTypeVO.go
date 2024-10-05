package mess_vo

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-05 20:26
 */
type AIInfoVO struct {
	UsageId    uint64 // 自增ID
	UserId     string // 用户ID
	ModelId    uint64 // AI模型ID
	ModelName  string // AI模型名称
	QueryCount int64  // 使用剩余次数
}
type AIInfosVo struct {
	AIInfos []AIInfoVO
}

type AIInfosUsVo struct {
	DeOk        bool   //是否成功被使用
	UserId      string // 用户ID
	DeModelId   uint64 //被使用模型id
	DeModelName string // AI模型名称
	QueryCount  int64  // 使用剩余次数
}
type PrizeVO struct {
	Id         int64
	StrategyId int64
	AwardId    int64
	AwardTitle string
	Sort       int64
}
type PrizesVO struct {
	PrizeVOs []PrizeVO
}
