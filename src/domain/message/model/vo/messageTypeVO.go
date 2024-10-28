package mess_vo

import "time"

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-05 20:26
 */
type AIInfoVO struct {
	UsageId    uint64 `json:"usage_id"`    // 自增ID
	UserId     string `json:"user_id"`     // 用户ID
	ModelId    uint64 `json:"model_id"`    // AI模型ID
	ModelName  string `json:"model_name"`  // AI模型名称
	QueryCount int64  `json:"query_count"` // 使用剩余次数
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
type HistoryVO struct {
	Id         int32
	Desc       string
	CreateTime time.Time
}
