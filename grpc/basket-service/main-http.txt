// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/url"
// 	"os"
// 	"strconv"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// 	"github.com/rabbitmq/amqp091-go"
// 	"github.com/redis/go-redis/v9"
// )

// type Basket struct {
// 	ID    uuid.UUID    `json:"id"`
// 	Items []BasketItem `json:"items"`
// }

// type BasketItem struct {
// 	ID        uuid.UUID `json:"id"`
// 	ProductId uuid.UUID `json:"catalogId"`
// 	Price     float32   `json:"price"`
// 	Quantity  uint      `json:"quantity"`
// }

// func main() {
// 	app := fiber.New()

// 	ch := connectToRabbitMQ()
// 	redis := connectToRedis()

// 	app.Get("/api/basket/:id", func(c *fiber.Ctx) error {
// 		id, err := uuid.Parse(c.Params("id"))
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Could not parse id to UUID"})
// 		}

// 		value, err := redis.Get(c.Context(), id.String()).Result()
// 		if err != nil {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Record not found"})
// 		}

// 		var basket Basket
// 		if err := json.Unmarshal([]byte(value), &basket); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Record invalid"})
// 		}

// 		return c.Status(fiber.StatusOK).JSON(basket)
// 	})

// 	app.Get("/api/basket/:id/checkout", func(c *fiber.Ctx) error {
// 		id, err := uuid.Parse(c.Params("id"))
// 		if err != nil {
// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Could not parse id to UUID"})
// 		}

// 		value, err := redis.Get(c.Context(), id.String()).Result()
// 		if err != nil {
// 			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Record not found"})
// 		}

// 		var basket Basket
// 		if err := json.Unmarshal([]byte(value), &basket); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Record invalid"})
// 		}

// 		q, err := ch.QueueDeclare(
// 			"orders", // name
// 			true,     // durable
// 			false,    // delete when unused
// 			false,    // exclusive
// 			false,    // no-wait
// 			nil,      // arguments
// 		)
// 		if err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Unable to checkout basket"})
// 		}

// 		if err := ch.PublishWithContext(c.Context(),
// 			"",     // exchange
// 			q.Name, // routing key
// 			false,  // mandatory
// 			false,  // immediate
// 			amqp091.Publishing{
// 				Body: []byte(value),
// 			}); err != nil {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Unable to checkout basket"})
// 		}

// 		return c.SendStatus(fiber.StatusNoContent)
// 	})

// 	app.Listen(":3000")
// }

// func connectToRedis() *redis.Client {
// 	server := os.Getenv("REDIS_HOST")
// 	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
// 	failOnError(err, "Could parse REDIS_PORT to int")
// 	pass := os.Getenv("REDIS_PASSWORD")

// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     fmt.Sprintf("%s:%d", server, port),
// 		Password: pass,
// 		DB:       0,
// 	})

// 	return rdb
// }

// func connectToRabbitMQ() *amqp091.Channel {
// 	user := os.Getenv("RABBITMQ_USERNAME")
// 	pass := os.Getenv("RABBITMQ_PASSWORD")
// 	host := os.Getenv("RABBITMQ_HOST")
// 	port := os.Getenv("RABBITMQ_PORT")

// 	u := &url.URL{
// 		Scheme: "amqp",
// 		User:   url.UserPassword(user, pass),
// 		Host:   fmt.Sprintf("%s:%s", host, port),
// 	}

// 	conn, err := amqp091.Dial(u.String())
// 	failOnError(err, "Failed to connect to RabbitMQ")
// 	defer conn.Close()

// 	ch, err := conn.Channel()
// 	failOnError(err, "Failed to open a channel")
// 	defer ch.Close()

// 	return ch
// }

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }
