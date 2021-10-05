package conf

import "fmt"

var CurrentConf *Config = new(Config)

// hold a reference to the config stored in database
var AppConfig Config

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
	KeycloakRealm     string
	KeycloakSecret    string
	KeycloakClientID  string
}

func (c Config) String() string {
	return fmt.Sprintf(`{
		LastReadEvent : %v,
		EventStoreHost : %v,
		ArangoHost : %v,
		ArangoDb : %v,
		KeycloakHost : %v,
		KeycloakRealm: %v,
		}`,
		c.LastReadEvent,
		c.EventstoreHost,
		c.ArangoHost,
		c.ArangoDb,
		c.KeycloakHost,
		c.KeycloakRealm,
	)
}
