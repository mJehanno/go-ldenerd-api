//
package goldmanager

import (
	"github.com/arangodb/go-driver"
	"github.com/mjehanno/goldenerd/database"
)

//
func getCoinsCollection() *driver.Collection {
	db := *database.GetDb()

	if found, err := db.CollectionExists(database.DbContext, "coins"); err != nil {
		panic(err)
	} else if !found {
		db.CreateCollection(database.DbContext, "coins", nil)
	}

	col, err := db.Collection(database.DbContext, "coins")

	if err != nil {
		panic(err)
	}

	return &col
}

//Return current amount of Gold
func GetCurrentGoldAmount() Coin {
	db := *database.GetDb()
	col := *getCoinsCollection()
	var coin Coin

	if len, err := col.Count(database.DbContext); err != nil {
		panic(err)
	} else if len == 0 {
		col.CreateDocument(database.DbContext, Coin{})
	}

	query := "FOR d IN coins LIMIT 1 RETURN d"
	cursor, err := db.Query(database.DbContext, query, nil)
	if err != nil {
		panic(err)
	}

	defer cursor.Close()
	for {

		meta, err := cursor.ReadDocument(database.DbContext, &coin)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}
		coin.Id = meta.Key
	}

	return coin
}

//Update current amount of gold
func UpdateGoldAmount(c Coin) {
	col := *getCoinsCollection()

	_, err := col.UpdateDocument(database.DbContext, c.Id, c)

	if err != nil {
		panic(err)
	}
}
