package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/reddaemon/calendarsqlqueue/config"
	"github.com/reddaemon/calendarsqlqueue/logger"
	"github.com/reddaemon/calendarsqlqueue/queue"
)

var (
	counter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "sender_notifications_total",
			Help: "Sent notifications counter",
		},
	)
)

func init() {
	prometheus.MustRegister(counter)
}

func main() {
	go runPrometheusServer()
	NotifConfigPath := flag.String("config", "config.yml", "path to config file")
	flag.Parse()
	c, err := config.GetConfig(*NotifConfigPath)
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
}

func runPrometheusServer() {
	log.Println("run prometheus exporter server")
	http.Handle("/metrics", promhttp.Handler())
	_ = http.ListenAndServe(":9187", nil)
}
