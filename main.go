package main

import (
	"s2p-api/config"
	"s2p-api/core"
	"s2p-api/database"
	"s2p-api/services"

	log "github.com/sirupsen/logrus"

	"net/http"
)

func main() {
	log.Infoln("      _____                        _ ")
	log.Infoln("     / __  \\                      (_)")
	log.Infoln(" ___ `' / /'_ __ ______ __ _ _ __  _ ")
	log.Infoln("/ __|  / / | '_ \\______/ _` | '_ \\| |")
	log.Infoln("\\__ \\./ /__| |_) |    | (_| | |_) | |")
	log.Infoln("|___/\\_____/ .__/      \\__,_| .__/|_|")
	log.Infoln("           | |              | |      ")
	log.Infoln("           |_|              |_|      \n\n")

	log.Infoln("[Configuration]: Getting variables")

	if err := config.Load(); err != nil {
		log.Fatalln("[Configuration]: Could not possible to get the configuration variables", err)
	}

	services.Initialize()

	database.Connect()

	if err := core.Initialize(); err != nil {
		log.Fatalln("[Schema-Manager]: Could not possible to create the schema", err)
	} else {
		log.Infoln("[Schema-Manager]: Schema was generated successfully")
	}
	address := config.HTTP.Address + ":" + config.HTTP.Port
	log.Infoln("Starting server on:", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Errorln(err)
	}
}
