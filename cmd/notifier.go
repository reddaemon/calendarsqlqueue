package cmd

import (
	"log"

	"github.com/reddaemon/calendarsqlqueue/config"
	"github.com/reddaemon/calendarsqlqueue/logger"
	"github.com/reddaemon/calendarsqlqueue/queue"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(notifier)
}

var notifier = &cobra.Command{
	Use:   "sender",
	Short: "sender",
	Long:  "sender",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.GetConfig()
		if err != nil {
			log.Fatalf("unable to load config: %v", err)
		}
		l, err := logger.GetLogger(c)
		if err != nil {
			log.Fatalf("unable to load logger")
		}

		amqpq := queue.GetConnection(c)

		ch, err := amqpq.Channel()

		if err != nil {
			log.Fatalf("Failed to open amqp channel")
		}
		defer ch.Close()
		err = ch.ExchangeDeclare(
			c.Broker["exchange"],
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			log.Fatalf("Unable to declare exchange")
		}

		q, err := ch.QueueDeclare(
			c.Broker["queue"],
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Unable to declare queue")
		}

		err = ch.QueueBind(
			q.Name,
			"",
			c.Broker["exchange"],
			false,
			nil)
		if err != nil {
			l.Fatal("Unable to bind queue")
		}

		msgs, err := ch.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil)
		if err != nil {
			l.Fatal("Unable to register consumer")
		}

		always := make(chan bool)

		go func() {
			for m := range msgs {
				log.Printf("[x] %s", m.Body)
			}
		}()

		log.Printf("To exit press CTRL + C")
		<-always
	},
}
