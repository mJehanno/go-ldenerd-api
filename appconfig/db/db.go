package db

import (
	"github.com/arangodb/go-driver"
	"github.com/mjehanno/go-ldenerd-api/appconfig/conf"
	"github.com/mjehanno/go-ldenerd-api/database"
)

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
func GetConfigFromDb() *conf.Config {
	db := *database.GetDb()
	col := *getConfigCollection()

	if len, err := col.Count(database.DbContext); err != nil {
		panic(err)
	} else if len == 0 {
		col.CreateDocument(database.DbContext, conf.Config{LastReadEvent: 0})
	}

	query := "FOR d IN config LIMIT 1 RETURN d"

	cursor, err := db.Query(database.DbContext, query, nil)
	if err != nil {
		panic(err)
	}

	defer cursor.Close()
	for {
		meta, err := cursor.ReadDocument(database.DbContext, &conf.AppConfig)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}
		conf.AppConfig.Id = meta.Key
	}

	return &conf.AppConfig
}

// Update stored config with given object
func UpdateConfig(c conf.Config) {
	col := *getConfigCollection()

	_, err := col.UpdateDocument(database.DbContext, c.Id, c)

	if err != nil {
		panic(err)
	}
}
