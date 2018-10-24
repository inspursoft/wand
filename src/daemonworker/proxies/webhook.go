package proxies

import (
	"fmt"
	"log"
	"net/http"

	"github.com/inspursoft/wand/src/daemonworker/dao"
	"github.com/inspursoft/wand/src/daemonworker/models"
	"github.com/inspursoft/wand/src/daemonworker/utils"
)

var jsonHeader = http.Header{
	"content-type": []string{"application/json"},
}

func (p *ProxiedHandler) Webhook(resp http.ResponseWriter, req *http.Request, data []byte) {
	event := req.Header.Get("X-Gogs-Event")
	log.Printf("Intercepting event: %s webhook ...\n", event)
	var payload models.CustomWebhookPayload
	payload.GogsURL = p.GogitsBaseURL
	payload.APIURL = fmt.Sprintf("%s/api/v1", p.GogitsBaseURL)
	payload.MasterURL = p.JenkinsMasterURL
	payload.NodeIP = p.JenkinsNodeIP
	payload.RegistryURL = p.KVMRegistryURL
	trigger, err := payload.AdaptToCustom(event, data)
	if err != nil {
		log.Printf("Failed to convert payload with event: %s, error: %+v\n", event, err)
	}
	if trigger {
		dao.NewBuildConfig(payload.RepoName, payload.Username).SetPayload(&payload).Update()
		log.Printf("raw: %+v\n\ncustom: %+v\n", string(data), utils.PrettyPrintJSON(payload))
		utils.SimplePostRequestHandle(fmt.Sprintf("%s/generic-webhook-trigger/invoke", p.JenkinsMasterURL), jsonHeader, payload)
	}
}
