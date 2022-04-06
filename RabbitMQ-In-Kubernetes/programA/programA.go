package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"fmt"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(RandInt(65, 90))
	}
	return string(bytes)
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func publishing(n int) (res int, err error) {
	// Connection
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//Channel Connection
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	//Exchange Declaration
	err = ch.ExchangeDeclare(
		"NumberStoreA", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	FailOnError(err, "Failed to declare an exchange")

	// Queue Declaration
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	//Consume for Queue
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	//CorrId Creation
	corrId := RandomString(32)

	err = ch.Publish(
		"NumberStoreA", // exchange
		"rpc_queue",    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	FailOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.Atoi(string(d.Body))
			FailOnError(err, "Failed to convert body to integer")
			break
		}
	}

	return
}

/*
Prediction function
-If comes a produced number which from queue,
-Consumer function working for the prediction.
*/
func prediction(n int) bool {

	for i := 0; i < 5; i++ {

		x := RandInt(0, 9)

		if n == x {
			fmt.Println("The random number is predicted! Congrulations!")
			return true
		}

	}
	return false
}

/*
	Publisher Function. The Work pipline is:
 1-) Create a random number,
 2-) The number delivered to Consumer of Program B.
 3-) Consumer of Program B received a number and return score B.
 4*) At the end of function, the result of ScoreB print at the screen.

*/
func publisher() {
	// Thanks to this for loop, the machine deliver a random number to the consumer every 3 seconds.
	for {
		rand.Seed(time.Now().UTC().UnixNano())

		n := RandInt(0, 9)

		log.Printf(" [x] The Random number from ProgramA is: (%d)", n)
		res, err := publishing(n)
		FailOnError(err, "Failed to handle RPC request")

		log.Printf(" [.] Point Store A is: %d", res)
		time.Sleep(3 * time.Second)
	}
}

/*
	Consumer Function. The Work pipline is:
 1-) Read a queue,
 2-) predict a random number
 3-) If Prediction is achieved, Score value increment 1, Otherwise break a loop

*/
func consumer() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"NumberStoreB", // name
		"direct",       // typea
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	FailOnError(err, "Failed to declare an exchange")
	q, err := ch.QueueDeclare(
		"rpc_queue2", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	// The value defines the max number of unacknowledged deliveries that are permitted on a channel.
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	FailOnError(err, "Failed to set QoS")

	// The function provide a binding with some specifications
	err = ch.QueueBind(
		q.Name,         // queue name
		"rpc_queue2",   // routing key
		"NumberStoreB", // exchange
		false,
		nil,
	)
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}
	var scoreA int
	go func() {
		for d := range msgs {
			n, err := strconv.Atoi(string(d.Body))
			FailOnError(err, "Failed to convert body to integer")

			log.Printf(" [.] Prediction(%d)", n)
			response := prediction(n)

			if response {
				scoreA += 1

			} else {

			}

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte(strconv.Itoa(scoreA)),
				})
			FailOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
func main() {
	//Main part of the Program A. We use a Go Routine function for the Publisher.

	go publisher()
	time.Sleep(1 * time.Second)
	consumer()

}

func bodyFrom(args []string) int {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "30"
	} else {
		s = strings.Join(args[1:], " ")
	}
	n, err := strconv.Atoi(s)
	FailOnError(err, "Failed to convert arg to integer")
	return n
}
