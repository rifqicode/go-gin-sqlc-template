package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gitlab.gift.id/lv2/loyalty/cmd/server"
	logger "gitlab.gift.id/lv2/loyalty/pkg/logger/zap"
)

var rootCmd = &cobra.Command{
	Use:   "loyalty",
	Short: "Launch service 🚀",
	Long:  `Launch TADA Loyalty service`,
}

var RunServerCmd = &cobra.Command{
	Use:     "run-server",
	Short:   "Running the server",
	Example: "loyalty run-server",
	Run: func(cmd *cobra.Command, args []string) {
		appServer, err := server.NewApp()

		if err != nil {
			logger.Error("Failed to initialize loyalty server", map[string]interface{}{"error": err.Error()})
			return
		}

		go func() {
			quitSign := make(chan os.Signal, 1)
			signal.Notify(quitSign, syscall.SIGINT, syscall.SIGTERM)

			<-quitSign
			logger.Info("Received shutdown signal, initiating graceful shutdown", nil)

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()

			if err := appServer.Stop(ctx); err != nil {
				logger.Error("Failed to gracefully shut down loyalty server", map[string]interface{}{
					"error":   err.Error(),
					"timeout": "15s",
				})
			} else {
				logger.Info("Loyalty server shut down successfully", nil)
			}
		}()

		logger.Info("Starting loyalty server", nil)
		if err := appServer.Run(); err != nil {
			logger.Error("Loyalty server encountered an error and stopped", map[string]interface{}{"error": err.Error()})
		}
	},
}

func init() {
	rootCmd.AddCommand(RunServerCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Failed to execute loyalty command", map[string]interface{}{
			"error":   err.Error(),
			"command": rootCmd.Use,
		})
		os.Exit(1)
	}
}
