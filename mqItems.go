package onekmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queues struct {
	QueueMap map[string]*Queue `json:"queues"`
}
type Queue struct {
	Name        string                 `json:"Name"`
	IsDurable   bool                   `json:"IsDurable"`
	AutoDelete  bool                   `json:"AutoDelete"`
	IsExclusive bool                   `json:"IsExclusive"`
	NoWait      bool                   `json:"NoWait"`
	Arguments   map[string]interface{} `json:"Arguments"`
	RoutingKey  string                 `json:"RoutingKey"`
}

type Exchange struct {
	Raw ExchangeRaw `json:"exchange"`
}
type ExchangeRaw struct {
	Name       string                 `json:"Name"`
	Type       string                 `json:"Type"`
	IsDurable  bool                   `json:"IsDurable"`
	AutoDelete bool                   `json:"AutoDelete"`
	IsInternal bool                   `json:"IsInternal"`
	NoWait     bool                   `json:"NoWait"`
	Arguments  map[string]interface{} `json:"Arguments"`
}

type Consumers struct {
	ConsumerMap map[string]*Consumer `json:"consumers"`
}
type Consumer struct {
	QueueName         string     `json:"QueueName"`
	ConsumerName      string     `json:"ConsumerName"`
	AutoAck           bool       `json:"AutoAck"`
	IsExclusive       bool       `json:"IsExclusive"`
	NoLocal           bool       `json:"NoLocal"`
	NoWait            bool       `json:"NoWait"`
	Arguments         amqp.Table `json:"Arguments"`
	ConcurrentsNumber int        `json:"ConcurrentsNumber"`
}

type Producers struct {
	ProducerMap map[string]*Producer `json:"producers"`
}

type Producer struct {
	ExchangeName string `json:"ExchangeName"`
	RoutingKey   string `json:"RoutingKey"`
	DeliveryMode string `json:"DeliveryMode"`
	ContentType  string `json:"ContentType"`
}
