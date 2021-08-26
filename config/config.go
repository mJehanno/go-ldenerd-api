package config

type Config struct {
	Id                string
	LastReadEvent     int
	EventstoreHost    string
	EvenstoreDb       string
	EventstoreUser    string
	EvenstorePassword string
	ArangoDb          string
	ArangoUser        string
	ArangoPassword    string
	ArangoHost        string
	KeycloakHost      string
}
