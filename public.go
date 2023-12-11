package onekmq

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Initialize(path string) (Queues, Exchange) {
	// Open configuration file
	fileContent, err := os.ReadFile(path)

	if err != nil {
		FailOnError(errors.New("[File]"), "Couldn't open file "+path)
	}

	var queues Queues
	json.Unmarshal(fileContent, &queues)
	var exchange Exchange
	json.Unmarshal(fileContent, &exchange)

	return queues, exchange
}

func LaunchService(dialName string, exchange ExchangeRaw, queues Queues) {
	conn := NewMQConnection(dialName)
	defer conn.Close()
	mqChannel := NewMQChannel(conn)
	defer mqChannel.Close()

	loadExchange(mqChannel, exchange)

	for _, currQueue := range queues.QueueMap {
		loadQueues(mqChannel, *currQueue)
		bindQueue(mqChannel, *currQueue, exchange)
	}

	rabbitDraw := "=========================================\n"
	rabbitDraw += "            //       //      \n"
	rabbitDraw += "            ('>     ('>      \n"
	rabbitDraw += "            /rr     /rr     \n"
	rabbitDraw += "           *\\))_   *\\))_   \n"
	rabbitDraw += "========================================="

	fmt.Print("[x] Successfully launched RabbitMQ items\n")
	fmt.Print(rabbitDraw)

	var forever chan struct{}
	<-forever
}

// "amqp://guest:guest@localhost:5672/"
func NewMQConnection(dialName string) *amqp.Connection {
	con, err := amqp.Dial(dialName)
	FailOnError(err, "Failed to connect to RabbitMQ")

	return con
}
func NewMQChannel(con *amqp.Connection) *amqp.Channel {

	ch, err := con.Channel()
	FailOnError(err, "Failed to open a channel")

	return ch
}

func (queues *Queues) GetQueueByNameSubstr(nameSubstr string) *Queue {

	for key, queue := range queues.QueueMap {
		if strings.Contains(key, nameSubstr) {
			return queue
		}
	}

	return nil
}

func (consumers *Consumers) GetConsumerByNameSubstr(nameSubstr string) *Consumer {

	for key, consumer := range consumers.ConsumerMap {
		if strings.Contains(key, nameSubstr) {
			return consumer
		}
	}

	return nil
}

func Consume(
	configurationFile string,
	dialName string,
	consumerName string,
	queueName string,
	callbackProcessor func(string) bool) {

	// Retrieve structs via configuration file
	queues, _ := Initialize(configurationFile)
	consumers := initializeConsumer(configurationFile)

	// Validate queue full/partial name
	queue := queues.GetQueueByNameSubstr(queueName)
	if queue == nil {
		FailOnError(errors.New("[CONSUMER]"), "Failed to find queue "+queueName)
	} else {

		// Validate consumer/partial name
		consumer := consumers.GetConsumerByNameSubstr(consumerName)
		if consumer == nil {
			FailOnError(errors.New("[CONSUMER]"), "Failed to find consumer "+consumerName)
		} else {
			// Start a new channel and a new connection
			conn := NewMQConnection(dialName)
			mqChannel := NewMQChannel(conn)

			// Quality of Service declaration (for concurrency needs)
			err := mqChannel.Qos(
				consumer.ConcurrentsNumber,
				0,
				false,
			)
			FailOnError(err, "Failed to set RabbitMQ QoS")

			// Consumer declaration
			msgs, err := mqChannel.Consume(
				consumer.QueueName,
				consumer.ConsumerName,
				consumer.AutoAck,
				consumer.IsExclusive,
				consumer.NoLocal,
				consumer.NoWait,
				consumer.Arguments,
			)
			FailOnError(err, "Failed to register a consumer")

			// Use to acknowledge concurrent works
			var wg sync.WaitGroup

			// Processing queued messages
			go func() {
				for msg := range msgs {
					// Ackowledge a new concurrent worker
					wg.Add(1)

					// Ensure that the message will be right (since it can change using concurrency)
					msg := msg

					go func(msg amqp.Delivery) {
						defer wg.Done()
						isMessageProcessed := callbackProcessor(string(msg.Body))
						// Do not need to acknowledge the message if we're in an AutoAck context
						if !consumer.AutoAck {
							if isMessageProcessed {
								msg.Ack(false)
							} else {
								msg.Nack(false, false)
							}
						}
					}(msg)
				}
			}()

			log.Printf("\n[*] Waiting for messages. To exit press CTRL+C\n")

			// Wait for interrupt signal to gracefully exit
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, os.Interrupt)
			<-sig

			// Ensure every consumption has been completed
			wg.Wait()

			mqChannel.Close()
			conn.Close()

		}
	}
}
