package dto

import mess_vo "github.com/delyr1c/dechoric/src/domain/message/model/vo"

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description: 消息相关dto
 * @Date: 2024-10-04 19:40
 */
type AIInfoRequestDTO struct {
	UserId string `json:"user_id"` // 用户ID
}
type AIInfoResponseDTO struct {
	AiModels []mess_vo.AIInfoVO `json:"ai_infos"`
}
