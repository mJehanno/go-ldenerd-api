package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mjehanno/goldenerd/transaction"
	"github.com/vectorhacker/goro"
)

func main() {

	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		AppName:       "goldener",
	})

	// Default middleware config
	app.Use(logger.New())

	client := goro.Connect("http://127.0.0.1:2113", goro.WithBasicAuth("admin", ""))
	catchupSubscription := client.CatchupSubscription("transactions", 0) // sta
	writer := client.Writer("transactions")
	go func() {
		ctx := context.Background()

		transactions := catchupSubscription.Subscribe(ctx)

		for transaction := range transactions {
			fmt.Printf("%s\n", transaction.Event.Data)
		}
	}()

	ctx := context.Background()

	app.Get("api/transactions/history", func(c *fiber.Ctx) error {
		reader := client.FowardsReader("transactions")

		events, err := reader.Read(ctx, 0, 1)
		if err != nil {
			panic(err)
		}
		return c.JSON(events)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/api/transactions", func(c *fiber.Ctx) error {

		t := new(transaction.Transaction)
		if err := c.BodyParser(t); err != nil {
			fmt.Println("can't convert body")
			return err
		}
		obj, err := json.Marshal(&t)
		if err != nil {
			fmt.Println("can't convert to json")
			return err
		}

		event := goro.CreateEvent("transaction", json.RawMessage(obj), nil, 0)
		err = writer.Write(ctx, goro.ExpectedVersionAny, event)
		if err != nil {
			fmt.Println("can't write in eventstore")
			return (err)
		}
		return c.SendStatus(201)
	})

	app.Listen(":8080")
}
