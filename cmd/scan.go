package cmd

import (
	"time"

	"golang.org/x/net/context"

	"github.com/reddaemon/calendarsqlqueue/internal/domain/grpc/sender"
	"github.com/reddaemon/calendarsqlqueue/internal/domain/grpc/sender/queue"
	"github.com/reddaemon/calendarsqlqueue/models/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scan)
}

var scan = &cobra.Command{
	Use:   "scan",
	Short: "Start scanner server",
	Long:  "Start scanner server",
	Run: func(cmd *cobra.Command, args []string) {
		a := getApp()
		contextTimeout := time.Millisecond * 800
		q := queue.NewEventQueue(a)
		repo := storage.NewPsqlRepository(a)
		ucs := sender.NewEventUseCase(a, repo, q, contextTimeout)
		ctx := context.Background()

		t := time.Now()
		n := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
		d := n.Sub(t)

		for {
			err := ucs.GetEvent(ctx, time.Now())
			if err != nil {
				a.Logger.Error(err.Error())
			}
			time.Sleep(d)
		}
	},
}
