/*
Copyright © 2024 TARIQUE M NASRULLAH nasrullahtarique@gmail.com

*/

package cmd

import (
	"context"
	"fmt"
	"github.com/minishop/internal/delivery/rest"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts a http server",
	Long:  "Starts a http server",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		defer stop()

		srv := rest.New(ctx)

		var srvErr = make(chan error, 1)
		var err error

		go func() {
			fmt.Println("Serving on port :8080")
			srvErr <- srv.ListenAndServe()
		}()

		select {
		case err = <-srvErr:
			return err
		case <-ctx.Done():
			fmt.Println("received os signal, shutting down")
			stop()
		}

		err = srv.Shutdown(context.Background())
		return err
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
