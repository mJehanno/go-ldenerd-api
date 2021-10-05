package env

import (
	"os"

	"github.com/mjehanno/go-ldenerd-api/appconfig/conf"
)

// Returns a config ojbect filled with env variables
func GetConfigFromEnv() {
	conf.CurrentConf.KeycloakRealm = os.Getenv("GOLDENER_KEYCLOAK_REALM")
	conf.CurrentConf.KeycloakClientID = os.Getenv("GOLDENER_KEYCLOAK_CLIENTID")
	conf.CurrentConf.KeycloakSecret = os.Getenv("GOLDENER_KEYCLOAK_SECRET")
	conf.CurrentConf.KeycloakHost = os.Getenv("GOLDENER_KEYCLOAK_HOST")
	conf.CurrentConf.ArangoHost = os.Getenv("GOLDENER_ARANGO_HOST")
	conf.CurrentConf.ArangoDb = os.Getenv("GOLDENER_ARANGO_DB")
	conf.CurrentConf.ArangoUser = os.Getenv("GOLDENER_ARANGO_USER")
	conf.CurrentConf.ArangoPassword = os.Getenv("GOLDENER_ARANGO_PASSWORD")
	conf.CurrentConf.EventstoreHost = os.Getenv("GOLDENER_EVENTSTORE_HOST")
	conf.CurrentConf.EvenstoreDb = os.Getenv("GOLDENER_EVENTSTORE_STREAM")
	conf.CurrentConf.EventstoreUser = os.Getenv("GOLDENER_EVENTSTORE_USER")
	conf.CurrentConf.EvenstorePassword = os.Getenv("GOLDENER_EVENTSTORE_PASSWORD")
}
