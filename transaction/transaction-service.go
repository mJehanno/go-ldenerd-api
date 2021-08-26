package transaction

import (
	"sort"
	"strconv"

	"github.com/arangodb/go-driver"
	"github.com/mjehanno/goldenerd/database"
)

func getTransactionCollection() *driver.Collection {
	db := *database.GetDb()

	if found, err := db.CollectionExists(database.DbContext, "transactions"); err != nil {
		panic(err)
	} else if !found {
		db.CreateCollection(database.DbContext, "transactions", nil)
	}

	col, err := db.Collection(database.DbContext, "transactions")

	if err != nil {
		panic(err)
	}

	return &col
}

func GetAllTransactionHistory() []Transaction {
	transactions := []Transaction{}
	db := *database.GetDb()
	query := "FOR d IN transactions RETURN d"
	cursor, err := db.Query(database.DbContext, query, nil)
	if err != nil {
		panic(err)
	}

	defer cursor.Close()

	for {
		var t Transaction
		meta, err := cursor.ReadDocument(database.DbContext, &t)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			panic(err)
		}
		t.Id = meta.Key

		transactions = append(transactions, t)
	}

	sort.Slice(transactions, func(i, j int) bool {
		a, _ := strconv.Atoi(transactions[i].Id)
		b, _ := strconv.Atoi(transactions[j].Id)
		return a > b
	})
	return transactions
}

func AddTransaction(t Transaction) {
	col := *getTransactionCollection()

	col.CreateDocument(database.DbContext, t)
}
