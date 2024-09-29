package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"gitlab.gift.id/lv2/loyalty/cmd/server"
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
			fmt.Println("error initialize server", err)
		}

		go func() {
			quitSign := make(chan os.Signal, 1)
			signal.Notify(quitSign, syscall.SIGINT, syscall.SIGTERM)

			<-quitSign
			fmt.Println("shuting down server...")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()

			if err := appServer.Stop(ctx); err != nil {
				fmt.Println("error shutdown server: %v", err)
			}
		}()

		if err := appServer.Run(); err != nil {
			fmt.Println("error running server", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(RunServerCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
