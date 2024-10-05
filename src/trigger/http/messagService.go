package tr_http

import (
	"context"
	"net/http"

	"github.com/delyr1c/dechoric/src/api/dto"
	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	"github.com/gin-gonic/gin"
)

/*
 * @Author: deylr1c
 * @Email: linyugang7295@gmail.com
 * @Description:
 * @Date: 2024-10-05 22:09
 */
type MessageService struct {
	MessageRepo *infra_repository.MessageRepository
}

func NewMessageService(repository *infra_repository.MessageRepository) *MessageService {
	return &MessageService{
		MessageRepo: repository,
	}
}
func (mc *MessageService) GetUserAIInfo(c *gin.Context) {
	var req dto.AIInfoRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从请求中获取用户 ID
	userId := req.UserId
	ctx := context.Background()
	// 查询用户的 AI 信息
	aiInfo, err := mc.MessageRepo.QueryAIInfoByUserId(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 返回查询结果
	c.JSON(http.StatusOK, gin.H{"data": aiInfo})
}
