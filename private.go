package onekmq

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func loadExchange(ch *amqp.Channel, exInfos ExchangeRaw) {
	err := ch.ExchangeDeclare(
		exInfos.Name,
		exInfos.Type,
		exInfos.IsDurable,
		exInfos.AutoDelete,
		exInfos.IsInternal,
		exInfos.NoWait,
		exInfos.Arguments,
	)
	FailOnError(err, "Failed to create exchange "+exInfos.Name)
}

func loadQueues(ch *amqp.Channel, qInfos Queue) {
	_, err := ch.QueueDeclare(
		qInfos.Name,
		qInfos.IsDurable,
		qInfos.AutoDelete,
		qInfos.IsExclusive,
		qInfos.NoWait,
		qInfos.Arguments,
	)
	FailOnError(err, "Failed to declare a queue")
}

func bindQueue(ch *amqp.Channel, qInfos Queue, exInfos ExchangeRaw) {
	err := ch.QueueBind(
		qInfos.Name,
		qInfos.RoutingKey,
		exInfos.Name,
		false,
		nil,
	)
	FailOnError(err, "Failed to bind "+qInfos.Name+" queue")
}

func initializeConsumer(pathToConfiguration string) Consumers {
	// Open configuration file
	fileContent, err := os.ReadFile(pathToConfiguration)

	if err != nil {
		FailOnError(errors.New("[File]"), "Couldn't open file "+pathToConfiguration)
	}

	var consumers Consumers
	json.Unmarshal(fileContent, &consumers)

	return consumers
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
