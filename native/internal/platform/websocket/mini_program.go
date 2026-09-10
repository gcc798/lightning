package websocket

import (
	"bytes"
	"encoding/json"

	logging "github.com/gcc798/microservice-kit/internal/logger"
	"go.uber.org/zap"
)

// MiniProgramMessage 表示小程序发送到服务端的 WebSocket 消息。
type MiniProgramMessage struct {
	Type string `json:"type"`
}

// NewMiniProgramProtocolHandler 创建门禁小程序 WebSocket 协议处理器。
func NewMiniProgramProtocolHandler(logger logging.Logger) TextMessageHandler {
	return func(client *Client, message []byte) bool {
		var req MiniProgramMessage
		if err := json.Unmarshal(bytes.TrimSpace(message), &req); err != nil {
			return false
		}

		switch req.Type {
		case "ping":
			if err := client.WriteJSON(map[string]string{"type": "pong"}); err != nil {
				logger.Error("failed to send websocket pong",
					zap.Int64("userId", client.UserId),
					zap.Error(err))
			}
			return true
		case "pong":
			return true
		default:
			return false
		}
	}
}

// BuildMiniProgramHeartbeatMessage 构建门禁小程序服务端心跳消息。
func BuildMiniProgramHeartbeatMessage(client *Client) interface{} {
	return map[string]string{"type": "ping"}
}
