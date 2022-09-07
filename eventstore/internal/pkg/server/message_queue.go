package server

import (
	"context"
	"encoding/json"
	"eventstore/internal/pb"
	"eventstore/pkg/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (s EventStoreServer) PublishEvent(ctx context.Context, event *pb.Event) error {

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ch, err := s.rabbitmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,       // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
		ctx,
		exchangeName,    // exchange
		event.EventType, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonEvent),
		})
	if err != nil {
		return err
	}

	return nil
}

func (s EventStoreServer) ListenEvents(abort <-chan struct{}) (<-chan string, error) {

	ch, err := s.rabbitmq.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName,       // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}

	events, err := ch.Consume(
		queue.Name,   // queue
		"eventstore", // consumer
		true,         // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		return nil, err
	}

	channel := make(chan string)
	go func() {
		defer close(channel)
		for msg := range events {
			select {
			case channel <- string(msg.Body):
				log.PrintLn("transport", "amqp", "procedure", "ListenEvents", "msg", "received")
			case <-abort: // receive on closed channel can proceed immediately
				log.PrintLn("transport", "amqp", "procedure", "ListenEvents", "msg", "aborting")
				return
			}
		}
	}()

	return channel, nil
}
