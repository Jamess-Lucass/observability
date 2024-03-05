package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ctx = context.Background()

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"exchange-01", // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("%v", err)
	}

	http.HandleFunc("GET /sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Transfer-Encoding", "chunked")

		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		err = ch.QueueBind(
			q.Name,        // queue name
			"",            // routing key
			"exchange-01", // exchange
			false,
			nil,
		)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		go func() {
			<-r.Context().Done()
			log.Println("CLOSED")
			ch.QueueDelete(q.Name, false, false, false)
		}()

		for d := range msgs {
			msg := fmt.Sprintf("data: %s\n\n", string(d.Body))
			fmt.Fprintf(w, msg)

			w.(http.Flusher).Flush()
		}
	})

	http.HandleFunc("POST /producer", func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("message: %s", uuid.New().String())

		err := ch.PublishWithContext(ctx,
			"exchange-01", // exchange
			"",            // routing key
			false,         // mandatory
			false,         // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})
		if err != nil {
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		json.NewEncoder(w).Encode(fmt.Sprintf("published event: %s", msg))
	})

	fmt.Println("Starting web server on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
