package websocket

import (
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	logging "github.com/gcc798/lightning/internal/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// ErrHubClosed 表示 Hub 已停止，无法再投递消息。
var ErrHubClosed = errors.New("websocket hub 已关闭")

// Hub 是 WebSocket 连接管理中心。
type Hub struct {
	clients    map[int64]map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	quit       chan struct{}
	startOnce  sync.Once
	stopOnce   sync.Once
	wg         sync.WaitGroup
	mu         sync.RWMutex
	logger     logging.Logger
}

// Client WebSocket 客户端。
type Client struct {
	UserId       int64
	Conn         *websocket.Conn
	Hub          *Hub
	writeTimeout time.Duration
	readTimeout  time.Duration
	writeMu      sync.Mutex
	activeSeq    atomic.Uint64
}

// Message 推送消息。
type Message struct {
	UserId  int64       `json:"-"`
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	payload []byte      `json:"-"`
}

// NewHub 创建 WebSocket Hub。
func NewHub(logger logging.Logger) *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
		quit:       make(chan struct{}),
		logger:     logger,
	}
}

// Start 启动 Hub 处理循环，可重复调用但只生效一次。
func (h *Hub) Start() {
	h.startOnce.Do(func() {
		h.wg.Add(1)
		go h.run()
	})
}

// Close 停止 Hub 并等待处理循环退出、全部连接关闭，可重复调用。
func (h *Hub) Close() {
	h.stopOnce.Do(func() {
		close(h.quit)
	})
	h.wg.Wait()
}

func (h *Hub) run() {
	defer h.wg.Done()
	defer h.closeAll()

	for {
		select {
		case <-h.quit:
			return
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserId] == nil {
				h.clients[client.UserId] = make(map[*Client]bool)
			}
			h.clients[client.UserId][client] = true
			h.mu.Unlock()
			h.logger.Info("websocket client registered",
				zap.Int64("userId", client.UserId),
				zap.Int("totalConnections", len(h.clients[client.UserId])))

		case client := <-h.unregister:
			h.removeClient(client)

		case message := <-h.broadcast:
			h.sendToUser(message)
		}
	}
}

// removeClient 摘除并关闭连接，仅由处理循环调用。
func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	connections, ok := h.clients[client.UserId]
	if !ok || !connections[client] {
		h.mu.Unlock()
		return
	}
	delete(connections, client)
	if len(connections) == 0 {
		delete(h.clients, client.UserId)
	}
	h.mu.Unlock()

	client.Conn.Close()
	h.logger.Info("websocket client unregistered", zap.Int64("userId", client.UserId))
}

// closeAll 关闭全部残留连接，处理循环退出时调用。
func (h *Hub) closeAll() {
	h.mu.Lock()
	clients := h.clients
	h.clients = make(map[int64]map[*Client]bool)
	h.mu.Unlock()

	total := 0
	for _, connections := range clients {
		for client := range connections {
			client.Conn.Close()
			total++
		}
	}
	h.logger.Info("websocket hub stopped", zap.Int("closedConnections", total))
}

func (h *Hub) sendToUser(message *Message) {
	if message.payload == nil {
		payload, err := json.Marshal(message)
		if err != nil {
			h.logger.Error("failed to marshal websocket message",
				zap.Int64("userId", message.UserId),
				zap.Error(err))
			return
		}
		message.payload = payload
	}

	h.mu.RLock()
	targets := make([]*Client, 0, len(h.clients[message.UserId]))
	for client := range h.clients[message.UserId] {
		targets = append(targets, client)
	}
	h.mu.RUnlock()

	if len(targets) == 0 {
		h.logger.Debug("no websocket connections for user", zap.Int64("userId", message.UserId))
		return
	}

	for _, client := range targets {
		if err := client.WriteMessage(websocket.TextMessage, message.payload); err != nil {
			h.logger.Error("failed to send websocket message",
				zap.Int64("userId", message.UserId),
				zap.Error(err))
			h.removeClient(client)
		}
	}

	h.logger.Debug("websocket message sent",
		zap.Int64("userId", message.UserId),
		zap.String("type", message.Type),
		zap.Int("connections", len(targets)))
}

// SendToUser 发送消息给指定用户。
func (h *Hub) SendToUser(userId int64, msgType string, data interface{}) error {
	return h.enqueue(&Message{UserId: userId, Type: msgType, Data: data})
}

// SendJSONToUser 发送已组装好的 JSON 结构给指定用户。
func (h *Hub) SendJSONToUser(userId int64, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return h.enqueue(&Message{UserId: userId, payload: data})
}

// enqueue 投递消息，Hub 已关闭时立即返回 ErrHubClosed，避免调用方永久阻塞。
func (h *Hub) enqueue(message *Message) error {
	select {
	case <-h.quit:
		return ErrHubClosed
	case h.broadcast <- message:
		return nil
	}
}

// Register 注册客户端连接，Hub 已关闭时直接关闭连接。
func (h *Hub) Register(client *Client) {
	select {
	case <-h.quit:
		client.Conn.Close()
	case h.register <- client:
	}
}

// Unregister 注销客户端连接，Hub 已关闭时直接返回（连接由 closeAll 关闭）。
func (h *Hub) Unregister(client *Client) {
	select {
	case <-h.quit:
	case h.unregister <- client:
	}
}

// GetConnectionCount 获取指定用户的连接数。
func (h *Hub) GetConnectionCount(userId int64) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if connections, ok := h.clients[userId]; ok {
		return len(connections)
	}
	return 0
}

// GetTotalConnections 获取总连接数。
func (h *Hub) GetTotalConnections() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, connections := range h.clients {
		total += len(connections)
	}
	return total
}

// WriteMessage 向客户端写入 WebSocket 消息。
func (c *Client) WriteMessage(messageType int, payload []byte) error {
	return c.writeMessage(messageType, payload, true)
}

func (c *Client) writeMessage(messageType int, payload []byte, refreshReadDeadline bool) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.writeTimeout > 0 {
		if err := c.Conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
			return err
		}
	}
	if err := c.Conn.WriteMessage(messageType, payload); err != nil {
		return err
	}
	if refreshReadDeadline {
		c.markActive()
		return c.refreshReadDeadline()
	}
	return nil
}

// WriteJSON 向客户端写入 JSON 文本消息。
func (c *Client) WriteJSON(payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return c.WriteMessage(websocket.TextMessage, data)
}

func (c *Client) refreshReadDeadline() error {
	if c.readTimeout <= 0 {
		return nil
	}
	return c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout))
}

func (c *Client) markActive() {
	c.activeSeq.Add(1)
}

func (c *Client) activeSequence() uint64 {
	return c.activeSeq.Load()
}
