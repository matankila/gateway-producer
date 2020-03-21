package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/streadway/amqp"
)

var (
	conn *amqp.Connection
	err  error
)

func getBD(w http.ResponseWriter, req *http.Request) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"BD",  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	failOnError(err, "Failed to declare a queue")
	body := ""
	err = ch.Publish(
		"",    // exchange
		"BD",  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:     map[string]interface{}{"entry": "bitcore", "req": "get"},
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
	fmt.Fprintf(w, "get BD")
}

func createBD(w http.ResponseWriter, req *http.Request) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"BD",  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	failOnError(err, "Failed to declare a queue")
	body := "{\"config\": {\"users\": [\"IPG\"]}}"
	err = ch.Publish(
		"",    // exchange
		"BD",  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:     map[string]interface{}{"entry": "bitcore", "req": "post"},
			ContentType: "application/json",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
	fmt.Fprintf(w, "Create BD")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// TODO: use httprouter
func main() {
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//TODO: change it to /BD/:BD
	http.HandleFunc("/getBD", getBD)

	// TODO: change it to work with POST
	// TODO: change it to: /BD/:BD
	// TODO: send body param of request.
	http.HandleFunc("/newBD", createBD)
	http.ListenAndServe(":8090", nil)
}
