package cmd

import (
	"github.com/labstack/gommon/log"
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
		a := getApp()
		defer a.Amqp.Close()
		ch, err := a.Amqp.Channel()
		if err != nil {
			log.Fatalf("Failed to open amqp channel")
		}

		err = ch.ExchangeDeclare(
			a.Config.Broker.Exchange,
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
			a.Config.Broker.Queue,
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
			a.Config.Broker.Exchange,
			false,
			nil)
		if err != nil {
			log.Fatalf("Unable to bind queue")
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
			log.Fatalf("Unable to register consumer")
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
