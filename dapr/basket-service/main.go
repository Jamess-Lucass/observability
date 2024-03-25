package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	dapr "github.com/dapr/go-sdk/client"
)

type Basket struct {
	Id    uuid.UUID    `json:"id"`
	Items []BasketItem `json:"items"`
}

type BasketItem struct {
	Id        uuid.UUID `json:"id"`
	ProductId uuid.UUID `json:"catalogId"`
	Price     float64   `json:"price"`
	Quantity  uint      `json:"quantity"`
}

type CreateBasketRequest struct {
	Items []CreateBasketItemRequest `json:"items"`
}

type CreateBasketItemRequest struct {
	ProductId string `json:"catalogId"`
	Quantity  uint   `json:"quantity"`
}

type ProductResponse struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func main() {
	app := fiber.New()
	ctx := context.Background()

	client, err := dapr.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// basketSchema := schema.Object(map[string]schema.ISchema{
	// 	"Items": schema.Array(schema.Object(map[string]schema.ISchema{
	// 		"ProductId": schema.UUID().Refine(func(productId uuid.UUID) bool {
	// 			_, err := client.InvokeMethod(ctx, "product-service", fmt.Sprintf("/api/products/%s", productId.String()), "GET")
	// 			return err != nil
	// 		}),
	// 		"Quantity":  schema.Int().Gt(0),
	// 	})).Min(1),
	// })

	data, err := client.InvokeMethod(ctx, "product-service", "/api/products/72119506-89ef-4c0c-ace7-6cbd984bfc50", "GET")
		
	log.Printf("%v", err)
	log.Printf("%v", data)

	// app.Get("/api/basket/:id", func(c *fiber.Ctx) error {
	// 	id, err := uuid.Parse(c.Params("id"))
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Could not parse id to UUID"})
	// 	}

	// 	value, err := redis.Get(c.Context(), id.String()).Result()
	// 	if err != nil {
	// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Record not found"})
	// 	}

	// 	var basket Basket
	// 	if err := json.Unmarshal([]byte(value), &basket); err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Record invalid"})
	// 	}

	// 	return c.Status(fiber.StatusOK).JSON(basket)
	// })

	// app.Post("/api/basket", func(c *fiber.Ctx) error {
	// 	var request CreateBasketRequest
	// 	if err := c.BodyParser(&request); err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	// 	}

	// 	// todo: make better
	// 	if len(request.Items) <= 0 {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "At least one product is required"})
	// 	}

	// 	for _, item := range request.Items {
	// 		_, err := uuid.Parse(item.ProductId)
	// 		if err != nil {
	// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Product id must be a valid id"})
	// 		}

	// 		if item.Quantity <= 0 {
	// 			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Quantity must be greater than 1"})
	// 		}
	// 	}

	// 	basket := Basket{
	// 		Id: uuid.New(),
	// 	}

	// 	for _, item := range request.Items {
	// 		uri := fmt.Sprintf("%s/api/products/%s", os.Getenv("PRODUCT_SERVICE_BASE_URL"), item.ProductId)
	// 		res, err := http.Get(uri)
	// 		if err != nil {
	// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not validate product"})
	// 		}
	// 		defer res.Body.Close()

	// 		if res.StatusCode != fiber.StatusOK {
	// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Invalid product"})
	// 		}

	// 		var product ProductResponse
	// 		if err := json.NewDecoder(res.Body).Decode(&product); err != nil {
	// 			return err
	// 		}

	// 		basketItem := BasketItem{
	// 			Id:        uuid.New(),
	// 			ProductId: uuid.MustParse(item.ProductId),
	// 			Price:     product.Price,
	// 			Quantity:  item.Quantity,
	// 		}

	// 		basket.Items = append(basket.Items, basketItem)
	// 	}

	// 	if err := redis.Set(c.Context(), basket.Id.String(), basket, 24*time.Hour).Err(); err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Could not validate product"})
	// 	}

	// 	return c.Status(fiber.StatusOK).JSON(basket)
	// })

	// app.Get("/api/basket/:id/checkout", func(c *fiber.Ctx) error {
	// 	id, err := uuid.Parse(c.Params("id"))
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Could not parse id to UUID"})
	// 	}

	// 	value, err := redis.Get(c.Context(), id.String()).Result()
	// 	if err != nil {
	// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Record not found"})
	// 	}

	// 	var basket Basket
	// 	if err := json.Unmarshal([]byte(value), &basket); err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Record invalid"})
	// 	}

	// 	q, err := ch.QueueDeclare(
	// 		"orders", // name
	// 		true,     // durable
	// 		false,    // delete when unused
	// 		false,    // exclusive
	// 		false,    // no-wait
	// 		nil,      // arguments
	// 	)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Unable to checkout basket"})
	// 	}

	// 	if err := ch.PublishWithContext(c.Context(),
	// 		"",     // exchange
	// 		q.Name, // routing key
	// 		false,  // mandatory
	// 		false,  // immediate
	// 		amqp091.Publishing{
	// 			Body: []byte(value),
	// 		}); err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Unable to checkout basket"})
	// 	}

	// 	return c.SendStatus(fiber.StatusNoContent)
	// })

	app.Listen(":8080")
}
