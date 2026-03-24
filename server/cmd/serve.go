package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/root-gg/skybook/common"
	"github.com/root-gg/skybook/metadata"
	"github.com/root-gg/skybook/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the SkyBook server",
	Run:   runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Make serve the default command when no subcommand is given
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		serveCmd.Run(cmd, args)
	}
}

func runServe(cmd *cobra.Command, args []string) {
	log := slog.Default()

	// Load config
	config, err := common.LoadConfig(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("No config file found, using defaults", "path", configPath)
			config = common.NewConfig()
		} else {
			log.Error("Failed to load config", "error", err)
			os.Exit(1)
		}
	}

	// Apply env var overrides
	config.ApplyEnvironment()

	if err := config.Validate(); err != nil {
		log.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	// Set log level
	if config.Server.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	log.Info("Starting SkyBook",
		"listen", fmt.Sprintf("%s:%d", config.Server.ListenAddress, config.Server.ListenPort),
	)

	// Initialize database
	backend, err := metadata.NewBackend(config.Database, log)
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer backend.Shutdown()

	// Create and start server
	srv := server.NewSkyBookServer(config, backend, log)
	go func() {
		if err := srv.Start(); err != nil {
			log.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Info("Shutting down", "signal", sig)

	if err := srv.Shutdown(); err != nil {
		log.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}

	log.Info("SkyBook stopped")
}
