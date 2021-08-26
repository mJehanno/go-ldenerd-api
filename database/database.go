//Package database manage everything related to database
package database

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

var currentConnection *driver.Client = nil
var DbContext = context.Background()

//Returns a database opened connection
func getConnexion() *driver.Client {
	if currentConnection == nil {
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{"http://localhost:8529"},
		})
		if err != nil {
			panic(err)
		}
		c, err := driver.NewClient(driver.ClientConfig{
			Connection: conn,
		})
		if err != nil {
			panic(err)
		}
		currentConnection = &c
	}
	return currentConnection
}

//Returns the app database
func GetDb() *driver.Database {
	conn := *getConnexion()

	db, err := conn.Database(DbContext, "goldener")

	if err != nil {
		panic(err)
	}

	return &db
}
