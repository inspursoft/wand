package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/inspursoft/wand/src/daemonworker/utils"
)

func (c *Handler) ResolveIcon(resp http.ResponseWriter, req *http.Request) {
	iconName := req.FormValue("name")
	if strings.TrimSpace(iconName) == "" {
		rendStatus(resp, http.StatusBadRequest, fmt.Sprintln("No icon name provided."))
		return
	}
	err := utils.DrawTag(resp, iconName)
	if err != nil {
		log.Printf("Failed to draw tag with icon name: %s\n", iconName)
		utils.DrawText(resp, "N/A")
	}
}
