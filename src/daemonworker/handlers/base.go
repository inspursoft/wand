package handlers

import (
	"log"
	"net/http"

	"github.com/inspursoft/wand/src/daemonworker/models"
)

const (
	uploadResourcePath = "/root/website"
)

type Handler struct {
	Cache *models.CachedStore
}

func rendStatus(resp http.ResponseWriter, statusCode int, message string) {
	log.Printf("%s", message)
	resp.WriteHeader(statusCode)
	resp.Write([]byte(message))
}
