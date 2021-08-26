package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mjehanno/goldenerd/config"
	goldmanager "github.com/mjehanno/goldenerd/gold-manager"
	"github.com/mjehanno/goldenerd/transaction"
	"github.com/vectorhacker/goro"
)

func main() {
	// Creating fiber app with config
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		AppName:       "goldener",
	})

	// Default middleware config
	app.Use(logger.New())
	//var cashedCoin goldmanager.Coin

	currentConf := config.GetConfigFromDb()

	fmt.Println("lastevent : ", currentConf.LastReadEvent)
	// Eventstore
	client := goro.Connect("http://127.0.0.1:2113", goro.WithBasicAuth("admin", ""))
	catchupSubscription := client.CatchupSubscription("transactions", int64(currentConf.LastReadEvent))
	writer := client.Writer("transactions")
	ctx := context.Background()

	go func() {
		transactions := catchupSubscription.Subscribe(ctx)

		t := new(transaction.Transaction)
		for tr := range transactions {
			// Unquoting received string to sanitize
			jsonInput, err := strconv.Unquote(string(tr.Event.Data))
			if err != nil {
				fmt.Println(err)
			}
			// Json Parsing
			err = json.Unmarshal([]byte(jsonInput), &t)
			if err != nil {
				fmt.Println(err)
			}
			transaction.AddTransaction(*t)
			coins := goldmanager.GetCurrentGoldAmount()

			if t.Type == transaction.Credit {
				for _, a := range t.Amount {
					switch a.Currency {
					case goldmanager.Copper:
						coins.Copper += a.Value
					case goldmanager.Silver:
						coins.Silver += a.Value
					case goldmanager.Electrum:
						coins.Electrum += a.Value
					case goldmanager.Gold:
						coins.Gold += a.Value
					case goldmanager.Platinum:
						coins.Platinum += a.Value
					}
				}
			} else {
				sort.Slice(t.Amount, func(i, j int) bool { return t.Amount[i].Currency < t.Amount[j].Currency })
				if coins.Copper >= t.Amount[0].Value {
					coins.Copper -= t.Amount[0].Value
				} else {

				}
				if coins.Silver >= t.Amount[1].Value {
					coins.Silver -= t.Amount[1].Value
				} else {

				}
				if coins.Electrum >= t.Amount[2].Value {
					coins.Electrum -= t.Amount[2].Value
				} else {

				}
				if coins.Gold >= t.Amount[3].Value {
					coins.Gold -= t.Amount[3].Value
				} else {

				}
				if coins.Platinum >= t.Amount[4].Value {
					coins.Platinum -= t.Amount[4].Value
				}
			}
			goldmanager.UpdateGoldAmount(coins)
			currentConf.LastReadEvent++
			config.UpdateConfig(currentConf)

		}
	}()

	api := app.Group("/api")

	gold := api.Group("/gold")

	gold.Get("/", func(c *fiber.Ctx) error {
		coins := goldmanager.GetCurrentGoldAmount()
		sum := 0
		if reflectedCoin := reflect.ValueOf(coins); reflectedCoin.Kind() == reflect.Struct {
			for i := 0; i < int(goldmanager.Limit); i++ {
				index := reflectedCoin.Field(i + 1).Int()
				sum += int(goldmanager.Convert(int(index), goldmanager.Currency(i), goldmanager.Copper))
			}
		}
		return c.JSON(struct {
			Gold int
		}{Gold: goldmanager.Convert(sum, goldmanager.Copper, goldmanager.Gold)})
	})

	gold.Get("/details", func(c *fiber.Ctx) error {
		coins := goldmanager.GetCurrentGoldAmount()
		coins.Id = ""
		return c.JSON(coins)
	})

	tr := api.Group("/transactions")

	tr.Get("/history", func(c *fiber.Ctx) error {
		transations := transaction.GetAllTransactionHistory()
		return c.JSON(transations)
	})

	/**
	*	TODO: Protect this route with keycloak.
	 */
	tr.Post("/", func(c *fiber.Ctx) error {

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
