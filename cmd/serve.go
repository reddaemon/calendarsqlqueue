package cmd

import "C"
import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/reddaemon/calendarsqlqueue/internal/database/postgres"

	"github.com/reddaemon/calendarsqlqueue/config"
	db2 "github.com/reddaemon/calendarsqlqueue/db"
	"github.com/reddaemon/calendarsqlqueue/internal/domain/grpc/server"
	"github.com/reddaemon/calendarsqlqueue/internal/domain/grpc/service"
	"github.com/reddaemon/calendarsqlqueue/logger"
	eventpb "github.com/reddaemon/calendarsqlqueue/protofiles"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	rootCmd.AddCommand(serve)
}

var serve = &cobra.Command{
	Use:   "serve",
	Short: "Start calendar server",
	Long:  `Start calendar server`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.GetConfig()
		if err != nil {
			log.Fatalf("unable to load config: %v", err)
		}

		l, err := logger.GetLogger(c)
		if err != nil {
			log.Fatalf("unable to load logger: %v", err)
		}

		db, err := db2.GetDb(c)

		if err != nil {
			log.Fatalf("unable to load db: %v", err)
		}

		lis, err := net.Listen("tcp", c.Host+":"+c.Port)
		if err != nil {
			l.Fatal(fmt.Sprintf("failed to listen %v", err))
		}
		grpcServer := grpc.NewServer()

		var contextTimeout time.Duration
		contextTimeout = time.Millisecond * 500
		repo := postgres.NewPsqlRepository(db)
		ucs := service.NewEventUsecase(repo, contextTimeout)
		eventpb.RegisterEventServiceServer(grpcServer, server.NewServer(ucs, l))

		err = grpcServer.Serve(lis)

		if err != nil {
			log.Fatal(err)
		}
	},
}
