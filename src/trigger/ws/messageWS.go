package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	infra_repository "github.com/delyr1c/dechoric/src/infrastructure/persistent/repository"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/xinggaoya/qwen-sdk/qwen"
)

// WebSocket 消息结构体
type WebSocketMessage struct {
	UserID    string `json:"user_id"`
	AimodelID int    `json:"model_id"`
	Message   string `json:"message"`
}
type WebSocketResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	End       bool   `json:"end"` // 添加结束标识
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的连接
	},
}

type WebSocketHandler struct {
	qwenClient        *qwen.Chat
	messageRepository *infra_repository.MessageRepository
}

func NewWebSocketHandler(apiKey string, messageRepository *infra_repository.MessageRepository) *WebSocketHandler {
	return &WebSocketHandler{
		qwenClient:        qwen.NewWithDefaultChat(apiKey),
		messageRepository: messageRepository,
	}
}

// HandleWebSocket 处理 WebSocket 聊天
func (wsh *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}
	defer conn.Close()

	for {
		// 读取客户端消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("读取消息失败: %v", err)
			break
		}

		// 解析 JSON 消息
		var wsMessage WebSocketMessage
		err = json.Unmarshal(message, &wsMessage)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("无效的消息格式，必须为 JSON 格式"))
			continue
		}

		// 处理 AI 模型逻辑
		aiID := uint64(wsMessage.AimodelID)
		aiInfo, _ := wsh.messageRepository.QueryAIInfoByUserId(context.Background(), wsMessage.UserID, aiID)
		if aiInfo.QueryCount <= 0 {
			response := WebSocketResponse{
				Message:   string("次数用尽，请及时充值"),
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				End:       false, // 根据逻辑设置是否结束
			}
			responseJSON, _ := json.Marshal(response)
			conn.WriteMessage(websocket.TextMessage, responseJSON)
			continue
		}
		wsh.messageRepository.UseAIByUserId(context.Background(), wsMessage.UserID, aiID)

		// 使用 QWEN 获取 AI 回复
		aiReply, err := wsh.getAIReply(wsMessage.Message)
		if err != nil {
			response := WebSocketResponse{
				Message:   string("获取 AI 回复失败"),
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				End:       false, // 根据逻辑设置是否结束
			}
			responseJSON, _ := json.Marshal(response)
			log.Printf("获取 AI 回复失败: %v", err)
			conn.WriteMessage(websocket.TextMessage, responseJSON)
			continue
		}

		// 创建响应消息并设置 end 字段
		response := WebSocketResponse{
			Message:   aiReply,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			End:       false, // 根据逻辑设置是否结束
		}

		// 根据业务逻辑决定是否结束
		if shouldEndConversation() {
			response.End = true // 当满足某些条件时将其设置为 true
		}

		// 将 AI 回复发送回客户端
		responseJSON, _ := json.Marshal(response)
		err = conn.WriteMessage(websocket.TextMessage, responseJSON)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
			break
		}
	}
}

func shouldEndConversation() bool {
	return true
}

// getAIReply 获取 AI 对话的回复
func (wsh *WebSocketHandler) getAIReply(userMessage string) (string, error) {
	messages := []qwen.Messages{
		{Role: qwen.ChatUser, Content: userMessage},
	}
	resp, err := wsh.qwenClient.GetAIReply(messages)
	if err != nil {
		return "", err
	}
	return resp.Output.Text, nil
}
