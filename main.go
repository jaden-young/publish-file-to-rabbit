package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/cenkalti/backoff"
	"github.com/streadway/amqp"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Value:  "localhost",
			EnvVar: "RABBIT_HOST",
			Usage:  "hostname of rabbitmq broker",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  5672,
			EnvVar: "RABBIT_PORT",
			Usage:  "port of rabbitmq broker",
		},
		cli.StringFlag{
			Name:   "user",
			Value:  "guest",
			EnvVar: "RABBIT_USER",
			Usage:  "username for rabbitmq broker",
		},
		cli.StringFlag{
			Name:   "password",
			Value:  "guest",
			EnvVar: "RABBIT_PASSWORD",
			Usage:  "password for rabbitmq broker",
		},
		cli.StringFlag{
			Name:   "queue",
			Value:  "eiffel",
			EnvVar: "QUEUE_NAME",
			Usage:  "name of rabbit queue",
		},
		cli.StringFlag{
			Name:   "file",
			Value:  "events.json",
			EnvVar: "EVENTS_FILE",
			Usage:  "file to publish. MUST be a single JSON array of objects.",
		},
		cli.IntFlag{
			Name:   "limit",
			Value:  1000,
			EnvVar: "EVENTS_LIMIT",
			Usage:  "maximum number of objects to read from `FILE` and send over RabbitMQ",
		},
	}
	app.Usage = ""
	app.Description = "Reads a file with a single array of JSON objects and publishes each object as a message to a RabbitMQ queue"
	app.Version = "0.0.2"
	app.Action = func(c *cli.Context) error {
		uri := &amqp.URI{
			Host:     c.String("host"),
			Port:     c.Int("port"),
			Username: c.String("user"),
			Password: c.String("password"),
			Scheme:   "amqp",
			Vhost:    "/",
		}

		var conn *amqp.Connection
		err := backoff.Retry(func() error {
			log.Printf("Dialing %s", uri.String())
			c, err := amqp.Dial(uri.String())
			if err != nil {
				return err
			}
			conn = c
			return nil
		}, backoff.NewExponentialBackOff())
		if err != nil {
			log.Fatalf("Unable to connect to %s", uri.String())
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			return err
		}
		defer ch.Close()

		q, err := ch.QueueDeclare(c.String("queue"), false, false, false, false, nil)
		if err != nil {
			return err
		}

		b, err := ioutil.ReadFile(c.String("file"))
		if err != nil {
			return err
		}

		var obj interface{}
		json.Unmarshal(b, &obj)

		events, ok := obj.([]interface{})
		if !ok {
			return errors.New("Error reading events file")
		}

		limit := c.Int("limit")
		log.Print("Sending events...")
		for i, event := range events {
			if i > limit {
				break
			}
			b, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				return err
			}
			ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        b,
				})
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
