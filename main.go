package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/golang-jwt/jwt"
	"github.com/mjehanno/goldenerd/auth"
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

	currentConf := config.GetConfigFromDb()
	currentConf.KeycloakRealm = "Goldener"
	currentConf.KeycloakClientID = "goldener"
	currentConf.KeycloakSecret = "3bb43fbe-e07e-4243-a53a-1d6bcdb1f3a7"

	fmt.Println("lastevent : ", currentConf.LastReadEvent)
	// Eventstore
	client := goro.Connect("http://127.0.0.1:2113", goro.WithBasicAuth("admin", ""))
	catchupSubscription := client.CatchupSubscription("transactions", int64(currentConf.LastReadEvent))
	writer := client.Writer("transactions")
	eventStoreCtx := context.Background()
	keycloackCtx := context.Background()

	go func() {
		transactions := catchupSubscription.Subscribe(eventStoreCtx)

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

	api.Post("/login", func(c *fiber.Ctx) error {
		var u auth.User
		if err := c.BodyParser(&u); err != nil {
			fmt.Println("can't convert body")
			return err
		}
		token, err := auth.GetClient().Login(keycloackCtx, currentConf.KeycloakClientID, currentConf.KeycloakSecret, currentConf.KeycloakRealm, u.Username, u.Password)
		if err != nil {
			log.Println(fmt.Errorf("error happened while login : %s", err))
		}
		return c.JSON(auth.Jwt{AccessToken: token.AccessToken, ExpiresIn: token.ExpiresIn, RefreshToken: token.RefreshToken})
	})

	api.Post("/refresh", func(c *fiber.Ctx) error {
		head := c.Get("Authorization")
		head = strings.Replace(head, "Bearer ", "", -1)
		headToken, _ := jwt.Parse(head, func(token *jwt.Token) (interface{}, error) {
			return token.Claims, nil
		})

		result, err := auth.GetClient().RetrospectToken(keycloackCtx, headToken.Raw, currentConf.KeycloakClientID, currentConf.KeycloakSecret, currentConf.KeycloakRealm)
		if err != nil || !*result.Active {
			return c.SendStatus(401)
		}

		var t auth.Jwt
		if err := c.BodyParser(&t); err != nil {
			fmt.Println("can't convert body")
			return err
		}

		token, _ := auth.GetClient().RefreshToken(keycloackCtx, t.RefreshToken, currentConf.KeycloakClientID, currentConf.KeycloakSecret, currentConf.KeycloakRealm)

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

	/**
	*	TODO: Protect this route with keycloak.
	 */
	tr.Post("/", func(c *fiber.Ctx) error {
		head := c.Get("Authorization")
		head = strings.Replace(head, "Bearer ", "", -1)
		token, _ := jwt.Parse(head, func(token *jwt.Token) (interface{}, error) {
			return token.Claims, nil
		})

		result, err := auth.GetClient().RetrospectToken(keycloackCtx, token.Raw, currentConf.KeycloakClientID, currentConf.KeycloakSecret, currentConf.KeycloakRealm)
		if err != nil || !*result.Active {
			return c.SendStatus(401)
		}

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
		err = writer.Write(eventStoreCtx, goro.ExpectedVersionAny, event)
		if err != nil {
			fmt.Println("can't write in eventstore")
			return (err)
		}
		return c.SendStatus(201)
	})

	app.Listen(":8080")
}
