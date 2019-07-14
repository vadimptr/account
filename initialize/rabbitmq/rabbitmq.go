package rabbitmq

import (
	"account-sync/models"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

var Client *amqp.Connection
var Channel *amqp.Channel

const DefaultUrl = "amqp://zwdkijew:FgA3Ilyct6--rfo1zfXIMFRlNpO6OC5j@whale.rmq.cloudamqp.com/zwdkijew"
const InputExchangeName = "input_balance_change_exchange"
const InputQueueName = "input_balance_change_queue"
const OutputExchangeName = "output_balance_change_exchange"
const OutputQueueName = "output_balance_change_queue"

func init() {
	url := os.Getenv("CLOUDAMQP_URL")
	if url == "" {
		url = DefaultUrl
	}
	connectToAmqp(url)
}

func connectToAmqp(url string) {
	fmt.Printf("Connection to amqp... ")
	var err error
	Client, err = amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	Channel, err = Client.Channel()
	if err != nil {
		panic(err)
	}
	err = Channel.Qos(1,0,false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[success]\n")
	setupAmqp()
}

func setupAmqp() {
	setupQueue(InputExchangeName, InputQueueName)
	setupQueue(OutputExchangeName, OutputQueueName)
}

func setupQueue(exchange, queue string) {
	err := Channel.ExchangeDeclare(exchange, "topic", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	_, err = Channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	err = Channel.QueueBind(queue, queue, exchange, false, nil)
	if err != nil {
		panic(err)
	}
}

func EnqueueMessage(message *models.OutputMessage) error {
	var bytes []byte
	var err error
	bytes, err = json.Marshal(message)
	if err != nil {
		return err
	}
	err = Channel.Publish(OutputExchangeName, OutputQueueName, true, false,	amqp.Publishing{Body: bytes,})
	return err
}

func Listen() (<-chan amqp.Delivery, error) {
	return Channel.Consume(InputQueueName, "", false, false, false, false, nil)
}
