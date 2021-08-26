//Package config manage the app configuration
package config

import (
	"github.com/arangodb/go-driver"
	"github.com/mjehanno/goldenerd/database"
)

var appConfig Config

// Returns db's config collection
func getConfigCollection() *driver.Collection {
	db := *database.GetDb()
	if found, err := db.CollectionExists(database.DbContext, "config"); err != nil {
		panic(err)
	} else if !found {
		db.CreateCollection(database.DbContext, "config", nil)
	}

	col, err := db.Collection(database.DbContext, "config")
	if err != nil {
		panic(err)
	}

	return &col
}

// Returns a config object with stored config
func GetConfigFromDb() Config {
	db := *database.GetDb()
	col := *getConfigCollection()

	if len, err := col.Count(database.DbContext); err != nil {
		panic(err)
	} else if len == 0 {
		col.CreateDocument(database.DbContext, Config{LastReadEvent: 0})
	}

	query := "FOR d IN config LIMIT 1 RETURN d"

	cursor, err := db.Query(database.DbContext, query, nil)
	if err != nil {
		panic(err)
	}

	defer cursor.Close()
	for {
		meta, err := cursor.ReadDocument(database.DbContext, &appConfig)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}
		appConfig.Id = meta.Key
	}

	return appConfig
}

// Returns a config ojbect filled with env variables
func GetConfigFromEnv() Config {
	return appConfig
}

// Update stored config with given object
func UpdateConfig(c Config) {
	col := *getConfigCollection()

	_, err := col.UpdateDocument(database.DbContext, c.Id, c)

	if err != nil {
		panic(err)
	}
}
