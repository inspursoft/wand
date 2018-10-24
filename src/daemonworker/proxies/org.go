package proxies

import (
	"log"
	"net/http"
	"net/url"

	"github.com/inspursoft/wand/src/daemonworker/dao"
)

func (p *ProxiedHandler) CreateOrg(resp http.ResponseWriter, req *http.Request, data []byte) {
	log.Printf("Intercepting org creation: %s\n", req.URL.Path)
	form, _ := url.ParseQuery(string(data))
	orgName := form.Get("org_name")
	dao.AddOrUpdateUserAccess(orgName, "", 1)
}
