package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/undeadpelmen/webrobot-robot/internal/api"
	"github.com/undeadpelmen/webrobot-robot/internal/cli"
	"github.com/undeadpelmen/webrobot-robot/internal/config"
	"github.com/undeadpelmen/webrobot-robot/internal/hardware"
	"github.com/undeadpelmen/webrobot-robot/internal/interfaces"
	"github.com/undeadpelmen/webrobot-robot/internal/robot"
	"github.com/undeadpelmen/webrobot-robot/internal/websocket"
)

type Application struct {
	config       *config.Config
	logger       zerolog.Logger
	robotService interfaces.RobotController
	hardware     interfaces.HardwareDriver
	apiHandler   *api.Handler
	wsServer     *websocket.Server
	cliService   *cli.Service
}

type Flags struct {
	ConfigPath   string
	Debug        bool
	TestMode     bool
	EnableCLI    bool
	EnableHTTP   bool
	EnableWS     bool
	Interactive  bool
	WebSocketURL string
}

func NewApplication() *Application {
	return &Application{}
}

func (app *Application) ParseFlags() *Flags {
	flags := &Flags{}

	flag.StringVar(&flags.ConfigPath, "config", "", "Path to configuration file")
	flag.BoolVar(&flags.Debug, "debug", false, "Enable debug logging")
	flag.BoolVar(&flags.TestMode, "test", false, "Enable test mode with mock hardware")
	flag.BoolVar(&flags.EnableCLI, "cli", false, "Enable CLI interface")
	flag.BoolVar(&flags.EnableHTTP, "http", false, "Enable HTTP API")
	flag.BoolVar(&flags.EnableWS, "websocket", false, "Enable WebSocket server")
	flag.BoolVar(&flags.Interactive, "i", true, "Enable interactive CLI")
	flag.StringVar(&flags.WebSocketURL, "ws-url", "", "WebSocket server URL (for client mode)")

	flag.Parse()

	return flags
}

func (app *Application) Initialize(flags *Flags) error {
	if err := app.loadConfig(flags); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := app.setupLogging(flags); err != nil {
		return fmt.Errorf("failed to setup logging: %w", err)
	}

	if err := app.setupHardware(flags); err != nil {
		return fmt.Errorf("failed to setup hardware: %w", err)
	}

	if err := app.setupRobot(); err != nil {
		return fmt.Errorf("failed to setup robot service: %w", err)
	}

	if err := app.setupServices(); err != nil {
		return fmt.Errorf("failed to setup services: %w", err)
	}

	return nil
}

func (app *Application) loadConfig(flags *Flags) error {
	cfg, err := config.Load(flags.ConfigPath)
	if err != nil {
		return err
	}

	if flags.TestMode {
		cfg.Hardware.TestPins = true
	}

	if flags.EnableCLI {
		cfg.Robot.EnableCLI = true
	}

	if flags.EnableHTTP {
		cfg.Server.EnableHTTP = true
	}

	if flags.EnableWS {
		cfg.Server.EnableWS = true
	}

	app.config = cfg
	return nil
}

func (app *Application) setupLogging(flags *Flags) error {
	if err := app.config.EnsureLogDirectory(); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile, err := os.OpenFile(app.config.Logging.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	var writers io.Writer
	if app.config.Logging.Console {
		writers = io.MultiWriter(os.Stdout, logFile)
	} else {
		writers = logFile
	}

	level := app.config.GetLogLevel()
	if flags.Debug {
		level = zerolog.DebugLevel
	}

	app.logger = zerolog.New(writers).With().Timestamp().Logger().Level(level)
	app.logger.Info().Msg("Logging initialized")

	return nil
}

func (app *Application) setupHardware(flags *Flags) error {
	app.logger.Info().Str("driver", app.config.Hardware.Driver).Msg("Initializing hardware driver")

	driver, err := hardware.NewL298NDriver(app.config.Hardware, app.logger)
	if err != nil {
		return err
	}

	app.hardware = driver
	return nil
}

func (app *Application) setupRobot() error {
	app.logger.Info().Msg("Initializing robot service")

	robotService := robot.NewService(app.hardware, app.logger)
	app.robotService = robotService

	return nil
}

func (app *Application) setupServices() error {
	app.logger.Info().Msg("Setting up application services")

	app.apiHandler = api.NewHandler(app.robotService, app.logger)
	app.wsServer = websocket.NewServer(app.robotService, app.logger)
	app.cliService = cli.NewService(app.robotService, app.logger)

	return nil
}

func (app *Application) Run(ctx context.Context) error {
	app.logger.Info().Msg("Starting application")

	if err := app.robotService.Initialize(ctx); err != nil {
		return fmt.Errorf("failed to initialize robot service: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if app.config.Server.EnableHTTP {
		go app.startHTTPServer(ctx)
	}

	if app.config.Server.EnableWS {
		go app.startWebSocketServer(ctx)
	}

	if app.config.Robot.EnableCLI {
		go app.startCLI()
	}

	app.waitForShutdown(ctx)
	return nil
}

func (app *Application) startHTTPServer(ctx context.Context) {
	addr := fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.HTTPPort)
	app.logger.Info().Str("addr", addr).Msg("Starting HTTP server")

	router := app.apiHandler.SetupRoutes()

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		app.logger.Info().Msg("Shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		app.logger.Error().Err(err).Msg("HTTP server failed")
	}
}

func (app *Application) startWebSocketServer(ctx context.Context) {
	addr := fmt.Sprintf("%s:%d", app.config.Server.Host, app.config.Server.WebSocketPort)
	app.logger.Info().Str("addr", addr).Msg("Starting WebSocket server")

	if err := app.wsServer.Start(ctx, addr); err != nil {
		app.logger.Error().Err(err).Msg("WebSocket server failed")
	}
}

func (app *Application) startCLI() {
	app.logger.Info().Msg("Starting CLI interface")

	if app.config.Robot.AutoConnect {
		app.cliService.StartInteractive()
	} else {
		app.cliService.StartNonInteractive()
	}
}

func (app *Application) waitForShutdown(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		app.logger.Info().Str("signal", sig.String()).Msg("Received shutdown signal")
	case <-ctx.Done():
		app.logger.Info().Msg("Context cancelled, shutting down")
	}

	app.shutdown(ctx)
}

func (app *Application) shutdown(ctx context.Context) {
	app.logger.Info().Msg("Shutting down application")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.robotService.Shutdown(shutdownCtx); err != nil {
		app.logger.Error().Err(err).Msg("Failed to shutdown robot service")
	}

	app.logger.Info().Msg("Application shutdown complete")
}
