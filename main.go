package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/streadway/amqp"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host",
			Value:  "localhost",
			EnvVar: "AMQP_HOST",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  5672,
			EnvVar: "AMQP_PORT",
		},
		cli.StringFlag{
			Name:   "user",
			Value:  "guest",
			EnvVar: "AMQP_USER",
		},
		cli.StringFlag{
			Name:   "password",
			Value:  "guest",
			EnvVar: "AMQP_PASSWORD",
		},
		cli.StringFlag{
			Name:   "queue",
			Value:  "eiffel",
			EnvVar: "QUEUE_NAME",
		},
		cli.StringFlag{
			Name:   "file",
			Value:  "events.json",
			EnvVar: "EVENTS_FILE",
		},
	}
	app.Action = func(c *cli.Context) error {
		uri := &amqp.URI{
			Host:     c.String("host"),
			Port:     c.Int("port"),
			Username: c.String("user"),
			Password: c.String("password"),
			Scheme:   "amqp",
			Vhost:    "/",
		}
		log.Println(uri.String())
		conn, err := amqp.Dial(uri.String())
		if err != nil {
			return err
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

		for i, event := range events {
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
			log.Printf("Sent event #%d", i)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
