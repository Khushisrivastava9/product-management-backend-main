package queue

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type Queue struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewQueue() (*Queue, error) {
	conn, err := amqp.Dial(os.Getenv("QUEUE_URL"))
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"image_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &Queue{
		conn:    conn,
		channel: channel,
		queue:   queue,
	}, nil
}

func (q *Queue) AddToQueue(imageURL string) error {
	err := q.channel.Publish(
		"",
		q.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(imageURL),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Added image URL to queue: %s", imageURL)
	return nil
}

func (q *Queue) Close() {
	q.channel.Close()
	q.conn.Close()
}
