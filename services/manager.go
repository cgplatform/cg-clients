package services

import (
	"s2p-api/services/mail"

	log "github.com/sirupsen/logrus"
)

func Initialize() {
	log.Infoln("[Services]: Starting services")
	mail.Initialize()
}
