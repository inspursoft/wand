package proxies

import (
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/inspursoft/wand/src/daemonworker/dao"
	"github.com/inspursoft/wand/src/daemonworker/gogs"
)

func (p *ProxiedHandler) SignUp(resp http.ResponseWriter, req *http.Request, data []byte) {
	log.Printf("Intercepting user login: %s\n", req.URL.Path)
	form, _ := url.ParseQuery(string(data))
	username := form.Get("user_name")
	password := form.Get("password")
	time.Sleep(time.Second * 2)
	accessToken, err := gogs.NewGogsHandler(p.GogitsBaseURL, username, "").CreateAccessToken(password)
	if err != nil {
		log.Printf("Failed to get access token: %+v\n", err)
		return
	}
	dao.AddOrUpdateUserAccess(username, accessToken.Sha1, 0)
	log.Printf("Created access token with username: %s, access token: %s\n", username, accessToken.Sha1)
}
