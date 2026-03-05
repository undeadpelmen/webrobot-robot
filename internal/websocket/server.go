package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/undeadpelmen/webrobot-robot/internal/interfaces"
)

type Server struct {
	robotService interfaces.RobotController
	logger       zerolog.Logger
	upgrader     websocket.Upgrader
	clients      map[*websocket.Conn]bool
	clientsMu    sync.RWMutex
	broadcast    chan []byte
	register     chan *websocket.Conn
	unregister   chan *websocket.Conn
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type MoveMessage struct {
	Direction string `json:"direction"`
	Speed     int    `json:"speed"`
}

type StatusMessage struct {
	Status string `json:"status"`
	Speed  int    `json:"speed"`
}

func NewServer(robotService interfaces.RobotController, logger zerolog.Logger) *Server {
	return &Server{
		robotService: robotService,
		logger:       logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logger.Info().Str("addr", addr).Msg("Starting WebSocket server")

	http.HandleFunc("/ws", s.handleWebSocket)

	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	go s.broadcastLoop()

	go func() {
		<-ctx.Done()
		s.logger.Info().Msg("Shutting down WebSocket server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	return server.ListenAndServe()
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to upgrade connection")
		return
	}

	s.clientsMu.Lock()
	s.clients[conn] = true
	s.clientsMu.Unlock()

	s.logger.Info().Msg("WebSocket client connected")

	go s.readPump(conn)
	go s.writePump(conn)
}

func (s *Server) readPump(conn *websocket.Conn) {
	defer func() {
		s.unregister <- conn
		conn.Close()
	}()

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.Error().Err(err).Msg("WebSocket error")
			}
			break
		}

		if err := s.handleMessage(conn, msg); err != nil {
			s.logger.Error().Err(err).Str("type", msg.Type).Msg("Failed to handle message")
		}
	}
}

func (s *Server) writePump(conn *websocket.Conn) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-s.broadcast:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *Server) handleMessage(conn *websocket.Conn, msg Message) error {
	switch msg.Type {
	case "move":
		return s.handleMove(conn, msg.Payload)
	case "stop":
		return s.handleStop(conn)
	case "status":
		return s.handleStatus(conn)
	case "set_speed":
		return s.handleSetSpeed(conn, msg.Payload)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func (s *Server) handleMove(conn *websocket.Conn, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var moveMsg MoveMessage
	if err := json.Unmarshal(data, &moveMsg); err != nil {
		return err
	}

	if moveMsg.Speed == 0 {
		moveMsg.Speed = s.robotService.GetSpeed()
	}

	if err := s.robotService.Move(moveMsg.Direction, moveMsg.Speed); err != nil {
		s.sendError(conn, "move_failed", err.Error())
		return err
	}

	s.sendSuccess(conn, "move_successful")
	return nil
}

func (s *Server) handleStop(conn *websocket.Conn) error {
	if err := s.robotService.Stop(); err != nil {
		s.sendError(conn, "stop_failed", err.Error())
		return err
	}

	s.sendSuccess(conn, "stop_successful")
	return nil
}

func (s *Server) handleStatus(conn *websocket.Conn) error {
	status := s.robotService.Status()
	speed := s.robotService.GetSpeed()

	response := Message{
		Type: "status_response",
		Payload: StatusMessage{
			Status: status,
			Speed:  speed,
		},
	}

	return conn.WriteJSON(response)
}

func (s *Server) handleSetSpeed(conn *websocket.Conn, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var req struct {
		Speed int `json:"speed"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		return err
	}

	if err := s.robotService.SetSpeed(req.Speed); err != nil {
		s.sendError(conn, "set_speed_failed", err.Error())
		return err
	}

	s.sendSuccess(conn, "speed_set_successful")
	return nil
}

func (s *Server) sendError(conn *websocket.Conn, errorType, message string) {
	response := Message{
		Type: "error",
		Payload: map[string]string{
			"error":   errorType,
			"message": message,
		},
	}
	conn.WriteJSON(response)
}

func (s *Server) sendSuccess(conn *websocket.Conn, message string) {
	response := Message{
		Type: "success",
		Payload: map[string]string{
			"message": message,
		},
	}
	conn.WriteJSON(response)
}

func (s *Server) broadcastLoop() {
	for {
		select {
		case client := <-s.register:
			s.clientsMu.Lock()
			s.clients[client] = true
			s.clientsMu.Unlock()

		case client := <-s.unregister:
			s.clientsMu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				client.Close()
			}
			s.clientsMu.Unlock()

		case message := <-s.broadcast:
			s.clientsMu.RLock()
			for client := range s.clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					client.Close()
					delete(s.clients, client)
				}
			}
			s.clientsMu.RUnlock()
		}
	}
}
