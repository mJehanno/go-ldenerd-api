//Package database manage everything related to database
package database

import (
	"context"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/mjehanno/go-ldenerd-api/appconfig/conf"
)

var currentConnection *driver.Client = nil
var DbContext = context.Background()

//Returns a database opened connection
func getConnexion() *driver.Client {
	if currentConnection == nil {
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{conf.CurrentConf.ArangoHost},
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

	db, err := conn.Database(DbContext, conf.CurrentConf.ArangoDb)

	if err != nil {
		panic(err)
	}

	return &db
}
