package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/streadway/amqp"
)

var (
	conn   *amqp.Connection
	err    error
	client *redis.Client
)

type response struct {
	IsError bool
	Error   string
	Result  interface{}
}

func getBD(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := response{}
	keyToSearch := "bd-" + ps.ByName("businessDomainName")
	val, err := client.Get(keyToSearch).Result()
	if err != nil {
		resp.IsError = true
		resp.Error = fmt.Sprintf("Failed to fetch from DB, %s", err)
	} else {
		resp.Result = fmt.Sprintf("key: %s, value: %s", keyToSearch, val)
	}

	json.NewEncoder(w).Encode(resp)
}

func createBD(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	resp := response{}
	keyToSet := "bd-" + ps.ByName("businessDomainName")
	err := client.Set(keyToSet, "pending", 0).Err()
	if err != nil {
		resp.IsError = true
		resp.Error = fmt.Sprintf("Failed to set DB: %s", err)
	}

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
	body, err := ioutil.ReadAll(r.Body)
	err = ch.Publish(
		"",    // exchange
		"BD",  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers:     map[string]interface{}{"entry": ps.ByName("businessDomainName"), "request": r.Method},
			ContentType: "application/json",
			Body:        body,
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
	json.NewEncoder(w).Encode(resp)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	router := httprouter.New()
	router.GET("/cd/business-domains/:businessDomainName", getBD)
	router.POST("/cd/business-domains/:businessDomainName", createBD)
	log.Fatal(http.ListenAndServe(":8090", router))
}
