package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/mjehanno/go-ldenerd-api/appconfig/conf"
	"github.com/mjehanno/go-ldenerd-api/appconfig/db"
	"github.com/mjehanno/go-ldenerd-api/appconfig/env"
	"github.com/mjehanno/go-ldenerd-api/auth"
	goldmanager "github.com/mjehanno/go-ldenerd-api/gold-manager"
	"github.com/mjehanno/go-ldenerd-api/transaction"
	"github.com/vectorhacker/goro"
)

var eventStoreClient *goro.Client
var catchupSubscription goro.Subscriber
var writer goro.Writer
var eventStoreCtx context.Context
var keycloackCtx context.Context

func main() {
	app, config := appInit()
	conf.CurrentConf = config

	go transactionEventHandler()

	api := app.Group("/api")

	api.Post("/login", func(c *fiber.Ctx) error {
		var u auth.User
		if err := c.BodyParser(&u); err != nil {
			fmt.Println("can't convert body")
			return fiber.NewError(fiber.StatusBadRequest, "no login/password provided")
		}
		token, err := auth.GetClient().Login(keycloackCtx, conf.CurrentConf.KeycloakClientID, conf.CurrentConf.KeycloakSecret, conf.CurrentConf.KeycloakRealm, u.Username, u.Password)
		if err != nil {
			log.Println(fmt.Errorf("error happened while login : %s", err))
			return fiber.ErrBadRequest
		}
		aut := auth.Jwt{AccessToken: token.AccessToken, ExpiresIn: token.ExpiresIn, RefreshToken: token.RefreshToken, TokenType: token.TokenType}
		return c.JSON(aut)
	})

	api.Post("/refresh", func(c *fiber.Ctx) error {
		headToken := c.Get("Authorization")
		headToken = strings.Replace(headToken, "Bearer ", "", -1)
		result, err := auth.GetClient().RetrospectToken(keycloackCtx, headToken, conf.CurrentConf.KeycloakClientID, conf.CurrentConf.KeycloakSecret, conf.CurrentConf.KeycloakRealm)
		if err != nil || !*result.Active {
			return c.SendStatus(401)
		}

		var t auth.Jwt
		if err := c.BodyParser(&t); err != nil {
			fmt.Println(err.Error())
			return fiber.NewError(fiber.StatusBadRequest, "Requests Data not formated correctly")
		}
		token, err := auth.GetClient().RefreshToken(keycloackCtx, t.RefreshToken, conf.CurrentConf.KeycloakClientID, conf.CurrentConf.KeycloakSecret, conf.CurrentConf.KeycloakRealm)
		if err != nil {
			fmt.Println(err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Error in the server")
		}
		return c.JSON(token)
	})

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

	tr.Post("/", func(c *fiber.Ctx) error {
		head := c.Get("Authorization")
		head = strings.Replace(head, "Bearer ", "", -1)

		result, err := auth.GetClient().RetrospectToken(keycloackCtx, head, conf.CurrentConf.KeycloakClientID, conf.CurrentConf.KeycloakSecret, conf.CurrentConf.KeycloakRealm)
		if err != nil || !*result.Active {
			return c.SendStatus(401)
		}
		t := new(transaction.Transaction)
		if err := c.BodyParser(t); err != nil {
			fmt.Println(err.Error())
			return fiber.NewError(fiber.StatusBadRequest, "Requests Data not formated correctly")
		}

		obj, err := json.Marshal(&t)
		if err != nil {
			fmt.Println(err.Error())
			return fiber.NewError(fiber.StatusBadRequest, "Requests Data not formated correctly")
		}

		coins := goldmanager.GetCurrentGoldAmount()
		incomingGold := transaction.ConvertSumOfAmountToCoin(t.Amount)

		_, err = transaction.Align(coins, incomingGold, t.Type)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Not enought coins")
		}

		event := goro.CreateEvent("transaction", json.RawMessage(obj), nil, 0)
		err = writer.Write(eventStoreCtx, goro.ExpectedVersionAny, event)
		if err != nil {
			fmt.Println(err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "can't write to database")
		}
		return c.SendStatus(201)
	})

	app.Listen(":8000")
}

//Init application
func appInit() (*fiber.App, *conf.Config) {
	// Creating fiber app with config
	app := fiber.New(fiber.Config{
		CaseSensitive: true,
		AppName:       "goldener",
	})

	// Default middleware config
	app.Use(logger.New())
	env.GetConfigFromEnv()

	conf.CurrentConf.LastReadEvent, conf.CurrentConf.Id = db.GetConfigFromDb().LastReadEvent, db.GetConfigFromDb().Id
	fmt.Println(conf.CurrentConf)

	//fmt.Println("lastevent : ", conf.CurrentConf.LastReadEvent)
	// Eventstore
	eventStoreClient = goro.Connect(conf.CurrentConf.EventstoreHost, goro.WithBasicAuth(conf.CurrentConf.EventstoreUser, ""))
	catchupSubscription = eventStoreClient.CatchupSubscription(conf.CurrentConf.EvenstoreDb, int64(conf.CurrentConf.LastReadEvent))
	writer = eventStoreClient.Writer(conf.CurrentConf.EvenstoreDb)
	eventStoreCtx = context.Background()
	keycloackCtx = context.Background()

	return app, conf.CurrentConf
}

func transactionEventHandler() {
	transactions := catchupSubscription.Subscribe(eventStoreCtx)

	t := new(transaction.Transaction)

	for tr := range transactions {
		if len(tr.Event.Data) > 0 {
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
			currentGold := goldmanager.GetCurrentGoldAmount()

			incomingGold := transaction.ConvertSumOfAmountToCoin(t.Amount)

			coins, err := transaction.Align(currentGold, incomingGold, t.Type)
			if err != nil {
				panic(err)
			}
			goldmanager.UpdateGoldAmount(coins)
			conf.CurrentConf.LastReadEvent++

			db.UpdateConfig(*conf.CurrentConf)
		}
	}
}
