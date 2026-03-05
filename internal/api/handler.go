package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/undeadpelmen/webrobot-robot/internal/interfaces"
)

type Handler struct {
	robotService interfaces.RobotController
	logger       zerolog.Logger
}

type MoveRequest struct {
	Direction string `json:"direction" binding:"required"`
	Speed     int    `json:"speed"`
}

type StatusResponse struct {
	Status string `json:"status"`
	Speed  int    `json:"speed"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewHandler(robotService interfaces.RobotController, logger zerolog.Logger) *Handler {
	return &Handler{
		robotService: robotService,
		logger:       logger,
	}
}

func (h *Handler) SetupRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(h.loggingMiddleware())

	api := router.Group("/api/v1")
	{
		api.GET("/status", h.getStatus)
		api.POST("/move", h.move)
		api.POST("/stop", h.stop)
		api.GET("/speed", h.getSpeed)
		api.PUT("/speed", h.setSpeed)
	}

	router.GET("/health", h.health)

	return router
}

func (h *Handler) getStatus(c *gin.Context) {
	status := h.robotService.Status()
	speed := h.robotService.GetSpeed()

	c.JSON(http.StatusOK, StatusResponse{
		Status: status,
		Speed:  speed,
	})
}

func (h *Handler) move(c *gin.Context) {
	var req MoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid move request")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	if req.Speed == 0 {
		req.Speed = h.robotService.GetSpeed()
	}

	if err := h.robotService.Move(req.Direction, req.Speed); err != nil {
		h.logger.Error().Err(err).Str("direction", req.Direction).Int("speed", req.Speed).Msg("Failed to move robot")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "move_failed",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info().Str("direction", req.Direction).Int("speed", req.Speed).Msg("Robot moved successfully")
	c.JSON(http.StatusOK, gin.H{"message": "move successful"})
}

func (h *Handler) stop(c *gin.Context) {
	if err := h.robotService.Stop(); err != nil {
		h.logger.Error().Err(err).Msg("Failed to stop robot")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "stop_failed",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info().Msg("Robot stopped successfully")
	c.JSON(http.StatusOK, gin.H{"message": "stop successful"})
}

func (h *Handler) getSpeed(c *gin.Context) {
	speed := h.robotService.GetSpeed()
	c.JSON(http.StatusOK, gin.H{"speed": speed})
}

func (h *Handler) setSpeed(c *gin.Context) {
	var req struct {
		Speed int `json:"speed" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().Err(err).Msg("Invalid speed request")
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	if err := h.robotService.SetSpeed(req.Speed); err != nil {
		h.logger.Error().Err(err).Int("speed", req.Speed).Msg("Failed to set speed")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "set_speed_failed",
			Message: err.Error(),
		})
		return
	}

	h.logger.Info().Int("speed", req.Speed).Msg("Speed set successfully")
	c.JSON(http.StatusOK, gin.H{"message": "speed set successful"})
}

func (h *Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "robot-api",
	})
}

func (h *Handler) loggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
