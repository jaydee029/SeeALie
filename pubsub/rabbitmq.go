package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type QueueType int
type Acktype int

const (
	DurableQueue QueueType = iota
	TransientQueue
)

const (
	Ack Acktype = iota
	NackRequeue
	NackDiscard
)

func (pb *PubSub) PublishGob(exchange, key string, val any) error {
	ch, err := pb.conn.Channel()

	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err = enc.Encode(val)

	if err != nil {
		return err
	}

	return ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/gob",
		Body:        buffer.Bytes(),
	})
}

func (pb *PubSub) DeclareAndBind(exchange, queueName, key string, simpleQueueType QueueType) (*amqp.Channel, amqp.Queue, error) {
	pubchannel, err := pb.conn.Channel()

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	table := amqp.Table{}
	table["x-dead-letter-exchange"] = "SeeALie_dlx"
	pubqueue, err := pubchannel.QueueDeclare(queueName,
		simpleQueueType == DurableQueue,
		simpleQueueType != DurableQueue,
		simpleQueueType != DurableQueue,
		false, table)

	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = pubchannel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return pubchannel, pubqueue, nil
}

func (pb *PubSub) SubscribeGob(exchange, queueName, key string, simpleQueueType QueueType, handler func(val any) Acktype) error {

	return pb.subscribe(exchange, queueName, key, simpleQueueType, handler, func(msg []byte) (any, error) {
		buffer := bytes.NewBuffer(msg)
		decoder := gob.NewDecoder(buffer)
		var data any
		err := decoder.Decode(&data)
		return data, err
	})
}

func (pb *PubSub) subscribe(exchange, queueName, key string, simpleQueueType QueueType, handler func(val any) Acktype, unmarshaller func([]byte) (any, error)) error {

	subch, queue, err := pb.DeclareAndBind(exchange, queueName, key, simpleQueueType)

	if err != nil {
		return err
	}

	deliverychan, err := subch.Consume(queue.Name, "", false, false, false, false, nil)

	if err != nil {
		return err
	}

	go func() {
		defer subch.Close()
		for k := range deliverychan {
			data, err := unmarshaller(k.Body)
			if err != nil {
				log.Println(err)
			}
			switch handler(data) {
			case Ack:
				k.Ack(false)
			case NackRequeue:
				k.Nack(false, true)
			case NackDiscard:
				k.Nack(false, false)
			}
		}
	}()
	return nil
}
